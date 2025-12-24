# Go Upload Service

File upload microservice using Cloudflare R2 with Go.

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Configure `config.yaml` with your R2 credentials

3. Run the server:
```bash
go run main.go
```

## Build

```bash
go build -o go-upload
./go-upload
```

## API Endpoints

- `GET /health` - Health check
- `POST /api/upload/` - Upload file
- `DELETE /api/upload/` - Delete file

## Project Structure

```
go-upload/
├── config/          # Configuration management
├── controllers/     # HTTP handlers
├── services/        # Business logic
├── routes/          # Route definitions
├── utils/           # Utility functions (R2)
├── models/          # Data models
└── main.go          # Application entry point
```
