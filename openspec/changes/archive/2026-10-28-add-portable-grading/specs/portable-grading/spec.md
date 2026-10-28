## Purpose

为仅准备Docker Compose的教师部署环境提供可复现自动验收，在隔离资源中覆盖当前已交付能力并生成真实可追溯的HTML与JSON报告；不得把缺失或失败的检查计为通过，不依赖宿主开发运行时。

## ADDED Requirements

### Requirement: G1 Docker-only execution
系统 SHALL 通过sh scripts/grade.sh在Linux Docker环境运行当前骨架全部自动验收，宿主不需要Python、Node、Go或Chrome；构建工具和浏览器均在容器。

#### Scenario: P01 Valid environment normal
- **GIVEN** 克隆仓库、填全部.env且Docker Compose可用
- **WHEN** 执行sh scripts/grade.sh
- **THEN** 静态/历史秘密、单元/race/vet、前端测试/构建、strict、漏洞和真实HTTP/DB故障/浏览器用例逐项执行；通过时退出0。

#### Scenario: P02 Invalid configuration failure
- **GIVEN** .env缺必填项或密码非法
- **WHEN** 执行验收
- **THEN** 退出非零并生成基础设施失败报告，不展示凭证，不把未运行项写成成功。

### Requirement: G2 Truthful reports
系统 SHALL 每次产生独立报告目录，包含自包含report.html、report.json和脱敏证据；列出版本、范围、逐项结果、判据、耗时及错误。任何必测失败或未运行 MUST 总体非零，未来业务N/A不计通过。

#### Scenario: P03 Case failure report
- **GIVEN** 一条必测断言被破坏
- **WHEN** 验收运行结束
- **THEN** 报告该项失败，其他可运行检查继续，总体失败并非零。

#### Scenario: P04 Dependency unavailable
- **GIVEN** 镜像下载或漏洞公告依赖不可达
- **WHEN** 构建/扫描超时或报错
- **THEN** 记录失败或阻塞并给出重试步骤，不输出“无漏洞”或成功，仍生成HTML。

#### Scenario: P05 Missing evidence
- **GIVEN** 验收阶段没有产生预期必测记录
- **WHEN** 汇总
- **THEN** 补一条失败/阻塞，不能因结果列表为空而判通过。

### Requirement: G3 Resource and credential isolation
验收 SHALL 不发布测试端口、不挂载Docker socket，使用唯一Plus项目与自建卷；结束仅清理自建资源，报告保留，正式部署和GD不受影响。报告 MUST 不包含真实配置秘密。

#### Scenario: P06 Denied operations boundary
- **GIVEN** 骨架已启动
- **WHEN** 匿名访问未登记教师资源、未知API和附件
- **THEN** 404，无数据泄露；该检查不冒充业务JWT权限验收。

#### Scenario: P07 Coexistence normal
- **GIVEN** 正式Plus已占用3001，GD可能占用3000
- **WHEN** 自动验收启动、停库、停API、清理
- **THEN** 无端口冲突，不停止其他项目，正式Plus仍健康，报告不含密码。
