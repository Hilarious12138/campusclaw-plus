# 固定AI评测集

`cases.json`固定12条输入与判据。认证、上下文、工具和故障可以用可重复夹具验证；语义质量、模型泄露与提示词抵抗必须在配置真实供应商后另跑，保留模型/模板/Skill版本、实际输出与判定。不能把夹具的固定回答记录成真实模型12/12通过。

已有自动入口：`cases/api/conversations/POST.py`、`cases/api/mcp/servers/POST.py`、`cases/api/retrieval/search/POST.py`及适配器单测。未逐项实测的条目保留未完成；`reports/sec08-experiment.md`须说明实际运行环境与范围。
