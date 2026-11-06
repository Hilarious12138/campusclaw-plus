## MODIFIED Requirements

### Requirement: P2 Neutral acceptance outputs

**验收命名与判定一致**

自动验收 SHALL 保留原有成功、失败、缺阶段判断及原有15个必测标识，并按新增能力追加必测。报告范围、标题和目录使用通用项目命名，CI同步收集新路径。

#### Scenario: S03 Run acceptance
- **GIVEN** 前置条件：Docker环境和配置有效。
- **WHEN** 操作：执行`sh scripts/grade.sh`。
- **THEN** 结果：产物写入`reports/acceptance-*`，HTML/JSON标注实际受验能力与夹具/真实依赖边界，CI能收集。
- **AND** 只有必测齐全且通过时，才总体成功。

#### Scenario: S04 Fail acceptance
- **GIVEN** 前置条件：环境缺失或必测失败。
- **WHEN** 操作：生成报告。
- **THEN** 结果：仍按原契约非零退出；未执行项记失败/阻塞，不因命名变化改判。
