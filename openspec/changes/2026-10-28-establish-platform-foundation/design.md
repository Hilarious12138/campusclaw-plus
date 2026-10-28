## Context

Plus新仓只有首次说明提交，无业务源码；见proposal。用户已要求完成第7课四步，实施边界为骨架，未来教学目标写README路线图。

## Goals / Non-Goals

设计提供可观测、可复现、可扩展的壳；不引入JWT伪实现、空业务表、模型Key或GD数据。无权限失败以未登记业务路由默认404覆盖；真实业务JWT/资源越权留对应未来变更。

## Decisions

### D1 分层与依赖倒置

backend/cmd/api装配配置、数据库适配器、service、httpapi；domain定义Dependency.Ping(ctx)与健康结果；service依赖domain端口。repository包装database/sql MySQL Ping；httpapi只处理方法、参数、JSON和日志。替代“全部写main”无法用静态检查证明边界，故不选。

### D2 最小接口与统一响应

GET /health为就绪、GET /live为存活；未知/API/附件404，不把所有路径fallback到index。health/live任何query400，非GET405。code与HTTP均明确：OK/200、INVALID_ARGUMENT/400、NOT_FOUND/404、METHOD_NOT_ALLOWED/405、DEPENDENCY_UNAVAILABLE/503。错误error.message有下一步，不泄露驱动原文。

健康JSON：{code,data:{status,version,time,dependencies:{mysql:up|down}},request_id}；依赖失败同时error:{message}。request_id与响应头一致，合法字符[A-Za-z0-9_-]长度1..64，非法则随机生成。替代“只回ok:true”不便定位依赖与版本，不选。

### D3 配置与真实依赖

API必填APP_ENV/APP_VERSION/HTTP_ADDR/DB_ADDR/DB_NAME/DB_USER/DB_PASSWORD/DB_TIMEOUT_MS/SHUTDOWN_TIMEOUT_SECONDS；Compose另需WEB_PORT/DB_ROOT_PASSWORD。DB_TIMEOUT_MS范围100..10000，停机1..10秒，秘密至少24位且不记原值。后续无新增业务配置时不能因未填模型Key阻塞骨架。

MySQL8创建空独立schema；Go启动不要求DB已连通，health如实503并可恢复；连接池不重建服务。替代SQLite会改变用户既定栈，静态假探针会假报健康，均不选。

### D4 日志、部署和前端

Go slog单行JSON，注入module/request_id；Nginx也用escape=json记录请求（duration_seconds为秒，Go duration_ms为毫秒），统一传播X-Request-ID。只映射web→localhost:WEB_PORT，/api与health/live反代Go，/uploads固定404。所有镜像用本机已缓存的上游摘要固定；Go镜像跑非root，凭证只runtime env，无.env COPY。

React18 TS Vite只做工程状态页，stateFromResponse纯函数支持四态单测，真实Chrome验证页面；Vite同源代理web，生产Nginx路径一致。未知静态路径404，不做伪SPA路由。health探测超时和优雅停机可注入测试；信号context控制Shutdown，先排空再关闭DB。

### D5 验证与课程适配

make check用Docker独立测试容器，不依赖本机Go；前端npm ci/状态测试/类型检查/构建，Go test-race/vet，Python静态结构/秘密历史扫描，OpenSpec strict。make acceptance以Plus独立Compose执行匿名健康、参数/方法/未知路由、停DB、恢复、JSON日志、停止耗时，并真实Chrome验证。禁止操作GD服务。

正常/失败/边界标识与S01..S16写入docs/lesson-07/verification.md。另做分层/放行变异验证与干净clone验证；外部公告审计独立make audit，联网失败不写通过。第一版本没有前一运行版本，只验证当前标签重启，不假造回滚。

郭地三角色自审，独立同伴/异机不能完成就记差异；空白Harness另启无上下文进程验证，没执行前不勾选。首次提交README/.gitignore，其后规则/OpenSpec提交先于实现；PR经CI后合并。

## Risks / Trade-offs

- 课程写slash命令而本机Codex生成6个项目Skill → 使用真实openspec CLI和Skill；validate以CLI strict执行，记录映射，不编造命令。
- 课件无权限要求与纯骨架冲突 → 以默认拒绝未注册资源体现边界，不引入假登录；未来业务变更补完整401/403/跨班404场景。
- 依赖扫描/镜像下载受网络影响 → 缓存摘要复用、记录网络失败，不能假通过。
- 新项目无前一版本/真实业务数据 → 回滚与数据恢复不适用部分写N/A，未来首次业务上线前补演练。

## Migration Plan

独立Compose项目campusclaw-plus，建议端口3001。配置必填，docker compose up --build -d --wait后检查/health。数据卷只在本项目命名，down默认保留；GD不改。首版本无业务迁移；回到同提交/标签重建验证可复现，不使用down -v删除用户数据。

## Implementation adjustments

- npm审计发现Vite7.2.0的高危项，按公告固定到7.3.7；Go1.27.1标准库有可达漏洞，固定到1.27.3及镜像摘要，govulncheck改为兼容的1.9.0。
- MySQL配置若经FormatDSN再Open会丢失Logger，故使用NewConnector/OpenDB，真实断连验证请求日志保持JSON。
- check/audit各用独立node_modules卷，避免并行npm ci互相覆盖；浏览器本机依赖另装。
- 初次骨架涉及多文件工程配置/测试，不是新增整套教学业务；生成OpenSpec技能较长但不重复载入无关内容。
