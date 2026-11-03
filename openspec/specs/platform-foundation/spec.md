# 平台工程基线（platform-foundation）

## Purpose

为Plus提供独立于GD的可复现起点。配置错误要明确失败，公开接口保持最小范围，依赖状态和请求日志可追踪，后续教学能力在此基础上增量交付。

## Requirements

### Requirement: F1 Reproducible isolated startup

**独立配置与启动**

系统 SHALL 满足以下要求：
- 仅按README说明，就能完成clone、填写本机`.env`并用一条Compose命令启动。
- 配置样例只写键名，真实秘密不能出现在全部Git历史中。
- Plus不能复用GD运行资源；必填运行配置缺失或非法 MUST 启动失败。

#### Scenario: S01 Clean startup normal
- **GIVEN** 前置条件：在新目录clone，已填全部有效配置，Docker可用。
- **WHEN** 操作：执行`docker compose up --build -d --wait`。
- **THEN** 结果：web/api/db健康，只有web发布本机配置端口，Plus使用独立卷。
- **AND** GET `/health`在2秒内返回200/OK。

#### Scenario: S02 Missing or out of range config failure
- **GIVEN** 前置条件：其他配置有效。
- **WHEN** 操作：删除任一API必填项，或将DB_TIMEOUT_MS设为0/10001。
- **THEN** 结果：API非零退出，不监听端口。
- **AND** 日志只指出配置键，不包含原始值。

#### Scenario: S03 Timeout boundary
- **GIVEN** 前置条件：全部必填项有效。
- **WHEN** 操作：DB_TIMEOUT_MS分别设为100和10000。
- **THEN** 结果：两值均被接受；99和10001明确拒绝，不能退回默认值运行。

### Requirement: F2 Minimal public API and default deny

**最小公开接口与默认拒绝**

系统 SHALL 只显式公开GET `/health`和GET `/live`：
- 未知API、附件访问均返回404；非GET方法返回405。
- health/live不接受查询参数。
- 当前骨架不实现登录，不能宣称业务角色鉴权已完成。

#### Scenario: S04 Invalid filter failure
- **GIVEN** 前置条件：系统已就绪。
- **WHEN** 操作：GET `/health?class_id=other`。
- **THEN** 结果：400/INVALID_ARGUMENT，无额外依赖详情或数据，也无数据库写入。

#### Scenario: S05 Unregistered protected operation denied
- **GIVEN** 前置条件：匿名用户没有业务授权。
- **WHEN** 操作：GET `/api/teachers/materials`。
- **THEN** 结果：404/NOT_FOUND，不返回教师数据。
- **AND** 这只验证未登记业务默认不可访问，不替代未来JWT/班级授权测试。

#### Scenario: S06 Missing resource failure
- **GIVEN** 前置条件：系统运行中。
- **WHEN** 操作：GET `/uploads/guessed.txt`或GET `/api/missing`。
- **THEN** 结果：404，不返回文件、前端index或私有数据。

#### Scenario: S07 Unsupported method failure
- **GIVEN** 前置条件：系统运行中。
- **WHEN** 操作：POST `/health`。
- **THEN** 结果：405/METHOD_NOT_ALLOWED，Allow包含GET，没有写入副作用。

### Requirement: F3 Observable readiness and liveness

**区分存活与就绪**

系统 SHALL 允许匿名读取`/health`，返回版本、UTC时间和MySQL连通状态。依赖正常用200/OK，不可达用503/DEPENDENCY_UNAVAILABLE；`/live`仅表示进程存活。

响应 MUST 不含凭证、内部地址或原始驱动错误。

#### Scenario: S08 Healthy dependency normal
- **GIVEN** 前置条件：MySQL真实可达。
- **WHEN** 操作：不带Authorization请求GET `/health`。
- **THEN** 结果：code=OK，data.version等于配置版本，time为UTC RFC3339，dependencies.mysql=up。
- **AND** request_id非空，与X-Request-ID相同。

#### Scenario: S09 Dependency unreachable failure
- **GIVEN** 前置条件：API已启动，Plus MySQL被停止。
- **WHEN** 操作：请求GET `/health`和GET `/live`。
- **THEN** 结果：health在配置超时加1秒容差内返回503/DEPENDENCY_UNAVAILABLE及mysql=down；live仍为200/OK。
- **AND** 恢复DB后，health恢复200。

### Requirement: F4 Traceable logs and graceful shutdown

**日志关联与有界停机**

系统 SHALL 为请求生成或校验标识，在响应和单行JSON日志中关联。日志包含时间、级别、模块、事件、方法、路径、状态和耗时。

秘密与query/body MUST 不进入日志。收到SIGTERM后，按配置在≤10秒内排空已有请求并退出。

#### Scenario: S10 Correlated sanitized log normal
- **GIVEN** 前置条件：系统运行，请求X-Request-ID=verify-platform。
- **WHEN** 操作：GET `/health`。
- **THEN** 结果：响应和请求日志保留该ID；日志可逐行解析JSON，必要字段齐全。
- **AND** 日志无凭证或查询正文。

#### Scenario: S11 Unsafe request id boundary
- **GIVEN** 前置条件：系统运行中。
- **WHEN** 操作：请求ID含非法字符或长度超过64。
- **THEN** 结果：服务端生成新合法ID，不回显原值；长度1和64的合法ID可保留。

#### Scenario: S12 Drain active request normal
- **GIVEN** 前置条件：有未完成、但可在停机预算内结束的请求。
- **WHEN** 操作：服务收到停止信号。
- **THEN** 结果：停止接收新连接，让已有请求完成，关闭依赖后退出。
- **AND** 超出预算时强制结束，不无限等待。

### Requirement: F5 Complete homepage states

**完整的首页反馈**

系统 SHALL 提供中文工程首页，从同源`/health`读取状态。页面要有加载、正常、错误和空状态，并提供重试。未知页面返回404，支持响应式和键盘操作，不能把未实现功能显示为可用。

#### Scenario: S13 Ready page normal
- **GIVEN** 前置条件：health返回200/OK及有效数据。
- **WHEN** 操作：打开`/`。
- **THEN** 结果：显示CampusClaw Plus、服务已就绪、当前版本及MySQL状态。

#### Scenario: S14 Loading empty error boundary
- **GIVEN** 前置条件：health分别处于等待、数据缺失、网络故障或503。
- **WHEN** 操作：页面读取状态。
- **THEN** 结果：显示对应的加载、空或错误文案；错误/空状态可重试。
- **AND** 不无限加载，不伪装成正常。

### Requirement: F6 Executable standards and truthful evidence

**规则可执行，证据真实**

系统 SHALL 将五类规范放入Harness可读取的规则文件，将场景流程放入项目Skill。一键门禁失败即非零。验证记录包含能力与Scenario对应、命令结果和单人例外，不能假冒独立评审。

#### Scenario: S15 Verification and isolation normal
- **GIVEN** 前置条件：规则、spec、tasks和骨架已落库。
- **WHEN** 操作：执行`make check`与`make acceptance`。
- **THEN** 结果：单元、静态、构建、strict和真实Compose验证有可定位结果，不修改GD。
- **AND** 未完成项明确列出，不提前勾选。

#### Scenario: S16 Rule violation rejected failure
- **GIVEN** 前置条件：正常门禁通过。
- **WHEN** 操作：临时让服务层导入HTTP/数据库，或把未知API改为成功放行。
- **THEN** 结果：静态门禁或负向用例失败；还原后重新通过，证明检查能发现相应缺陷。
