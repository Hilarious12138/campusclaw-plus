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
- make check：规则/全历史秘密、Go格式/race/vet、前端3条状态测试及类型/构建、OpenSpec strict。最终补丁镜像结果待下方更新。
- make acceptance：实际Compose/MySQL/HTTP、未知入口、日志、停机和Chrome；补丁前完整通过一次，补丁版重跑结果待更新。
- make audit：初始发现Vite高危与Go1.27.1标准库漏洞；已锁定Vite7.3.7、Go1.27.3及govulncheck1.9.0，最终扫描结果待更新，不能把初始失败写成通过。
- 规则生效：service暂加database/sql import，静态检查失败；暂把未知API改200/OK，负向HTTP测试失败；已恢复源码。脱敏记录reports/mutation.json。
- 原始本机报告与真实Chrome图片在reports/（忽略Git）；正式可提交证据仅含脱敏状态/统计，不含.env或完整Docker环境。

## 三角色记录与单人差异

| 工程职责 | 本人产出/验收 |
| --- | --- |
| 产品负责人 | 定位/边界、F1-F6及S01-S16、明确不做本课业务；学期目标保留作业、错题、学情。 |
| 开发者 | 独立分层骨架、真实MySQL依赖、统一响应、默认拒绝、JSON日志与Compose。 |
| 测试验收者 | 正常/失败/边界、真实停库恢复、变异检查、干净clone、独立上下文检查与证据。 |

- 同一人三角色不等同三名成员。独立同伴评审/另一台机器由其他成员从零启动，本人无法代替；本项目按用户明确要求作单人适配，课程是否接受由教师判断，原要求不标通过。
- 本机CLI空白会话因账号不支持所配置模型失败，未声称通过；另在Codex应用中新建无历史只读会话（01a18017-b55d-7612-acf1-1a4a2ba1415c），只给路径与change进行检查，结果待补。
- 干净clone使用新目录、独立凭证/Compose项目/端口，不是异机；执行结果待补。
- 首个发布无上一个运行版本，不假称生产双版本回滚。README提供当前基线新目录重启；骨架无业务表，业务备份恢复N/A，未来上线前必须补。
- 骨架多文件基础设施与测试不属于一次生成整套业务；不引入GD业务、模型、Qdrant或MCP。

## 本次发现与修复

1. npm审计发现Vite7.2.0高危项，升级并锁定7.3.7。
2. 旧govulncheck与Go1.27不兼容，升级1.9.0；进一步发现标准库可达漏洞，修复工具链到1.27.3。
3. 真实断连出现驱动非JSON日志，原因是FormatDSN丢失Logger，改NewConnector/OpenDB后真实故障日志验证通过。
4. 并发npm ci共享node_modules导致验收失败，check/audit改独立卷，本机浏览器依赖分开，重跑真实门禁。

## 最终交付状态

完成后填入实际提交、PR/CI、归档路径、骨架标签及干净clone结果；未完成不提前勾选。
