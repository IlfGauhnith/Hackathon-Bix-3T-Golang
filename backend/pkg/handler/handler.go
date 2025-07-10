package handler

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/pkg/config"
	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/pkg/logger"
	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/pkg/model"
	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/pkg/repository"
	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/pkg/service"
	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
)

// HealthHandler godoc
// @Summary      Check API health
// @Description  Returns a simple "API is running" status.
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.HealthResponse
// @Router       /health [get]
func HealthHandler(c *gin.Context) {
	logger.Log.Info("HealthHandler")
	c.JSON(http.StatusOK, gin.H{"status": "API is running"})
}

// UploadHandler godoc
// @Summary      Upload CSV file and detect data divergences
// @Description  Accepts a multipart‐form CSV file, optional X-Batch-Size header overrides the default batch size.
// @Tags         CSV
// @Accept       multipart/form-data
// @Produce      application/json
// @Param        X-Batch-Size  header  int  false  "Override batch size for CSV splitting"
// @Param        file          formData  file  true  "CSV file to upload and verify"
// @Success      200           {object}  map[string][]model.Divergence
// @Failure      400           {object}  gin.H  "Bad request"
// @Failure      500           {object}  gin.H  "Internal server error"
// @Router       /upload [post]
func UploadHandler(cfg *config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {
		defer util.TimeTrack(time.Now(), "UploadHandler")
		logger.Log.Info("UploadHandler")

		// Parse optional X-Batch-Size header (must be 1–1000)
		batchSize := cfg.BatchSize
		if hdr := c.GetHeader("X-Batch-Size"); hdr != "" {
			if v, err := strconv.Atoi(hdr); err != nil {
				logger.Log.Warnf("Invalid X-Batch-Size header %q: %v (using default %d)", hdr, err, cfg.BatchSize)
			} else if v < 1 || v > 1000 {
				logger.Log.Warnf("X-Batch-Size %d out of allowed range [1,1000] (using default %d)", v, cfg.BatchSize)
			} else {
				batchSize = v
				logger.Log.Infof("Overriding batch size to %d via header", batchSize)
			}
		}

		// Retrieve file from multipart form
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}
		f, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
			return
		}
		defer f.Close()

		// Unmarshal CSV rows into structs
		var rows []*model.CSVRecord
		if err := gocsv.Unmarshal(f, &rows); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid CSV format"})
			return
		}

		total := len(rows)
		totalBatches := (total + batchSize - 1) / batchSize

		divergencesChan := make(chan []model.Divergence, totalBatches)
		var wg sync.WaitGroup
		client := repository.NewExternalAPIClient(cfg)

		// limit to N concurrent workers
		// semaphore sized by config, defaulting to number of CPUs
		// goroutines is expected to do some I/O bound-work (network calls),
		// so quadriplicate the concurrency to allow bufferization
		sem := make(chan struct{}, cfg.MaxConcurrency*4)
		for i := 0; i < totalBatches; i++ {
			sem <- struct{}{} // Acquire a semaphore slot
			wg.Add(1)         // Increment wait group for this batch

			start := i * batchSize
			end := start + batchSize
			if end > total {
				end = total
			}
			batch := rows[start:end]
			pageNum := cfg.StartPage + i

			// In a production environment, I would move that goroutine to dedicated worker pool microservice
			// to avoid blocking the main API thread.
			go func(batch []*model.CSVRecord, page int) {

				defer wg.Done()
				defer func() { <-sem }() // release

				logger.Log.Infof("Processing batch %d/%d (records %d–%d)", i+1, totalBatches, start+1, end)

				// Fetch external products for this page
				apiResp, err := client.GetProducts(pageNum, batchSize)
				if err != nil {
					logger.Log.Errorf("Error fetching external products for batch %d: %v", i, err)
					return
				}

				// Convert pointers to values
				arr := make([]model.CSVRecord, len(batch))
				for i, r := range batch {
					arr[i] = *r
				}
				divs, err := service.CompareBatch(arr, apiResp)
				if err == nil {
					divergencesChan <- divs
				}
				if err != nil {
					logger.Log.Errorf("Error comparing batch %d: %v", i, err)
					return
				}
			}(batch, pageNum)
		}

		wg.Wait()
		close(divergencesChan)

		// Collect all divergences
		var all []model.Divergence
		for part := range divergencesChan {
			all = append(all, part...)
		}

		logger.Log.Infof("Completed all %d batches, found %d divergences", totalBatches, len(all))
		c.JSON(http.StatusOK, gin.H{"divergences": all})
	}
}

// UploadHandlerSequential godoc
// @Summary      Upload CSV file and detect data divergences (sequential)
// @Description  Accepts a multipart‐form CSV file, processes each batch one after another (no concurrency), compares records against the external products API, and returns any mismatches.
// @Tags         CSV
// @Accept       multipart/form-data
// @Produce      application/json
// @Param        X-Batch-Size  header  int  false  "Override batch size for CSV splitting"
// @Param        file  formData  file  true  "CSV file to upload and verify"
// @Success      200   {object}  map[string][]model.Divergence  "Key `divergences` with an array of divergence objects"
// @Failure      400   {object}  gin.H                         "Bad request (e.g. missing file or invalid CSV format)"
// @Failure      500   {object}  gin.H                         "Internal server error (e.g. external API failure)"
// @Router       /upload-seq [post]
func UploadHandlerSequential(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Measure total handler time
		defer util.TimeTrack(time.Now(), "UploadHandlerSequential")
		logger.Log.Info("UploadHandlerSequential started")

		// Parse optional X-Batch-Size header (must be 1–1000)
		batchSize := cfg.BatchSize
		if hdr := c.GetHeader("X-Batch-Size"); hdr != "" {
			if v, err := strconv.Atoi(hdr); err != nil {
				logger.Log.Warnf("Invalid X-Batch-Size header %q: %v (using default %d)", hdr, err, cfg.BatchSize)
			} else if v < 1 || v > 1000 {
				logger.Log.Warnf("X-Batch-Size %d out of allowed range [1,1000] (using default %d)", v, cfg.BatchSize)
			} else {
				batchSize = v
				logger.Log.Infof("Overriding batch size to %d via header", batchSize)
			}
		}

		// Retrieve file from multipart form
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}
		f, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
			return
		}
		defer f.Close()

		// Parse CSV into slice of pointers
		var rows []*model.CSVRecord
		if err := gocsv.Unmarshal(f, &rows); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid CSV format"})
			return
		}

		total := len(rows)
		totalBatches := (total + batchSize - 1) / batchSize

		client := repository.NewExternalAPIClient(cfg)
		var allDivergences []model.Divergence

		// Loop over each batch sequentially
		for i := 0; i < totalBatches; i++ {
			start := i * batchSize
			end := start + batchSize
			if end > total {
				end = total
			}
			batch := rows[start:end]
			pageNum := cfg.StartPage + i

			logger.Log.Infof("Processing batch %d/%d (records %d–%d)", i+1, totalBatches, start+1, end)

			// Fetch external products for this batch/page
			apiResp, err := client.GetProducts(pageNum, batchSize)
			if err != nil {
				logger.Log.Errorf("Failed to fetch API page %d: %v", pageNum, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "external API failure"})
				return
			}

			// Convert []*CSVRecord to []CSVRecord
			recValues := make([]model.CSVRecord, len(batch))
			for j, r := range batch {
				recValues[j] = *r
			}

			// Compare and collect divergences
			divs, err := service.CompareBatch(recValues, apiResp)
			if err != nil {
				logger.Log.Errorf("Comparison error on batch %d: %v", i+1, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "comparison failure"})
				return
			}
			allDivergences = append(allDivergences, divs...)
		}

		// Return all collected divergences
		logger.Log.Infof("Completed all %d batches, found %d divergences", totalBatches, len(allDivergences))
		c.JSON(http.StatusOK, gin.H{"divergences": allDivergences})
	}
}
