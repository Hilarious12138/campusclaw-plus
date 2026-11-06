#!/bin/sh
# Portable acceptance entry: POSIX shell + Docker Compose only. No host language runtimes/browser.
set -eu
cd "$(dirname "$0")/.."
DOCKER=${DOCKER_BIN:-docker}
stamp=$(date -u +%Y%m%dT%H%M%SZ)
PROJECT="campusclaw-plus-grade-$(date +%s)-$$"
export GRADE_IMAGE="$PROJECT-grader"
export GRADE_REPORT_DIR="$PWD/reports/acceptance-$stamp-$$"
mkdir -p "$GRADE_REPORT_DIR"
CONFIG_FILE=${GRADE_ENV_FILE:-.env}
failed=0
image_ready=0
resources=0
stage=initialization
dc() { "$DOCKER" compose --env-file "$CONFIG_FILE" -p "$PROJECT" -f compose.grade.yaml "$@"; }
phase() { dc run --rm --no-deps grader "$@" >>"$GRADE_REPORT_DIR/runner.log" 2>&1 || failed=1; }
finish() {
 result=$?
 trap - EXIT INT TERM
 set +e
 if [ "$result" -ne 0 ]; then failed=1; fi
 if [ "$resources" -eq 1 ]; then
 dc down -v --remove-orphans >>"$GRADE_REPORT_DIR/infra.log" 2>&1
 if [ "$?" -ne 0 ]; then failed=1; stage=cleanup; printf '\nCleanup failed; isolated project: %s\n' "$PROJECT" >>"$GRADE_REPORT_DIR/infra.log"; fi
 fi
 if [ "$image_ready" -eq 1 ]; then
   "$DOCKER" run --rm --network none --env-file "$CONFIG_FILE" -e HOST_FAILURE="$failed" -e HOST_FAILURE_REASON="$stage" -v "$GRADE_REPORT_DIR:/reports" -v "$PWD/.git:/work/.git:ro" "$GRADE_IMAGE" finalize >>"$GRADE_REPORT_DIR/runner.log" 2>&1
   if [ "$?" -ne 0 ]; then failed=1; fi
 fi
 if [ ! -s "$GRADE_REPORT_DIR/report.html" ]; then
   failed=1
   printf '<!doctype html><meta charset="utf-8"><title>CampusClaw Plus 验收失败</title><h1>基础设施检查失败</h1><p>阶段：%s。未执行的必测项不能记为通过。</p><p>请确认.env、Docker Compose与镜像网络，查看infra.log后重新运行。宿主无需安装Python、Node或Chrome。</p>' "$stage" >"$GRADE_REPORT_DIR/report.html"
   printf '{"status":"failed","stage":"%s","reason":"infrastructure failure; required cases not executed","exit_code":1,"checks":[' "$stage" >"$GRADE_REPORT_DIR/report.json"
   separator=''
   for id in C01 C02 C03 C04 C05 C06 C07 C08 B01 B02 B03 B04 A01 A02 A03 R01 R02 R03 R04 R05 R06 R07; do
     printf '%s{"id":"%s","status":"blocked","reason":"grading infrastructure unavailable"}' "$separator" "$id" >>"$GRADE_REPORT_DIR/report.json"
     separator=,
   done
   printf ']}\n' >>"$GRADE_REPORT_DIR/report.json"
   printf '<p>未执行的必测项（均为阻塞）：C01-C08、R01-R07，共22项，不计通过。</p>' >>"$GRADE_REPORT_DIR/report.html"
 fi
 printf '\n验收报告：%s/report.html\n机器报告：%s/report.json\n退出码：%s\n' "$GRADE_REPORT_DIR" "$GRADE_REPORT_DIR" "$failed"
 exit "$failed"
}
trap finish EXIT
trap 'exit 1' INT TERM
[ -f "$CONFIG_FILE" ] || { printf 'Missing .env; copy .env.example and fill documented keys.\n' >"$GRADE_REPORT_DIR/infra.log"; exit 1; }
stage=fixture-configuration
TEST_FIXTURE_API_KEY=$("$DOCKER" run --rm node:24-alpine@sha256:ebfe2f90462722a7a4de65e91990e97fe0d401c70e0e762c5b53302f905ec1c1 node -e 'process.stdout.write(require("node:crypto").randomBytes(32).toString("hex"))')
export TEST_FIXTURE_API_KEY
stage=build
printf 'Building isolated grading tools and application containers...\n'
dc build >>"$GRADE_REPORT_DIR/infra.log" 2>&1
image_ready=1
stage=checks
resources=1
phase checks
stage=startup
printf 'Starting isolated stack without host ports...\n'
dc up -d --wait --wait-timeout 180 db api web broker qdrant fixture mcp-demo dispatcher worker >>"$GRADE_REPORT_DIR/infra.log" 2>&1
stage=healthy
phase healthy
stage=business
phase business
stage=concurrency
dc stop worker dispatcher >>"$GRADE_REPORT_DIR/infra.log" 2>&1
phase concurrency
dc start worker dispatcher >>"$GRADE_REPORT_DIR/infra.log" 2>&1
stage=resilience
phase resilience
stage=broker-outage
dc stop broker >>"$GRADE_REPORT_DIR/infra.log" 2>&1
dc run --rm --no-deps -e BROKER_PHASE=down grader broker >>"$GRADE_REPORT_DIR/runner.log" 2>&1 || failed=1
dc up -d --wait --wait-timeout 90 broker >>"$GRADE_REPORT_DIR/infra.log" 2>&1
dc run --rm --no-deps -e BROKER_PHASE=recovered grader broker >>"$GRADE_REPORT_DIR/runner.log" 2>&1 || failed=1
stage=async-recovery
for delay in 10000 30000; do
 dc stop worker >>"$GRADE_REPORT_DIR/infra.log" 2>&1
 dc run --rm --no-deps -e ASYNC_PHASE=prepare -e ASYNC_DELAY_MS="$delay" grader async >>"$GRADE_REPORT_DIR/runner.log" 2>&1 || failed=1
 dc start worker >>"$GRADE_REPORT_DIR/infra.log" 2>&1
 dc run --rm --no-deps -e ASYNC_PHASE=complete -e ASYNC_DELAY_MS="$delay" grader async >>"$GRADE_REPORT_DIR/runner.log" 2>&1 || failed=1
 dc exec -T broker rabbitmqctl purge_queue campusclaw.work >>"$GRADE_REPORT_DIR/infra.log" 2>&1 || failed=1
 dc kill worker >>"$GRADE_REPORT_DIR/infra.log" 2>&1
 dc start worker >>"$GRADE_REPORT_DIR/infra.log" 2>&1
 dc run --rm --no-deps -e ASYNC_PHASE=replayed -e ASYNC_DELAY_MS="$delay" grader async >>"$GRADE_REPORT_DIR/runner.log" 2>&1 || failed=1
done
phase async-summary
stage=dependency-failure
dc stop db >>"$GRADE_REPORT_DIR/infra.log" 2>&1
phase down
dc up -d --wait --wait-timeout 180 db api web broker qdrant fixture mcp-demo dispatcher worker >>"$GRADE_REPORT_DIR/infra.log" 2>&1
phase recovered
stage=shutdown
api_id=$(dc ps -q api)
dc stop api >>"$GRADE_REPORT_DIR/infra.log" 2>&1
api_exit=$("$DOCKER" inspect --format '{{.State.ExitCode}}' "$api_id")
dc logs --no-color --no-log-prefix api >"$GRADE_REPORT_DIR/api.log" 2>&1
dc logs --no-color --no-log-prefix web >"$GRADE_REPORT_DIR/web.log" 2>&1
dc run --rm --no-deps -e API_EXIT_CODE="$api_exit" grader stopped >>"$GRADE_REPORT_DIR/runner.log" 2>&1 || failed=1
dc logs --no-color --no-log-prefix worker dispatcher broker >"$GRADE_REPORT_DIR/queue.log" 2>&1
stage=completed
