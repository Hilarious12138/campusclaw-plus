## Context

现有product.md、backend.md、security.md、testing.md、operations.md五份规则；后三类正文暂位于AGENTS.md。上一条直接修改AGENTS.md的操作已撤回，当前仅规划，尚不实施拆分。

## Decisions

1. 保留前五份；新增AI规则文件为`.agents/rules/ai.md`、数据为`data.md`、评审与完成定义为`review.md`。分别迁移现有条款，删除原重复正文。
2. AGENTS.md保留八份索引、适用条件和“先四文档及概览、敲定后实施”的入口约束。AI/工具变更读ai，模型/数据迁移读data，实施验收读review；其他规范按实际影响读取。
3. 当前change的design列出涉及规范、缺口、检查方式；tasks列具体完善任务。交付记录只记本次完成项和剩余缺口，不新增独立待办体系。
4. 八份文件存在仅表示基础规范齐全。条款的实现与验证仍按业务逐次完善；不自动继承GD完成状态。

## Stack / Interfaces / Data

仅文档组织变更，无运行接口、配置、数据模型、数据库迁移或技术栈变动。继续React18/TypeScript/Vite、Go net/http、MySQL8、Nginx/Compose。

## Verification

E1/S01：检查八份文件、AGENTS索引与迁移条款，避免重复；E2/S02：核对四份文件先存在且tasks未勾选；E2/S03：核对缺口有位置及检查方式，未验证项不标完成。规划执行strict；实际实施后跑既有门禁，结果如实记录。

## 本次授权与状态

用户已明确授权将“每次开发先四文档、概览敲定后实施”写成常驻规约。本次仅更新AGENTS.md、openspec/config.yaml及本change的对应描述；八份规范的拆分仍待敲定，不随本次请求实施。验证采用内容核对、git diff --check与OpenSpec strict，不运行无关业务测试。


## 2026-11-02 评分核对与实施授权

用户已授权完善AI、数据、评审三类；本次先更新四文档，再实施规则整理及讲解稿，不修改业务或验收程序。

| 评分项 | 现状缺口 | 本次措施 / 判据 |
| --- | --- | --- |
| A 落位与可加载30 | 五份规则由AGENTS引用，后三类未拆；不能声称规则目录自动全部加载 | AGENTS列八份路径/适用条件；从Plus根目录启动Harness；SKILL.md核对name/description，规则正文通过入口按需读取 |
| B 规范覆盖20 | 五独立文件+三类混写 | 拆ai/data/review，迁移已有要求并补可执行检查；八份非空且有检查方式 |
| C 可判定30 | AI评测、数据副本、DoD仍需更具体 | 每类写触发、动作、判据；未来业务用例标待验证，不以当前15项替代；抽查三条能说明如何失败 |
| D 组织与自洽20 | 重复正文、部分文档引用README旧位置 | AGENTS留索引与核心边界，修正部署说明引用；design→Scenario→任务→证据链；一份讲解稿复用证据、不复制规则 |

技术栈不变。Codex使用根AGENTS.md与.agents/skills；SKILL.md须name/description。八份普通Markdown规则通过AGENTS显式索引，不把globs/alwaysApply虚构为Codex加载机制，也不为另一Harness复制一套规则。

依据：https://learn.chatgpt.com/docs/agent-configuration/agents-md 、https://learn.chatgpt.com/docs/build-skills （2026-11-02核对）。若课堂要求换用支持globs/alwaysApply的Harness，另确认适配；本次不宣称跨Harness验证。

## 展示与验证设计

新增docs/lesson-07/showcase.md作为唯一讲解稿：约8分钟，README定位→AGENTS入口→八份规则→Skill元数据→一条规则的Scenario/测试→报告和版本/PR。报告说明运行日期与提交，不声称今日重跑；备选用原始报告，现场不等待首次镜像下载。

检查：八文件非空、每份含检查方式、AGENTS八路径存在、迁移条款不重留、讲解稿链接有效、strict与静态门禁。AI/数据业务尚未实现，评测/恢复/跨班授权均待后续能力验证；不勾未来工作。GitHub同步须走PR及check/grade，当前先完成本地整理。
