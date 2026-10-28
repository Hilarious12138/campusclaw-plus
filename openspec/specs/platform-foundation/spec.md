# platform-foundation Specification

## Purpose
为CampusClaw Plus提供独立于GD的可复现工程基线：运行配置缺失即失败，工程接口最小公开，必要依赖就绪情况可观测，日志可追踪且不泄露秘密，后续教学能力在此基础增量交付。

## Requirements

### Requirement: F1 Reproducible isolated startup
系统 SHALL 仅用README说明完成clone、填本机.env、一条Compose命令启动；配置样例仅键名，秘密不进入全部Git历史；Plus运行资源不复用GD。所有运行必填项缺失/非法 MUST 启动失败。

#### Scenario: S01 Clean startup normal
- **GIVEN** 新目录clone并填全部有效配置，Docker可用
- **WHEN** 执行docker compose up --build -d --wait
- **THEN** web/api/db健康，只有web发布本机配置端口，独立Plus卷存在，GET /health在2秒内返回200/OK。

#### Scenario: S02 Missing or out of range config failure
- **GIVEN** 其他配置有效
- **WHEN** 删除任一API必填项或将DB_TIMEOUT_MS设为0/10001
- **THEN** API进程非零退出且不监听；日志仅指出配置键、不包含原始值。

#### Scenario: S03 Timeout boundary
- **GIVEN** 全部必填项有效
- **WHEN** DB_TIMEOUT_MS分别设为100与10000
- **THEN** 配置接受；99与10001明确拒绝，不带默认值运行。

### Requirement: F2 Minimal public API and default deny
系统 SHALL 只显式公开GET /health及GET /live，未知API与附件访问均404，非GET方法405，health/live不接受查询参数。骨架不实现登录，不宣称完成业务角色鉴权。

#### Scenario: S04 Invalid filter failure
- **GIVEN** 系统就绪
- **WHEN** GET /health?class_id=other
- **THEN** 400/INVALID_ARGUMENT，无额外依赖详情或数据，无数据库写入。

#### Scenario: S05 Unregistered protected operation denied
- **GIVEN** 匿名用户无任何业务授权
- **WHEN** GET /api/teachers/materials
- **THEN** 404/NOT_FOUND，不返回教师数据；骨架未登记的业务能力默认不可访问，此场景不替代未来JWT/班级授权测试。

#### Scenario: S06 Missing resource failure
- **GIVEN** 系统运行
- **WHEN** GET /uploads/guessed.txt或GET /api/missing
- **THEN** 404，不返回文件、前端index或私有数据。

#### Scenario: S07 Unsupported method failure
- **GIVEN** 系统运行
- **WHEN** POST /health
- **THEN** 405/METHOD_NOT_ALLOWED，Allow包含GET，没有写入副作用。

### Requirement: F3 Observable readiness and liveness
系统 SHALL 匿名返回/health的版本、UTC时间和MySQL连通状态；必要依赖正常200/OK，依赖不可达503/DEPENDENCY_UNAVAILABLE，/live只表示进程存活。响应 MUST 不含凭证、内部地址或原始驱动错误。

#### Scenario: S08 Healthy dependency normal
- **GIVEN** MySQL真实可达
- **WHEN** 不带Authorization访问GET /health
- **THEN** code=OK，data.version为配置版本，time为UTC RFC3339，dependencies.mysql=up，request_id非空且与X-Request-ID相同。

#### Scenario: S09 Dependency unreachable failure
- **GIVEN** API已启动，Plus MySQL被停止
- **WHEN** GET /health和GET /live
- **THEN** health在配置超时加1秒容差内返回503/DEPENDENCY_UNAVAILABLE与mysql=down，live返回200/OK；恢复DB后health恢复200。

### Requirement: F4 Traceable logs and graceful shutdown
系统 SHALL 每个请求生成或校验请求标识，在响应与单行JSON日志中关联；字段包含时间、级别、模块、事件、方法、路径、状态、耗时。秘密和query/body MUST 不进入日志；收到SIGTERM在配置≤10秒内排空已有请求并退出。

#### Scenario: S10 Correlated sanitized log normal
- **GIVEN** 系统运行且请求X-Request-ID=verify-lesson07
- **WHEN** GET /health
- **THEN** 响应与请求日志均为该ID，日志可逐行解析JSON，必要字段齐备，无凭证及查询正文。

#### Scenario: S11 Unsafe request id boundary
- **GIVEN** 系统运行
- **WHEN** 请求ID含非法字符或长度超过64
- **THEN** 服务端生成新的合法ID，不回显原值；长度1和64的合法ID可保留。

#### Scenario: S12 Drain active request normal
- **GIVEN** 存在未完成但可在停机预算内完成的请求
- **WHEN** 服务接收停止信号
- **THEN** 不再接收新连接，已有请求正常完成，关闭依赖后退出；超出预算强制结束，不无限等待。

### Requirement: F5 Complete homepage states
系统 SHALL 提供中文工程首页，读取同源/health并显示加载/正常/错误/空状态，提供重试操作；未知页面404，响应式与键盘操作可用，不显示尚未实现的教学功能为可用。

#### Scenario: S13 Ready page normal
- **GIVEN** health返回200/OK的有效数据
- **WHEN** 打开/
- **THEN** 显示CampusClaw Plus、服务已就绪、当前版本和MySQL状态。

#### Scenario: S14 Loading empty error boundary
- **GIVEN** health请求分别处于等待、返回缺失数据、网络故障或503
- **WHEN** 页面读取状态
- **THEN** 分别显示规定的加载、空、错误文案，错误/空可重试；不得无限加载或伪装正常。

### Requirement: F6 Executable standards and truthful evidence
系统 SHALL 将五类规范落在Harness读取的规则文件，按场景流程落在项目Skill；一键门禁失败即非零，验证记录含能力-Scenario对应、命令结果与单人例外，不假冒独立评审。

#### Scenario: S15 Verification and isolation normal
- **GIVEN** 规则、spec、tasks和骨架已落库
- **WHEN** 执行make check与make acceptance
- **THEN** 单元/静态/构建/strict及实际Compose验证有可定位结果，不改GD；未完成项明确列出而不预勾选。

#### Scenario: S16 Rule violation rejected failure
- **GIVEN** 正常门禁通过
- **WHEN** 暂时让服务层导入HTTP/数据库，或将未知API改为成功放行
- **THEN** 静态门禁或负向用例失败；还原后重新通过，证明规则确实能发现缺陷。
