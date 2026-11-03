# CampusClaw Plus · Harness 常驻规则

## 每次变更
- 此仓库是独立Plus项目，GD仅参考；不得改GD或复制其.env、Git元数据、数据卷及完成状态。
- 每次开发前（新增、修复、重构、配置和规范调整均适用），必须先在`openspec/changes/YYYY-MM-DD-kebab-case/`写齐`proposal.md`、`design.md`、`specs/<能力>/spec.md`、`tasks.md`。先读README与相关规则；先交付四份文档和简短概览（范围、技术栈/设计、规范待完善项、验收方式），与用户修改敲定后才实施。不得先改实现再补文档；计划阶段tasks全为`[ ]`，实际验证后才勾选。已获确认范围内的连续工作更新同一change，不重复建文档或重复索取确认；新增范围先更新规划并敲定。已归档变更不代表后续开发自动获准。
- 常驻安全/测试约束见下；代码涉及前端读`.agents/rules/product.md`，后端读`backend.md`，任意代码读`security.md`与`testing.md`，构建/部署读`operations.md`。规范文件中的“检查”是门禁要求，不能只读标题。
- 验证时使用`.agents/skills/verify-change/SKILL.md`；不读无关旧changes与长报告，减少上下文。新能力每次建YYYY-MM-DD-kebab-case变更。

## 不可绕过的边界
- 默认拒绝，逐资源校验角色、班级和归属。公开接口须明确登记；当前只有/、/health、/live，骨架没有登录和业务资源。
- 密码/令牌/模型及MCP密钥仅后端环境变量；.env不得入Git/镜像/前端/日志；.env.example仅键名。Secrets扫描覆盖全部Git历史，发现后先轮换，不靠删文件假装修复。
- Go入口→HTTP→服务→领域端口；数据库适配器实现领域端口，只在cmd装配。业务层不import数据库/HTTP或执行SQL；测试用替身不得冒充真实运行证据。
- 未通过验证不勾tasks、不归档；发现阻塞如实记录。合并前`make check`必过，开发机真实运行另跑`make acceptance`，老师交付必须通过仅Docker的`sh scripts/grade.sh`。后续新增能力同步扩展容器用例/Scenario映射；不适用不计通过。故意移除授权或引入反向依赖时，门禁必须失败。

## 规范读取索引

根AGENTS是常驻入口；下列普通Markdown文件通过本索引按任务读取，不宣称目录中文件自动全部加载。操作流程见`.agents/skills/verify-change/SKILL.md`。

| 规范 | 路径 | 何时读取 |
| --- | --- | --- |
| 产品设计 | `.agents/rules/product.md` | 页面、交互、文案、产品场景 |
| 后端开发 | `.agents/rules/backend.md` | 后端接口、分层、查询 |
| 安全开发 | `.agents/rules/security.md` | 任意实现或配置变更 |
| 测试 | `.agents/rules/testing.md` | 用例、修复、验收 |
| 部署运维 | `.agents/rules/operations.md` | 构建、部署、日志、回滚 |
| AI应用 | `.agents/rules/ai.md` | 模型、提示词、检索、Skill、MCP |
| 数据 | `.agents/rules/data.md` | 数据模型、迁移、材料索引、留存/删除 |
| 评审与完成定义 | `.agents/rules/review.md` | 实施、验证、合并与交付 |

每次合入的审核记录见`docs/reviews/`，按review规则进行七类/七维核查；不得把AI辅助或CI称为独立同学交叉评审。

八份基础规范齐全不等于所有业务已实现/已验证。当前change需标涉及规范、待完善项及检查方式；只加载相关正文。单人三职责不等于独立同伴评审；未执行项如实保留。
