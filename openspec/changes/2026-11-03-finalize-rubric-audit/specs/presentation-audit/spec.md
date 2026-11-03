## Purpose

把原评分表与下一阶段审核AI代码的课件要求逐项映射到仓库文件和实际证据，给出当前满足程度、缺口及可直接讲述的最终稿，避免将规划、静态检查或AI辅助核对冒充已完成的业务实践与独立人类评审。

## ADDED Requirements

### Requirement: V1 Itemized course audit

清单 SHALL 覆盖A-D全部问题、八类规范、七类对象、七维核查、三分法及三层审核机制，每项有状态和证据或补齐动作。

#### Scenario: S01 Trace rubric requirements
- **GIVEN** 教师查看评分核对表
- **WHEN** 按A/B/C/D逐项核对
- **THEN** 每个评分问题都有仓库路径/检查依据，未做的加载复验、人工评审和业务测试未被计通过。

#### Scenario: S02 Trace next-stage practice
- **GIVEN** 查看课件对照部分
- **WHEN** 核查七类、七维、三分法和审核流程
- **THEN** 每项能定位规则/模板/实际记录，规范就绪与实际实践分别标识，缺同学peer不会写流程全部完成。

### Requirement: V2 Final outcome-focused script

最终讲稿 SHALL 围绕做了什么、为什么做、后续如何审核，包含页面提示、实际已验证案例和未来业务案例的明确边界，不以实现细节为主。

#### Scenario: S03 Present final script
- **GIVEN** 用户直接开仓库讲解
- **WHEN** 按约6分钟讲稿顺序展开
- **THEN** 能说明项目价值、四步基线、八类规范用途、后续阶段审核重点和交付条件，证据入口可打开。

#### Scenario: S04 Preserve evidence limits
- **GIVEN** 讲稿引用自动报告或重复提交案例
- **WHEN** 说明当前成果
- **THEN** 报告明确受验提交/范围，未来业务标计划，不将AI辅助或CI称为本人/同学已审阅。

### Requirement: V3 Verifiable documentation handoff

本次 SHALL 执行内容/链接与既有门禁检查，变异检查结束恢复实现；GitHub推送前说明更新内容与原因，发布状态不得提前写成成功。

#### Scenario: S05 Verify and publish
- **GIVEN** 清单和讲稿完成
- **WHEN** 执行验证与交付
- **THEN** 链接有效，strict/适用门禁有实际结果，变异后实现无残留改动，README无姓名/填写表，PR未通过前不称已合并。
