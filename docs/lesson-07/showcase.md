# 第7课仓库展示路线（约8分钟）

展示日期与实际课次安排为准。这里按评分A30/B20/C30/D20组织证据，不预估教师得分。只展示CampusClaw Plus，GD是经验来源，不是本项目已实现业务。

## 课前准备

- 确认GitHub main已含本次规范：若只在本地或PR中，先如实指出分支，不把未同步文件说成主干已有。
- 开仓库README、AGENTS、八份rules、verify-change Skill、骨架spec及真实报告几个页签。GitHub报告HTML可下载后浏览器打开；JSON可直接查看。
- 提前部署并跑`sh scripts/grade.sh`，打开本次输出report.html。现场不等首次镜像/Chromium下载；历史报告备用但必须讲明日期与提交。
- 不展示.env内容、访问令牌、完整环境或含凭证的终端命令。

## 0:00—0:40 项目和边界

打开[README](../../README.md)。

> CampusClaw Plus面向班级教学，后续连接资料检索、学习与作业订正。本课交付工程基线，不展示尚未实现的登录、模型或作业业务。我是单人项目，以产品、开发、测试三种职责组织工作，不能等同于三位独立成员。

## 0:40—2:00 A：规则放在哪里、如何加载（30分）

打开[AGENTS.md](../../AGENTS.md)，指“每次变更”和“规范读取索引”。

> 根AGENTS是常驻入口：先四份规划文档、给概览、确认后实施。八份规则按影响读取，避免每轮把全部正文塞进上下文。仅放在rules目录不代表自动生效，所以入口有明确路径和读取条件。

再打开[verify-change/SKILL.md](../../.agents/skills/verify-change/SKILL.md)。指出开头`name`、`description`与正文检查流程。

> 我用Codex：AGENTS提供常驻约束，Skill的name/description用于发现、正文在触发后读取。截图里的globs/alwaysApply是需要结合Harness解释的字段，我没有把它们伪装成Codex八份普通Markdown的自动加载机制。

官方依据：[AGENTS](https://learn.chatgpt.com/docs/agent-configuration/agents-md)、[Skills](https://learn.chatgpt.com/docs/build-skills)。从Plus仓库根目录启动新的Harness；可以让它列出读取来源，检查是否指出AGENTS和适用规范。文件存在/元数据有效不等于已经在另一工具验证加载。

## 2:00—3:10 B：八类规范覆盖（20分）

打开`.agents/rules/`，快速指八个文件，不逐字朗读。

| 文件 | 一句话说明 |
| --- | --- |
| [product.md](../../.agents/rules/product.md) | 页面、四态、术语、交互基线 |
| [backend.md](../../.agents/rules/backend.md) | 分层、错误码、参数、事务与配置 |
| [security.md](../../.agents/rules/security.md) | 默认拒绝、逐资源授权、秘密与出站边界 |
| [testing.md](../../.agents/rules/testing.md) | 正常/失败/边界、缺陷信号和一键门禁 |
| [operations.md](../../.agents/rules/operations.md) | Compose、health、日志、停机与恢复 |
| [ai.md](../../.agents/rules/ai.md) | 提示词版本、模型边界、引用与12题评测 |
| [data.md](../../.agents/rules/data.md) | 归属、迁移、数据分级、留存及副本删除 |
| [review.md](../../.agents/rules/review.md) | 规划顺序、完成定义、评审门禁与证据 |

> 前五类是基础，后三类结合AI教学平台补充。每份都有“怎么检查、怎样算失败”。AI与数据业务还未实现，评测与恢复不冒充已通过。

## 3:10—5:20 C：抽三条可重放判据（30分）

1. **后端分层**：打开backend.md首条，再看[scripts/check.py](../../scripts/check.py)的依赖检查及[mutation.py](../../scripts/mutation.py)。服务层导入database/sql会使门禁失败；不是只写“保持良好分层”。已有故意违规验证记录见[verification.md](verification.md)。现场无需手改正式代码。
2. **依赖故障**：打开[骨架spec](../../openspec/specs/platform-foundation/spec.md)和[grader.mjs](../../scripts/grader.mjs)的R04/R05：停测试MySQL后health为503、mysql=down、live仍200，恢复后health为200/up。脚本只操作隔离项目。
3. **AI未找到依据**：打开ai.md相应条款：无命中不能伪造材料/切片引用，并有6正常+6负向评测要求。这是后续业务的可判定标准，当前不能用骨架15项通过来替代AI评测。

> 规范的关键是可以改坏、可以被检查发现；尚未实施的能力也写明未来怎样判，而不提前填“通过”。

## 5:20—6:40 D：规约、规则和证据如何串起来（20分）

打开[本次proposal](../../openspec/changes/archive/2026-10-28-maintain-eight-engineering-standards/proposal.md)，接着切design→spec→tasks，解释各自只管动因/方案/验收/执行。

> 一次变更一套四文件；规则是跨变更约束，spec是能力验收标准，Skill是按场景执行流程，报告是实际证据。完成后归档、长期spec同步。未做项不勾选，不复制多套正文。

展示骨架归档、Git标签`v0.0.1-skeleton`及`v0.0.2-grading`，再打开GitHub PR和Actions。不要将“本地已改”当成“GitHub已合并”。

## 6:40—8:00 验收结果与收尾

展示最新一次真实报告；备用[历史HTML](teacher-report.html)/[历史JSON](teacher-report.json)是2026-10-28 Linux骨架15项通过，必须指出报告中的source_commit和时间。

> 老师克隆、按部署文档配置后，只需Docker Compose及sh运行sh scripts/grade.sh。报告有逐项判据、通过/失败/阻塞和截图，缺证据、网络失败或清理失败均不会判成功。自动检查只能覆盖已编码的判据，不能替代人类评审或证明未来教学业务完成。

结束语：本次展示证明规范有入口、有八类内容、有可判定检查、有规约与证据关联。后续按同一流程增量开发业务。

## 老师追问的简答

- **八份是不是都自动读？** 根AGENTS先读，再按适用条件读取相关规则；不是自动加载整个目录。
- **AI测试通过了吗？** 未实现教学AI，不宣称通过；已给出实施后必须满足的评测判据。
- **为什么没多人评审？** 单人项目；自审、新AI会话与独立人类peer review分别记录，不混算。
- **能证明规则有效吗？** 指静态检查、故意违规失败记录和真实依赖故障报告；若现场复验，使用已有隔离脚本，失败就解释失败，不演假结果。
- **为什么不现场完整构建？** 可以运行，但首轮依赖下载耗时且依赖网络；提前生成的原始报告必须可追溯到展示提交。历史结果只证明历史提交。
