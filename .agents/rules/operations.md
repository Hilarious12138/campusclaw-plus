# 部署运维规范

| 约定 | 检查方式 |
| --- | --- |
| 多阶段Dockerfile，npm ci锁依赖，镜像固定摘要；同一镜像跨local/test/production配置，只有本机Nginx端口发布。 | `make check`构建；docker compose ps与配置检查，无API/DB宿主映射。 |
| .env必填配置缺失即失败；Plus默认建议3001，独立Compose项目名campusclaw-plus与卷，不改GD服务/卷。 | 空env失败；真实clone用不同项目/端口启动。 |
| /live仅进程存活，/health为就绪，免认证、含version/time/必要依赖状态；DB失败503，探测不泄露凭证。 | 停Plus DB，live200/health503；恢复后health200。 |
| SIGTERM先停止接收并在10秒配置内排空，服务返回后关闭DB；不使用log.Fatal绕过defer。 | 单元在途请求测试+容器stop耗时/退出码。 |
| Go与Nginx每请求一行JSON：time/level/module/event/request_id/method/path/status/duration_ms；Go回X-Request-ID。启动/关闭日志用system标识。 | JSON解析全部请求日志；请求头与记录ID一致；无Authorization/body/query/秘密。 |
| 指标预留请求量/错误率/延迟/依赖状态；本课以日志统计，不部署监控栈。未来连续3次health503或5分钟错误率>5%由郭地响应。 | 本课手工故障事件可定位；正式告警留后续变更，不假称已部署。 |
| 发布：固定提交构建→必要迁移→新容器就绪→切流；回滚指定前一已验收标签，先确认迁移向后兼容。首版本无前一可运行标签，不造“已回滚”证据。 | README同标签重启/恢复演练；以后双版本演练耗时入docs。 |
| 业务上线后每日备份、保留7天，RPO24h/RTO30min；删除时核对MySQL/向量/对象/缓存/备份。骨架无业务表，数据恢复演练记N/A。 | 有数据后必须真实恢复，不能只备份命令存在就勾选。 |
