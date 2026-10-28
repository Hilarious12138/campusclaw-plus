# CampusClaw Plus

面向班级教学的智能教学平台，连接讲义检索、辅助学习、作业批改与错题订正。项目基于 CampusClaw-GD 的经验独立建设。

## 当前进度

已完成工程骨架：React 前端、Go 分层后端、MySQL、Nginx、Docker Compose、健康检查、结构化日志和自动验收。

教学登录、知识库检索、多 AI 应用、提示词模板、Skill、MCP、作业与学情功能将按 OpenSpec 分阶段交付，目前尚未实现。

## 部署与验收

运行环境为 Docker Engine / Docker Desktop 与 Compose v2。部署配置见[部署文档](docs/deployment.md)。

```sh
docker compose up --build -d --wait
```

自动验收仅需 Docker Compose 和系统 sh，无需安装宿主开发运行时：

```sh
sh scripts/grade.sh
```

验收会创建隔离环境并自动清理测试容器与卷，生成 `reports/teacher-*/report.html` 和 `report.json`。报告包含逐项判据、结果和截图；必测失败、未执行或清理失败均返回非零退出码。当前验收覆盖工程骨架，后续业务随变更扩展。

## 项目结构

- `frontend/`：React 18、TypeScript、Vite。
- `backend/`：Go、标准库 net/http、MySQL 数据访问。
- `deploy/`：容器构建与 Nginx 配置。
- `openspec/`：能力规约与变更归档。
- `scripts/`：检查、测试和自动验收程序。
- `docs/`：部署说明与验证记录。

开发前阅读 [AGENTS.md](AGENTS.md)，按规约实施并通过 `make check`。验证证据见[第 7 课记录](docs/lesson-07/verification.md)。
