# 工程规范与开发流程（engineering-standards）

## Purpose

用八份可按需读取的规范指导开发。先规划、与用户敲定，再实施；只记录当前变更涉及的完善项，避免重复正文和无关上下文。

## Requirements

### Requirement: E1 Eight discoverable standards

**八类规范有入口**

仓库 SHALL 提供八份独立规范：产品设计、后端、安全、测试、部署运维、AI应用、数据、评审与完成定义。AGENTS.md MUST 链接各文件并说明何时读取。同一条款只保留一份权威正文。

#### Scenario: S01 Discover relevant standards
- **GIVEN** 前置条件：用户确认八份规范整理任务，且任务已实施。
- **WHEN** 操作：Harness从AGENTS.md查找本次涉及的规范。
- **THEN** 结果：八文件存在，能找到适用规则与检查方式；迁出条款不在AGENTS重复保存。

### Requirement: E2 Plan and maintain before implementation

**所有变更先规划**

新增、修复、重构、配置和规范调整 SHALL 都遵守以下顺序：
1. 写齐当前change的`proposal.md`、`design.md`、`specs/<能力>/spec.md`、`tasks.md`。
2. 给出范围、技术栈/设计、规范待完善项和验收方式的概览。
3. 用户敲定后再实施。禁止先改实现再补文档；用户敲定前 MUST 不实施计划中的代码或规则变更。

已确认范围内沿用同一change，不重复建文档或索取确认；新增范围先更新规划并敲定。未完成任务保持`[ ]`。规范缺口写在当前change，不另建平行待办。

#### Scenario: S02 Planning awaits confirmation
- **GIVEN** 前置条件：用户要求新增或调整规范。
- **WHEN** 操作：首次交付开发概览。
- **THEN** 结果：四份规划文件可查，tasks全部未勾选，计划中的规则拆分和业务开发尚未实施。

#### Scenario: S03 Unverified clauses remain pending
- **GIVEN** 前置条件：某条规范尚无实现或验证证据。
- **WHEN** 操作：更新设计和交付记录。
- **THEN** 结果：写明待完善项及检查方式；不把文件存在或GD的验证算作Plus已验证。

### Requirement: E3 Observable supplemental standards

**补充规范也要能判断结果**

AI、数据、评审规范 SHALL 分别说明适用条件、操作约束、检查方式和失败判据。未实现业务的检查 MUST 标待验证，不计入当前骨架通过结果。

#### Scenario: S04 Review supplemental standards
- **GIVEN** 前置条件：三份补充规范已整理。
- **WHEN** 操作：核查AI评测、数据删除副本和完成定义。
- **THEN** 结果：能找到12条AI评测及拒答判据、数据副本/恢复目标、Done和阻断条件。
- **AND** 未执行的业务验证仍明确标为待验证。

### Requirement: E4 Repository documentation traceability

**项目资料可追溯**

项目 SHALL 为README、规划、工程规范、部署和验证记录提供清晰导航。不能虚构Harness加载机制或验收结果，非产品材料不作为公开入口。

#### Scenario: S05 Navigate project records
- **GIVEN** 前置条件：维护者打开README。
- **WHEN** 操作：沿链接查看规范与验证记录。
- **THEN** 结果：能找到八份规范、AGENTS入口、Skill元数据、Scenario和验证摘要。
- **AND** 结果有时间/范围，未来业务不标通过。


### Requirement: E5 Complete change lifecycle navigation
AGENTS.md SHALL 提供简短全流程导航：读取规范、四文档规划并确认、严格校验、按任务实施、审核与测试、修复复验、独立评审或如实记录未完成、同步长期规约并归档、经PR门禁合入。详细要求 MUST 引用现有评审规则和验证技能，不重复复制规范正文。

#### Scenario: S06 Follow lifecycle from entry
- **GIVEN** Harness首次读取AGENTS.md。
- **WHEN** 查找一次变更从规划到交付的顺序。
- **THEN** 可定位上述各阶段，并可沿链接找到评审规则与验证技能；确认后先实施，验证后再同步归档。

#### Scenario: S07 Keep incomplete review explicit
- **GIVEN** 变更尚无独立人类评审。
- **WHEN** 按导航记录审核和交付状态。
- **THEN** 独立评审保留未完成或有来源的例外；AI自审和CI不计为独立评审，不改变既有合并权限。
