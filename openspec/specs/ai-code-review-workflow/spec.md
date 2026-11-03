# ai-code-review-workflow Specification

## Purpose
将项目文档、阶段规划和AI代码审核流程统一组织，明确已交付成果、建设原因、后续教学能力和风险重点；通过固定评审结构，使每次变更的规约、审核、验证与合并证据可追溯，同时如实区分自动检查、人工审阅和独立评审者交叉评审。

## Requirements

### Requirement: R1 Outcome-oriented documentation

项目说明 SHALL 准确列出已交付成果、建设目的与后续审核计划；未交付业务不得表述为已实现。

#### Scenario: S01 Present actual progress
- **GIVEN** 使用者打开仓库README
- **WHEN** 查看项目说明与路线图
- **THEN** 能说明已有基线及其用途、下一阶段审核闭环、教学业务仍为计划，历史报告可定位日期与提交。

### Requirement: R2 Risk-linked roadmap

项目规划 SHALL 将教学阶段与权限、事务、幂等、并发、一致性、失败及查询风险对应，并为每阶段写明完成证据，不使用虚假交付日期。

#### Scenario: S02 Plan next increments
- **GIVEN** 查看路线图
- **WHEN** 选择下一阶段
- **THEN** 能找到预期能力、为什么优先、人工审核重点与测试/评审证据，未实现阶段标计划。

### Requirement: R3 Recorded AI-code review

评审规则与模板 SHALL 覆盖七类审核对象、七个核对维度、逐文件可用/待修改/风险无效三分结论，高风险部分记录“查了什么/依据/结论”；交叉评审状态如实记录。

#### Scenario: S03 Review one change
- **GIVEN** 一次AI生成代码变更
- **WHEN** 进行审核
- **THEN** 七类及七维逐项有结论或有理由的不适用，逐文件归类，事务/并发/一致性不可无理由跳过，待修改和风险无效阻断合入。

#### Scenario: S04 No independent peer
- **GIVEN** 单人项目尚无另一位评审者复核
- **WHEN** 归档当前文档自审或规划未来代码评审
- **THEN** 明确独立交叉评审未完成、流程适配待维护者确认；自审/AI/CI均不标为独立peer review通过。

### Requirement: R4 Navigable repository evidence

README SHALL 提供规划、部署、工程规范与验证入口；规则、评审模板和实际记录 SHALL 位于仓库可访问路径。推送GitHub前须向用户说明更新内容及原因。

#### Scenario: S05 Navigate and publish
- **GIVEN** 文档改动已准备
- **WHEN** 检查并同步GitHub
- **THEN** 本地链接有效，先说明改动及目的，PR保留检查记录，无姓名/本机填写教程重新进入README，未通过门禁不称已合并。
