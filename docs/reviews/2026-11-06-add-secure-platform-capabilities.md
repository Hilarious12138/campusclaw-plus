# 安全能力转换与异步任务复核

受验代码717c14c569c67b8e02426a995159de0ee5f5d7c3；规划4fb6f39先于实现。固定干净克隆的make check退出0，22项通过。后续文档提交只补充此代码的证据。990968c52b9d5302cc20dcdf4206dcf4947b0603的GitHub工作流另行通过，见下表。

## 七类文件与结论

| 文件组 | 核对结果与依据 |
|---|---|
| internal/httpapi、cmd | 可用：身份上下文来自Bearer与数据库，状态/机器码按OpenAPI，B01匿名/学生/跨班与故障回归 |
| internal/service | 可用：领域端口依赖，材料事务接收、技能选择、逐资源授权；B01与分层扫描 |
| internal/repository、迁移 | 可用但补测项保留：审计事务、幂等、索引版本、lease/fence；B01/B04/A01—03；升级回滚未实演 |
| 公共函数、解析、密码/JWT | 可用：转换来源摘要齐，哈希/算法/issuer/audience/期限及切片/SSE测试；不继承GD完成声明 |
| 错误处理与配置 | 可用：非法配置启动失败，依赖故障不伪装空命中；总超时与HTTP头超时一致已修复 |
| 安全与隐私、Compose/Nginx | 可用：秘密仅后端env、无上传静态目录、仅web映射，扫描/漏洞审计与浏览器转义通过 |
| frontend、apitest、scripts | 可用：实际Chromium截图和API/DB/故障证据；fixture明确标注，缺阶段/清理失败使验收失败 |

结论针对上述代码组，不表示每条未来需求全部完成。大范围来源转换原因及分阶段提交见design；锁文件/生成OpenAPI单列，不以生成量掩盖业务变动。

## 七维结论

需求覆盖：当前机制22项通过，尚有tasks明确未验项，未编覆盖百分比。需求一致性：AI API按用户要求关闭。失败路径：401/403/404/409/429、审计写失败、数据库与broker不可达有实际断言。查询性能：只测受理P95与健康延迟，不外推生产容量。重复代码：旧单包转换分层，未做量化重复率测量。漏洞：npm/Go最终无漏洞，实际旧依赖红记录保留。高风险：权限/并发/事务/一致性分别见reports/risk-*.md和请求追踪。

## 实际返工

AMQP旧版本漏洞升级；SQL认证插件和Rabbit交换机权限修复；浏览器选择器修复；30秒头超时与重投ACK等待修复。固定717c14c5重新执行成功。运行中改shell导致过一次无效失败，之后冻结克隆，未删除失败记录或修改判据。

## 限制及准入

AI辅助自审与脚本检查均已记录，独立人类评审、异机启动和真实模型实验未完成。单人三职责不等同三人。当前change继续开放；未归档、未打完成标签、未将未完成项算通过。真实AI凭证无需填写。自动门禁不能替代人工评审；合并前仍须相应PR检查和维护者处理未完成项。

## 范围、职责与检查入口

记录日期：2026-11-07，Asia/Shanghai。安全change为`2026-11-06-add-secure-platform-capabilities`，队列change为`2026-11-04-add-reliable-async-tasks`；本次补记录沿用前者任务5.3/5.5，不另建重复规划。

AI参与了来源转换、代码与测试生成及本记录整理。当前审核为AI辅助源码复核和自动测试核查；维护者仍须处理准入结论，本文没有独立人类评审者签名。范围是相对`origin/main`的变更文件；逐文件清单区分直接复核与测试证据，不宣称全部人工逐行审核。

依据：当前change的R1—R8/S01—S18、异步change的A1—A8，以及根AGENTS和八份工程规范。证据入口：[场景映射](../verification/secure-platform-scenarios.md)、[验收摘要](../verification/secure-platform-result.json)、[来源摘要](../verification/gd-conversion.md)、[漏洞修复历史](../verification/security-fix-history.json)。

## 场景完成状态

| 场景 | 状态与判定边界 |
|---|---|
| S01—S02 来源与分层 | 已验证源码来源、GD隔离、原测试及Plus检查；不包含GD真实模型集成 |
| S03—S05 认证与资源授权 | B01及JWT/哈希测试通过；批改服务未实现，不将404当批改安全证据 |
| S06—S07 SQL/XSS/跨站 | B01/B02通过；无认证Cookie，Cookie型CSRF机制不适用 |
| S08—S09 审计与回滚 | B01核对七列、范围及审计失败时无业务半写 |
| S10—S11 安全门禁 | 人工秘密fixture红绿、实际依赖红绿有记录；公告源故障出现过非零，不改为成功 |
| S12 异步任务 | A01—A03/B04覆盖受理、停机恢复、重复、租约及死信；完整异步change仍有未重放窗口 |
| S13 检索 | fixture下三模式、MySQL正文回查与权限机制通过；真实嵌入未验证 |
| S14 模板/Skill/会话 | 机制已验证；真实模型12题评测未执行，整体为部分完成 |
| S15 MCP | 真实MCP发现、调用、断开与越权拒绝已验证 |
| S16 对抗实验 | 未执行，用户暂缓真实AI；评分业务也未实现 |
| S17—S18 交付与复测 | Docker机制及既有修复证据已记录；独立人类/异机及完整实验仍未完成 |

## 两条实际链路的复核

正常：教师Bearer经`httpapi/platform.go`识别身份→`service/tasks.go`检查角色、UTF-8/2MiB/文件名/幂等键→repository事务写材料/任务/outbox/审计→202→dispatcher确认投递→worker租约执行→MySQL发布版本。B01/A01/A02核对HTTP、任务终态和数据库结果，不只看页面。

失败：学生上传在服务层拒绝403；跨班读取由repository班级条件产生404。配置写入遇到审计失败，事务回滚并返回503，B01核对数据库不变。旧worker租约过期后提交返回冲突，B04核对最终结果仅一份。

## 高风险复核明细

| 风险 | 源码/对象 | 判据与规约 | 证据与限制 |
|---|---|---|---|
| 权限 | Platform身份识别、Resources.Save、Materials.Create、MCP.Register/Connect及SQL范围 | R2/S04—05：身份来自Bearer/数据库，学生不能写，跨班不可读 | B01请求追踪与DB断言；不依赖前端隐藏按钮 |
| 并发 | ClaimTask、CompleteTask、tasks唯一键 | 异步A3/A4：失效fence不得提交，同键不重复任务 | B04真实MySQL20并发；双worker全部窗口尚待补，不推导全覆盖 |
| 事务 | CreateMaterial、SaveResource、审计/outbox | R4/S09及异步A2：多写原子，事务内无外部网络 | B01故障回滚、A02 broker恢复；数据库升级回滚未实演 |
| 一致性 | 确认投递、ACK、索引版本与MySQL回查 | 异步A4/A7：重投不重复生效、不可见半版本 | A01 ACK前kill与唯一结果；部分索引崩溃窗口未重放 |
| 外部调用 | MCP白名单、工具发现缓存、AI开关 | R7/S15：只允许本班已连接工具，数据不当系统指令 | 真实MCP用例通过；真实模型输出安全没有结论 |
| 性能 | 材料202受理与分页 | 异步20次/组P95<1000ms；分页上限100 | 两组12.85ms/3.46ms；不是处理完成耗时，也不是生产压测 |

## 提交与检查对应

| 对象 | 实际结果 | 可追溯入口 |
|---|---|---|
| 717c14c569c67b8e02426a995159de0ee5f5d7c3 | 干净副本make check退出0；Docker Linux arm64，22/22 | secure-platform-result.json记录原报告SHA256；本地原始报告不入Git |
| 990968c52b9d5302cc20dcdf4206dcf4947b0603 | GitHub verify已完成success，两次运行 | [运行42255197631](https://github.com/Hilarious12138/campusclaw-plus/actions/runs/42255197631)、[运行42255141482](https://github.com/Hilarious12138/campusclaw-plus/actions/runs/42255141482) |
| 本次文档补全 | 以提交后的CI为准，不预填通过 | [草稿PR11](https://github.com/Hilarious12138/campusclaw-plus/pull/11) |

## 逐文件清单

“可用（自动证据）”只表示其适用机器检查通过，不等于独立人工审阅。“待复核”不是已发现漏洞，而是仍欠指定验证。生成物和锁文件单列。

| 文件 | 结论 | 依据 / 后续动作 |
|---|---|---|
| `.agents/rules/ai.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `.agents/rules/backend.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `.agents/rules/data.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `.agents/rules/operations.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `.agents/rules/product.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `.agents/rules/review.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `.agents/rules/security.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `.agents/rules/testing.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `.env.example` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `.github/workflows/check.yml` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `.gitignore` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `AGENTS.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `Makefile` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `README.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `apitest/async_cases.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/broker_cases.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/browser.mjs` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/cases/api/audit-events/GET.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/cases/api/auth/login/POST.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/cases/api/conversations/POST.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/cases/api/materials/POST.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/cases/api/mcp/servers/POST.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/cases/api/retrieval/search/POST.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/cases/api/skills/POST.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/cases/openapi.json/GET.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/common.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/evals/README.md` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/evals/cases.json` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/fixtures/server.mjs` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/resilience_cases.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `apitest/run_all.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/cmd/api/main.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/cmd/job/main.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/cmd/mcp-server/main.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/cmd/mcp-server/main_test.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/go.mod` | 可用（依赖/生成物检查） | C03/C04/C06/C07及接口契约用例；不按手写业务计规模 |
| `backend/go.sum` | 可用（依赖/生成物检查） | C03/C04/C06/C07及接口契约用例；不按手写业务计规模 |
| `backend/internal/config/config.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/config/config_test.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/domain/ai.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/domain/platform.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/domain/tasks.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/httpapi/ai.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/httpapi/handler.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/httpapi/materials.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/httpapi/platform.go` | 可用（源码复核+测试） | 本次直接读取；七类/高风险表与B01/B04/A01—03；独立人类评审待补 |
| `backend/internal/httpapi/platform_test.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/auth.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/conversations.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/credentials.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/credentials_test.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/index.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/index_test.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/mcp.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/mcp_remote.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/model.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/mysql.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/password.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/password_test.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/platform-schema.sql` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/rabbit.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/resources.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/retrieval.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/seed_apps.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/stream.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/stream_test.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/tasks.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/repository/tasks_integration_test.go` | 可用（源码复核+测试） | 本次直接读取；七类/高风险表与B01/B04/A01—03；独立人类评审待补 |
| `backend/internal/service/auth.go` | 可用（源码复核+测试） | 本次直接读取；七类/高风险表与B01/B04/A01—03；独立人类评审待补 |
| `backend/internal/service/chunks.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/service/chunks_test.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/service/conversations.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/service/health.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/service/mcp.go` | 可用（源码复核+测试） | 本次直接读取；七类/高风险表与B01/B04/A01—03；独立人类评审待补 |
| `backend/internal/service/resources.go` | 可用（源码复核+测试） | 本次直接读取；七类/高风险表与B01/B04/A01—03；独立人类评审待补 |
| `backend/internal/service/retrieval.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/service/skill_import.go` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `backend/internal/service/tasks.go` | 可用（源码复核+测试） | 本次直接读取；七类/高风险表与B01/B04/A01—03；独立人类评审待补 |
| `compose.grade.yaml` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `compose.test.yaml` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `compose.yaml` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `deploy/api.Dockerfile` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `deploy/fixture.Dockerfile` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `deploy/grade.Dockerfile` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `deploy/nginx.conf` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `deploy/rabbitmq.conf` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `deploy/web.Dockerfile` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `docs/deployment.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `docs/reviews/2026-11-06-add-secure-platform-capabilities.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `docs/roadmap.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `docs/verification/gd-conversion.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `docs/verification/platform-foundation.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `docs/verification/secure-platform-result.json` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `docs/verification/secure-platform-scenarios.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `docs/verification/security-fix-history.json` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `frontend/package-lock.json` | 可用（依赖/生成物检查） | C03/C04/C06/C07及接口契约用例；不按手写业务计规模 |
| `frontend/package.json` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `frontend/public/openapi.json` | 可用（依赖/生成物检查） | C03/C04/C06/C07及接口契约用例；不按手写业务计规模 |
| `frontend/scripts/swagger.mjs` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `frontend/src/Workspace.tsx` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `frontend/src/api.ts` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `frontend/src/main.tsx` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `frontend/src/styles.css` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `mcp-server/README.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `openspec/changes/2026-11-04-add-reliable-async-tasks/design.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `openspec/changes/2026-11-04-add-reliable-async-tasks/proposal.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `openspec/changes/2026-11-04-add-reliable-async-tasks/specs/async-task-processing/spec.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `openspec/changes/2026-11-04-add-reliable-async-tasks/tasks.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `openspec/changes/2026-11-06-add-secure-platform-capabilities/design.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `openspec/changes/2026-11-06-add-secure-platform-capabilities/proposal.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `openspec/changes/2026-11-06-add-secure-platform-capabilities/specs/platform-foundation/spec.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `openspec/changes/2026-11-06-add-secure-platform-capabilities/specs/portable-grading/spec.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `openspec/changes/2026-11-06-add-secure-platform-capabilities/specs/project-documentation/spec.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `openspec/changes/2026-11-06-add-secure-platform-capabilities/specs/secure-platform-capabilities/spec.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `openspec/changes/2026-11-06-add-secure-platform-capabilities/tasks.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `reports/gate-red-green.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `reports/retrospective.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `reports/review-02.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `reports/review-03.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `reports/risk-authz.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `reports/risk-concurrency.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `reports/risk-transaction.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `reports/sec08-experiment.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `reports/security-baseline.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `reports/security-fix-log.md` | 可用（文档/配置检查） | strict、分层/秘密/配置检查及证据交叉核对；未完成项仍以tasks为准 |
| `scripts/audit.sh` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `scripts/build-openapi.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `scripts/check-security.sh` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `scripts/check.sh` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `scripts/clean_clone.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `scripts/gate-evidence.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `scripts/grade.sh` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `scripts/grader.mjs` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `scripts/integration.sh` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `scripts/npm-audit.mjs` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `scripts/security_scan.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `scripts/test-env.mjs` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `scripts/test_grader.mjs` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |
| `scripts/test_security_scan.py` | 可用（自动证据） | 717c14c5对应C/B/A/R测试范围；不声称本次全部逐行重读 |

## 待补项与最终准入结论

- 真实模型与嵌入、12题评测和四步AI安全实验：用户要求暂缓；保留开关关闭。恢复须用户明确要求并配置服务端环境变量，不能用fixture结论替代。
- 双worker竞争、确认丢失、索引崩溃与数据库升级回滚等尚未完整重放：按异步tasks补齐原判据后复核，不放宽期望。
- 独立人类交叉评审：未完成，未获得等价豁免；补齐条件是另一位实际审阅者记录范围、意见和处理。AI自审/CI只作替代检查资料，不能自行宣告等价。
- 异机复验：GitHub Linux CI提供另一运行环境证据，但不等于另一位成员按README独立复现；后者仍未完成。

责任职责：项目维护者处理剩余验证、评审意见与准入；当前仅授权补记录和提交GitHub。结论为“已验证机制可继续迭代，整个变更尚未完成”；PR保持草稿，不合并、不归档、不打完成标签。没有发现可据此记录为“风险无效”的文件；未验证风险已逐项列出，不把空白结论填成通过。
