# 后端开发规范

适用：Go服务、接口与数据访问。身份、配置、材料、检索、任务与审计均按本规范分层；验证范围见当前change。

## 分层与依赖

**要求**：`backend/cmd/api`负责装配。`internal/httpapi`处理协议与参数，service编排业务，domain声明健康与依赖端口，repository只实现数据库适配。

依赖方向为httpapi→service→domain、repository→domain。service不得导入`net/http`或`database/sql`。

检查：`make check`检查import；故意加入反向依赖时，检查必须失败。

## 接口路径

**要求**：业务接口使用`/api/<resources>`，不另开`/api/v1`。工程探活与登录/文档是显式公开入口；其他已注册业务要求Bearer，未注册的`/api/*`返回404。新增业务先明确鉴权要求。

检查：核对路由与接口清单；拒绝请求不能返回数据或产生副作用。

## 响应与错误码

**要求**：JSON字段使用snake_case（小写单词用下划线连接），结构为`{code,data?,error?,request_id}`。同时明确HTTP状态和业务code，不能只靠其中一个判断结果。

错误码：

- 参数非法：400/INVALID_ARGUMENT。
- 资源不存在：404/NOT_FOUND。
- 方法不支持：405/METHOD_NOT_ALLOWED。
- 依赖不可用：503/DEPENDENCY_UNAVAILABLE。
- 内部故障：500/INTERNAL_FAILURE。

健康成功用200/OK。新增业务还需登记401/UNAUTHENTICATED、403/FORBIDDEN、409/CONFLICT，不能统一写成“系统错误”。

检查：错误码集中在httpapi；测试同时检查code、必要字段和副作用。响应不得包含原始SQL或堆栈。

## 参数校验

**要求**：入口只接受白名单中的参数；服务层还要检查业务条件。`/health`、`/live`不接受query，带筛选参数返回400，不能忽略非法筛选后返回全量。

检查：发送`?class_id=other`、重复及未知参数，确认被拒绝。

## 查询与事务

**要求**：SQL只写在repository，使用参数绑定。事务中只做数据库操作，不调用外部服务。排序最后用唯一id兜底。分页limit默认20、上限100，超限返回400。

检查：业务实现后覆盖注入、并发和相同排序值的分页。不得用仅探活通过替代业务SQL检查。

## 幂等

**要求**：GET应幂等。幂等是指重复请求不会多产生业务结果。新增POST须声明幂等键与唯一约束，或者说明为何不提供幂等。

检查：重复发送写请求并核对数据库，不能仅靠页面禁止重复点击。

## 配置

**要求**：运行地址、版本、环境、探测和停机时限从env读取。必填项缺失或越界必须启动失败，不能悄悄使用默认值。Compose可固定基础设施镜像和端口契约，用户运行参数仍须外置。

检查：单测逐个删除配置；Compose使用空配置时必须失败。

## 异步与审计

任务接收用202/ACCEPTED，只读查询不触发执行。业务记录、任务与outbox同事务，发布/外部调用在事务外；消费者完成提交后ACK。写配置/业务与审计同事务，审计失败回滚，返回503。检查：并发同键、旧fence、审计/outbox不可写、stop/start及ACK前kill均核对数据库。
