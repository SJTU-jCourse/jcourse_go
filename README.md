# jcourse_go: 选课社区 2.0 后端

## 依赖

* Golang >= 1.21

## 编译

```shell
go build
```
构建成功后，项目下将出现`jcourse_go`可执行文件。

### 配置环境变量
可以在项目根目录下创建一个`.env`文件，内容包含必要的环境变量配置：
```sh
export REDIS_HOST=127.0.0.1
export REDIS_PORT=6379

export POSTGRES_HOST=127.0.0.1
export POSTGRES_PORT=5432
export POSTGRES_USER=jcourse
export POSTGRES_PASSWORD=jcourse
export POSTGRES_DBNAME=jcourse

export VECTORDB_HOST=127.0.0.1
export VECTORDB_PORT=5433
export VECTORDB_USER=jcourse
export VECTORDB_PASSWORD=jcourse
export VECTORDB_DBNAME=jcourse

export CSRF_KEY=your-csrf-key
export SESSION_SECRET=your-session-secret

export SMTP_HOST=127.0.0.1
export SMTP_PORT=25
export SMTP_SENDER=
export SMTP_USERNAME=
export SMTP_PASSWORD=

# 在Debug模式下，注册用户无需验证邮箱验证码
export DEBUG=true
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
go run cmd/migrate/migrate.go
# 从data/2024-2025-1.csv 导入课程数据
go run cmd/import_from_jwc/import_from_jwc.go
```
注：`data/2024-2024-1.csv`中为虚拟的测试数据，可以参考格式进行修改，但是请不要将真实的课程信息加入Git仓库。

### 启动后端
```sh
./jcourse
```
后端成功启动后在**8080**端口运行。

## 代码结构

* `\cmd`：可编译成单独执行文件的部分
* `\constant`：常量定义
* `\dal`：数据库链接
* `\handler`：后端接口承载，调用业务函数，但是不执行具体业务逻辑
* `\middleware`：HTTP 服务的中间件
* `\model`：模型定义
  - `\dto`：前后端交互模型
  - `\po`：DB 存储模型
  - `\domain`：业务领域模型
  - `\converter`：转换函数
* `\pkg`：与业务无关的通用库
* `\repository`：存储层查询，屏蔽存储细节，供业务调用
* `\rpc`：外部调用
* `\service`：执行具体的业务逻辑
* `\task`：异步任务
  * `\server`：异步任务单独的可执行文件
* `\util`：与业务无关的工具方法