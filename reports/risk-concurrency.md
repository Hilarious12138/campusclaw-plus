# 并发验证

验证时间：2026-11-07（Asia/Shanghai）；受验代码：`717c14c569c67b8e02426a995159de0ee5f5d7c3`。

`make check`退出0，Docker Linux arm64独立环境22/22通过。MySQL、RabbitMQ、Qdrant、MCP及Chromium为真实容器；对话与嵌入为明确标注的测试替身。真实AI、评分、独立人类评审、异机复验未完成。原始HTML/JSON与截图在本地`reports/acceptance-20261106T163702Z-99671/`，不进入Git。

B04：20个真实MySQL并发同幂等键只建一任务；旧租约失效后旧worker提交拒绝，新worker仅生效一份。A01：结果已持久化、ACK前kill消费者，重投不重复写结果。429最多3次，永久401仅1次，A03核对持久DLQ与无假成功。未将单元测试替身称为真实队列。
