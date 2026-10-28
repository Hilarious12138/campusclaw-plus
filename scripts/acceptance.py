#!/usr/bin/env python3
import json,os,subprocess,time,urllib.request,urllib.error
from pathlib import Path
ROOT=Path(__file__).resolve().parents[1];os.chdir(ROOT)
DOCKER=os.environ.get('DOCKER_BIN','docker');PROJECT=os.environ.get('PLUS_PROJECT','campusclaw-plus')
if not PROJECT.startswith('campusclaw-plus'):raise SystemExit('Only Plus projects permitted')
cmd=[DOCKER,'compose','-p',PROJECT]
values=dict(line.split('=',1) for line in (ROOT/'.env').read_text().splitlines() if '=' in line)
base='http://127.0.0.1:'+values['WEB_PORT'];records=[]
(ROOT/'reports').mkdir(exist_ok=True)
def run(*args):return subprocess.run(cmd+list(args),check=True,capture_output=True,text=True).stdout
def get(path,status,code=None,request_id=None):
 req=urllib.request.Request(base+path,headers={'X-Request-ID':request_id} if request_id else {})
 start=time.monotonic()
 try:r=urllib.request.urlopen(req,timeout=12)
 except urllib.error.HTTPError as e:r=e
 with r:
  body=r.read();actual=r.status;rid=r.headers.get('X-Request-ID')
 elapsed=(time.monotonic()-start)*1000
 assert actual==status,(path,actual)
 v=None
 if code:
  v=json.loads(body);assert v['code']==code,(path,v['code']);assert v['request_id']==rid and rid
 records.append({'path':path,'http':actual,'code':code,'elapsed_ms':round(elapsed,2)})
 return v,elapsed
print('Starting isolated Plus stack',flush=True)
run('up','--build','-d','--wait','--wait-timeout','180')
try:
 v,elapsed=get('/health',200,'OK','verify-lesson07');assert elapsed<=2000
 assert v['data']['dependencies']['mysql']=='up' and v['data']['version']==values['APP_VERSION']
 assert v['data']['time'].endswith('Z')
 get('/health?class_id=other',400,'INVALID_ARGUMENT');get('/api/teachers/materials',404,'NOT_FOUND');get('/api/missing',404,'NOT_FOUND');get('/uploads/guessed.txt',404);get('/missing-page',404)
 req=urllib.request.Request(base+'/health',method='POST')
 try:urllib.request.urlopen(req)
 except urllib.error.HTTPError as e:assert e.code==405 and json.load(e)['code']=='METHOD_NOT_ALLOWED'
 else:raise AssertionError('POST /health was incorrectly allowed')
 print('Normal and rejected HTTP paths passed; stopping Plus DB',flush=True)
 run('stop','db')
 try:
  v,elapsed=get('/health',503,'DEPENDENCY_UNAVAILABLE');assert elapsed<=int(values['DB_TIMEOUT_MS'])+1000;assert v['data']['dependencies']['mysql']=='down';get('/live',200,'OK')
 finally:run('start','db')
 for i in range(45):
  try:get('/health',200,'OK');break
  except (AssertionError,urllib.error.URLError):time.sleep(1)
 else:raise AssertionError('DB recovery timeout')
 logs=run('logs','--no-color','--no-log-prefix','api')
 for secret_key in ['DB_PASSWORD','DB_ROOT_PASSWORD']:
  assert values[secret_key] not in logs,'secret in logs'
 parsed=[json.loads(line) for line in logs.splitlines() if line.strip()]
 request_logs=[v for v in parsed if v.get('event')=='request']
 assert any(v.get('request_id')=='verify-lesson07' for v in request_logs)
 for v in request_logs:
  assert all(k in v for k in ['time','level','module','request_id','method','path','status','duration_ms'])
  assert '?' not in v['path']
 # Only nginx access records are JSON; native stderr diagnostics, if any, are not represented as request records.
 nginx=run('logs','--no-color','--no-log-prefix','web')
 nlogs=[json.loads(line) for line in nginx.splitlines() if line.startswith('{')]
 assert any(v.get('request_id')=='verify-lesson07' for v in nlogs)
 for v in nlogs:
  assert all(k in v for k in ['time','level','module','event','request_id','method','path','status','duration_seconds'])
  assert '?' not in v['path']
 for k in ['DB_PASSWORD','DB_ROOT_PASSWORD']:assert values[k] not in nginx,'secret in nginx logs'
 api_id=run('ps','-q','api').strip();db_id=run('ps','-q','db').strip()
 for cid in [api_id,db_id]:
  ports=json.loads(subprocess.check_output([DOCKER,'inspect','--format','{{json .HostConfig.PortBindings}}',cid],text=True));assert not ports,'private port exposed'
 print('Recovery, JSON logs and private ports passed; checking shutdown',flush=True)
 start=time.monotonic();run('stop','api');duration=time.monotonic()-start;assert duration<=int(values["SHUTDOWN_TIMEOUT_SECONDS"])
 state=json.loads(subprocess.check_output([DOCKER,'inspect','--format','{{json .State}}',api_id],text=True));assert state['ExitCode']==0,'ungraceful stop'
 run('up','-d','--wait','--wait-timeout','60','api','web')
 env=dict(os.environ,PLUS_BASE_URL=base)
 subprocess.run([os.environ.get('NODE_BIN','node'),'scripts/browser.mjs'],env=env,check=True)
 result={'status':'passed','project':PROJECT,'checks':records,'json_request_logs':len(request_logs),'shutdown_seconds':round(duration,2),'browser':'reports/browser.json','scope':'real Compose/MySQL/HTTP; browser 503 is injected UI case; no teaching business implemented'}
 (ROOT/'reports/acceptance.json').write_text(json.dumps(result,ensure_ascii=False,indent=2))
 print(json.dumps({'status':'passed','checks':len(records),'shutdown_seconds':round(duration,2)},ensure_ascii=False))
finally:
 # Recover on failure as well; leave the user's Plus environment up, never touch GD.
 subprocess.run(cmd+['start','db'],capture_output=True)
 subprocess.run(cmd+['start','api','web'],capture_output=True)
