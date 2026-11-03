# 公开项目文档（project-documentation）

## Purpose

公开仓库以独立产品提供用途、状态、规划、部署、规范和验证信息。非产品资料独立保存；报告使用通用项目命名，同时保留实际结果和历史记录的边界。

## Requirements

### Requirement: P1 Product-facing public tree

**公开树只放项目资料**

当前发布树 SHALL 只保留项目文档和代码。移出外部材料前先留本机副本；README、页面、规范及导航不能残留过时的非产品标识。

#### Scenario: S01 Inspect public documentation
- **GIVEN** 前置条件：打开当前主干。
- **WHEN** 操作：查看README、页面、工程规范和OpenSpec。
- **THEN** 结果：内容围绕项目用途、规划和维护，内部链接指向实际存在的公开文件。

#### Scenario: S02 Preserve external materials
- **GIVEN** 前置条件：部分原文件将退出发布范围。
- **WHEN** 操作：迁移文件。
- **THEN** 结果：仓库外副本与原文件摘要一致，当前树不再提供相关入口。
- **AND** 不复制秘密，不重写Git历史。

### Requirement: P2 Neutral acceptance outputs

**验收命名与判定一致**

自动验收 SHALL 保留原有成功、失败、缺阶段判断及15个必测标识。报告范围、标题和目录使用通用项目命名，CI同步收集新路径。

#### Scenario: S03 Run acceptance
- **GIVEN** 前置条件：Docker环境和配置有效。
- **WHEN** 操作：执行`sh scripts/grade.sh`。
- **THEN** 结果：产物写入`reports/acceptance-*`，HTML/JSON标注项目基础能力范围，CI能收集。
- **AND** 只有必测齐全且通过时，才总体成功。

#### Scenario: S04 Fail acceptance
- **GIVEN** 前置条件：环境缺失或必测失败。
- **WHEN** 操作：生成报告。
- **THEN** 结果：仍按原契约非零退出；未执行项记失败/阻塞，不因命名变化改判。

### Requirement: P3 Honest project records

**不改写验证事实**

项目 SHALL 保留适用的工程规划、规范和验证事实。不能把原始结果改写成新验收；未实现业务、未完成评审仍须明确。整理当前树不表示历史已擦除。

#### Scenario: S05 Read verification history
- **GIVEN** 前置条件：历史证据和当前验收并存。
- **WHEN** 操作：查阅项目记录。
- **THEN** 结果：能辨认提交与范围，原始副本未被篡改，业务状态和评审边界未被夸大。
