## 1. 骨架（仅本组产生应用代码）

- [ ] 1.1 建立domain/service/repository/httpapi与cmd装配，verify：分层静态检查及依赖替身单测通过。
- [ ] 1.2 必填env校验、MySQL探测及health/live/default-deny，verify：S02-S09正常/失败/边界测试通过。
- [ ] 1.3 请求标识、单行JSON日志、优雅停机，verify：S10-S12测试及实际日志核对通过。
- [ ] 1.4 React工程首页与四态，verify：状态单测、tsc/Vite构建及实际浏览器页面通过。
- [ ] 1.5 Compose/Nginx/多阶段镜像与配置示例，verify：匿名health与独立卷、单入口、缺配置失败，GD未改。

## 2. 规范与规约验证

- [ ] 2.1 五类规则+AGENTS+verify-change技能与四份规划齐全，verify：项目技能校验和OpenSpec strict通过，规范提交早于应用实现。
- [ ] 2.2 实现一键check/acceptance/audit与CI，verify：正常/失败/边界按S01-S16映射，非零阻断，审计结果真实记录。
- [ ] 2.3 规则生效验证，verify：临时注入反向依赖及未知API放行使门禁失败，还原后通过。
- [ ] 2.4 空白上下文检查，verify：新Harness进程只给仓库路径和change名能定位规则、按任务验证；记录证据，不冒充独立同伴。

## 3. 交付验证与记录

- [ ] 3.1 真实Compose、DB故障恢复、停机与浏览器验证，verify：make acceptance通过，记录匿名health/日志与超时结果。
- [ ] 3.2 干净clone从README复现，verify：独立临时目录/Compose项目/端口运行health成功；明确“同机”不冒充异机。
- [ ] 3.3 仓库/PR/CI与提交历史检查，verify：首次提交仅README/.gitignore，保护规则先于代码、CI通过、全历史秘密扫描无原值。
- [ ] 3.4 docs/lesson-07/verification.md记录两列表、三角色与限制，verify：逐项真实结果，独立同伴和异机差异明确；适用门禁通过后归档与v0.0.1-skeleton，交付验证不循环计入本任务。
