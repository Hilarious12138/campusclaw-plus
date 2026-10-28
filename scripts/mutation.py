import subprocess, pathlib, os, json
root=pathlib.Path(__file__).resolve().parents[1];os.chdir(root)
p=root/'backend/internal/service/rule_probe.go'
try:
 p.write_text('package service\nimport _ "database/sql"\n')
 r=subprocess.run(['python3','scripts/check.py'],capture_output=True,text=True)
 assert r.returncode!=0 and 'layer IO violation' in r.stderr+r.stdout
 print('reverse dependency correctly blocked')
finally:p.unlink(missing_ok=True)
p=root/'backend/internal/httpapi/handler.go';old=p.read_text()
try:
 new=old.replace('reply(404, "NOT_FOUND",','reply(200, "OK",',1);assert new!=old;p.write_text(new)
 r=subprocess.run([os.environ.get('DOCKER_BIN','docker'),'run','--rm','-v',str(root/'backend')+':/work','-w','/work','-v','campusclaw-plus-go-cache:/go/pkg','-v','campusclaw-plus-go-build:/root/.cache/go-build','golang:1.27.3-alpine@sha256:d200d2dd76b03a99b52651fd1325aff7dcd35422c35c01bac56571b0507f3deb','go','test','./internal/httpapi','-run','TestFailurePathsS04S05S06S07','-count=1'],capture_output=True,text=True)
 assert r.returncode!=0 and 'FAIL' in r.stdout,r.stdout+r.stderr
 print('unregistered API allowance correctly rejected by negative tests')
finally:p.write_text(old)
(root/'reports/mutation.json').write_text(json.dumps({'reverse_dependency':'blocked','unregistered_api_success':'tests_failed_as_expected','source_restored':True},indent=2))
