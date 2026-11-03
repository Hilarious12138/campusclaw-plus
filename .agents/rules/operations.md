# 部署运维规范

适用：构建、部署、日志、监测、发布与恢复。当前无业务表，业务备份恢复和正式告警仍待后续实施。

## 构建与环境

**要求**：使用多阶段Dockerfile、`npm ci`和固定摘要的镜像。同一镜像用配置区分local/test/production，只有本机Nginx端口对外发布。

`.env`必填项缺失即失败。Plus建议端口3001，独立Compose项目名campusclaw-plus及数据卷，不操作GD服务或卷。

检查：`make check`构建；查看docker compose ps与配置，确认API/DB无宿主映射。空env失败；干净clone使用独立项目/端口。

## 健康与停机

**要求**：`/live`表示进程存活，`/health`表示依赖就绪。两者免认证；health含version/time及必要依赖状态，DB失败返回503，不泄露凭证。

收到SIGTERM后先停止接收请求，在配置的10秒内排空已有请求，再关闭DB。不用`log.Fatal`跳过defer清理。

检查：停测试DB时live200、health503；恢复后health200。用在途请求单测及容器stop验证耗时和退出码。

## 结构化日志

**要求**：Go与Nginx每请求一行JSON，字段包括time、level、module、event、request_id、method、path、status与耗时。

Go耗时用duration_ms；Nginx用duration_seconds记录原生request_time，单位不能混用。Go返回X-Request-ID；启动/关闭事件使用system标识。日志不记Authorization、body、query或秘密。

检查：逐行解析请求日志，核对字段、请求头与记录ID，扫描敏感内容。

## 监测与告警

**要求**：预留请求量、错误率、延迟和依赖状态指标。当前用日志统计，不部署监控栈。未来连续3次health503，或5分钟错误率超过5%，由维护者响应。

检查：现阶段故障事件可从日志定位。正式告警未部署，不标为已完成。

## 发布与回滚

**要求**：按固定提交构建，执行必要迁移，等待新容器就绪后切流。回滚到前一已验收标签，先确认迁移向后兼容。首版本无前一可运行标签，不虚构回滚结果。

检查：按`docs/deployment.md`重启/恢复同标签版本。后续双版本演练记录实际耗时。

## 备份与恢复

**要求**：业务上线后每日备份、保留7天。RPO（可接受的数据丢失时间）为24h，RTO（恢复耗时目标）为30min。删除时核对MySQL、向量、对象、缓存和备份副本。

检查：有业务数据后做真实恢复；只有备份命令不算完成。当前无业务表，恢复演练不适用。
