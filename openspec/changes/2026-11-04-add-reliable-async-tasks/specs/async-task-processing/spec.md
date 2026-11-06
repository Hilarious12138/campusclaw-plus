# 可靠异步任务能力（计划）

## ADDED Requirements

### Requirement: A1 Accept without waiting

服务 SHALL 在校验与持久化完成后返回202、task_id和status_url，不在请求内执行解析/嵌入。

#### Scenario: S01 Measure acceptance
- **GIVEN** 前置条件：已具备认证/材料前置能力，真实DB和broker；执行适配器分别延迟10秒、30秒。
- **WHEN** 操作：每组上传20份合法不同材料，单独记录文件接收完成后的处理耗时。
- **THEN** 结果：各组接收阶段P95<1000ms且小于执行延迟10%；响应不含执行结果，完整上传耗时另记。延迟替身与真实嵌入证据分开。

### Requirement: A2 Persist and publish reliably

材料版本、任务与outbox MUST 同事务提交，队列发布在事务外并等待路由/持久化确认。

#### Scenario: S02 Recover broker outage
- **GIVEN** 前置条件：DB可用，broker停止，本班非终态任务少于100。
- **WHEN** 操作：教师提交合法材料，再恢复broker。
- **THEN** 结果：返回202且task/outbox存在；恢复后30秒内开始执行，无需再次提交；确认丢失重发不丢任务。

#### Scenario: S03 Roll back failed transaction
- **GIVEN** 前置条件：材料文件已接收，注入任务或outbox插入失败。
- **WHEN** 操作：提交请求。
- **THEN** 结果：返回503，无已提交材料版本/任务/outbox半写，不调用外部执行；孤儿文件由清理机制处理。

### Requirement: A3 Resume interrupted processing

任务 SHALL 持久保存并由独立worker执行；停止与崩溃后可恢复，旧租约不能提交结果。

#### Scenario: S04 Stop and resume worker
- **GIVEN** 前置条件：真实broker/DB正常、worker停止。
- **WHEN** 操作：提交并查任务，再启动worker。
- **THEN** 结果：停止期间为pending且202成功；启动后30秒内running，使用10秒执行夹具时90秒内succeeded；无重新提交。

#### Scenario: S05 Recover expired lease
- **GIVEN** 前置条件：任务running，worker持有效租约。
- **WHEN** 操作：kill worker并启动另一实例，等待旧租约到期，再模拟旧fence提交。
- **THEN** 结果：任务在租约过期后30秒内被接管，attempts递增；旧提交被拒绝，最后仅一个结果。ACK前kill已成功任务时只增delivery_count、不增加执行attempts。

### Requirement: A4 Make business effects idempotent

系统 MUST 以业务版本唯一键和结果唯一约束去重，不能承诺外部调用恰好一次。

#### Scenario: S06 Replay and compete
- **GIVEN** 前置条件：相同教师/班级/材料版本/处理版本，两个worker。
- **WHEN** 操作：同一幂等键并发请求20次，再重复投递同task消息。
- **THEN** 结果：返回相同task_id；task_results该task仅1条，切片/向量无重复可见版本；成功后重投不改变结果摘要。

#### Scenario: S07 Reject conflicting key
- **GIVEN** 前置条件：已有同用户同班同类型幂等键及内容摘要。
- **WHEN** 操作：使用同键提交不同内容。
- **THEN** 结果：409/IDEMPOTENCY_CONFLICT，不新增任务、不覆盖原结果。

### Requirement: A5 Bound retry and dead letters

临时失败 SHALL 最多执行3次，间隔5秒/10秒加0—1秒抖动；永久失败直接failed。终态与死信待发记录同事务保存。

#### Scenario: S08 Retry transient errors
- **GIVEN** 前置条件：外部适配器连续返回429，注入可控时钟。
- **WHEN** 操作：执行至第三次失败并恢复broker。
- **THEN** 结果：attempts=3、failed、脱敏错误可查，死信确认发布后DLQ可见；无第四次执行。broker不可达时显示待发布，不谎报已入DLQ。

#### Scenario: S09 Reject permanent errors
- **GIVEN** 前置条件：任务输入无效或执行前发现权限已撤销。
- **WHEN** 操作：worker领取任务。
- **THEN** 结果：一次执行后failed，不调用嵌入或自动重试；畸形消息隔离，不阻断下一条正常消息。

### Requirement: A6 Authorize status and input

任务 MUST 逐资源授权，班级仅取认证身份。GET只读，拒绝时不得泄露或产生副作用。

#### Scenario: S10 Deny invalid access
- **GIVEN** 前置条件：存在A班材料任务及B班账号。
- **WHEN** 操作：匿名查任务、学生上传、B班查A任务及随机不存在任务。
- **THEN** 结果：分别401、403、404、404；两个404响应结构一致，不返回源输入，不写任务/结果。

#### Scenario: S11 Validate and limit
- **GIVEN** 前置条件：教师已登录，本班非终态任务达到100。
- **WHEN** 操作：提交未知class_id参数、超2MiB文件、不支持扩展名，再提交合法新任务。
- **THEN** 结果：前三种各400且无任务；容量请求429；GET反复查询不改变attempts/outbox或触发执行。

### Requirement: A7 Publish complete index versions

材料处理 SHALL 只发布完整且授权有效的索引版本，MySQL存正文，Qdrant只存向量和标识。

#### Scenario: S12 Interrupt index publication
- **GIVEN** 前置条件：旧版本可检索，新版本正写向量。
- **WHEN** 操作：中途kill再恢复worker。
- **THEN** 结果：处理中仍只能检索旧有效版本；恢复后新版本完整生效，确定性id无重复；MySQL发布失败不露出半成品，向量回查必须满足班级与有效版本。

### Requirement: A8 Observe and reproduce

系统 SHALL 提供任务状态、队列/执行指标及Docker自动验收，缺证据不能算成功。

#### Scenario: S13 Produce evidence
- **GIVEN** 前置条件：隔离Plus测试Compose环境，三条核心判据与负向用例已注册。
- **WHEN** 操作：运行make check及sh scripts/grade.sh。
- **THEN** 结果：回归执行并输出HTML/JSON、耗时原始值、DB结果核对、进程恢复和去重证据；reports含并发/事务专项记录；必测失败或缺证据退出非零，秘密不进报告。

#### Scenario: S14 Show task lifecycle
- **GIVEN** 前置条件：材料页有pending/running/succeeded/failed样本。
- **WHEN** 操作：查看、切换页面及模拟查询网络故障。
- **THEN** 结果：每2秒只读轮询、终态或离页停止；四态文案明确，错误可重新查询，不自动重新创建任务。
