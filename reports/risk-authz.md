# 逐资源授权验证

验证时间：2026-11-07（Asia/Shanghai）；受验代码：`717c14c569c67b8e02426a995159de0ee5f5d7c3`。

`make check`退出0，Docker Linux arm64独立环境22/22通过。MySQL、RabbitMQ、Qdrant、MCP及Chromium为真实容器；对话与嵌入为明确标注的测试替身。真实AI、评分、独立人类评审、异机复验未完成。原始HTML/JSON与截图在本地`reports/acceptance-20261106T163702Z-99671/`，不进入Git。

手段：apitest用真实HTTP与MySQL核对。对象：认证、材料、任务、检索、Skill、MCP、会话、审计。判据：匿名401、学生写403、跨班与资源不存在均404；请求自带班级不改变服务端范围。B01通过，逐请求状态/机器码/request_id见api/request-trace.json；日志不保存Bearer或正文。B02检查浏览器内存令牌和转义DOM。未知接口404本身不证明业务授权。
