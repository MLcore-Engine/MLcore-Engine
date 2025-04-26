# API响应结构文档

## 整体响应结构

后端API响应遵循统一的多层嵌套结构，便于前端处理各类数据和状态。结构共有4层嵌套：

### 第一层：axios响应对象

前端使用axios请求后端时，得到的顶层响应对象包含多个属性：

```typescript
interface AxiosResponse {
  data: APIResponse<T>;  // 后端返回的实际数据
  status: number;        // HTTP状态码
  statusText: string;    // HTTP状态描述
  headers: any;          // 响应头
  config: any;           // 请求配置
}
```

因此，当前端收到响应时，通过`response.data`访问后端返回的APIResponse对象。

### 第二层：APIResponse<T> (response.data)

这是后端返回的标准响应格式，包含以下字段：

```typescript
interface APIResponse<T> {
  success: boolean;  // 请求是否成功
  message: string;   // 状态消息或错误信息
  data: T;           // 实际返回的数据，泛型T根据不同API而变化
}
```

所以，在前端代码中使用`response.data.data`来访问实际的业务数据。

### 第三层：ProjectResponse (response.data.data)

在项目相关API中，T通常是ProjectResponse类型，包含分页信息和项目列表：

```typescript
interface ProjectResponse {
  limit: number;      // 每页限制数量
  page: number;       // 当前页码
  total: number;      // 总项目数
  projects: Project[]; // 项目列表
}
```

前端通过`response.data.data.projects`访问项目列表数组。

### 第四层：Project[] (response.data.data.projects)

项目数组，每个项目对象包含项目基本信息：

```typescript
interface Project {
  id: string;          // 项目ID
  name: string;        // 项目名称
  description: string; // 项目描述
  createTime: string;  // 创建时间
  updateTime: string;  // 更新时间
  creator: string;     // 创建者ID
  users: User[];       // 项目成员列表
  // 其他项目属性...
}
```

### 第五层：User[] (response.data.data.projects[i].users)

项目成员列表，包含每个成员的信息：

```typescript
interface User {
  id: string;         // 用户ID
  username: string;   // 用户名
  email: string;      // 邮箱
  role: string;       // 在项目中的角色
  joinTime: string;   // 加入项目时间
  // 其他用户属性...
}
```

## 数据访问流程图

```
response                 -> axios响应对象
  └── response.data      -> APIResponse<T> (success, message, data)
      └── response.data.data -> ProjectResponse (limit, page, total, projects)
          └── response.data.data.projects      -> Project[] (项目列表)
              └── response.data.data.projects[i].users -> User[] (成员列表)
```

## 统一响应规范

所有API接口必须遵循以下响应规范，确保前端能够统一处理响应数据。

### 后端DTO规范

后端使用以下标准结构体来构建API响应：

#### 1. 基础响应结构

```go
// 基础响应结构
type BaseResponse struct {
    Success bool        `json:"success"`             // 请求是否成功
    Message string      `json:"message,omitempty"`   // 状态消息或错误信息
    Data    interface{} `json:"data,omitempty"`      // 实际返回的数据
}

// 成功响应
type SuccessResponse struct {
    Success bool        `json:"success" example:"true"`
    Message string      `json:"message" example:"操作成功"`
    Data    interface{} `json:"data,omitempty"`
}

// 错误响应
type ErrorResponse struct {
    Success bool        `json:"success" example:"false"`
    Message string      `json:"message" example:"发生错误"`
    Data    interface{} `json:"data,omitempty"`
}
```

#### 2. 分页数据结构

所有分页数据必须包含以下标准字段：

```go
// 通用分页数据结构
type PagedData struct {
    Total int64 `json:"total"` // 总记录数
    Page  int   `json:"page"`  // 当前页码
    Limit int   `json:"limit"` // 每页限制数
}
```

#### 3. 业务领域DTO

每个业务领域使用专门的DTO类型：

- **项目相关**：`ProjectDTO`, `ProjectResponse`, `ProjectListData`等
- **用户相关**：`UserDTO`, `LoginResponse`, `UsersListData`等
- **数据集相关**：`DatasetDTO`, `DatasetResponse`, `DatasetsListData`等
- **笔记本相关**：`NotebookDTO`, `NotebookResponse`, `NotebooksListData`等
- **训练任务相关**：`TrainingJobDTO`, `TrainingJobResponse`等
- **模型部署相关**：`TritonDeployDTO`, `TritonDeployResponse`等

### 响应实现规范

所有API处理函数必须遵循以下模式：

#### 1. 成功响应

```go
c.JSON(http.StatusOK, SuccessResponse{
    Success: true,
    Message: "操作成功消息",
    Data:    someDTO或SomeListData,  // 使用适当的DTO或列表数据
})
```

#### 2. 错误响应

```go
c.JSON(errorHttpStatusCode, ErrorResponse{
    Success: false,
    Message: "错误消息",
    // Data字段通常为nil
})
```

#### 3. 分页列表响应

```go
c.JSON(http.StatusOK, SuccessResponse{
    Success: true,
    Message: "获取列表成功",
    Data: SomeListData{
        Items: []SomeDTO{...},  // 列表项数组
        PagedData: PagedData{
            Total: total,
            Page:  page,
            Limit: limit,
        },
    },
})
```

## 实际DTO示例

以下是几个常用DTO的实际定义：

### 项目DTO

```go
// 项目DTO
type ProjectDTO struct {
    ID          uint        `json:"id"`
    Name        string      `json:"name"`
    Description string      `json:"description"`
    CreatedAt   time.Time   `json:"created_at,omitempty"`
    UpdatedAt   time.Time   `json:"updated_at,omitempty"`
    Users       []MemberDTO `json:"users,omitempty"`
}

// 项目列表数据
type ProjectListData struct {
    Projects []ProjectDTO `json:"projects"`
    PagedData
}
```

### 用户DTO

```go
// 用户DTO
type UserDTO struct {
    Id          uint      `json:"id"`
    Username    string    `json:"username"`
    DisplayName string    `json:"display_name"`
    Role        int       `json:"role"`
    Status      int       `json:"status"`
    Email       string    `json:"email"`
    CreatedAt   time.Time `json:"created_at,omitempty"`
}
```

## 完整JSON示例

```json
// 完整axios响应对象中的data属性
{
  "success": true,
  "message": "获取项目列表成功",
  "data": {              // response.data.data
    "projects": [        // response.data.data.projects
      {
        "id": 1,
        "name": "MLcore引擎开发",
        "description": "机器学习核心引擎开发项目",
        "created_at": "2023-05-15T08:00:00Z",
        "updated_at": "2023-06-20T10:30:00Z",
        "users": [       // response.data.data.projects[0].users
          {
            "userId": 123,
            "username": "张三",
            "email": "zhangsan@example.com",
            "role": 1
          },
          {
            "userId": 456,
            "username": "李四",
            "email": "lisi@example.com",
            "role": 2
          }
        ]
      }
    ],
    "total": 25,
    "page": 1,
    "limit": 10
  }
}
```

## 前端处理示例

```typescript
import axios from 'axios';
import { getProjects } from '../api/projectAPI';

// 使用自定义API函数获取项目列表
const fetchProjectsWithAPI = async () => {
  try {
    // getProjects内部已处理axios响应，直接返回APIResponse
    const response = await getProjects(1, 10);
    
    // response已是APIResponse<ProjectResponse>类型
    if (response.success) {
      // 通过response.data访问ProjectResponse
      const { projects, total } = response.data;
      console.log(`共有${total}个项目`);
      
      // 访问第一个项目的成员
      if (projects.length > 0) {
        const firstProjectUsers = projects[0].users;
        console.log(`第一个项目有${firstProjectUsers.length}个成员`);
      }
    } else {
      console.error('获取项目失败:', response.message);
    }
  } catch (error) {
    console.error('API请求异常:', error);
  }
};

// 直接使用axios获取项目列表
const fetchProjectsWithAxios = async () => {
  try {
    // axios返回的是最外层响应
    const response = await axios.get('/api/projects?page=1&limit=10');
    
    // 通过response.data访问APIResponse
    if (response.data.success) {
      // 通过response.data.data访问ProjectResponse
      const projectsData = response.data.data;
      console.log(`共有${projectsData.total}个项目`);
      
      // 通过response.data.data.projects访问项目列表
      const projects = projectsData.projects;
      
      // 访问第一个项目的成员
      if (projects.length > 0) {
        // 通过response.data.data.projects[0].users访问成员列表
        const firstProjectUsers = projects[0].users;
        console.log(`第一个项目有${firstProjectUsers.length}个成员`);
      }
    } else {
      console.error('获取项目失败:', response.data.message);
    }
  } catch (error) {
    console.error('API请求异常:', error);
  }
};
```

## 命名规范说明

根据项目规范：

1. 所有ID字段在前端使用小驼峰命名法(camelCase)，如`id`、`userId`等
2. 后端结构体字段使用PascalCase，JSON标签使用camelCase
3. API交互中所有字段统一使用camelCase

## 严格遵守规范

**重要提示：**所有API开发人员必须严格遵守以下规则：

1. **响应结构统一**：所有接口返回必须使用`SuccessResponse`或`ErrorResponse`
2. **状态字段完整**：每个响应必须包含`success`和`message`字段
3. **类型一致性**：同一类型的字段在不同DTO间必须保持一致的类型和命名
4. **分页数据标准**：所有分页列表必须包含标准的`total`、`page`和`limit`字段
5. **错误码统一**：API错误必须使用统一的HTTP状态码和业务状态码

## 接口设计提示

为简化前端访问，建议:

1. API模块封装时剥离外层axios响应，直接返回`response.data`(APIResponse)
2. 考虑在API函数中直接解构数据，返回`response.data.data`或更深层次的数据
3. 使用TypeScript类型定义确保类型安全，减少访问错误 