## Purpose

在保持八份工程规范原有职责、约束强度、技术标识和验收判据的前提下，用分节、短句和明确的要求/检查结构提高阅读效率，减少重复背景与术语负担，使维护者和AI都能按同一份正文理解并执行。

## ADDED Requirements

### Requirement: R1 Readable organization

八份规范 SHALL 按主题分节，以清楚的适用范围、要求和检查方式组织；密集多动作长句须拆分，必要术语首次出现时简短解释。

#### Scenario: S01 Locate an actionable rule
- **GIVEN** 改写后的任一规范
- **WHEN** 读取某个主题
- **THEN** 能分别找到何时适用、必须执行的动作和检查方式，不依赖跨越多个表格长单元格理解要求。

### Requirement: R2 Preserve existing constraints

改写 MUST 保留既有约束的条件、强度、技术标识、数值边界与失败判据，不得以简化语言为由删减约束或新增未确认的要求。

#### Scenario: S02 Compare old and new rules
- **GIVEN** 改写前后两版正文
- **WHEN** 按原条款逐项核对
- **THEN** 每项约束都有对应的新表述，命令/路径/字段/状态码/限额一致，禁止项未被弱化为建议。

#### Scenario: S03 Keep future scope explicit
- **GIVEN** 某条要求涉及尚未实现的业务
- **WHEN** 改写适用范围或状态说明
- **THEN** 仍明确实施后需要满足的要求，不将未实现或未验证事项写成当前完成，不适用项保留原因。

### Requirement: R3 Keep one authoritative version

八份文件的路径和职责 SHALL 不变，继续由AGENTS按需索引；不得另建一套平行规范，亦不得为解释术语无限扩写。

#### Scenario: S04 Read repository guidance
- **GIVEN** 从AGENTS读取某类规范
- **WHEN** 按原索引打开文件
- **THEN** 文件仍存在且是唯一权威正文，内部引用有效，改写保持精简且不增加新的运行依赖。
