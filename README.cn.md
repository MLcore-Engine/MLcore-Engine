# 机器学习平台

[![Go Version](https://img.shields.io/badge/Go-1.20+-blue.svg)]()
[![React Version](https://img.shields.io/badge/React-18.0+-blue.svg)]()
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)]()

## 项目简介

这是一个基于 Kubernetes 的机器学习平台，提供完整的模型开发、训练和部署流程。采用标准的 Go 后端和 React 前端架构，便于快速上手和二次开发。

## 技术栈

### 后端
- 框架：Gin + GORM
- 存储：MinIO (对象存储) + MySQL
- 容器：Docker + Kubernetes
- 模型服务：Triton Inference Server
- 监控：Prometheus + Grafana

### 前端
- 框架：React 18 + React Router
- UI：Semantic UI React
- 状态管理：Redux Toolkit
- 请求：Axios
- 图表：Recharts

## 快速开始

### 环境要求
- Go 1.20+
- Node.js 16+
- Docker 20+
- Kubernetes 1.20+
- MySQL 8.0+

### 本地开发

1. 克隆项目
```bash
git clone https://github.com/your-org/ml-platform.git
cd ml-platform
```

2. 启动后端
```bash
cd backend
go mod tidy
go run main.go
```

3. 启动前端
```bash
cd frontend
npm install
npm start
```

## API 文档

### 核心数据结构

#### 训练任务
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

#### API 返回格式
```go
type Response struct {
    Code    int         `json:"code"`    // 状态码
    Message string      `json:"message"` // 提示信息
    Data    interface{} `json:"data"`    // 数据
    Total   int64      `json:"total"`   // 总数(列表接口)
}
```

### 主要 API 接口

#### 训练任务
- `POST /api/v1/training-jobs`: 创建训练任务
- `GET /api/v1/training-jobs`: 获取训练任务列表
- `GET /api/v1/training-jobs/:id`: 获取训练任务详情
- `PUT /api/v1/training-jobs/:id`: 更新训练任务
- `DELETE /api/v1/training-jobs/:id`: 删除训练任务

## 扩展开发

### 添加新的存储后端
实现 `Storage` 接口：
```go
type Storage interface {
    Upload(ctx context.Context, bucket, object string, reader io.Reader) error
    Download(ctx context.Context, bucket, object string) (io.ReadCloser, error)
    Delete(ctx context.Context, bucket, object string) error
}
```

### 添加新的模型服务
实现 `ModelServer` 接口：
```go
type ModelServer interface {
    Deploy(ctx context.Context, model *Model) error
    Undeploy(ctx context.Context, modelID string) error
    GetStatus(ctx context.Context, modelID string) (*ModelStatus, error)
}
```

## 学习资料

### Go 语言学习
1. [Go 官方文档](https://golang.org/doc/)
2. [Go by Example](https://gobyexample.com/)
3. [Gin 框架文档](https://gin-gonic.com/docs/)
4. [GORM 文档](https://gorm.io/docs/)

### React 学习
1. [React 官方文档](https://reactjs.org/docs/getting-started.html)
2. [React Router 文档](https://reactrouter.com/docs/en/v6)
3. [Redux Toolkit 文档](https://redux-toolkit.js.org/)
4. [Semantic UI React 文档](https://react.semantic-ui.com/)

## 贡献指南

欢迎提交 Pull Request 或 Issue。

## 许可证

Apache License 2.0