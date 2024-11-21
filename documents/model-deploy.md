# Triton 模型推理部署文档

MLcore-Engine 采用的是NV的triton框架作为推理引擎，采用deploy+service方式部署到kubernetes集群中， Triton提供了丰富的部署方式和模型格式支持。

svc-template:

```yaml

#service中没有暴露metric相关的接口，如果需要使用metric相关的接口  
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: "2024-11-15T01:33:49Z"
  labels:
    app: root-triszbm9
  name: root-triszbm9
  namespace: triton-serving
  resourceVersion: "473059"
  uid: cb86471e-33d6-4b8a-8314-c4708c3050c4
spec:
  clusterIP: 10.99.106.127
  clusterIPs:
  - 10.99.106.127
  externalTrafficPolicy: Cluster
  internalTrafficPolicy: Cluster
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - name: http-triton
    nodePort: 30236
    port: 8000
    protocol: TCP
    targetPort: 8000
  - name: grpc-triton
    nodePort: 32046
    port: 8001
    protocol: TCP
    targetPort: 8001
  selector:
    app: root-triszbm9
  sessionAffinity: None
  type: NodePort

```

deploy-template

```yaml

apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "1"
  creationTimestamp: "2024-11-15T01:33:48Z"
  generation: 1
  labels:
    app: root-triszbm9
  name: root-triszbm9
  namespace: triton-serving
  resourceVersion: "534794"
  uid: 8a39bc1b-3a9a-4d08-b8fc-7bd50c691184
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: root-triszbm9
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: root-triszbm9
    spec:
      containers:
      - args:
        - --model-repository=/model
        - --allow-gpu-metrics=false
        - --strict-model-config=false
        command:
        - tritonserver
        image: 192.168.12.121:5005/traning/triton:24.10-py3
        imagePullPolicy: IfNotPresent
        name: root-triszbm9
        ports:
        - containerPort: 8000
          name: http-triton
          protocol: TCP
        - containerPort: 8001
          name: grpc-triton
          protocol: TCP
        - containerPort: 8002
          name: metrics-triton
          protocol: TCP
        resources:
          limits:
            cpu: "2"
            memory: "4294967296"
            nvidia.com/gpu: "0"
          requests:
            cpu: "2"
            memory: "4294967296"
            nvidia.com/gpu: "0"
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30

```
模板中已经将模型打入到镜像中，所以没有使用volume挂载模型，如果需要使用volume挂载模型，可以参考如下配置：    

```yaml
volumes:
- name: model-volume
  persistentVolumeClaim:
    claimName: model-pvc
```
也可以根据安装文档使用S3加载模型，参考：https://docs.nvidia.com/deeplearning/triton-inference-server/user-guide/docs/user_guide/model_repository.html#s3

### Triton Inference Server 常用功能与参数配置

| **功能**          | **参数**                                    | **说明**                                                                                     |
|-------------------|---------------------------------------------|--------------------------------------------------------------------------------------------|
| **HTTP 请求支持**  | `--allow-http=true`                        | 启用 HTTP 协议                                                                              |
| **gRPC 请求支持**  | `--allow-grpc=true`                        | 启用 gRPC 协议                                                                              |
| **启用 gRPC SSL**  | `--grpc-use-ssl=true`                      | 开启 gRPC 的 SSL 加密，确保请求数据安全                                                      |
| **动态模型更新模式**| `--model-control-mode=explicit`            | 设置为 `explicit` 模式，通过 API 动态加载和卸载模型                                           |
| **启用性能指标**   | `--allow-metrics=true`                     | 启用 Prometheus 风格的性能指标，便于监控系统性能                                             |
| **缓存功能**       | `--cache-config=local,size=104857600`      | 开启本地缓存，并设置缓存大小（单位：字节，例如 100MB）                                        |
| **跟踪功能**       | `--trace-config=triton,file=/tmp/trace.json`| 开启请求跟踪功能，将结果输出到指定文件路径 `/tmp/trace.json`                                 |
| **严格健康检查**   | `--strict-readiness=true`                  | 启用严格健康检查，仅在所有模型加载完成时返回就绪状态                                         |
| **关闭未使用的协议**| `--allow-http=false` 或 `--allow-grpc=false`| 关闭未使用的协议以减少暴露的攻击面                                                          |
| **禁用自动配置**   | `--disable-auto-complete-config=true`      | 禁用模型的自动补全配置，要求手动提供完整的模型配置文件                                        |

## 示例启动命令

```bash
tritonserver \
    --model-repository=/path/to/model_repo \
    --allow-http=true \
    --allow-grpc=true \
    --grpc-use-ssl=true \
    --model-control-mode=explicit \
    --allow-metrics=true \
    --cache-config=local,size=104857600 \
    --trace-config=triton,file=/tmp/trace.json \
    --strict-readiness=true \
    --disable-auto-complete-config=true
```
### 相关代码路径

#### 后端代码
- 控制器: `controller/triton_deploy.go`
  - 负责处理模型部署的 HTTP 请求
  - 包含创建、更新、删除和列表等操作

- 数据模型: `model/triton_deploy.go`
  - 定义了 TritonDeploy 数据结构
  - 包含数据库操作相关方法

#### 前端代码
- 页面组件: `web/src/pages/ModelDeployment.js`
  - 模型部署的主页面
  - 显示部署统计和操作入口

- 侧边栏组件: `web/src/components/sidebars/ModelDeploySidebar.js`
  - 模型部署相关的导航菜单
  - 包含部署列表和创建部署等链接

Triton官方文档： https://docs.nvidia.com/deeplearning/triton-inference-server/user-guide/docs/contents.html
备注：有必要在模型部署前使用model-analyzer工具对模型进行优化 https://github.com/triton-inference-server/tutorials/tree/main/Conceptual_Guide/Part_3-optimizing_triton_configuration