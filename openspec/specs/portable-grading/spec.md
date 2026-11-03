# 便携自动验收（portable-grading）

## Purpose

让仅有Docker Compose的使用者在隔离环境中复验已交付能力，取得可追溯的HTML/JSON报告。宿主不需要开发运行时；失败或缺少检查都不能计为通过。

## Requirements

### Requirement: G1 Docker-only execution

**宿主只需Docker**

系统 SHALL 在Linux Docker环境通过`sh scripts/grade.sh`运行全部当前骨架验收。Python、Node、Go、Chrome和构建工具均在容器中，宿主无需安装。

#### Scenario: P01 Valid environment normal
- **GIVEN** 前置条件：已clone、填好全部`.env`，Docker Compose可用。
- **WHEN** 操作：执行`sh scripts/grade.sh`。
- **THEN** 结果：逐项执行静态/历史秘密、单元/race/vet、前端测试/构建、strict、漏洞及真实HTTP/DB故障/浏览器检查。
- **AND** 全部通过时退出0。

#### Scenario: P02 Invalid configuration failure
- **GIVEN** 前置条件：`.env`缺必填项或密码非法。
- **WHEN** 操作：执行验收。
- **THEN** 结果：非零退出并生成基础设施失败报告，不展示凭证，不把未运行项写成成功。

### Requirement: G2 Truthful reports

**报告不能假成功**

系统 SHALL 为每次运行创建独立目录，保存自包含的report.html、report.json和脱敏证据。自包含表示无需外部网页资源即可查看。

报告写明版本、范围、逐项结果、判据、耗时和错误。任一必测失败或未执行 MUST 总体非零；未来业务不适用项不计通过。

#### Scenario: P03 Case failure report
- **GIVEN** 前置条件：一条必测断言被破坏。
- **WHEN** 操作：验收结束。
- **THEN** 结果：该项记失败，其他可运行检查继续；总体失败并非零退出。

#### Scenario: P04 Dependency unavailable
- **GIVEN** 前置条件：镜像下载或漏洞公告依赖不可达。
- **WHEN** 操作：构建/扫描超时或报错。
- **THEN** 结果：记失败或阻塞，给出重试步骤，仍生成HTML。
- **AND** 不输出成功或“无漏洞”。

#### Scenario: P05 Missing evidence
- **GIVEN** 前置条件：某阶段没有生成预期的必测记录。
- **WHEN** 操作：汇总结果。
- **THEN** 结果：补记失败或阻塞，不能因列表为空而判通过。

### Requirement: G3 Resource and credential isolation

**隔离资源，保护凭证**

验收 SHALL 使用唯一Plus项目及自建卷，不发布测试端口、不挂载Docker socket。结束只清理自建资源，保留报告，不影响正式部署或GD。报告 MUST 不包含真实配置秘密。

#### Scenario: P06 Denied operations boundary
- **GIVEN** 前置条件：骨架已启动。
- **WHEN** 操作：匿名访问未登记教师资源、未知API或附件。
- **THEN** 结果：404且无数据泄露；不把该检查当作业务JWT权限验收。

#### Scenario: P07 Coexistence normal
- **GIVEN** 前置条件：正式Plus占用3001，GD可能占用3000。
- **WHEN** 操作：自动验收启动、停库、停API并清理。
- **THEN** 结果：无端口冲突，不停止其他项目，正式Plus仍健康，报告不含密码。
