# CampusClaw Plus

面向班级教学的智能教学平台，连接讲义检索、辅助学习、作业批改与错题订正。

## 项目状态

已具备React前端、Go分层后端、MySQL、Nginx与Docker Compose工程基线，支持健康检查、结构化日志和隔离自动验收。

账号认证、班级知识库、多AI应用、提示词模板、Skill、MCP、作业与学情功能按[项目规划](docs/roadmap.md)分阶段交付，目前尚未实现。

## 部署与验收

依赖Docker Engine / Docker Desktop与Compose v2。配置说明见[部署文档](docs/deployment.md)。

```sh
docker compose up --build -d --wait
```

自动验收只需Docker Compose与系统sh：

```sh
sh scripts/grade.sh
```

验收创建独立环境，执行检查后清理测试容器与卷，保留`reports/acceptance-*/report.html`和`report.json`。报告包含逐项判据、结果与截图；必测失败、未执行或清理失败均返回非零退出码。当前覆盖平台基础能力，后续业务随变更扩展。

## 项目结构

- `frontend/`：React 18、TypeScript、Vite。
- `backend/`：Go、标准库net/http、MySQL适配。
- `deploy/`：镜像构建与Nginx配置。
- `openspec/`：能力规约与变更归档。
- `scripts/`：检查、测试和自动验收。
- `docs/`：规划、部署、审核与验证记录。

## 工程流程

开发前阅读[AGENTS.md](AGENTS.md)，先完成变更规划，再实施与验证。八类工程规范按职责读取；每次交付包含能力、审核记录与验证证据。

- [项目规划](docs/roadmap.md)
- [评审模板](docs/reviews/template.md)
- [验证记录](docs/verification/platform-foundation.md)
