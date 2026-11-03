## Why

公开仓库应作为独立软件项目提供用途、状态、规划、部署、工程规范和验证证据。当前存在非产品材料及非标准命名，影响对外一致性，需要迁移并统一表达，同时保留实际验证结果与未来能力边界。

## What Changes

- 外部用途材料先复制到仓库之外，再从当前发布树移出；产品文档与工程审核记录按职责整理。
- 统一README、路线图、八类规范、OpenSpec和当前页面的项目表述。
- 验证记录归入docs/verification，报告标题、范围与输出目录改为通用验收命名，保留grade入口及失败语义。
- 更新导航和引用，保留实际结果/提交对应，不修改旧原始报告冒充新结果。

## Capabilities

### New Capabilities
- `project-documentation`: 面向产品与维护者的公开文档边界及验收命名。

### Modified Capabilities
以项目化表述整理engineering-standards、ai-code-review-workflow和platform-foundation文档；运行能力不变。

## Non-goals

不开发教学业务、不删除个人副本、不更换技术栈、不改认证/数据逻辑、不重写Git历史或移除旧标签，不将旧证据编辑成新的通过结果。

## Alternatives

放弃只改README：页面、规范和报告仍会存在不一致。也不为隐藏旧内容而强推或重建仓库；当前发布树整理与历史清理分开。

## Impact

涉及文档迁移、说明文本、测试请求标识、验收产物路径及CI收集路径。对外入口仍为Docker Compose与sh scripts/grade.sh，main的check/grade门禁保持。独立评审未完成仍如实记录，由项目维护者管理例外。
