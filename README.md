# jcourse_go: 选课社区 2.0 后端

## 依赖

* Golang >= 1.21

## 编译

```shell
go build
```

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