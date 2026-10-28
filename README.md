# CampusClaw Plus

基于 CampusClaw-GD 经验独立建设的班级智能教学平台，串起“讲义检索 → 辅助学习 → 作业批改与错题订正”。

当前为第7课项目启动；首次提交只建立说明和忽略规则，下一变更交付可启动骨架。后续才逐步复用GD能力并拓展应用、模板、Skill、MCP、作业、错题与学情；当前不具备这些业务。

## 本机运行（骨架变更交付后）

依赖 Git、Docker Desktop（Compose v2）；安装时联网，启动后不调用模型。复制 `.env.example` 为本机 `.env`，按下表填值，再执行：

```sh
cp .env.example .env
# 在本机编辑 .env，填写所有键；不要提交或发送真实值。
docker compose up --build -d --wait
```

| 键 | 本机如何填写 |
| --- | --- |
| WEB_PORT | 未占用端口，建议3001；GD的3000不受影响 |
| APP_ENV | local，测试环境用test |
| APP_VERSION | 本课填v0.0.1-skeleton |
| HTTP_ADDR | 容器内监听地址，填:8080 |
| DB_ADDR | Compose内部MySQL地址，填db:3306 |
| DB_NAME | 独立数据库名，填campusclaw_plus |
| DB_USER | 独立普通账号，填campusclaw_plus |
| DB_PASSWORD | 自己生成至少24位字母数字随机值，不复用GD |
| DB_ROOT_PASSWORD | 另生成至少24位字母数字随机值，与普通账号不同 |
| DB_TIMEOUT_MS | 单次依赖探测上限，填1500 |
| SHUTDOWN_TIMEOUT_SECONDS | 优雅停机秒数，填10 |

浏览器访问 `http://127.0.0.1:3001/`；匿名访问 `http://127.0.0.1:3001/health`，正常HTTP 200、code=OK、dependencies.mysql=up；故障HTTP 503、code=DEPENDENCY_UNAVAILABLE。若WEB_PORT另填端口，URL同步替换。

停止：`docker compose down`。不要加 `-v`，该参数删除本项目数据卷。Plus的Compose项目、数据库、凭证与GD隔离。API与MySQL不发布宿主端口。

## 角色与范围

郭地单人承担产品负责人、开发者、测试验收者。产品业务角色仍为教师/学生。本人自审不冒充独立同伴评审；同机干净clone不冒充异机验证。

只有骨架首页 `/`、免认证 `/health` 与 `/live` 属于本课；未来按OpenSpec逐项交付教学能力。GD保持冻结，代码如被复用须标来源并在Plus重新验证，不复制其.env、数据库或Git历史。
