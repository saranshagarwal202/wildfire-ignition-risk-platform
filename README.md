# wildfire-ignition-risk-platform

wildfire-risk-platform/
├── README.md
├── LICENSE
├── .gitignore
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.work
├── go.work.sum
│
├── api/
│   └── proto/
│       ├── services.proto
│       ├── generate.go
│       └── generated/
│           └── (generated protobuf files will go here)
│
├── services/
│   ├── api-gateway/
│   │   ├── main.go
│   │   ├── handlers/
│   │   │   └── risk_assessment.go
│   │   ├── middleware/
│   │   │   └── cors.go
│   │   ├── config/
│   │   │   └── config.go
│   │   └── Dockerfile
│   │
│   ├── orchestrator/
│   │   ├── main.go
│   │   ├── server/
│   │   │   └── orchestrator.go
│   │   ├── database/
│   │   │   └── queries.go
│   │   ├── kafka/
│   │   │   └── producer.go
│   │   ├── config/
│   │   │   └── config.go
│   │   └── Dockerfile
│   │
│   ├── infrastructure/
│   │   ├── main.go
│   │   ├── server/
│   │   │   └── infrastructure.go
│   │   ├── osm/
│   │   │   └── overpass.go
│   │   ├── config/
│   │   │   └── config.go
│   │   └── Dockerfile
│   │
│   ├── topography/
│   │   ├── main.go
│   │   ├── server/
│   │   │   └── topography.go
│   │   ├── usgs/
│   │   │   └── dem_downloader.go
│   │   ├── config/
│   │   │   └── config.go
│   │   └── Dockerfile
│   │
│   ├── downloader/
│   │   ├── main.go
│   │   ├── kafka/
│   │   │   ├── consumer.go
│   │   │   └── producer.go
│   │   ├── satellite/
│   │   │   └── landsat.go
│   │   ├── config/
│   │   │   └── config.go
│   │   └── Dockerfile
│   │
│   ├── spark-processor/
│   │   ├── main.py
│   │   ├── processor/
│   │   │   ├── __init__.py
│   │   │   ├── raster_ops.py
│   │   │   ├── risk_calculator.py
│   │   │   └── spatial_join.py
│   │   ├── kafka_handler/
│   │   │   ├── __init__.py
│   │   │   ├── consumer.py
│   │   │   └── producer.py
│   │   ├── weather/
│   │   │   ├── __init__.py
│   │   │   └── openweather.py
│   │   ├── config/
│   │   │   ├── __init__.py
│   │   │   └── config.py
│   │   ├── requirements.txt
│   │   └── Dockerfile
│   │
│   └── ingestion/
│       ├── main.go
│       ├── kafka/
│       │   └── consumer.go
│       ├── database/
│       │   └── writer.go
│       ├── config/
│       │   └── config.go
│       └── Dockerfile
│
├── shared/
│   ├── database/
│   │   ├── connection.go
│   │   ├── migrations/
│   │   │   └── 001_initial_schema.sql
│   │   └── models/
│   │       └── job.go
│   ├── kafka/
│   │   ├── client.go
│   │   └── topics.go
│   ├── config/
│   │   └── env.go
│   └── utils/
│       ├── geojson.go
│       └── logger.go
│
├── data/
│   ├── dem/
│   │   └── .gitkeep
│   ├── landsat/
│   │   └── .gitkeep
│   └── processed/
│       └── .gitkeep
│
├── scripts/
│   ├── setup-db.sh
│   ├── kafka-topics.sh
│   └── test-api.sh
│
├── deployments/
│   ├── docker/
│   │   ├── postgres.yml
│   │   ├── kafka.yml
│   │   └── spark.yml
│   └── k8s/
│       └── (kubernetes manifests for production)
│
├── docs/
│   ├── api/
│   │   └── postman_collection.json
│   ├── architecture/
│   │   └── diagrams/
│   └── setup/
│       └── development.md
│
└── tests/
    ├── integration/
    │   └── api_test.go
    ├── unit/
    │   └── (unit test files)
    └── testdata/
        └── sample_aoi.geojson