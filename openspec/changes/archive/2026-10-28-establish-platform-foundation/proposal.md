## Why

CampusClaw Plus面向教师和学生，解决讲义、辅助学习与作业闭环分散、结果难追溯的问题。本课先建立可被Harness读取、可验证且独立于GD的工程基线，后续按同一规约增量交付，不每轮重建约定。

## What Changes

- 编号交付：F1独立仓库与可复现配置；F2分层骨架与默认拒绝；F3免认证存活/就绪探测；F4单行JSON日志与优雅停机；F5四状态首页；F6规则/技能/门禁与真实验证记录。
- 新仓首次提交只有README/.gitignore，主干先设禁止删除/强推；后续PR与CI合并。
- 五类规范置于.agents/rules，由AGENTS指定读取；验证流程为项目Skill，规则不只停在对话里。
- 实现前先strict，实施只跑骨架任务，不引入GD业务、模型、Qdrant或MCP。
- 单人三角色由郭地承担；同机干净clone与空白Harness检查作为可执行验证，不冒充他人异机peer review。课程原项保留差异说明。

## Capabilities

### New Capabilities
- `platform-foundation`: 配置校验、分层/默认拒绝、健康状态、日志停机、四状态首页及验收门禁。

### Modified Capabilities
无。Plus独立仓库，GD规约不代表Plus已实现。

## Non-goals

本课不做登录、角色业务、班级表、讲义上传、模型/检索、模板、Skill执行、MCP、作业、错题、学情或公网部署。这些属于学期目标，README路线图保留。

放弃“全量复制GD后补规范”：这会把旧依赖、凭证风险与业务故障混入骨架，难以独立验收；改为后续分批复用。放弃“先做静态假健康”：必须真实查询MySQL。

## Impact

新增backend/frontend/deploy/scripts与Compose；MySQL仅创建空schema，不建业务表。增加Go MySQL驱动与React/Vite构建依赖，锁定版本；Nginx唯一宿主入口，建议3001，独立Plus卷，不修改GD。
