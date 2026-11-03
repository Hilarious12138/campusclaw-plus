# 最终核对清单：评分标准 × 下一阶段实践

核对日期：2026-11-03（Asia/Shanghai）；基线e883f5d，规划提交9d69412。本表依据用户给出的A30/B20/C30/D20评分图，以及[下一阶段AI代码审核课件](https://devops.hello1023.com/小组项目/小组项目下一阶段重点-AI代码审核.html)。不代教师评分，也不把八份规范存在等同于所有业务通过。

状态口径：**已核对**指本次检查到明确文件/代码/记录；**规范就绪/实践待补**指要求写入但还没有对应业务验证；**待补**指尚无完成证据；**不适用**必须给理由。代码阅读和文档检查由AI辅助执行，不能当作项目负责人已完成人工审核或同学peer review。

## 一、A：落位与可加载（30分）

| 编号 / 逐项问题 | 实际证据 | 结论 / 补齐动作 |
| --- | --- | --- |
| A1 常驻入口是否在Harness读取位置 | [AGENTS.md](../../AGENTS.md)位于Plus根目录；说明从Plus根目录启动 | 位置已核对；不能从父目录启动就假定子项目入口已加载 |
| A2 分类规则是否有明确读取路径与条件 | AGENTS八份索引与适用条件，目标路径实际存在 | 已核对；普通rules目录不宣称自动全部加载 |
| A3 frontmatter是否符合实际Harness | [SKILL.md](../../.agents/skills/verify-change/SKILL.md)含name/description，check.py检查该元数据 | 已核对现有格式；Codex AGENTS为Markdown，未把globs/alwaysApply套到不支持的入口上 |
| A4 是否有总纲和操作流程 | AGENTS每次变更/安全边界/索引；verify-change步骤 | 已核对；常驻约束与按场景流程分开 |
| A5 是否有实质内容并实际生效 | 八规则均非空且带检查方式；原骨架有空白会话记录；本次执行静态和变异检查 | 部分具实证；新增八份规则在新Harness中的独立读取/遵循复验仍待补，不能用旧五类复核顶替 |

## 二、B：八类规范覆盖（20分）

| 编号 / 规范 | 文件与具体内容 | 本次结论 |
| --- | --- | --- |
| B1 产品 | [product.md](../../.agents/rules/product.md)：用户/场景、路由、四态、术语、视觉/表单 | 规范已核对；业务页面/导入尚待对应能力 |
| B2 后端 | [backend.md](../../.agents/rules/backend.md)：依赖方向、响应、错误码、参数、幂等与事务 | 规范已核对；骨架分层/错误接口有测试，业务写入待补 |
| B3 安全 | [security.md](../../.agents/rules/security.md)：默认拒绝、归属、密钥、输入输出/附件/依赖 | 规范已核对；历史扫描/默认拒绝适用，真正登录/跨班授权未实现 |
| B4 测试 | [testing.md](../../.agents/rules/testing.md)：三类用例、隔离、机器码/副作用/泄露判定 | 规范已核对；本次重跑既有门禁，不把文档检查当业务覆盖 |
| B5 运维 | [operations.md](../../.agents/rules/operations.md)：Compose、health、JSON日志、停机/恢复目标 | 规范已核对；骨架有历史Linux报告，业务备份/恢复待补 |
| B6 AI | [ai.md](../../.agents/rules/ai.md)：调用边界、提示词版本、引用、注入、12题评测 | 规范就绪/实践待补；Plus未实现教学模型调用 |
| B7 数据 | [data.md](../../.agents/rules/data.md)：归属、迁移、敏感分级、留存及副本删除 | 规范就绪/实践待补；当前只有MySQL探测，没有业务表/向量 |
| B8 评审 | [review.md](../../.agents/rules/review.md)：Done、七类/七维、三分法、准入/例外 | 规则与[模板](../reviews/template.md)已核对；人工/同学评审未完成 |

## 三、C：可判定性（30分）

| 编号 / 问题 | 可重放判据与证据 | 结论 |
| --- | --- | --- |
| C1 能否判断违反分层 | service导入database/sql时check.py应非零；[mutation.py](../../scripts/mutation.py)负责探针并恢复 | 本次已重放并恢复通过，见第七节 |
| C2 失败是否会被当成成功 | 未登记API若改200/OK，[handler_test.go](../../backend/internal/httpapi/handler_test.go)负向测试必须失败 | 本次已重放并恢复通过，见第七节 |
| C3 是否有前置/操作/可观测结果 | [platform-foundation spec](../../openspec/specs/platform-foundation/spec.md)有S01-S16，覆盖参数、拒绝、资源、依赖 | 已核对结构；骨架404不冒充JWT权限测试 |
| C4 是否有明确数值和信号 | 请求ID1/64边界、依赖超时边界、HTTP与code、数据/日志判据；AI规则要求12题 | 骨架有测试；12题AI评测目前只有要求，无实际完成证据 |
| C5 无证据、跳过或造成功能否挡住 | [test_grader.mjs](../../scripts/test_grader.mjs)覆盖断言失败/缺阶段/清理失败/缺env；grade程序非零 | 已有可执行检查，本次make check重跑；不能由此推导无其他缺陷 |

## 四、D：组织与自洽（20分）

| 编号 / 问题 | 对应位置 | 结论 / 补齐动作 |
| --- | --- | --- |
| D1 规则是否分层 | AGENTS入口→分类规则→场景Skill；spec/任务/证据各自独立 | 已核对；不复制八份正文到讲稿或清单 |
| D2 与规约提案是否挂钩 | openspec/changes/archive保存proposal/design/spec/tasks；长期spec有能力规则 | 已核对；当前变更由独立四文档先规划 |
| D3 能否追溯实现与证据 | [verification.md](verification.md)、PR/Actions、[review记录](../reviews/2026-11-03-finalize-rubric-audit.md) | 已核对来源；历史报告仅证明其提交，未补造人类签字 |
| D4 是否写“怎么检查”且无空壳 | 八规则检查列、脚本命令、Scenario、实际记录 | 已核对；template中的待填是复制入口，不算完成的评审证据 |
| D5 准入规则是否与实际状态一致 | 仓库单人审批0；课件要求独立同学评审；现有文档明确差异 | **部分满足**：不能声称课程流程全部完成；后续业务须补评审或获教师明确认可的适配，不能由CI替代 |

## 五、课件第一层：七类对象逐段读

已读取现有handler、service、repository、lifecycle、main及前端状态映射。以下是AI辅助的骨架检查范围，不是独立人类审核。

| 对象 | 本次读到什么 / 依据 | 已做与待补 |
| --- | --- | --- |
| 接口入口 | [handler.go](../../backend/internal/httpapi/handler.go)检查路径、方法、query，响应含code/request_id | 骨架已有实现/测试；业务参数、登录和归属待实现后核对 |
| 业务逻辑 | [health.go](../../backend/internal/service/health.go)将依赖状态映射就绪结果，设置有界超时 | 骨架范围可核对；教学状态流转待业务 |
| 数据库操作 | [mysql.go](../../backend/internal/repository/mysql.go)仅连接/Ping，无业务SQL | 事务、唯一约束和写入并发不适用本骨架，不能据此说已验证 |
| 公共函数 | 生命周期、健康判定、前端状态各有明确职责；重试/错误判断有测试 | 已检查当前职责和局部复用；无正式重复比例测量，不编“0%重复” |
| 错误处理 | 路径/方法/参数/依赖各有明确失败码，配置/停机有测试 | 已核对骨架；业务约束冲突和模型故障待补 |
| 安全与隐私 | 未注册入口默认拒绝、query不入请求日志、配置值不进样例、前端无危险输出写法 | 已有静态/测试；业务越权、提示注入和凭证撤销仍待补 |
| 性能 | DB连接池/超时有界；无列表/循环业务SQL | 循环查询、索引与分页专项暂无适用业务，不能称数据库性能审核完成 |

审核四问已进入[roadmap](../roadmap.md)：要保证什么→靠什么保证→怎样会失效→怎样修复/证明。最终稿以未来重复提交作例，不假称现有作业功能。

## 六、课件第二、三层：打分、机制与准入

### 七个核对维度

| 维度 | 已落在哪里 | 当前满足程度 |
| --- | --- | --- |
| 需求覆盖率 | template逐Scenario已实现/部分实现/未实现；spec映射 | 表与骨架证据可用；无业务覆盖率统计 |
| 与需求的一致性 | template与每次spec/design、实际文件记录 | 文档与骨架范围已核对；业务状态和字段尚待实现 |
| 失败路径覆盖 | handler/config/报告测试，模板故障/冲突栏 | 已有骨架失败验证；教学写入路径待补 |
| 查询性能风险 | backend/data规则、template查询栏 | 核查要求就绪；业务查询工具/EXPLAIN证据未具备 |
| 重复代码比例 | template要求具体范围与测量口径 | 核查要求就绪；没有量化测量结果 |
| 漏洞风险 | security规则、[audit.sh](../../scripts/audit.sh)、CI | 自动扫描入口存在；本次CI按实际结果记录，不能承诺整库无漏洞 |
| 高风险场景 | template分别记录权限/并发/事务/一致性及依据 | 骨架无业务写入理由明确；业务专项与人工结论待补 |

### 结论与固定机制

| 要求 | 证据 / 动作 | 结论 |
| --- | --- | --- |
| 每文件分可用/待修改/风险无效 | review规则与template；本次文档记录逐文件分类 | 形式与动作就绪；不是所有文件都已经被人类判可用 |
| 待修改重审，风险无效重写 | review明文阻断条件 | 规则就绪；未来缺陷实例须留真实处理链 |
| 工具批量扫描 | check/audit/grade，CI执行；SQL专项没有工具 | **部分满足**；现有工具覆盖有限，业务查询扫描待补 |
| 人工重点复核 | 模板记录查了什么/规约依据/结论 | **待补人工执行**；AI读取代码不算本人或同学逐行审核 |
| 发现反哺规则 | verification记录旧日志故障、分层/默认拒绝探针；规则已有相应检查 | 有历史与本次检查证据；未制造新的业务缺陷案例 |
| AI初稿说明范围 | change proposal/design与PR记录 | 本次有记录，纯文档未生成业务代码 |
| 自己逐段审核 | 本次AI辅助记录，没有项目负责人确认 | **待补** |
| 同学交叉评审 | 当前无另一位同学的意见/签字 | **待补**，需教师适配或同学参与 |
| 验证、证据归档后合入 | 既有PR与check/grade流程；单人审批0不等于peer通过 | 自动门禁已落地，课程完整准入**不能宣称全部满足** |

## 七、本次实际检查与展示口径

- 四份规划先于本次文档实施；不修改GD，不实现新教学业务，不改技术栈。
- 本次已实际通过：八规则非空/检查列、AGENTS索引、Skill name/description格式、全部本地链接、README边界、评分/课件项目覆盖、strict、静态/全历史检查、make check（退出0）。`mutation.py`退出0：反向依赖被阻断、错误放行导致指定测试失败；finally还原后静态和指定HTTP测试重新通过，`git diff -- backend frontend`为空。变异原始结果位于忽略Git的reports/mutation.json；本表和本次审核记录保留脱敏摘要。
- 历史[teacher-report.json](teacher-report.json)/[HTML](teacher-report.html)生成于2026-10-28，source_commit写在报告中。此前骨架Linux15项通过，不等于教学业务或独立评审通过。
- 实施版本bb6e08c的[check/grade均成功](https://github.com/Hilarious12138/campusclaw-plus/actions/runs/41628412280)，本次[PR6](https://github.com/Hilarious12138/campusclaw-plus/pull/6)可追溯；归档后的最终版本仍须门禁通过，未通过前不称已合并。
- 尚待补齐：本次新增规则的新Harness遵循复验、项目负责人逐段审核、同学peer或教师认可适配、业务SQL/并发/事务/一致性/AI专项实证。已写进规范和路线图；未来能力触发对应任务，不另造一套待办台账。

## 八、课堂怎么用

按[最终讲稿](showcase.md)打开README→AGENTS→roadmap→评审模板→核对表/真实证据。主讲“已做什么、为什么、下一步”，细节只作追问。不讲得分预测，不说全项目已完成。
