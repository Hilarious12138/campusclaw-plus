# CampusClaw Plus

面向班级教学的智能教学平台，连接讲义检索、辅助学习、作业批改与错题订正。

## 项目状态

采用React、Go分层后端、MySQL、Nginx与Docker Compose，提供教师/学生身份、本班材料、可溯源检索、可靠异步任务及操作审计。

多应用、提示词、Skill、MCP与流式会话已有配置和适配代码。真实AI API接入暂缓；AI_ENABLED=false时无需填写AI地址或密钥，可使用材料与关键字检索。自动验收明确使用模型/嵌入测试替身，不能据此宣称真实模型效果或提示词安全已验证。

作业评分、学情、生产备份恢复与独立人类交叉评审仍未完成，见[项目规划](docs/roadmap.md)和[验收映射](docs/verification/secure-platform-scenarios.md)。

## 部署与验收

依赖Docker Engine / Docker Desktop与Compose v2。配置说明见[部署文档](docs/deployment.md)。

```sh
docker compose up --build -d --wait
```

自动验收只需Docker Compose与系统sh。使用已有配置：

```sh
sh scripts/grade.sh
```

无须填写AI配置的隔离测试可直接运行`sh scripts/integration.sh`，自动生成临时测试配置并在结束时删除。

验收创建独立环境，执行检查后清理测试容器与卷，保留`reports/acceptance-*/report.html`和`report.json`。报告包含逐项判据、结果与截图；必测失败、未执行或清理失败均返回非零退出码。覆盖工程、业务安全、真实队列故障与浏览器检查；报告明确标注替身与未验证范围。

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

安全检查入口：`make check-security`。秘密扫描或漏洞失败、公告不可达最终重试失败均阻断；原始结果与前后对照保存在`reports/`。
