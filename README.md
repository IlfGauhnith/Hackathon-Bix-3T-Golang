# Hackathon Bix 3T Golang

<details>
<summary><strong> Challenge </strong></summary>

### Description

This project addresses the data discrepancy issues faced by a major retail network by enabling the comparison of a large internal CSV file (up to 500,000 rows) with data retrieved from an external API. It focuses on performance, scalability, and delivering reliable results to support accurate business decision-making.

### Features

- Upload and validation of CSV files
- Data retrieval from an external API
- Efficient comparison between data sources
- Exposure of results via REST API
- Export of results in CSV and JSON formats

</details>

### Env
- ```{{BACKEND_PORT}}``` Backend port to listen to 
### Local Execution

```bash
# Repo clone
git clone https://github.com/IlfGauhnith/Hackathon-Bix-3T-Golang.git

cd Hackathon-Bix-3T-Golang

# Build and execute using docker compose
docker-compose up --build
```

### API Documentation
The full OpenAPI (Swagger) specification is available at ```http://localhost:{{PORT}}/docs/openapi```