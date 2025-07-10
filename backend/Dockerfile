# Use a Node.js Debian base image so we can install Go tools
FROM node:20-bullseye AS final

# Install essential system tools and Go build dependencies
RUN apt-get update && apt-get install -y \
    bash git curl tar gcc g++ pkg-config nano tree

# Define environment variables for Go installation
ENV GO_VERSION=1.24.4
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV GOBIN=/usr/local/bin
ENV PATH=$GOBIN:$GOROOT/bin:$GOPATH/bin:$PATH

# Download and install Go manually
RUN curl -LO https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# Install redoc-cli globally using npm (for generating HTML API docs)
RUN npm install -g redoc-cli

# Set the working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY backend/go.mod backend/go.sum ./src/
WORKDIR /app/src
RUN go mod download all

# Install Swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy the rest of the project
COPY /backend ./

# Generate swagger.yaml & docs.go
RUN swag init \
    --generalInfo ./main.go \
    --dir cmd,pkg/handler \
    --exclude pkg/config,docs \
    --output ../docs \
    --parseDependency
# Generate Swagger html if ReDoc
RUN redoc-cli bundle ../docs/swagger.yaml -o ../docs/docs.html

# Build application binaries
RUN go build -o /app/backend ./cmd/

# Set back to /app for CMD
WORKDIR /app

# Run the built binary
CMD ["./backend"]  