## 1. 基础设施骨架（先执行）

- [ ] 1.1 确认设计范围和前置业务的独立change；verify：认证/班级/上传未验收前不开放材料任务，材料与批改验收差异明确。
- [ ] 1.2 增加broker、dispatcher、worker入口与领域端口；verify：Compose仅Nginx对外，镜像摘要/依赖许可/配置校验和JSON日志符合原规范。
- [ ] 1.3 增量迁移tasks/outbox/attempts/results与唯一键；verify：空库/升级/重复迁移与回滚方案可复核，无事务外调。

## 2. 任务内核

- [ ] 2.1 实现事务发件箱与publisher confirm；verify：A2/S02—S03，断broker、确认丢失和写库失败均无丢任务/半写。
- [ ] 2.2 实现worker、租约/fence和恢复；verify：A3/S04—S05，stop/start、kill与旧worker提交结果均符合判据。
- [ ] 2.3 实现业务/结果唯一约束、重试及死信outbox；verify：A4—A5/S06—S09，两worker并发、重复、429与永久失败均核对数据库及DLQ。

## 3. 材料业务接入（前置能力通过后）

- [ ] 3.1 接入教师上传202、任务只读查询和逐资源授权；verify：A1/S01、A6/S10—S11，响应不含执行结果，401/403/404/409/429及副作用可判定。
- [ ] 3.2 接入解析/切片/嵌入与索引版本发布；verify：A7/S12，崩溃不暴露半版本，向量不含正文，禁止跨班回查。
- [ ] 3.3 实现材料页任务状态、轮询与OpenAPI；verify：A8/S14，终态/离页停止，四态与错误重查，Swagger可试调。

## 4. 验证与交付

- [ ] 4.1 按URL补apitest/cases并接入make check与Docker grade；verify：A8/S13，三判据真实DB/broker运行，延迟适配器与真实嵌入分开记录，20次/两组P95含原始值。
- [ ] 4.2 产出并发/事务与运行证据；verify：reports/risk-concurrency.md和risk-transaction.md含手段/对象/判据/结果，HTML/JSON缺证据失败，docs保留脱敏摘要。
- [ ] 4.3 更新八规范涉及条款与项目规划；verify：按design缺口表逐项回填，不建平行正文，独立评审/异机未做就保留未完成。
- [ ] 4.4 strict、check/grade、审核后归档并打async-202标签；verify：同提交证据、PR门禁与标签对应；未验证不勾选、不归档、不标业务完成。
