# Log Ingestion Service

This project is a data ingestion service that collects logs from a public API, processes them, and stores them in a cloud-native storage solution.

## Features

- Fetches data from JSONPlaceholder API
- Transforms data by adding metadata
- Stores data in MongoDB (cloud-native storage)
- Provides a REST API to retrieve ingested data
- Tracks the latest successful data ingestion
- Containerized with Docker
- Comprehensive test coverage

## Architecture

The service follows a modular architecture with the following components:

- **Fetcher**: Responsible for retrieving data from external APIs
- **Transformer**: Processes and enriches the raw data
- **Storage**: Handles data persistence
- **API**: Exposes endpoints to interact with the service
- **Tracker**: Monitors ingestion progress

## Getting Started

### Prerequisites

- Go 1.20+
- Docker and Docker Compose
- MongoDB (for local development without Docker)

### Setup Instructions

1. Clone the repository:

```bash
git clone https://github.com/tiwariayush700/log-ingestion-service.git
cd log-ingestion-service
```

2. Set up environment variables (or use the defaults in docker-compose.yml):

```bash
cp .env.example .env
```

3. Build and run with Docker Compose:

```bash
docker-compose up --build
```

### Running Tests

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## API Endpoints

- `GET /api/logs`: Retrieve all ingested logs
- `GET /api/logs/:id`: Retrieve a specific log by ID
- `GET /api/status`: Get the latest ingestion status

## Cloud Deployment

### AWS Deployment

To deploy to AWS ECS:

1. Create an ECR repository:

```bash
aws ecr create-repository --repository-name log-ingestion-service
```

2. Build and push the Docker image:

```bash
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <your-aws-account-id>.dkr.ecr.us-east-1.amazonaws.com
docker build -t <your-aws-account-id>.dkr.ecr.us-east-1.amazonaws.com/log-ingestion-service:latest .
docker push <your-aws-account-id>.dkr.ecr.us-east-1.amazonaws.com/log-ingestion-service:latest
```

3. Create an ECS cluster, task definition, and service using the AWS console or CLI.

### Google Cloud Run Deployment

To deploy to Google Cloud Run:

1. Build and push the Docker image:

```bash
gcloud auth configure-docker
docker build -t gcr.io/<your-project-id>/log-ingestion-service:latest .
docker push gcr.io/<your-project-id>/log-ingestion-service:latest
```

2. Deploy to Cloud Run:

```bash
gcloud run deploy log-ingestion-service \
  --image gcr.io/<your-project-id>/log-ingestion-service:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars="MONGO_URI=<your-mongodb-uri>"
```

## Database Schema

### EnrichedPost Collection

| Field       | Type     | Description                           |
|-------------|----------|---------------------------------------|
| _id         | ObjectID | MongoDB document ID                   |
| userId      | int      | User ID from the original post        |
| postId      | int      | Post ID from the original post        |
| title       | string   | Post title                            |
| body        | string   | Post body                             |
| ingested_at | datetime | UTC timestamp of ingestion            |
| source      | string   | Source identifier                     |

### IngestStatus Collection

| Field     | Type     | Description                           |
|-----------|----------|---------------------------------------|
| _id       | ObjectID | MongoDB document ID                   |
| timestamp | datetime | UTC timestamp of ingestion attempt    |
| success   | boolean  | Whether the ingestion was successful  |
| count     | int      | Number of records ingested            |
| error     | string   | Error message (if any)                |

## Design Decisions and Trade-offs

### Storage Choice: MongoDB

I chose MongoDB as the cloud-native storage solution for the following reasons:

1. **Flexibility**: MongoDB's document model is well-suited for semi-structured log data, allowing for easy schema evolution.
2. **Scalability**: MongoDB can scale horizontally to handle growing data volumes.
3. **Cloud-native options**: MongoDB Atlas provides a fully managed cloud service with automatic scaling and backups.
4. **Query capabilities**: MongoDB's query language is powerful for log analysis and filtering.

Trade-offs:
- MongoDB may not be as efficient as specialized time-series databases for pure time-series data.
- For extremely high write throughput, a solution like Apache Kafka might be more appropriate.

### Architecture Decisions

1. **Modular Design**: The service is built with clear separation of concerns, making it easy to extend or replace components.
2. **Periodic Ingestion**: Data is ingested at regular intervals rather than streaming, which is appropriate for the JSONPlaceholder API.
3. **Stateful Tracking**: The service tracks ingestion status to provide visibility and enable recovery from failures.

### Challenges and Improvements

The most challenging aspects of the implementation were:

1. **Error Handling**: Ensuring robust error handling across all components, especially for network and database operations.
2. **Testing**: Creating comprehensive tests that cover edge cases without being brittle.

With more time, I would improve:

1. **Metrics and Monitoring**: Add Prometheus metrics for better observability.
2. **Rate Limiting**: Implement adaptive rate limiting for API requests.
3. **Data Validation**: Add more sophisticated validation of incoming data.
4. **Pagination**: Implement pagination for the API endpoints to handle large datasets.
5. **Authentication**: Add authentication for the API endpoints.