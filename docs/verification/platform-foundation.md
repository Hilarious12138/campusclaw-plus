# 平台基础能力验证记录

本文件描述工程基线和验证范围，不代表全部产品能力已实现。原始机器摘要见[platform-foundation.json](platform-foundation.json)；该摘要保持原始内容，受验提交b2d448b，日期2026-10-28。

## 已验证范围

| 能力 | 证据与判据 |
| --- | --- |
| 独立配置与启动 | 配置缺失/越界拒绝，独立Compose项目/卷；同机干净clone启动验证通过 |
| 最小公开接口 | health/live明确登记；未知API/附件404、非法query400、非GET方法405 |
| 就绪与存活 | 真实MySQL正常时health200/up；停库后health503/down而live200；恢复后health200 |
| 日志与生命周期 | JSON请求日志、响应关联标识、秘密不入日志；在途排空与超预算停机测试 |
| 前端状态 | 加载、正常、空、错误/重试；浏览器响应式与键盘检查 |
| 工程约束 | 分层违规/错误放行变异会被发现，恢复后重新通过；静态/历史秘密及依赖审计 |

原始摘要含实际HTTP结果、耗时、工具版本和限制。受验时npm及Go审计通过不表示当前或未来依赖永远无漏洞，应按每次CI重新扫描。

## 自动验收

```sh
sh scripts/grade.sh
```

宿主只需Docker Compose与sh。验收镜像包含构建、检查和浏览器工具，使用不发布端口的唯一测试项目；真实停库/恢复及API退出仅作用于该项目。结果输出为`reports/acceptance-*/report.html`与`report.json`。

历史版本15个必测覆盖静态/历史、配置、Go race/vet、前端、strict、漏洞、报告判定、真实HTTP/数据库故障/恢复、浏览器、日志与停机。未执行/失败/清理失败均非零，不以空列表判通过。报告内有受验提交、环境、范围和证据；每次运行按实际结果判断。

已有Linux运行证据：[基础验收](https://github.com/Hilarious12138/campusclaw-plus/actions/runs/40617258614)、[工程文档整理验收](https://github.com/Hilarious12138/campusclaw-plus/actions/runs/41445814906)。这些链接证明相应历史提交，不能代表当前版本已重新验收。

## 过程与限制

- 早期审计发现依赖/工具链问题，更新后重跑；驱动日志格式、配置校验、停机预算与负向断言缺口均修复后验证，原始过程由Git提交/PR追溯。
- 同机干净clone不是异机复现；AI辅助静态核对不是独立人类评审。两项独立验证仍待完成。
- 该历史基线受验时没有业务表、登录、Qdrant、模型或MCP调用；新增业务证据见[安全与转换映射](secure-platform-scenarios.md)，不覆盖改写原始摘要。
- 原始个人材料已单独保存；公开记录只保留项目适用范围。当前文件整理不改变旧提交或标签内容。

## 项目文档与验收命名更新

2026-11-03，实施提交f542bba通过[Linux check/grade](https://github.com/Hilarious12138/campusclaw-plus/actions/runs/41633637659)。15项必测通过，CI产物名acceptance-report、输出目录acceptance-前缀及通用报告范围均已检查；实际浏览器截图角标为工程基线。范围只涵盖基础能力，不包含尚未实现的业务。最终归档版本由PR门禁继续验证。
