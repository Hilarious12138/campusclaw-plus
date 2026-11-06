# 安全基线

验证时间：2026-11-07（Asia/Shanghai）；受验代码：`717c14c569c67b8e02426a995159de0ee5f5d7c3`。

`make check`退出0，Docker Linux arm64独立环境22/22通过。MySQL、RabbitMQ、Qdrant、MCP及Chromium为真实容器；对话与嵌入为明确标注的测试替身。真实AI、评分、独立人类评审、异机复验未完成。原始HTML/JSON与截图在本地`reports/acceptance-20261106T163702Z-99671/`，不进入Git。

全部Git历史/配置/分层扫描通过；npm公告审计0漏洞；govulncheck无漏洞。公告获取失败会阻断，不算通过。C01/C06/C07是本次实际信号。
