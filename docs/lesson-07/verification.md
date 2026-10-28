# 第7课交付与验证记录

日期：2026-10-28（Asia/Shanghai）。项目：CampusClaw Plus，郭地单人项目；独立仓库 https://github.com/Hilarious12138/campusclaw-plus 。GD未纳入本次变更。

## 范围与流程

本课仅工程骨架，教学目标与后续范围见README。本仓首次提交880dbdb只有README.md与.gitignore；之后7adcb2b单独提交规则/Skill/OpenSpec，先strict再实现，无GD.env或业务数据迁入。

主干规则main-protection（26003569）在代码前创建，禁止删除与强推；首次说明提交后启用PR合并，单人所需批准数0。CI形成后再启用check必需状态。最终状态与PR链接在交付段更新，不将“配置计划”记成生效。

课件来源：
- [第7课总览](https://devops.hello1023.com/课件/第7课-课件-小组项目启动/index.html)
- 同目录pages/repo.html、spec.html、rules.html、skeleton.html已逐页核对。
- std-product/backend/security/testing/ops/ai/data/review及case-files、case-check已核对。五类常驻规范分别落.agents/rules，由AGENTS指向；补充AI/数据/DoD集中于AGENTS，不复制课件全文。

## 能力与Scenario（两列表）

| 能力 | Scenario及对应证据 |
| --- | --- |
| F1 独立可复现配置 | S01真实Compose与干净clone；S02逐键缺失/非法，config_test；S03 99/100/10000/10001边界。 |
| F2 公开最小面/默认拒绝 | S04未知参数400；S05未登记教师API404；S06未知API/附件404；S07 POST health405；handler_test与acceptance。 |
| F3 存活/就绪 | S08匿名health返回UTC时间/版本/MySQL真实up；S09停Plus DB后health503/live200，恢复后200。 |
| F4 日志/停机 | S10响应和日志verify-lesson07关联；S11非法/1/64/65长度ID；S12在途请求排空单测与容器退出。 |
| F5 四态页面 | S13真实Chrome正常页；S14前端正常/失败/边界单测，Chrome注入503检查错误/重试，不将其称真实DB故障截图。 |
| F6 规范与证据 | S15一键check/acceptance、S16分层违规/默认放行两次变异失败；空白会话、CI与例外单列。 |

无权限类失败在骨架中验证“未登记业务资源不可访问”，不冒充业务登录或跨班授权；未来业务change须有401/403/跨班404逐资源场景。四类失败信号为参数400、默认拒绝404、资源不存在404、依赖503。

## 已执行的验证

- OpenSpec 1.13.0 init --tools codex --language zh完成；Codex生成6个Skill，未提供课件同名slash命令。strict使用真实CLI；不改全局profile。项目verify-change经官方quick_validate.py校验通过（临时PyYAML环境）。
- make check：规则/全历史秘密、Go格式/race/vet、前端3条状态测试及类型/构建、OpenSpec strict。最终Go1.27.3补丁镜像已通过，前端3条状态测试、Go race/vet与strict均通过。
- make acceptance：实际Compose/MySQL/HTTP、未知入口、日志、停机和Chrome；最终补丁版及空白会话反馈修复后均已通过：9条GET探测、POST拒绝、JSON日志和浏览器；最新正常health约18ms，DB不可达约1500ms，容器正常退出约0.11秒。
- make audit：初始发现Vite高危与Go1.27.1标准库漏洞；已锁定Vite7.3.7、Go1.27.3及govulncheck1.9.0，最终npm全部依赖0漏洞，Go扫描No vulnerabilities found；初始失败与修复分开保留。
- 规则生效：service暂加database/sql import，静态检查失败；暂把未知API改200/OK，负向HTTP测试失败；已恢复源码。脱敏记录reports/mutation.json。
- 原始本机报告与真实Chrome图片在reports/（忽略Git）；正式可提交证据仅含脱敏状态/统计，不含.env或完整Docker环境。

## 三角色记录与单人差异

| 工程职责 | 本人产出/验收 |
| --- | --- |
| 产品负责人 | 定位/边界、F1-F6及S01-S16、明确不做本课业务；学期目标保留作业、错题、学情。 |
| 开发者 | 独立分层骨架、真实MySQL依赖、统一响应、默认拒绝、JSON日志与Compose。 |
| 测试验收者 | 正常/失败/边界、真实停库恢复、变异检查、干净clone、独立上下文检查与证据。 |

- 同一人三角色不等同三名成员。独立同伴评审/另一台机器由其他成员从零启动，本人无法代替；本项目按用户明确要求作单人适配，课程是否接受由教师判断，原要求不标通过。
- 本机CLI空白会话因账号不支持所配置模型失败；随后在Codex应用新建无历史只读会话（01a18017-b55d-7612-acf1-1a4a2ba1415c），输入仅提供Plus路径、change名与只读验收任务，不带本聊天。
- 该会话确认规则可发现、分层/命名/默认拒绝符合；初评同时指出记录文件/提交尚未就绪、根密码规则缺校验、停机预算偏宽/缺超预算用例、POST误放行漏判和Nginx字段未全查。记录/提交在会话检查后补齐；代码四项缺口已修复并重跑check/acceptance，不将初评写成无问题通过。
- 根密码已用DB入口验证≥24字符且不同于普通密码，Python正常/失败测试；生命周期增加超预算强制结束单测；真实停机断言按配置秒数；POST意外成功必失败；Nginx字段及秘密均检查。
- 新Harness静态检查不等于另一个人运行；未执行第二次独立复核，不虚构“二次review全部通过”。
- 干净clone已通过：受验提交b2d448b，脚本新建临时目录、随机凭证、独立Compose项目campusclaw-plus-repro-*与端口3002，启动约14.61秒，匿名health 200/OK与MySQL up；随后仅清理其新建卷与目录。这是同机干净clone，不是异机。
- 首个发布无上一个运行版本，不假称生产双版本回滚。README提供当前基线新目录重启；骨架无业务表，业务备份恢复N/A，未来上线前必须补。
- 骨架多文件基础设施与测试不属于一次生成整套业务；不引入GD业务、模型、Qdrant或MCP。

## 本次发现与修复

1. npm审计发现Vite7.2.0高危项，升级并锁定7.3.7。
2. 旧govulncheck与Go1.27不兼容，升级1.9.0；进一步发现标准库可达漏洞，修复工具链到1.27.3。
3. 真实断连出现驱动非JSON日志，原因是FormatDSN丢失Logger，改NewConnector/OpenDB后真实故障日志验证通过。
4. 并发npm ci共享node_modules导致验收失败，check/audit改独立卷，本机浏览器依赖分开，重跑真实门禁。

## 最终交付状态

- 实现提交：39a5b46；空白会话反馈修复提交：b2d448b。首个规范提交与首次说明提交见上文。
- PR：https://github.com/Hilarious12138/campusclaw-plus/pull/1 。39a5b46与b2d448b的GitHub Actions均实际成功；修复版记录为https://github.com/Hilarious12138/campusclaw-plus/actions/runs/40615067678 。归档后提交继续执行同一门禁。
- 生效主干规则：main-protection，禁止删除/强推、PR合并、必需GitHub Actions check（integration_id=15368），单人批准数0。
- 骨架归档目标：openspec/changes/archive/2026-10-28-establish-platform-foundation/；长期规约：openspec/specs/platform-foundation/spec.md。归档后strict再验。
- 标签目标：v0.0.1-skeleton，主干CI通过并合并后创建，不能将目标写成已存在。最终标签存在性另由Git远端核对。
- 课程原始“独立peer review/另一台机器复现”仍不标完成；单人适用任务验收不代表教师已批准例外。

## 依赖用途与许可（本次直接依赖）

| 依赖 | 用途 | 许可审查 |
| --- | --- | --- |
| React18/react-dom | 工程状态页 | MIT，保留上游声明 |
| Vite7.3.7/TypeScript及类型包 | 构建与静态类型 | Vite MIT、TypeScript Apache-2.0、类型MIT |
| Playwright1.63.0 | 本机Chrome验收 | Apache-2.0，非业务依赖 |
| go-sql-driver/mysql1.10.1 | MySQL适配器 | MPL-2.0，未修改其源文件 |
| edwards25519 1.2.0 | MySQL驱动认证依赖 | BSD-3-Clause |
| MySQL/Go/Nginx/Node上游镜像 | 本课部署构建 | 固定摘要，遵循各镜像上游声明；MySQL GPL-2.0运行组件已显式记录 |

本仓未选新LICENSE，不代表可以移除依赖许可或将全部上游代码声明为自有。

## 可复验命令与限制

- 规范生效：`python3 scripts/mutation.py`，临时文件/代码在finally还原，返回失败即不可称门禁有效。
- 干净复现：先提交代码，再`python3 scripts/clean_clone.py`；本机3002需空闲，脚本用全新随机凭证/隔离Compose项目，完毕仅删除脚本新建卷与目录。
- 空白会话：新建无历史Harness，给仓库路径与日期变更名，要求仅凭仓库核对规则可发现性、任务可执行性与当前实现，禁止修改/停机；记录缺口并修复，不预填通过。

脱敏机器证据：`evidence.json`。实际Chrome截图：`homepage-ready.png`、`homepage-error.png`；错误页使用503注入，仅验证UI，真实停库结果单列在机器证据。
