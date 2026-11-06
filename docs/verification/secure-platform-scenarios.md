# 安全与转换验收映射

以下是验收入口，不表示未执行项通过。原始值以同提交 `report.json` 与各阶段记录为准。

| 要求 / 场景 | 自动检查入口 | 观察信号 |
| --- | --- | --- |
| R1 S01—02 | gd-conversion.md、check.py、Go测试 | 来源SHA/摘要、独立测试、分层 |
| R2 S03 | auth/login/POST.py、credentials_test.go | 哈希、令牌校验/撤销、数据库与存储 |
| R2 S04—05 | skills/POST.py、materials/POST.py、conversations/POST.py | 401/403/404、DB无副作用 |
| R3 S06—07 | retrieval/search/POST.py、skills/POST.py、browser.mjs | SQL参数绑定、转义DOM、Origin/无Bearer |
| R4 S08—09 | audit-events/GET.py、skills/POST.py | 七列、本班过滤、审计失败业务回滚 |
| R5 S10—11 | security_scan.py、gate-evidence.py、audit.sh | 非零/零、实际漏洞公告与升级复验 |
| R6 S12 | async_cases.py、grade.sh | 20次两组P95、stop/start、ACK前kill、唯一结果 |
| R6 S13 | retrieval/search/POST.py、index_test.go | 三模式、正文回查、双侧过滤、无依据 |
| R7 S14 | conversations/POST.py、browser.mjs | 实际快照、模板/Skill版本、SSE、引用 |
| R7 S15 | mcp/servers/POST.py、mcp-server测试 | 真实工具、断开/越权拒绝、审计 |
| R7 S16 | 用户已暂缓AI接入，后续配置后执行固定实验 | 不能以fixture回答证明提示词抗攻击 |
| R8 S17—18 | grade.sh、固定reports | 证据、修复同场景、缺阶段失败 |

本阶段作业评分业务、真实外部模型/向量服务的供应商验证、独立人类同伴评审和异机复验尚未完成。测试中的嵌入/对话fixture明确标识；MySQL、RabbitMQ、Qdrant、MCP及Chromium均为真实独立容器。fixture结果不能推导真实模型质量或安全性。
