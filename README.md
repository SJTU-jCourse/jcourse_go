# jcourse_go: 选课社区 2.0 后端

## 依赖

* Golang >= 1.24

## 编译

```shell
go build
```
构建成功后，项目下将出现`jcourse_go`可执行文件。

### 配置文件

项目使用两种配置方式：

#### 1. API 服务配置 (YAML)
编辑 `config/config.yaml` 文件：

```yaml
db:
  host: localhost
  port: 5432
  user: jcourse
  password: jcourse
  db_name: jcourse_v2
  debug: true

server:
  debug: true
  port: 8888

sqlite:
  path: ""

redis:
  addr: "localhost:6379"
  username: ""
  password: ""
  db: 0

smtp:
  host: "127.0.0.1"
  port: 25
  username: ""
  password: ""
  sender: ""
```

#### 2. Worker 服务配置 (.env)
在项目根目录下创建 `.env` 文件：

```sh
# 数据库配置
POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5432
POSTGRES_USER=jcourse
POSTGRES_PASSWORD=jcourse
POSTGRES_DBNAME=jcourse_v2

# Redis 配置
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_USERNAME=
REDIS_PASSWORD=

# 安全配置
CSRF_KEY=your-csrf-key
SESSION_SECRET=your-session-secret

# 邮件配置
SMTP_HOST=127.0.0.1
SMTP_PORT=25
SMTP_SENDER=
SMTP_USERNAME=
SMTP_PASSWORD=

# 调试模式
DEBUG=true
```

### 启动Docker容器
运行以下命令启动数据库和Redis容器：
```sh
docker compose up -d
```

### 初始化数据
执行以下命令以初始化数据库：
```sh
# 创建数据库表
go run script/migrate/migrate.go
# 从data/2024-2025-1.csv 导入课程数据
go run script/import_from_jwc/import_from_jwc.go
```
注：`data/2024-2024-1.csv`中为虚拟的测试数据，可以参考格式进行修改，但是请不要将真实的课程信息加入Git仓库。

### 启动后端
```sh
# 启动 API 服务（使用 YAML 配置）
./jcourse_go

# 启动异步任务处理服务（使用 .env 配置）
./worker
```
后端成功启动后在**8888**端口运行（API 服务）。

## 代码结构

项目采用清洁架构（Clean Architecture）和领域驱动设计（DDD）模式：

### 可执行程序 (`cmd/`)
* `api/`：API 服务主程序（HTTP 服务器）
* `worker/`：异步任务处理服务

### 内部代码 (`internal/`)
* `app/`：应用容器和依赖注入
* `config/`：配置管理
* `application/`：应用层（CQRS 模式）
  * `command/`：命令处理（写操作）
  * `query/`：查询处理（读操作）
  * `vo/`：视图对象
* `domain/`：领域层（核心业务逻辑）
  * `auth/`：认证领域
  * `course/`：课程领域
  * `email/`：邮件领域
  * `event/`：事件领域
  * `notification/`：通知领域
  * `rating/`：评分领域
  * `reaction/`：反应领域
  * `review/`：评论领域
  * `statistic/`：统计领域
  * `user/`：用户领域
  * `shared/`：共享领域
* `infrastructure/`：基础设施层
  * `dal/`：数据访问层（数据库连接）
  * `entity/`：数据库实体定义
  * `repository/`：仓储实现
  * `rpc/`：远程过程调用
* `interface/`：接口层
  * `controller/`：HTTP 控制器
  * `handler/`：处理器（如 LLM 处理）
  * `middleware/`：HTTP 中间件
  * `router/`：路由配置
  * `task/`：异步任务
    * `asynq/`：基于 Asynq 的任务队列
    * `base/`：任务基础组件
    * `lock/`：分布式锁
    * `biz/`：业务相关的异步任务
      * `ping/`：ping 任务
      * `statistic/`：统计任务
* `service/`：服务层

### 通用库 (`pkg/`)
* `apperror/`：应用错误处理
* `util/`：工具函数

### 脚本文件 (`script/`)
* `import_from_jwc/`：从教务系统导入数据
* `import_from_v1/`：从 v1 版本导入数据
* `migrate/`：数据库迁移
* `load/`：数据加载和扩展
  * `teacher/`：教师数据扩展
  * `trainingplan/`：培养方案数据扩展

### 配置文件 (`config/`)
* `config.yaml`：API 服务的主配置文件（YAML 格式）
* `.env`：Worker 服务的环境变量配置文件

## 架构特点

项目采用清洁架构（Clean Architecture）和领域驱动设计（DDD）模式：

1. **依赖方向**：外层依赖内层，内层不依赖外层
2. **CQRS 模式**：读写分离，命令（Command）处理写操作，查询（Query）处理读操作
3. **领域驱动**：每个业务领域独立，包含实体、值对象、仓储接口等
4. **分层架构**：
   - **Domain 层**：核心业务逻辑，不依赖任何框架
   - **Application 层**：应用服务，协调领域对象
   - **Infrastructure 层**：技术基础设施，如数据库、缓存等
   - **Interface 层**：对外接口，如 HTTP 控制器、中间件等

## 技术栈

* **Web 框架**：Gin
* **ORM**：GORM
* **数据库**：PostgreSQL (with pgVector extension), SQLite
* **缓存**：Redis
* **任务队列**：Asynq
* **配置管理**：YAML
* **邮件**：gomail
* **中文分词**：gse
* **Docker**：容器化部署

## 数据库支持

项目支持多种数据库：
* **PostgreSQL**：主数据库，支持向量搜索（pgVector）
* **SQLite**：轻量级数据库，适用于开发和测试
