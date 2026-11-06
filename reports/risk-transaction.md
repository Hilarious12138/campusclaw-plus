# 事务与一致性验证

验证时间：2026-11-07（Asia/Shanghai）；受验代码：`717c14c569c67b8e02426a995159de0ee5f5d7c3`。

`make check`退出0，Docker Linux arm64独立环境22/22通过。MySQL、RabbitMQ、Qdrant、MCP及Chromium为真实容器；对话与嵌入为明确标注的测试替身。真实AI、评分、独立人类评审、异机复验未完成。原始HTML/JSON与截图在本地`reports/acceptance-20261106T163702Z-99671/`，不进入Git。

B01：审计表/发件箱故障下拒绝写入，核对业务库无半写。A02：停止broker后上传仍202，恢复后可靠发件箱自动补投，任务完成。材料正文、任务与outbox同事务；外部嵌入/Qdrant调用在事务外，最终按租约与索引版本发布。此处没有验证数据库备份恢复目标，不能推导灾备已完成。
