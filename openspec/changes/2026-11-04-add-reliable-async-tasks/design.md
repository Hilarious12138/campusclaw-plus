# 异步任务设计（待确认）

## Context / 实施边界

当前Plus只有工程基线。顺序为基础设施骨架→任务内核→已验收的认证/班级/材料前置能力→材料索引接入→运行验收。本次先完成设计，不改业务代码。若验收要求必须使用作业批改，则在P3补作业领域变更；不能用材料演示冒充批改已实现。

## Decisions / 技术栈

沿用React18 + TypeScript + Vite、Go net/http、MySQL8、Nginx/Compose。新增RabbitMQ（实施时选受支持版本并锁镜像摘要，凭证仅env）及Go官方客户端amqp091-go（实施前审查版本/许可）。业务仍依赖领域端口，RabbitMQ/MySQL/Qdrant/嵌入调用在适配器中；worker复用服务层，不通过HTTP自调。

Compose增加broker、dispatcher、worker，与api同网；仅Nginx映射宿主端口。RabbitMQ数据卷独立，AMQP与管理端口不对外发布。第一期1个worker、prefetch=1，可用两个实例做竞争验收。单节点持久卷不代表高可用或磁盘损坏零丢失。

## 数据流

1. API鉴权并验证教师、本班材料、类型/大小和持久文件引用；上传正文先写共享持久附件卷。失败文件定期清理，不把文件IO放SQL事务中。
2. 一个MySQL事务写材料版本、task与outbox；提交后返回202。数据库失败返回503，不给虚假task_id。broker暂不可达仍可接收已落库任务；容量超限拒绝。
3. dispatcher每1秒扫描到期outbox，用短事务/租约竞争领取，事务外发布持久消息，mandatory路由并等publisher confirm；确认后才标已投递。确认丢失可重复发，禁止先标成功。连接故障持续退避并报警，不悄悄删outbox。
4. worker从RabbitMQ取任务id，回MySQL读取可信类型、归属与输入版本，重新检查资源仍可执行。领取租约后在事务外解析/嵌入/写向量；最后一个MySQL事务记录结果并置succeeded，提交后ACK。
5. 页面查任务状态；只有完整索引版本生效才可检索，处理中仍显示旧的有效版本或“处理中”。

消息只含schema_version=1、event_id、task_id、request_id；不传正文、密码或可自行指定的班级。班级来源是任务创建时的认证身份，worker不得信任消息中的归属。

## 接口一览（均为计划）

沿用统一{code,data?,error?,request_id}。所有业务接口Bearer鉴权；class_id不作请求参数，未知/重复参数返回400。

| 接口 | 权限与结果 |
| --- | --- |
| POST /api/materials（multipart） | 教师；Idempotency-Key必填，1—128 ASCII字符；校验并持久文件后202/ACCEPTED，data含material_id/task_id/status/status_url，不含切片/向量执行结果 |
| GET /api/tasks/{id} | 只读；教师可查本班材料任务，学生只查其可读材料的安全进度；200/OK含status、stage、attempts、next_run_at、脱敏error_code、result_ref与时间，不含源输入/密钥/他人正文 |
| GET /api/materials/{id} | 依材料授权查状态；文件仍由鉴权接口读取，不静态公开附件 |

401/UNAUTHENTICATED：无效令牌；403/FORBIDDEN：学生上传；404/NOT_FOUND：他班与不存在一致；400/INVALID_ARGUMENT：类型、大小、参数非法；409/IDEMPOTENCY_CONFLICT：同键不同内容；503/DEPENDENCY_UNAVAILABLE：DB/附件不可写；429/TASK_CAPACITY_EXCEEDED：本班非终态任务达到100。POST每次仍须先授权再返回旧任务。GET不能触发投递或重试。

初版材料只接受UTF-8 txt/md，最大2MiB，不压缩；PDF/Office另行设计解析安全边界。文件接收耗时不算异步消除的耗时。

## 数据模型与幂等

- tasks：id、class_id、created_by、type、resource_id、input_version、input_hash、pipeline_version、idempotency_key、status、stage、attempts、delivery_count、next_run_at、lease_owner/lease_until、fence、error_code、created_at/updated_at/started_at/finished_at、result_ref。
- 唯一键(class_id,type,resource_id,input_version,pipeline_version)挡住同一业务重复任务；请求幂等键按(class_id,created_by,type,idempotency_key)唯一，保存规范化内容摘要；相同内容返回同task_id，不同内容409。并发请求靠唯一约束，不靠先查。
- task_outbox：event_id主键、task_id、event_kind、generation、available_at、published_at、租约与publish_attempts；唯一(task_id,event_kind,generation)，涵盖初始投递、重试、恢复、死信事件。
- task_results：task_id唯一、结果版本引用及摘要；task_attempts记录每次实际执行和错误分类。结果写入/任务终态同一事务，不使用INSERT IGNORE吞掉无关SQL错误。
- knowledge_chunks正文只在MySQL；同材料版本/pipeline/序号唯一；Qdrant用确定性chunk_id upsert，不存正文。先建立完整待发布版本，再以MySQL当前版本指针原子发布；检索必须回查MySQL的有效版本与班级。崩溃残留向量不可成为可见结果，恢复幂等覆盖，孤儿由有界清理任务删除。

## 状态、重试与恢复

pending→running→succeeded；暂时失败running→pending，附next_run_at；永久失败或执行次数耗尽→failed。终态不直接修改，首版不开放人工重放。任务界面显示待处理/执行中/已完成/已失败；stage使用queued/parsing/chunking/embedding/publishing，非伪造百分比。

默认配置：总执行最多3次，重试等待5秒、10秒并加0—1秒抖动；单次执行上限120秒，单次嵌入调用30秒；SDK内部自动重试关闭，避免叠乘。超时/429/5xx重试；输入坏、权限失效、密钥无效直接失败。索引分批可续做，超120秒按同版本重试，不提前发布。所有值env集中配置、校验范围；.env.example只键名。Qdrant/嵌入参数由知识库前置变更最终确认，未配置不得伪成功。

重试：将pending状态、下次时间与重试outbox在同一事务写好后ACK原消息；到期再由dispatcher投递，不循环NACK立即requeue。failed与死信outbox同事务；dispatcher确认发布到持久DLQ后标记。代理不可用时死信待发布仍可查，不能声称已入DLQ。畸形/未知版本消息隔离死信并记录脱敏原因。

租约150秒，每10秒续租；恢复器每5秒检查过期running，按attempts决定重新排队或失败。领取时递增fence，最终提交须匹配owner/fence和有效租约，旧worker结果不能覆盖新结果。未结束事务不得做broker/模型调用。成功后ACK前kill：重投发现succeeded直接ACK，delivery_count递增、attempts不递增；执行中kill：租约过期后重跑，attempts递增。外部模型可能重复调用及计费，业务只保证一个可见结果。

## 界面与运维

材料页展示文件与任务状态，2秒轮询、终态/离页停止，加载/空/错误可区分，故障提供重新查询按钮而非重复创建。无新增独立任务管理台。

API健康保留DB必需依赖；broker断开单独标degraded，DB可用且未超容量仍能202。dispatcher就绪查DB+broker，worker查DB+broker和对应执行依赖；live只表示存活。worker停止领取后最多10秒排空，未完成取消上下文、不ACK，依租约恢复。结构化日志记录task_id、request_id、event_id、type、stage、attempts、耗时与错误码，不写正文。观测pending数量/最老等待、outbox积压、broker ready/unacked、执行耗时、重试与死信数；最老pending>60秒或死信>0报警记录，真实通知通道另行接入。

## 验收与取舍

三条核心判据见spec A1/A3/A4。20次提交测P95，测试执行延迟分别10秒/30秒；文件上传完成后的接收阶段P95<1秒且<执行延迟10%，总上传耗时另记。没有旧同步系统，使用明确标识的延迟适配器做可重复机制验收；真实嵌入另测，不伪称旧链路提速。

真实RabbitMQ/MySQL容器验证stop/start、broker重启、ACK前kill、两worker并发、DB回滚、outbox确认丢失、过期租约fence；单测使用注入时钟/适配器。脚本按URL放apitest/cases，make check纳入容器回归，一键Docker grade生成HTML/JSON和截图，覆盖不可执行即失败。原始报告reports/，脱敏摘要docs/verification；并发与事务报告分别输出reports/risk-concurrency.md、reports/risk-transaction.md。完成后按已验收提交打async-202标签；现在不生成假报告或标签。

部署先增表/索引→broker/dispatcher/worker→启用入口；旧字段/表不删除。回滚先禁入口、保留任务/outbox/卷，再回退兼容镜像；不能删卷当作回滚。

## 规范待完善与检查方式

| 原规范 | 实施时补充 | 检查 |
| --- | --- | --- |
| 产品 | 任务四态、轮询、处理中材料可见性 | 浏览器四态/离页停止轮询 |
| 后端 | 202/429/409、幂等键、worker依赖端口 | HTTP/并发约束和import检查 |
| 安全 | 任务逐资源授权、worker重验、broker最小权限 | A/B班和权限撤销、秘密扫描 |
| 测试 | 三判据、真实故障注入、性能口径 | 容器回归与报告缺证据失败 |
| 运维 | broker持久卷、恢复/积压与worker停机 | stop/kill/restart与时限 |
| AI | 嵌入超时/重试归属/重复成本 | 429/超时及调用次数核对 |
| 数据 | 任务/outbox/fence、结果与索引版本原子发布 | 回滚、唯一键、旧worker写入拒绝 |
| 评审 | 并发/事务证据与三判据映射 | Scenario→脚本→数据库→报告 |

## 实施调整

UTF-8 txt/md原文直接写MySQL LONGTEXT，与材料/任务/outbox同事务持久化，避免引入共享附件卷与孤儿文件；下载仍走材料鉴权接口。原文上限2MiB不变。原规划的文件卷/孤儿文件清理不适用于本期无独立附件存储，后续二进制附件另行设计。源文件接收后至事务提交的性能口径不变。任务内核拆为端口、SQL适配、AMQP适配、执行服务和独立入口，跨文件超过400行以保留完整故障边界；分组测试提交，不整体复制旧服务。
