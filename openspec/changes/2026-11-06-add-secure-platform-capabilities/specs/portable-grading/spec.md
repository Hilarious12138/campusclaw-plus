## MODIFIED Requirements

### Requirement: G1 Docker-only execution

**宿主只需Docker**

系统 SHALL 在Linux Docker环境通过`sh scripts/grade.sh`运行已实现能力的全部适用验收。Python、Node、Go、Chrome和构建工具均在容器中，宿主无需安装。

#### Scenario: P01 Valid environment normal
- **GIVEN** 前置条件：已clone、填好全部`.env`，Docker Compose可用。
- **WHEN** 操作：执行`sh scripts/grade.sh`。
- **THEN** 结果：逐项执行静态/历史秘密、单元/race/vet、前端测试/构建、strict、漏洞及真实HTTP/DB/队列故障、业务接口和浏览器检查；模型夹具与真实供应商证据分开。
- **AND** 全部通过时退出0。

#### Scenario: P02 Invalid configuration failure
- **GIVEN** 前置条件：`.env`缺必填项或密码非法。
- **WHEN** 操作：执行验收。
- **THEN** 结果：非零退出并生成基础设施失败报告，不展示凭证，不把未运行项写成成功。
