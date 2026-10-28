# CampusClaw Plus

基于 CampusClaw-GD 经验独立建设的班级智能教学平台，串起“讲义检索 → 辅助学习 → 作业批改与错题订正”。

当前交付第7课工程骨架；首次提交只建立说明和忽略规则，后续变更实现了分层、配置、健康检查和日志。后续才逐步复用GD能力并拓展应用、模板、Skill、MCP、作业、错题与学情；当前不具备这些业务。

## 本机运行

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

## 老师验收（Linux / Docker）

克隆本仓库，复制并按上表填好.env；只需要Docker Engine + Compose v2和系统sh，不需要宿主Python、Node、Go或Chrome。建议Docker分配4核、4GiB以上内存，首次构建需联网下载镜像/工具与漏洞公告。

```sh
sh scripts/grade.sh
```

可先按“本机运行”部署正式栈再验收，也可直接运行上述命令。验收会另建不发布端口的唯一项目与测试卷，真实执行停库恢复/停API，不碰正式Plus或GD。完成自动清理测试容器/卷，保留报告。macOS若docker不在PATH可在命令前指定DOCKER_BIN；Linux通常无需指定。

报告位置：`reports/teacher-<UTC时间>-<进程号>/report.html`，直接用浏览器打开；`report.json`供自动评分程序读取。HTML内含逐项判据、通过/失败/阻塞、耗时、错误详情与内嵌真实截图，无外部网页依赖。总体通过退出0；必测失败、缺阶段、基础设施或清理失败退出非零。镜像下载或公告源失败不会被写成“通过”，修复网络后重新运行即可。

当前15项覆盖第7课骨架：规则/全历史秘密、配置边界、Go race/vet、前端四态与构建、strict、漏洞、匿名health、默认拒绝、Linux Chromium、真实DB故障/恢复、JSON日志及正常退出。未交付教学业务和独立人类peer review列为不适用，不能据此宣称学期平台全部完成。后续每个业务change必须扩展老师验收用例及Scenario映射。

## 验证与开发

本机运行只需Git和Docker；开发/验收另需Python 3、Make、Node 24、OpenSpec 1.13.0；浏览器验证使用本机Chrome。Mac若docker不在PATH，命令前加 `DOCKER_BIN=/Applications/Docker.app/Contents/Resources/bin/docker`。没有本机Go也能验收。

```sh
npm install --global @fission-ai/openspec@1.13.0
make check
make audit
# 需要Chrome与Node；先安装前端测试依赖，Docker测试后已有node_modules时不用重复安装。
npm ci --prefix frontend
make acceptance
```

`make check`：规则/全历史秘密扫描、Go格式/race/vet、前端正常/失败/边界状态测试及构建、OpenSpec strict。`make audit`：npm全部依赖高危扫描、Go漏洞扫描；网络失败同样非零，不能算通过。`make acceptance`：真实Compose/HTTP/MySQL故障恢复、匿名health、日志/私有端口、优雅停机及Chrome响应式/错误重试；会短暂停止**Plus自己的**DB/API并恢复，不操作GD。报告和截图在忽略Git的 `reports/`；可提交的脱敏结论在 `docs/lesson-07/verification.md`。

前端开发：`cd frontend && npm ci && npm run dev`；Vite从项目本机.env读取WEB_PORT，仅将该非秘密端口用于同源代理，不将数据库配置注入前端；端口变更无需改代码。

贡献流程：先读AGENTS.md和相关规则，建立日期变更并strict，再按骨架任务组实施；调用项目 `verify-change` 技能验证。当前Codex用 `$openspec-propose` / `$openspec-apply-change` / `$openspec-archive-change`；课件 `/opsx-validate` 的实际校验使用 `openspec validate <change> --strict`。无需更改全局OpenSpec profile。

### 配置轮换与隔离

若本机凭证泄露，由郭地立即撤销/轮换并验证旧值失败，记录键名不记值。MySQL已有数据卷时修改.env不会自动改库内口令：先用管理员身份执行受控口令变更，再同步.env、重建API并检查health；保留业务数据，不用down -v当轮换。首次骨架没有GD数据，后续部署也不能复制GD凭证。

### 启停、回到基线与数据

`docker compose down`保留卷；恢复使用同一启动命令。回到已验收基线：在新的clone目录checkout `v0.0.1-skeleton`，填独立.env，使用不同WEB_PORT及 `docker compose -p campusclaw-plus-baseline up --build -d --wait`。首版本没有前一可运行发布，不宣称完成双版本生产回滚；业务迁移上线前再验回滚与备份恢复。

## 学期目标（尚未实现）

教师上传本班讲义，学生检索可追溯片段；教师配置AI应用、版本化提示词模板与材料/技能绑定，学生流式提问。Skill支持解题引导、作业批改、错题讲解，外部工具经授权MCP调用并审计。后续教师发布作业、学生截止前提交、教师复核/改分；扩展个人错题订正、本班学情反馈，最终完成HTTPS公网部署。

当前仅首页工程状态与健康接口可用。没有教学登录、业务表、Qdrant、模型或MCP调用；已有GD完成状态不等于Plus完成。后续按当前change之外的新变更逐项复用验证。
