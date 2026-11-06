# 验证中的实际修复

验证时间：2026-11-07（Asia/Shanghai）；受验代码：`717c14c569c67b8e02426a995159de0ee5f5d7c3`。

`make check`退出0，Docker Linux arm64独立环境22/22通过。MySQL、RabbitMQ、Qdrant、MCP及Chromium为真实容器；对话与嵌入为明确标注的测试替身。真实AI、评分、独立人类评审、异机复验未完成。原始HTML/JSON与截图在本地`reports/acceptance-20261106T163702Z-99671/`，不进入Git。

先后发现并修复：AMQP旧依赖可达漏洞、SQL客户端缺认证插件、Rabbit默认交换机权限、浏览器选择器与Origin、30秒HTTP头超时、重复消息ACK人为等待。曾在运行中修改验收shell导致输入偏移失败；之后改为固定提交干净克隆运行，失败日志保留。没有记录可靠的逐项修复工时，不补造数字。最终固定717c14c5代码通过22项。
