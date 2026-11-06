# Campus tools MCP

同仓独立进程，提供 `/health` 与 Bearer 鉴权的 `POST /mcp`。工具：`formula_hint`、`unit_convert`。

实现与测试位于 `backend/cmd/mcp-server/`，沿用后端锁定工具链；Compose 中的 `mcp-demo` 构建启动，不映射宿主端口。凭证仅通过运行时环境变量 `MCP_AUTH_TOKEN` 注入，主应用登记 `http://mcp-demo:8090` 后连接发现。
