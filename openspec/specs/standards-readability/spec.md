# 指导文档的表达方式（standards-readability）

## Purpose

在原有八份工程规范的基础上，统一AGENTS和长期能力规约的表达方式。通过分节、短句和中文提示提高阅读效率，同时保留路径、职责、约束强度、技术标识和验收依据，避免平行正文与无边界扩写。

## Requirements

### Requirement: R1 Readable organization

工程规范、AGENTS和长期能力规约 SHALL 使用清楚的适用范围、主题和短句组织内容。
- 工程规范区分要求与检查方式。
- AGENTS区分阅读、实施、边界和验证步骤。
- 能力规约保留OpenSpec英文骨架，用中文说明要求及验收步骤。
- 多动作长句须拆分，必要术语首次出现时简短解释。

#### Scenario: S01 Locate an actionable rule
- **GIVEN** 前置条件：改写后的规范、AGENTS或长期能力规约。
- **WHEN** 操作：读取某个主题或场景。
- **THEN** 结果：能找到适用条件、执行动作和检查结果，不依赖多动作长句理解要求。

### Requirement: R2 Preserve existing constraints

改写 MUST 保留既有约束的条件、强度、技术标识、数值边界与失败判据，不得借简化措辞删约束或新增未确认要求。长期spec还须保留Requirement/Scenario的原标题、编号与英文骨架。

#### Scenario: S02 Compare old and new rules
- **GIVEN** 前置条件：改写前后的正文。
- **WHEN** 操作：按原条款逐项核对。
- **THEN** 结果：每项约束都有对应表述；命令、路径、字段、状态码和限额一致。
- **AND** 禁止项未被弱化为建议，长期spec的原Requirement/Scenario标题和编号仍可定位。

#### Scenario: S03 Keep future scope explicit
- **GIVEN** 前置条件：某项要求涉及尚未实现的业务。
- **WHEN** 操作：改写其适用范围或状态说明。
- **THEN** 结果：仍说明未来实施时须满足的要求；未实现/未验证事项不写成当前完成，不适用项保留原因。

### Requirement: R3 Keep one authoritative version

既有规范、AGENTS和长期spec的路径与职责 SHALL 不变；八份工程规范继续由AGENTS按需索引。不得另建平行正文，也不得为解释术语无限扩写。

#### Scenario: S04 Read repository guidance
- **GIVEN** 前置条件：从仓库入口查找某类指导文件。
- **WHEN** 操作：按原索引或路径打开文件。
- **THEN** 结果：文件存在，仍是相应职责的唯一权威正文；内部引用有效，表达精简，不增加运行依赖。
