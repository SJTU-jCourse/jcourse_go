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
go run script/migrate/migrate.go
# 从data/2024-2025-1.csv 导入课程数据
go run script/import_from_jwc/import_from_jwc.go
```
注：`data/2024-2024-1.csv`中为虚拟的测试数据，可以参考格式进行修改，但是请不要将真实的课程信息加入Git仓库。

### 启动后端
```sh
./jcourse_go
```
后端成功启动后在**8080**端口运行。

## 代码结构

* `\cmd`：可编译成单独执行文件的部分
  * `\api`：API 服务主程序
* `\internal`：内部代码
  * `\constant`：常量定义
  * `\dal`：数据库链接
  * `\handler`：后端接口承载，调用业务函数，但是不执行具体业务逻辑
  * `\middleware`：HTTP 服务的中间件
  * `\model`：模型定义
    * `\dto`：前后端交互模型
    * `\po`：DB 存储模型
    * `\converter`：转换函数
  * `\repository`：存储层查询，屏蔽存储细节，供业务调用
  * `\rpc`：外部调用
  * `\service`：执行具体的业务逻辑
  * `\task`：异步任务
    * `\asynq`：基于 Asynq 的任务队列
    * `\lock`：分布式锁
    * `\biz`：业务相关的异步任务
* `\pkg`：与业务无关的通用库
* `\script`：脚本文件
  * `\codegen`：代码生成
  * `\gormgen`：GORM 模型生成
  * `\migrate`：数据库迁移
  * `\import_from_jwc`：从教务系统导入数据
  * `\import_from_v1`：从 v1 版本导入数据
  * `\refresh_db`：刷新数据库

---

## 为什么查询层只使用 GORM，而不直接使用 gorm/gen？

本项目采用 DDD + CQRS。对于 Query（读）侧，我们在 `internal/application/query` 中选择直接使用 GORM 的链式查询与 Scan 到自定义 VO 的方式，而不使用 gen 的强类型链式 API，主要有以下原因：

1) 读模型需要更强的 SQL 表达力与灵活性
- 课程列表/筛选中大量使用子查询、MAX 聚合出“最新学期”、EXISTS 子查询、以及对汇总结果再联结的写法；这些在 GORM 下以 Table/Joins/Where/Select 更直观，拼装成本低。
- 用 gen 也能实现，但表达会更冗长，链式类型约束使复杂 Join 与子查询的可读性下降，维护难度上升。

2) 更贴近 CQRS 的“只为展示塑形”目标
- Query 层直接扫描到轻量 struct/VO（只取必要字段），避免实体/关系的过度加载，降低 IO 与内存占用。
- gen 的典型用法更偏向 ORM 实体导航（Preload 等），容易引入多余字段；而我们需要的是极致按需的列选择与塑形。

3) 降低对代码生成器的耦合与运维成本
- 使用 gen 需要在模型或字段变更后重新生成代码；读侧的查询经常微调（新增筛选、调整聚合），使用 GORM 直接改 SQL 表达即可，避免生成/提交噪音。

4) 边界清晰：写侧/仓储可用 gen，读侧/应用查询用 GORM
- 项目仍保留 gen 生成的 query 包以支持仓储层等场景；但在应用层的读服务，为了应对跨表聚合与视图化塑形，我们刻意采用更自由的 GORM 查询。

5) 更好的可读性与跨角色协作
- 针对 DBA 或熟悉 SQL 的同学，直接映射到 SQL 语义的 GORM 语句更易沟通与 Review。

6) 性能与稳定性
- 读侧明确只选择必要列、聚合与排序，直接 Scan 到临时行结构，减少模型映射的额外开销。

综上，在读侧采用“GORM + 自定义 VO 扫描”的方式，能在满足 CQRS 前提下，以更低的心智负担实现复杂筛选与聚合逻辑；而 gen 仍可在需要强类型仓储访问的写侧继续发挥作用。文中所述实现可参考 `internal/application/query/course_query.go` 中的 GetCourseList/GetCourseDetail/GetCourseFilter。
