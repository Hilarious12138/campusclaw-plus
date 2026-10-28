## Why

老师可能在Linux机器上仅用Docker部署并依据自动报告评分；原本依赖宿主Python、Node、Chrome的开发验收不满足该交付方式。

## What Changes

- 提供sh scripts/grade.sh：宿主仅需Docker Compose、POSIX shell和已填写.env。
- 在独立项目/数据卷运行容器化静态、单元、构建、漏洞、真实接口/故障恢复与Chromium验收。
- 生成自包含HTML和JSON，逐项列出判据、结果、耗时、脱敏错误与截图；必测失败非零，不适用项不计为通过。
- 失败仍产出可读报告，清理仅脚本新建资源，不挂Docker socket给验收容器。

## Capabilities

### New Capabilities
- `portable-grading`: Linux Docker一键验收、隔离与机器可读/人可读报告。

### Modified Capabilities
无，保留开发用make check/acceptance；新增老师入口而不是让老师安装开发环境。

## Non-goals

不自动验收尚未实现的登录/检索/作业等学期业务，不伪造人类同伴评审；后续每项能力实现时同步扩展验收。放弃“老师安装本机Chrome/Node”：环境不一致且操作多；放弃Docker socket进容器：扩大主机控制面。

## Impact

新增grade.sh、Compose验收编排、Linux Chromium/Go/Node/Python工具镜像与报告程序；基础业务容器不变。首次构建需网络，公告服务不可达明确失败，不假报无漏洞。
