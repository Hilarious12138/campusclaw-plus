# 部署配置

依赖 Git、Docker Desktop（Compose v2）；安装时联网，启动后不调用模型。复制 `.env.example` 为本机 `.env`，按下表配置，再执行：

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


### 配置轮换与隔离

若本机凭证泄露，由维护者立即撤销/轮换并验证旧值失败，记录键名不记值。MySQL已有数据卷时修改.env不会自动改库内口令：先用管理员身份执行受控口令变更，再同步.env、重建API并检查health；保留业务数据，不用down -v当轮换。首次骨架没有GD数据，后续部署也不能复制GD凭证。

### 启停、回到基线与数据

`docker compose down`保留卷；恢复使用同一启动命令。回到已验收基线：在新的clone目录checkout `v0.0.1-skeleton`，填独立.env，使用不同WEB_PORT及 `docker compose -p campusclaw-plus-baseline up --build -d --wait`。首版本没有前一可运行发布，不宣称完成双版本生产回滚；业务迁移上线前再验回滚与备份恢复。

