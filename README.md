<p align="right">
   <strong>English</strong> | <a href="./README.cn.md">中文</a>
</p> 
# Machine Learning Platform

[![Go Version](https://img.shields.io/badge/Go-1.20+-blue.svg)]()
[![React Version](https://img.shields.io/badge/React-18.0+-blue.svg)]()
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)]()

## Overview

A Kubernetes-based machine learning platform providing comprehensive model development, training, and deployment workflows. Built with standard Go backend and React frontend architectures for easy adoption and secondary development.

## Tech Stack

### Backend
- Framework: Gin + GORM
- Storage: MinIO (Object Storage) + MySQL
- Container: Docker + Kubernetes
- Model Serving: Triton Inference Server
- Monitoring: Prometheus + Grafana

### Frontend
- Framework: React 18 + React Router
- UI: Semantic UI React
- State Management: Redux Toolkit
- HTTP Client: Axios
- Charts: Recharts

## Quick Start

### Prerequisites
- Go 1.20+
- Node.js 16+
- Docker 20+
- Kubernetes 1.20+
- MySQL 8.0+

### Local Development

1. Clone the repository
```bash
git clone https://github.com/your-org/ml-platform.git
cd ml-platform
```

2. Start backend
```bash
cd backend
go mod tidy
go run main.go
```

3. Start frontend
```bash
cd frontend
npm install
npm start
```

## API Documentation

### Core Data Structures

#### Training Job
```go
type TrainingJob struct {
    ID              uint      `json:"id"`
    UserID          uint      `json:"user_id"`
    ProjectID       uint      `json:"project_id"`
    Describe        string    `json:"describe"`
    Namespace       string    `json:"namespace"`
    Image           string    `json:"image"`
    ImagePullPolicy string    `json:"image_pull_policy"`
    RestartPolicy   string    `json:"restart_policy"`
    Args            []string  `json:"args"`
    MasterReplicas  int       `json:"master_replicas"`
    WorkerReplicas  int       `json:"worker_replicas"`
    GPUsPerNode     int       `json:"gpus_per_node"`
    CPULimit        string    `json:"cpu_limit"`
    MemoryLimit     string    `json:"memory_limit"`
    Status          string    `json:"status"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

#### API Response Format
```go
type Response struct {
    Code    int         `json:"code"`    // Status code
    Message string      `json:"message"` // Message
    Data    interface{} `json:"data"`    // Data
    Total   int64      `json:"total"`   // Total count (for list APIs)
}
```

### Main API Endpoints

#### Training Jobs
- `POST /api/v1/training-jobs`: Create training job
- `GET /api/v1/training-jobs`: Get training job list
- `GET /api/v1/training-jobs/:id`: Get training job details
- `PUT /api/v1/training-jobs/:id`: Update training job
- `DELETE /api/v1/training-jobs/:id`: Delete training job

## Extension Development

### Add New Storage Backend
Implement the `Storage` interface:
```go
type Storage interface {
    Upload(ctx context.Context, bucket, object string, reader io.Reader) error
    Download(ctx context.Context, bucket, object string) (io.ReadCloser, error)
    Delete(ctx context.Context, bucket, object string) error
}
```

### Add New Model Server
Implement the `ModelServer` interface:
```go
type ModelServer interface {
    Deploy(ctx context.Context, model *Model) error
    Undeploy(ctx context.Context, modelID string) error
    GetStatus(ctx context.Context, modelID string) (*ModelStatus, error)
}
```

## Learning Resources

### Go
1. [Go Official Documentation](https://golang.org/doc/)
2. [Go by Example](https://gobyexample.com/)
3. [Gin Framework Documentation](https://gin-gonic.com/docs/)
4. [GORM Documentation](https://gorm.io/docs/)

### React
1. [React Official Documentation](https://reactjs.org/docs/getting-started.html)
2. [React Router Documentation](https://reactrouter.com/docs/en/v6)
3. [Redux Toolkit Documentation](https://redux-toolkit.js.org/)
4. [Semantic UI React Documentation](https://react.semantic-ui.com/)

## Contributing

Pull requests and issues are welcome.

## License

Apache License 2.0