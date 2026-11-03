# AI代码审核流程（ai-code-review-workflow）

## Purpose

将项目说明、路线图、审核和验证记录关联起来。已交付成果、后续教学能力及风险重点要清楚；自动检查、人工审阅和独立交叉评审分别记录，不相互冒充。

## Requirements

### Requirement: R1 Outcome-oriented documentation

**如实说明项目进展**

项目说明 SHALL 列出已交付成果、建设目的及后续审核计划；未交付业务不能写成已实现。

#### Scenario: S01 Present actual progress
- **GIVEN** 前置条件：使用者打开README。
- **WHEN** 操作：查看项目说明及路线图。
- **THEN** 结果：能了解基线、用途和后续审核流程；教学业务仍标计划，历史报告可定位日期与提交。

### Requirement: R2 Risk-linked roadmap

**规划同时说明风险**

项目规划 SHALL 将各教学阶段与权限、事务、幂等、并发、一致性、失败及查询风险对应。每阶段写明完成证据，不编造交付日期。

#### Scenario: S02 Plan next increments
- **GIVEN** 前置条件：正在查看路线图。
- **WHEN** 操作：选择下一阶段。
- **THEN** 结果：能找到能力、优先原因、人工审核重点及测试/评审证据，未实现阶段标计划。

### Requirement: R3 Recorded AI-code review

**审核要有依据和结论**

评审规则与模板 SHALL 包含七类对象、七维核查，以及逐文件“可用/待修改/风险无效”结论。高风险部分记录查了什么、依据和结论。交叉评审状态必须真实。

#### Scenario: S03 Review one change
- **GIVEN** 前置条件：有一次AI生成代码的变更。
- **WHEN** 操作：进行审核。
- **THEN** 结果：七类/七维逐项有结论或有理由的不适用，文件逐个归类。
- **AND** 不无理由跳过事务、并发、一致性；待修改和风险无效阻断合入。

#### Scenario: S04 No independent peer
- **GIVEN** 前置条件：单人项目没有另一位评审者复核。
- **WHEN** 操作：归档文档自审或规划后续代码评审。
- **THEN** 结果：独立交叉评审标未完成，流程适配待维护者确认；自审、AI、CI不标为独立peer review通过。

### Requirement: R4 Navigable repository evidence

**资料和发布过程可查**

README SHALL 提供规划、部署、工程规范和验证入口。规则、评审模板和实际记录 SHALL 位于仓库可访问路径。推送GitHub前先向用户说明更新内容和原因。

#### Scenario: S05 Navigate and publish
- **GIVEN** 前置条件：文档改动已准备好。
- **WHEN** 操作：检查并同步GitHub。
- **THEN** 结果：本地链接有效，更新内容和目的已提前说明，PR保留检查记录。
- **AND** README没有重新加入姓名或本机填写教程；门禁通过前不称已合并。
