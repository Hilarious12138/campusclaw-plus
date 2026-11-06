# 门禁红绿验证

验证时间：2026-11-07（Asia/Shanghai）；受验代码：`717c14c569c67b8e02426a995159de0ee5f5d7c3`。

`make check`退出0，Docker Linux arm64独立环境22/22通过。MySQL、RabbitMQ、Qdrant、MCP及Chromium为真实容器；对话与嵌入为明确标注的测试替身。真实AI、评分、独立人类评审、异机复验未完成。原始HTML/JSON与截图在本地`reports/acceptance-20261106T163702Z-99671/`，不进入Git。

# 安全扫描门禁红绿证明

提交：717c14c569c67b8e02426a995159de0ee5f5d7c3

隔离人工标识fixture，不是真实凭证，也不加入Git历史。相同扫描器，相同期望：命中必须拒绝，删除后扫描通过。此记录只证明秘密扫描子门禁，完整依赖/安全门禁见当日基线。

红退出码：1
```json
{"check": "secret-scan", "findings": 1, "status": "failed"}
```
绿退出码：0
```json
{"check": "secret-scan", "findings": 0, "status": "passed"}
```

真实依赖红绿与人工扫描fixture是两类证据，不混用；真实依赖日志哈希见docs/verification/security-fix-history.json。
