## Purpose

让公开仓库作为独立产品提供清晰且一致的用途、进度、规划、部署、工程规范与验证信息；非产品资料与原始个人材料独立保存，验收产物使用通用项目命名，并保留真实结果和历史记录的明确边界。

## ADDED Requirements

### Requirement: P1 Product-facing public tree

当前发布树 SHALL 只保留项目相关文档和代码；不属于项目交付的外部材料须先保留本机副本再移出。README、页面、规范与导航不得残留过时的非产品标识。

#### Scenario: S01 Inspect public documentation
- **GIVEN** 打开当前主干
- **WHEN** 查看README、页面、工程规范与OpenSpec
- **THEN** 内容以项目用途、规划与维护为中心，所有内部链接指向实际存在的公开文件。

#### Scenario: S02 Preserve external materials
- **GIVEN** 部分原文件退出发布范围
- **WHEN** 进行迁移
- **THEN** 仓库外副本与原文件摘要一致，当前树不再提供相关入口；不复制秘密，不重写Git历史。

### Requirement: P2 Neutral acceptance outputs

自动验收 SHALL 保留原有成功/失败/缺阶段判断与15个必测标识，使用项目通用的报告范围、页面标题和产物目录，CI同步收集新路径。

#### Scenario: S03 Run acceptance
- **GIVEN** 有效Docker环境和运行配置
- **WHEN** 执行sh scripts/grade.sh
- **THEN** 产物写入reports/acceptance-*，HTML/JSON以项目基础能力标注范围，CI可收集，只有必测齐全且通过才总体成功。

#### Scenario: S04 Fail acceptance
- **GIVEN** 环境缺失或必测失败
- **WHEN** 生成报告
- **THEN** 仍按原契约退出非零，未执行项目保持失败/阻塞，不因命名整理改变判定。

### Requirement: P3 Honest project records

项目 SHALL 保留适用的工程规划、规范与验证事实；不得重写原始结果伪装成新验收，未实现业务与未完成评审保持明确。当前树整理不应被描述为历史已擦除。

#### Scenario: S05 Read verification history
- **GIVEN** 历史证据与当前验收并存
- **WHEN** 查阅项目记录
- **THEN** 能识别对应提交和范围，原始副本未被篡改，当前业务状态和评审边界不被夸大。
