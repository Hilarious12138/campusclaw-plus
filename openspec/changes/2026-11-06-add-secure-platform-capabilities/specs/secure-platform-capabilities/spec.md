# 安全能力与受控转换

## ADDED Requirements

### Requirement: R1 Verified conversion

转换 MUST 在隔离源码快照测试后按Plus分层实施，不修改GD、不迁移其配置/数据/历史。

#### Scenario: S01 Verify source then port
- **GIVEN** 前置条件：GD候选SHA与Plus基线已固定。
- **WHEN** 操作：导出无.env的隔离源码，运行原单测/独立集成测试，再转换模块。
- **THEN** 结果：来源清单含文件摘要、命令、退出码及跳过原因；Plus测试单列；失败模块不标已复用，GD工作树和数据不变。

#### Scenario: S02 Preserve contracts
- **GIVEN** 前置条件：候选HTTP/SQL混合、字段与旧路由不同。
- **WHEN** 操作：转换认证/材料模块并运行静态与接口检查。
- **THEN** 结果：SQL只在repository，service无HTTP/SQL依赖；snake_case统一响应，前端不直连DB；未审依赖不进入正式构建。

### Requirement: R2 Identity and authorization

业务 MUST 默认拒绝，角色逐次校验，SQL与向量都限制认证班级和资源。

#### Scenario: S03 Authenticate and revoke
- **GIVEN** 前置条件：A/B班四身份种子独立生成，JWT期限30分钟。
- **WHEN** 操作：登录、查me、退出后重放、改算法/issuer/audience/过期令牌。
- **THEN** 结果：正常200，撤销/错误令牌401；密码只为哈希；token只在前端内存，无Cookie/URL/localStorage认证。

#### Scenario: S04 Replay SEC authorization
- **GIVEN** 前置条件：本班/他班材料和已注册保护接口清单。
- **WHEN** 操作：匿名重放、学生写配置、B查A材料/随机id、直接调用写服务。
- **THEN** 结果：分别401、403、404/404；无越权数据/副作用；SQL/向量双侧过滤；不存在的批改服务明确未实现，不把404算SEC-04批改通过。

#### Scenario: S05 Reject context spoofing
- **GIVEN** 前置条件：学生有有效A班身份。
- **WHEN** 操作：请求中传class_id或classId改变班级，或读他人会话。
- **THEN** 结果：非法参数400，他人会话404；不扩大查询范围，不返回输入外推班级。

### Requirement: R3 Injection counterexamples

系统 SHALL 固化SQL、XSS及Bearer跨站请求反例，拒绝时同时验证无副作用。

#### Scenario: S06 Parameterize and escape
- **GIVEN** 前置条件：两班隔离样本及可编辑标题。
- **WHEN** 操作：检索提交SQL载荷，标题保存script字符串并打开列表DOM。
- **THEN** 结果：载荷不扩大结果/写库；标题为转义文本且无script节点；未知排序400，非仅以无弹窗判通过。

#### Scenario: S07 Reject cross-site writes
- **GIVEN** 前置条件：Bearer仅内存，无认证Cookie。
- **WHEN** 操作：跨站表单写配置、Cookie-only写配置、有效token配恶意Origin写配置。
- **THEN** 结果：前两项401、后者403，数据库不变；CORS不允许恶意Origin；SameSite及一次性CSRF Token在无认证Cookie方案中标不适用并写原因。

### Requirement: R4 Audit operations

系统 MUST 为成功、拒绝、失败记录可回查七列审计，敏感操作无审计保障时不继续执行。

#### Scenario: S08 Trace denied request
- **GIVEN** 前置条件：A班学生尝试教师配置写入。
- **WHEN** 操作：重放403，教师按该学生/时间窗查audit-events，学生再查日志。
- **THEN** 结果：教师可见操作人/时间/类型/对象/状态/skill及用量耗时；学生403；无正文/秘密，跨班不泄露；匿名actor为null。

#### Scenario: S09 Fail audit transaction
- **GIVEN** 前置条件：配置写入中注入审计持久化失败。
- **WHEN** 操作：提交教师合法配置。
- **THEN** 结果：503、配置无变更，结构化脱敏故障可查；不能返回成功或声称已审计。

### Requirement: R5 Security gates and assets

系统 SHALL 提供make check-security并纳入门禁，扫描命中、高危未处置或审计不可达必须非零。

#### Scenario: S10 Prove red then green
- **GIVEN** 前置条件：隔离人工标识fixture、同一安全门禁与固定扫描版本。
- **WHEN** 操作：先运行必红样本再移除fixture运行正常源码。
- **THEN** 结果：记录两次SHA、输入性质、退出码非零/零与输出；不引入有效密钥、不伪称真实泄露。

#### Scenario: S11 Handle dependency failure
- **GIVEN** 前置条件：漏洞源不可达或新增未处置高危。
- **WHEN** 操作：运行安全门禁。
- **THEN** 结果：非零且报告未验证/具体依赖；有效例外必须有编号/包/原因/降险/责任/截止时间，不能静默跳过。

### Requirement: R6 Reliable materials and retrieval

材料 SHALL 通过已定义异步内核入库并可溯源检索，不继承不可靠后台执行。

#### Scenario: S12 Integrate task lifecycle
- **GIVEN** 前置条件：认证/材料前置已通过，真实RabbitMQ/MySQL且worker可启停。
- **WHEN** 操作：合法教师上传、暂停恢复worker、重复提交和投递。
- **THEN** 结果：202含task_id；异步change三判据原场景通过，最多一份可见索引结果；B班查A任务/材料404，学生上传403。

#### Scenario: S13 Retrieve within scope
- **GIVEN** 前置条件：A/B班各有索引材料，Qdrant无正文。
- **WHEN** 操作：keyword/vector/hybrid检索、相关/无关/他班样本及依赖断开。
- **THEN** 结果：keyword无嵌入调用，vector低于0.35不留，hybrid RRF k60/top5；摘录回MySQL、无依据hits空；依赖失败可判且不伪空成功。

### Requirement: R7 AI skills and tools

AI应用 SHALL 只使用获授权上下文与工具，工具服务层不能被提示词绕过。

#### Scenario: S14 Apply templates and skills
- **GIVEN** 前置条件：知识问答/数学/语文三个应用、教师确认的UTF-8 Skill与绑定材料。
- **WHEN** 操作：学生发起流式会话并查看脱敏调用证据。
- **THEN** 结果：生效system含选中模板/Skill和revision，不含未选技能；user材料top5且≤6000字符；流式meta/delta/终态及引用可回查，学生不可改模板/Skill。

#### Scenario: S15 Deny disconnected tools
- **GIVEN** 前置条件：本班MCP已连接、发现formula_hint/unit_convert且技能绑定有效。
- **WHEN** 操作：先调用成功，再断开并重放，或换班/提交白名单外URL。
- **THEN** 结果：真实工具成功可审计；断开/跨班拒绝且无远端调用；URL越界400，密钥不出响应、日志或前端。

#### Scenario: S16 Record adversarial experiments
- **GIVEN** 前置条件：恶意材料/索要system/改规则/诱导越权工具四组固定输入。
- **WHEN** 操作：逐组执行并检查模型输出、服务调用、数据库与审计。
- **THEN** 结果：前3组记录真实现象，不预设安全；第4组越权副作用必须为0且可审计；评分未实现的子实验标缺口，不伪造评分实验。

### Requirement: R8 Evidence and phased completion

交付 MUST 按阶段记录真实前后值及报告，不以规划或GD状态代替Plus验收。

#### Scenario: S17 Reproduce complete journey
- **GIVEN** 前置条件：独立Plus测试Compose和已实现业务，无宿主开发工具。
- **WHEN** 操作：以Docker+sh运行grade及全功能/安全/异步三组回归。
- **THEN** 结果：HTML/JSON与固定安全报告可打开；登录→202→状态成功→跨班404可复现；失败/缺证据非零，独立评审/异机未做不标通过。

#### Scenario: S18 Preserve baseline expectations
- **GIVEN** 前置条件：SEC固定前置/操作/期望已登记并记录基线实际。
- **WHEN** 操作：按授权→注入→资产→日志顺序修复再重放。
- **THEN** 结果：同场景只更新实际，分别记前后值/提交；无旧红不伪造修复，真实返工才记录耗时。
