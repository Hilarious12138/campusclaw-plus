# engineering-standards Specification

## Purpose
以八份可被Harness按需读取的规范指导后续开发，并坚持先规划、共同敲定再实施；仅记录与当前变更相关的完善项，控制上下文长度。

## Requirements

### Requirement: E1 Eight discoverable standards

仓库 SHALL 提供产品设计、后端开发、安全开发、测试、部署运维、AI应用、数据、评审与完成定义八份独立规范；AGENTS.md MUST 链接并说明读取条件，同一条款只保留一份权威正文。

#### Scenario: S01 Discover relevant standards
- **GIVEN** 八份规范整理任务已获用户确认并实施
- **WHEN** Harness从AGENTS.md查找本次涉及的规范
- **THEN** 八个文件均存在，能找到适用规则及检查方式，迁移条款不在AGENTS重复保存。

### Requirement: E2 Plan and maintain before implementation

每次开发（新增、修复、重构、配置和规范调整）SHALL 先在当前change写齐`proposal.md`、`design.md`、`specs/<能力>/spec.md`、`tasks.md`，并给出范围、技术栈/设计、规范待完善项及验收方式的概览。禁止先实施再补文档。已敲定范围内的连续工作更新同一change，无需重复建文档或重复确认；新增范围须先更新规划并敲定。用户敲定前 MUST 不实施计划中的代码或规则变更；未完成任务保持[ ]。涉及的规范缺口须在当前change中标注，不增加平行待办文档。

#### Scenario: S02 Planning awaits confirmation
- **GIVEN** 用户要求新增或调整规范
- **WHEN** 首次提交该次开发概览
- **THEN** 当前change下四份规划文件可查，tasks全部未勾选，未进行计划中的规则拆分和业务开发。

#### Scenario: S03 Unverified clauses remain pending
- **GIVEN** 某条规范尚无实现或验证证据
- **WHEN** 更新当前change的设计及交付记录
- **THEN** 明确写出待完善项和检查方式，不将“文件存在”或GD验证结果算作Plus已验证。

### Requirement: E3 Observable supplemental standards

AI、数据和评审规范 SHALL 分别给出适用条件、操作约束、检查方式与失败判据；未实现业务的检查 MUST 标为待验证，不计入当前骨架通过结果。

#### Scenario: S04 Review supplemental standards
- **GIVEN** 三份补充规范已整理
- **WHEN** 核查AI评测、数据删除副本与完成定义
- **THEN** 能找到12条AI评测要求及拒答判据、数据副本清单与恢复目标、Done与阻断条件；未执行的业务验证明确标待验证。

### Requirement: E4 Repository presentation traceability

项目 SHALL 提供一份按A/B/C/D评分组织的仓库展示路线，链接规则、规约和真实证据；不得虚构Harness加载机制或当前验收成绩。

#### Scenario: S05 Walk through repository
- **GIVEN** 教师打开仓库讲解稿
- **WHEN** 沿展示顺序查看链接
- **THEN** 能定位八份规范、AGENTS入口、Skill元数据、Scenario及历史报告；报告时间/范围限制清楚，未来业务不标通过。
