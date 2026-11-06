# 安全与转换统筹设计（待批准）

## 当前事实与批准边界

2026-11-06只读盘点：Plus无业务，已有make check、audit、Docker grade与历史秘密检查；没有check-security或业务审计。GD HEAD为a86d1c0ec1635b24a09010fc06f29c3458cf80f4，工作树有未提交规划与output/tmp；本次不触碰。以下只是候选文件检查，不表示已跑测试、已安全或已可复用。

批准后先从GD该提交导出独立源码快照到仓库外隔离目录，禁用旧.env和卷；工作树非提交内容不自动纳入。记录文件摘要、Go/前端版本、原测试退出码/跳过项，再转换并跑Plus测试。失败模块先修复或重新实现，不以GD历史结果代替Plus结果。Go测试缓存、临时产物也只写隔离目录。

单次实施按下表小批提交；手写业务目标≤400行，确需超出须在当前任务记录说明拆分原因。批准此规划不等于允许跳过阶段门禁；新业务或契约变更先更新四文档。异步机制已有独立四文档，本变更引用而不复制正文，依赖前置完成后执行。

## 技术栈与复用决策

保留React18/TypeScript/Vite、Go net/http分层、MySQL8、Nginx/Compose；增加RabbitMQ实现可靠任务，Qdrant用于检索。模型/向量/MCP凭证重新配置到Plus，不能复用会话中或GD旧秘密。模型名/维度以批准的服务配置和探针为准，缺配置即失败，不猜测模型或生成替代回答。

| 阶段 | GD候选 | Plus转换与验收重点 |
| --- | --- | --- |
| T0 安全与来源 | 单测/apitest思路 | 独立源码快照测试、来源清单、秘密/依赖扫描、审计端口与固定SEC矩阵 |
| T1 身份/班级 | auth.go、password.go及测试 | Argon2id、JWT校验与撤销思路可候选复用；SQL移repository、鉴权移service/http适配；创建Plus独立issuer/audience与种子 |
| T2 材料/异步 | materials.go、chunks.go及测试 | 参数化SQL/文件授权/切片算法转换；10MiB旧限制收敛至已设计2MiB txt/md；移除旧索引后台执行，接task/outbox/worker |
| T3 检索 | retrieval.go、retrieval_scope.go及测试 | 默认hybrid、keyword不调用嵌入、vector阈值0.35、RRF k=60/top5；MySQL正文、Qdrant向量，双侧班级与授权材料过滤 |
| T4 AI与工具 | agents/skills/stream/conversations/mcp与测试、mcp-server | 三个独立内置应用（知识问答/数学引导/语文引导）、版本模板、Skill导入、流式引用与真实MCP；重做出站白名单和统一审计 |
| T5 总验收 | apitest按URL组织方式 | 重写Plus统一响应断言、独立A/B班样本、容器故障与固定报告，不搬GD报告当成绩 |

发现需重点复核的差异：GD main包混HTTP/SQL；部分响应camelCase；路由带classId；retrieval_scope缺上下文时默认放行（调用链还需审查，不直接认定漏洞）；MCP已有白名单/禁跳转但需复测；GD日志偏工具调用，不等于全操作审计。旧go.mod和前端锁文件不能覆盖Plus版本，新增依赖审计后固定。

## 身份与接口契约

JWT仅内存，HS256允许算法固定、密钥至少32字节随机值；issuer=campusclaw-plus、audience=campusclaw-plus-api，期限30分钟，校验exp/nbf/iat/sub/jti。退出撤销jti；敏感请求每次校验用户状态、成员资格与角色，失效身份默认拒绝。DB故障返回503，不能当匿名成功。刷新页面重新登录，首版不引入refresh Cookie。

POST /api/auth/login显式公开（其他业务默认鉴权）；GET /api/me、POST /api/auth/logout。登录失败统一401；每账号指纹+来源5次失败/15分钟后429，指纹不可明文记账号。种子用A/B班教师学生四角色样本，四个种子密码独立、互不相同，运行时注入并哈希，不在Git保存值。

每个登录身份只有一个当前认证班级；首版不开放切班。客户端class_id/classId均不能覆盖；未知/重复参数400。GD多班级路由转换成/api/materials、/api/retrieval/search、/api/retrieval/chunks/{id}、/api/agents、/api/skills、/api/conversations及其messages、/api/mcp/servers及connect/disconnect/tools、/api/audit-events。后续多班切换须新规约并验证成员资格。

业务响应统一{code,data?,error?,request_id}；401 UNAUTHENTICATED、403 FORBIDDEN、404 NOT_FOUND、400 INVALID_ARGUMENT、409 CONFLICT、429 RATE_LIMITED、503 DEPENDENCY_UNAVAILABLE。异步202 ACCEPTED等按异步design，不重定义容量码。跨班和不存在404同结构；本班检索不允许自选他班范围，结果只能含本班，故意提交class_id返回400而非忽略。资源读取/更新/删除在SQL条件直接限定授权范围，向量查询也过滤，回MySQL再次约束有效版本。教师配置写入逐次验角色；学生会话只许本人读写。文件通过材料鉴权接口读取，uploads不静态开放。

## SQL / XSS / 跨站请求

SQL参数绑定，排序白名单、分页默认20/上限100；' OR '1'='1作为普通检索值，不扩大结果或写库。标题/工具返回按文本渲染；检查DOM无script节点及HTML已转义，不能只看弹窗。

Bearer不会被浏览器自动附在跨站表单中：无Authorization表单应401且DB不变；有恶意Origin的写请求403，CORS仅精确同源、不允许任意带凭证跨域，生产同源。非浏览器请求无Origin时仍需Bearer及JSON内容类型/明确允许的multipart上传；登录浏览器Origin白名单、拒绝跨站表单和非JSON登录。拒绝Cookie-only及URL token认证。首版无认证Cookie，一次性CSRF Token和SameSite认证Cookie记“此认证设计不适用”，不是伪造已实现；若未来引入Cookie必须先补CSRF token+SameSite/Secure/HttpOnly后开放。

## AI / Skill / MCP边界

选中应用模板与Skill正文进system，材料摘录和工具结果作为user数据块；只传top5且总摘录≤6000字符，其他技能/全库不进入请求。应用/模板/Skill保存revision与实际快照。单次模型60秒、输出≤700 tokens、temperature=0、每轮最多2工具、首版流式中断不自动重试；外部失败返回可判错误，不制造引用。向量调用/重试使用异步约束；对话不排任务队列。

Skill纯文本UTF-8 md/txt≤64KiB，导入只解析/预览，教师确认后启用，不执行文件代码。MCP服务端env密钥引用；HTTP响应/前端/审计不得带回值。出站固定精确origin白名单，连接目标与DNS解析受控，拒绝重定向与白名单外/元数据地址；只例外允许配置的Compose内部MCP服务，不泛化禁止全部内网导致演示服务不可用。断开清缓存；保存技能和实际调用两处校验本班/已连接/已发现工具，工具调用仍经业务服务授权。

当前只提供只读formula_hint/unit_convert演示工具，无分数写工具。SEC-04覆盖所有已注册写服务直接调用与模型适配入口：匿名/学生/跨班均不能越权，不给模型直连DB。真正底层批改接口不存在，不能把404算“批改授权通过”，该子项保留待作业模块补齐。

## 七列操作审计

新增audit_events：actor_id（匿名null）、occurred_at UTC、operation、object_type/object_id、outcome与http_status、skill_id（无则null）、usage（未知unknown）与duration_ms；另有id/class_id/request_id/task_id/turn_id。七列语义允许用量与耗时同组；正常操作、401/403/404、工具失败均记。匿名不信任body的user_id/class_id，已认证越权记录其真实班级。

事件序列化为单行JSON存MySQL JSON列并配索引(class_id,occurred_at,id)、(class_id,actor_id,occurred_at,id)。只写白名单元数据，每文本字段≤256字符，单事件≤8KiB；转义CR/LF，不记录正文、提示词、Authorization或秘密。教师仅可按本班user/time/action/status查，学生403。每次调用的skill_id、turn_id与审计关联可回查。

配置等DB写操作与审计同事务，审计写失败回滚业务并503；拒绝请求审计写失败也503并结构化error日志，不返回伪造“403且已审计”。外部调用先写intent，完成再写结果；落结果失败标为未知可恢复，不谎报外部动作回滚。只读审计故障返回503，所有动作默认不在无审计保障下继续敏感执行。审计默认30天，备份7天，按数据规范留存删除。

## 安全门禁、证据与验收

make check-security=全历史/构建上下文/前端产物秘密扫描+依赖审计+安全用例；失败/审计源不可达非零。保留make audit为依赖审计入口，CI必需check内调用安全门禁；重复扫描合并但不减覆盖。升级可修高危；接受例外须记录CVE/包/版本/责任/补偿/截止时间，门禁显式消费有效例外，不隐性放行。

红绿门禁在隔离临时fixture/分支用人工测试标识或已知脆弱依赖快照演练，记录SHA/输出/退出码/测试性质，不往真实历史注入有效密钥，也不把合成命中当真实泄露。真实命中先撤销/轮换，再单独呈现历史清理影响供维护者批准；本计划不授权破坏性重写/通知其他人。项目不建设DB存模型密钥字段，密文存储条目记不适用，不为满足示例引入密钥托管。

SEC-01匿名、02归属、03角色、04服务层写权限、05资产、06依赖、07审计、08四步AI实验固定三段式。先基线再修，按同一输入/期望复测，只更新实际；无历史红则记新建基线/能力缺口，不伪造旧漏洞。实验：恶意材料、索要system、修改规则、诱导越权工具；前三项记录真实模型现象，不预填绿；第四项必须拒绝副作用且留审计。没有评分业务，评分偏移子实验保留缺口，不能以教学引导替代完整评分实验。

脚本apitest/cases按URL，正常/失败/边界均核对code、数据副作用和泄露。独立Docker环境执行make check、check-security、异步三判据与SQL/XSS/CSRF三件套；grade只需宿主Docker+sh，输出HTML/JSON、截图及以下固定reports文件：security-fix-log.md、security-baseline.md、risk-authz.md、gate-red-green.md、sec08-experiment.md、review-02.md、review-03.md、retrospective.md。reports默认忽略，原始报告不入Git；脱敏固定路径副本按白名单允许入库，docs只索引，下载完整报告也保留上述路径，不复制秘密。尚无实测不创建占位“通过报告”。

最终演示：登录→材料返回202→状态pending/running/succeeded→B班查A材料404；加教师回查学生403审计、MCP断开拒绝。班级材料演示与未完成作业评分分别注明。无人类同伴则至少两份AI/脚本证据且明确不是独立人类交叉评审；异机未测不声明已完成。返工只记真实发生的缺陷与耗时，未发生就如实记无，不编造一次。

## 与现有异步change的衔接

2026-11-04-add-reliable-async-tasks的四文档保留不动、未批准未实施。本变更先交付T0/T1及材料资源前置，再实施该change的队列/任务与材料处理；在异步spec的同一场景加身份/审计回归。最终本变更再验收检索/AI和综合安全。两个change分别按实际任务验收归档，不提前勾完任一个。

## 八规范待完善（实施时回填）

产品：登录/任务/审计四态与术语；后端：身份上下文/错误码/幂等；安全：SEC矩阵与Bearer跨站判据/密钥处置；测试：固定三段式与红绿门禁；运维：安全依赖、worker、报告与恢复；AI：12题/四步实验/出站；数据：身份/审计/索引版本迁移留存；评审：来源清单与权限/并发/事务证据。检查统一为Scenario→脚本→HTTP/DB/DOM/审计→报告，不新建平行规范正文。

## 实施记录

已批准按阶段推进。认证/审计首批跨越领域端口、SQL迁移、JWT适配、HTTP与容器配置，含必要新测试，超过400行；按各职责文件拆分，基础设施与单测单列，不整体复制GD。先固定规划，再测试来源与Plus，再提交首阶段；未完成的后续模块不提前标完成。

## AI接入暂缓

2026-11-07用户明确先完成测试替身验收，AI API暂不填写/接入。AI_ENABLED为必填开关；false时不要求任何模型/嵌入地址或密钥，材料仅发布MySQL切片供keyword检索，向量/混合和对话不伪装可用。true时必须配置完整供应商参数。Docker验收显式启用fixture，不能据其宣称真实AI已联调；后续不再新增AI接口实现，真实AI验证待单独恢复。嵌入超时env外置，生产建议30秒，延迟30秒的机制测试使用35秒调用上限/40秒任务上限，避免把依赖超时误当队列丢任务。

测试修复记录：HTTP总超时与ResponseHeaderTimeout必须一并按配置更新，防止30秒延迟实验仍被旧30秒头部超时截断。已补该回归。报告收集queue.log定位Worker/dispatcher/broker，而不记录原始凭证或正文。
