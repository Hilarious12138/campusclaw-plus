# 部署与配置

依赖Git、Docker Engine/Desktop及Compose v2。复制`.env.example`为本机`.env`并配置后执行：

```sh
docker compose up --build -d --wait
```

只有Nginx发布本机端口。API、MySQL、RabbitMQ、Qdrant及MCP只在Compose网内互通。`docker compose down`保留数据；不要用`down -v`作为普通停止、口令轮换或回滚。

## 必填配置

| 键 | 约束或建议 |
| --- | --- |
| WEB_PORT | 未占用端口，建议3001 |
| APP_ENV / APP_VERSION | local/test/production；准确部署版本 |
| APP_ORIGIN | 浏览器精确来源，例如http://localhost:3001；localhost与127.0.0.1不同 |
| HTTP_ADDR / DB_ADDR | :8080 / db:3306 |
| DB_NAME / DB_USER | 独立数据库与普通账号，建议campusclaw_plus |
| DB_PASSWORD / DB_ROOT_PASSWORD | 分别生成≥24位随机值且两者不同 |
| JWT_SECRET | 独立≥32字节随机值，不复用其他项目 |
| SEED_TEACHER_A_PASSWORD / SEED_STUDENT_A_PASSWORD | 各≥16位且与其他种子不同，仅首次初始化使用 |
| SEED_TEACHER_B_PASSWORD / SEED_STUDENT_B_PASSWORD | 各≥16位且与其他种子不同，仅首次初始化使用 |
| DB_TIMEOUT_MS / SHUTDOWN_TIMEOUT_SECONDS | 建议1500 / 10 |
| BROKER_ADDR / BROKER_USER / BROKER_PASSWORD | broker:5672 / 独立账号 / ≥24位随机口令 |
| QDRANT_URL | http://qdrant:6333 |
| MCP_AUTH_TOKEN / MCP_ALLOWED_BASE_URLS | 独立≥32位令牌 / http://mcp-demo:8090；应用和MCP容器令牌一致 |
| TASK_TIMEOUT_SECONDS / TASK_LEASE_SECONDS | 建议120 / 150，租约至少比任务上限多5秒 |
| TASK_MAX_ATTEMPTS / TASK_RETRY_BASE_MS / TASK_POLL_MS | 建议3 / 5000 / 1000 |
| AI_ENABLED | 当前建议false，无需填写AI接口 |

种子账号为teacher_a、student_a、teacher_b、student_b，各自使用对应密码。只初始化独立Plus数据库，不引入旧项目账号。不要公开种子密码，不以同一个口令覆盖不同角色/班级。

AI_ENABLED=false时，下列键可留空：EMBEDDING_BASE_URL、EMBEDDING_API_KEY、EMBEDDING_MODEL、EMBEDDING_DIM、EMBEDDING_TIMEOUT_SECONDS、CHAT_BASE_URL、CHAT_API_KEY、CHAT_MODEL。材料仍可解析切片并做关键字检索；向量/混合与AI对话明确不可用。未来启用前再配置完整AI参数并做真实供应商验证；现在不填写真实地址或密钥。

## 检查与验收

浏览器访问配置来源的`/`，API文档为`/swagger/`。`/health`免认证，DB正常返回200/OK；broker故障标degraded且任务仍可持久接收，DB故障503；`/live`仅表示进程存活。

```sh
make check-security
make check
sh scripts/grade.sh
```

`make check`是开发门禁，包含Go/前端/静态/strict、安全与隔离集成验收。`grade.sh`宿主只需Docker与sh：单独建无宿主端口的测试栈，明确用模型/嵌入fixture，真实运行MySQL、RabbitMQ、Qdrant、MCP与Chromium。测试覆盖空库初始化、角色/班级拒绝、审计/任务事务失败、并发唯一结果、停机/重投和真实浏览器。

报告在`reports/acceptance-*/report.html`、`report.json`；必测失败、缺阶段、清理失败均非零。10秒/30秒延迟是受控机制实验，嵌入调用上限在测试中设35秒，不能解读成真实AI效果或速度。`TEST_ACK_DELAY_MS`只允许test环境，生产拒绝该设置。

无本机配置的开发回归可运行`sh scripts/integration.sh`，它生成临时独立测试配置，结束删除配置并保留脱敏报告；不能将它的fixture配置用于正式运行。

## 数据迁移与回退

首次API启动执行幂等增量建表与独立种子；已有账号不会因env修改而自动覆盖密码。已有发布迁移不得编辑，后续只新增版本。当前基础表迁移仍处于本变更未发布阶段。

上线顺序：先兼容建表/索引，启动依赖和执行进程，再启用入口。回退先关闭新增写入口、保留任务/outbox/卷，再退兼容镜像。没有实际双版本回退或备份恢复的演练，不能声称生产就绪。

密钥泄露先撤销/轮换，验证旧值失效、新值生效；只记键名和结果。MySQL/RabbitMQ已有卷时，改env不会自动改服务器口令，须先执行受控服务端变更。JWT密钥轮换会使旧令牌失效；MCP轮换必须同步重建两端并复测。
