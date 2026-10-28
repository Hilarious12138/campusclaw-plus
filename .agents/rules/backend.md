# 后端开发规范

| 约定 | 检查方式 |
| --- | --- |
| backend/cmd/api装配；internal/httpapi处理协议/参数，service编排，domain声明健康与依赖端口，repository只实现数据库适配。依赖httpapi→service→domain、repository→domain；service不得导入net/http或database/sql。 | `make check`静态import检查；注入一段反向依赖必须失败。 |
| 业务路径未来沿用/api/<resources>，不双开/api/v1；本课仅显式公开GET /health、GET /live，任何/api/*未注册即404。新增业务必须先登记鉴权。 | 请求清单与路由比对；负向请求无数据与副作用。 |
| JSON使用snake_case，响应{code,data?,error?,request_id}；code与HTTP语义分别明确，健康成功200/OK，依赖失败503/DEPENDENCY_UNAVAILABLE。 | HTTP测试同时断言code、必要字段与副作用。 |
| 错误：400/INVALID_ARGUMENT，404/NOT_FOUND，405/METHOD_NOT_ALLOWED，503/DEPENDENCY_UNAVAILABLE，500/INTERNAL_FAILURE；未来业务登记401/UNAUTHENTICATED、403/FORBIDDEN、409/CONFLICT，禁止泛化“系统错误”。 | 错误码集中在httpapi，场景逐项覆盖；响应不含原始SQL/堆栈。 |
| 参数白名单在入口，服务层守业务不变量；本课health/live不接受query，任意筛选400，不静默返回全量。 | `?class_id=other`、重复/未知参数必须失败。 |
| SQL只在repository，使用绑定参数；事务只做数据库操作，不夹外调。排序末尾id唯一兜底，分页limit默认20/上限100，超限400。 | 业务出现时做注入/并发/分页同排序键测试，本课只有Ping。 |
| GET幂等，无业务写接口；新增POST必须声明幂等键与唯一约束或明确非幂等理由。 | 重放写请求核对DB，不能靠禁按钮代替。 |
| 运行地址、版本、环境和探测/停机时限从env读取，必填缺失/越界启动失败，无隐式回退。基础设施镜像/端口契约可在Compose固定，用户运行参数外置。 | 配置单测逐键删除；Compose空配置失败。 |
