# 🔥 Wildfire Ignition Risk & Impact Platform

**Status: 🚧 Active Development**

A cloud-native, event-driven platform for wildfire risk assessment that analyzes geographical areas using satellite imagery, topographical data, infrastructure mapping, and real-time weather conditions to generate comprehensive risk scores.

## 🎯 Project Overview

This platform provides on-demand wildfire risk analysis for user-defined geographical areas. It combines multiple data sources to create weighted risk scores that help identify areas and assets most vulnerable to wildfire ignition and spread.

### Key Features
- **Real-time Risk Assessment**: Analyze any geographical area on-demand
- **Multi-source Data Integration**: Satellite imagery, elevation data, infrastructure mapping, and weather APIs
- **Scalable Microservice Architecture**: Built with Go and Python microservices
- **Event-driven Processing**: Asynchronous job processing with Apache Kafka
- **Geospatial Analysis**: PostGIS database for spatial queries and analytics
- **Interactive Visualization**: Metabase dashboard for risk visualization

## 🏗️ Architecture

The system follows an event-driven microservice architecture with the following components:

```
User Request → API Gateway → Orchestrator → Data Services → Processing → Results
                     ↓              ↓              ↓           ↓
                  REST API     gRPC Services   Kafka Queue  PostGIS DB
```

### Data Flow
1. **User submits** a GeoJSON polygon via REST API
2. **API Gateway** validates and forwards to Orchestrator
3. **Orchestrator** coordinates data gathering from multiple services
4. **Infrastructure Service** fetches building/road data from OpenStreetMap
5. **Topography Service** downloads elevation data from USGS
6. **Downloader Service** obtains satellite imagery for vegetation analysis
7. **Spark Processor** performs geospatial analysis and risk calculations
8. **Ingestion Service** stores results in PostGIS database
9. **Metabase** provides visualization dashboard

## 🛠️ Technology Stack

### Backend Services
- **Go**: API Gateway, Orchestrator, Infrastructure, Topography, Downloader, Ingestion services
- **Python + PySpark**: Geospatial processing and risk calculations
- **gRPC**: Inter-service communication
- **Protocol Buffers**: Service contracts and message definitions

### Data & Messaging
- **PostgreSQL + PostGIS**: Geospatial database
- **Apache Kafka**: Event streaming and message queuing
- **Docker**: Containerization
- **Docker Compose**: Local development orchestration

### External Data Sources
- **USGS EarthExplorer**: Landsat satellite imagery
- **USGS 3DEP**: Digital elevation models
- **OpenStreetMap**: Infrastructure and asset data
- **OpenWeatherMap API**: Real-time weather conditions

## 📁 Project Structure

```
wildfire-ignition-risk-platform/
├── README.md                    # This file
├── docker-compose.yml          # Local development setup
├── go.mod                      # Go module definition
├── go.work                     # Go workspace configuration
│
├── api/                        # API definitions
│   └── proto/                  # Protocol Buffer definitions
│       ├── services.proto      # gRPC service contracts
│       └── generated/          # Generated protobuf code
│
├── services/                   # Microservices
│   ├── api-gateway/           # ✅ HTTP REST API entry point
│   ├── orchestrator/          # ✅ Job coordination service
│   ├── infrastructure/        # ✅ OpenStreetMap data fetcher
│   ├── topography/           # 🚧 USGS elevation data service
│   ├── downloader/           # 🚧 Satellite imagery downloader
│   ├── spark-processor/      # 🚧 Geospatial analysis engine
│   └── ingestion/            # 🚧 Database writer service
│
├── shared/                    # Shared utilities and models
│   ├── database/             # Database connection and models
│   ├── kafka/               # Kafka client utilities
│   ├── config/              # Configuration management
│   └── utils/               # Common utilities (logging, GeoJSON)
│
├── data/                     # Data storage directories
│   ├── dem/                 # Digital elevation models
│   ├── landsat/            # Satellite imagery
│   └── processed/          # Processed results
│
├── scripts/                 # Setup and utility scripts
├── docs/                   # Documentation and API specs
└── tests/                  # Test files
```

## 🚦 Current Status

### ✅ Completed Components
- **API Gateway**: HTTP REST endpoint for job submission
- **Orchestrator Service**: Job coordination and gRPC orchestration
- **Infrastructure Service**: OpenStreetMap data integration via Overpass API
- **Protocol Buffers**: Service contracts and message definitions
- **Database Schema**: PostGIS tables and spatial indexes
- **Shared Libraries**: Database connections, Kafka clients, utilities

### 🚧 In Progress
- **Topography Service**: USGS elevation data integration
- **Downloader Service**: Satellite imagery acquisition
- **Spark Processor**: Core geospatial analysis engine
- **Ingestion Service**: Database result writer
- **Docker Compose**: Complete development environment

### 📋 TODO List

#### High Priority
- [ ] Complete Topography Service implementation
- [ ] Implement Downloader Service for Landsat imagery
- [ ] Build Spark Processor for risk calculations
- [ ] Create Ingestion Service for result storage
- [ ] Set up complete Docker Compose environment
- [ ] Integrate all services with Kafka messaging

#### Medium Priority
- [ ] Add comprehensive error handling and logging
- [ ] Implement job status tracking and updates
- [ ] Add input validation and sanitization
- [ ] Create health check endpoints for all services
- [ ] Set up Metabase visualization dashboard
- [ ] Add API rate limiting and authentication

#### Low Priority
- [ ] Write comprehensive unit tests
- [ ] Add integration tests
- [ ] Create API documentation with Swagger
- [ ] Implement monitoring and metrics
- [ ] Add configuration management
- [ ] Create deployment scripts for production

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- Docker and Docker Compose
- Python 3.9+ (for Spark processor)
- OpenWeatherMap API key

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/wildfire-ignition-risk-platform.git
   cd wildfire-ignition-risk-platform
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Generate Protocol Buffer code**
   ```bash
   make proto-gen
   ```

4. **Start the development environment**
   ```bash
   docker-compose up -d
   ```

5. **Run database migrations**
   ```bash
   make migrate-up
   ```

### API Usage

Submit a wildfire risk assessment job:

```bash
curl -X POST http://localhost:8080/v1/wildfire-risk-jobs \
  -H "Content-Type: application/json" \
  -d '{
    "aoi": {
      "type": "Polygon",
      "coordinates": [[
        [-122.5, 37.8],
        [-122.4, 37.8],
        [-122.4, 37.7],
        [-122.5, 37.7],
        [-122.5, 37.8]
      ]]
    }
  }'
```

Response:
```json
{
  "job_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "status": "PENDING",
  "message": "Risk assessment job accepted. Use the job_id to track status."
}
```

## 🔧 Development Guide

### Understanding the Go Code Structure

#### 1. **Main Function (`main.go`)**
```go
func main() {
    // Load configuration
    config.LoadEnv()
    
    // Initialize database connection
    db, err := database.NewConnection()
    
    // Start HTTP server
    http.ListenAndServe(":8080", handler)
}
```
- Entry point of each microservice
- Loads environment variables
- Initializes dependencies (database, Kafka)
- Starts the service (HTTP server or gRPC server)

#### 2. **gRPC Services**
```go
type server struct {
    pb.UnimplementedOrchestratorServiceServer
    db *database.DB
}

func (s *server) CreateRiskAssessmentJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.CreateJobResponse, error) {
    // Business logic here
}
```
- Implements the Protocol Buffer service interface
- Contains business logic for each RPC method
- Uses dependency injection for database, Kafka clients

#### 3. **Database Models**
```go
type RiskAssessmentJob struct {
    JobID      uuid.UUID `json:"job_id" db:"job_id"`
    AOIPolygon string    `json:"aoi_polygon" db:"aoi_polygon"`
    JobStatus  JobStatus `json:"job_status" db:"job_status"`
}
```
- Struct tags define JSON serialization and database column mapping
- Uses UUID for unique job identification
- JobStatus is a custom type for type safety

#### 4. **Kafka Integration**
```go
writer := kafka.NewWriter(kafka.WriterConfig{
    Brokers: []string{"localhost:9092"},
    Topic:   "download.tasks",
})

err := writer.WriteMessages(context.Background(),
    kafka.Message{
        Key:   []byte(jobID),
        Value: messageBytes,
    },
)
```
- Kafka writers publish messages to topics
- Kafka readers consume messages from topics
- Messages are serialized as JSON

### Service Communication Flow

1. **REST → gRPC**: API Gateway converts HTTP requests to gRPC calls
2. **gRPC → Database**: Services query PostGIS for data storage/retrieval
3. **gRPC → Kafka**: Services publish messages for asynchronous processing
4. **Kafka → Processing**: Background services consume and process messages

## 📊 Risk Calculation Formula

The platform calculates wildfire risk using a weighted formula:

```
Overall Risk = (NDVI × 0.5) + (Slope × 0.4) + (Wind Factor × 0.1)
```

Where:
- **NDVI** (Normalized Difference Vegetation Index): Vegetation density from satellite imagery
- **Slope**: Terrain steepness from elevation data
- **Wind Factor**: Wind speed and direction from weather API

Each component is normalized to a 0.0-1.0 scale before applying weights.

## 🤝 Contributing

This is an active development project. Contributions are welcome!

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 External Resources

- [Protocol Buffers Documentation](https://developers.google.com/protocol-buffers)
- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)
- [Apache Kafka with Go](https://kafka.apache.org/documentation/)
- [PostGIS Documentation](https://postgis.net/documentation/)
- [USGS EarthExplorer](https://earthexplorer.usgs.gov/)
- [OpenStreetMap Overpass API](https://overpass-api.de/)
