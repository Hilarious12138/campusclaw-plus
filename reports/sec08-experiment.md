# AI安全实验状态

验证时间：2026-11-07（Asia/Shanghai）；受验代码：`717c14c569c67b8e02426a995159de0ee5f5d7c3`。

`make check`退出0，Docker Linux arm64独立环境22/22通过。MySQL、RabbitMQ、Qdrant、MCP及Chromium为真实容器；对话与嵌入为明确标注的测试替身。真实AI、评分、独立人类评审、异机复验未完成。原始HTML/JSON与截图在本地`reports/acceptance-20261106T163702Z-99671/`，不进入Git。

用户明确暂缓AI API。恶意材料、索要系统提示词、诱导修改规则、诱导越权工具四步真实模型实验未执行，不填写通过。固定输入见apitest的评测集。已验证的服务层授权/MCP断开拒绝可独立成立；fixture输出不能证明模型抗提示词攻击。评分业务不存在，评分偏移实验同样未完成。
