## Context

骨架已完成开发机验证。用户补充老师将以Linux Docker和自动报告评分，需要新增便携入口，覆盖当前F1-F6，不提前声称未来教学功能完成。

## Goals / Non-Goals

宿主只调用sh和docker compose；浏览器、Python、Go、Node、OpenSpec及扫描器均在grade镜像。无Docker socket挂载，不新增业务和公网发布。

## Decisions

- grade.sh在宿主编排阶段：容器检查→隔离服务启动→正常/浏览器→停DB→故障→恢复→日志/停机→汇总。只用shell/docker，不要求宿主curl/jq/python。
- compose.grade.yaml复用db/api定义，web不继承宿主ports；项目名带时间/PID，卷独立。grader只连本项目web，报告目录挂载；Git元数据只读用于历史扫描。
- Node报告程序每步捕获状态与耗时，阶段失败继续其他可行项；HTML内嵌CSS/截图，JSON给机器读。缺阶段补阻塞，infra构建失败用shell生成最小HTML/JSON，不以空列表成功。
- grader镜像基于固定Go Alpine工具链，复用GCC层并复制固定Node运行时，安装Chromium/Python/Git/CJK字体，锁OpenSpec与扫描器；不依赖宿主Chrome。首次下载允许耗时，失败报告保留构建日志，不隐藏失败。
- 验收真实DB失败与UI注入分开；shell只停止自己的项目。脚本最后清理测试卷/容器，保留reports，不删正式数据。
- Git history与报告比对实际DB秘密使用runtime env注入grader，不写镜像；报告错误集中脱敏。未来业务每个change增加对应自动用例，当前未来业务列N/A。

## Risks / Trade-offs

- 浏览器镜像较大、首次构建慢 → 缓存Docker层；报告明确网络前置。
- 公告服务波动 → 扫描失败记失败/阻塞并非零，不用忽略漏洞凑通过。
- 此前基础PR仍有收尾交付 → 先合并已验证骨架并打基线，再以独立变更交付便携验收。
