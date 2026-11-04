## ADDED Requirements

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
