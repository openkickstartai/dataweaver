# DataWeaver

Intelligent data integration and harmonization tool that automates schema matching, real-time transformation, and provides visual workflow design for diverse data sources.

## Features

- **Automatic Schema Detection**: Intelligently detects and analyzes data schemas from JSON, CSV, and database sources
- **Real-time Data Transformation**: High-performance transformation engine with configurable workflows
- **Visual Workflow Designer**: Web-based interface for creating and managing data transformation pipelines
- **Multi-source Connectors**: Built-in support for APIs, databases, and file formats
- **Data Validation**: Comprehensive data quality checks and validation rules
- **RESTful API**: Complete API for programmatic integration

## Tech Stack

- **Backend**: Go with Gin framework
- **Frontend**: JavaScript (React-based workflow designer)
- **Database**: PostgreSQL for metadata and workflow storage
- **API**: RESTful endpoints with JSON responses

## Quick Start

1. **Install Dependencies**
   ```bash
   go mod tidy
   ```

2. **Set Environment Variables**
   ```bash
   export DATABASE_URL="postgres://user:password@localhost/dataweaver?sslmode=disable"
   export PORT="8080"
   ```

3. **Run the Server**
   ```bash
   go run main.go
   ```

## API Endpoints

### Schema Detection
```bash
POST /api/v1/schemas/detect
{
  "data": "{\"name\": \"John\", \"age\": 30}",
  "format": "json"
}
```

### Workflow Management
```bash
# Create workflow
POST /api/v1/workflows
{
  "name": "User Data Transform",
  "description": "Transform user data format",
  "steps": [
    {
      "id": "step1",
      "type": "transform",
      "config": {
        "field_mapping": {
          "name": "full_name",
          "age": "user_age"
        }
      }
    }
  ]
}

# Execute transformation
POST /api/v1/transform
{
  "workflow_id": 1,
  "data": {
    "name": "John Doe",
    "age": 30
  }
}
```

## Architecture

- `main.go`: Application entry point and HTTP server setup
- `internal/api/`: REST API handlers and routing
- `internal/schema/`: Schema detection and analysis engine
- `internal/workflow/`: Workflow execution engine
- `internal/database/`: Database connection and models
- `internal/config/`: Configuration management

## Development

The project follows Go best practices with a clean architecture pattern. Each package has a specific responsibility and can be tested independently.

## License

MIT License