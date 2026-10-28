#!/usr/bin/env python3
"""Same-machine clean-clone reproducibility; independent disposable resources only."""
import json,os,secrets,shutil,subprocess,tempfile,time,urllib.request
from pathlib import Path
root=Path(__file__).resolve().parents[1];docker=os.environ.get('DOCKER_BIN','docker')
sha=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip()
temp=Path(tempfile.mkdtemp(prefix='campusclaw-plus-clone-'));clone=temp/'repo';project='campusclaw-plus-repro-'+secrets.token_hex(3);port=os.environ.get('PLUS_REPRO_PORT','3002')
cmd=[docker,'compose','-p',project];started=time.monotonic()
try:
 subprocess.run(['git','clone','--no-local',str(root),str(clone)],check=True,capture_output=True)
 subprocess.run(['git','checkout','--detach',sha],cwd=clone,check=True,capture_output=True)
 values={'WEB_PORT':port,'APP_ENV':'test','APP_VERSION':'v0.0.1-skeleton','HTTP_ADDR':':8080','DB_ADDR':'db:3306','DB_NAME':'campusclaw_plus','DB_USER':'campusclaw_plus','DB_PASSWORD':secrets.token_hex(24),'DB_ROOT_PASSWORD':secrets.token_hex(24),'DB_TIMEOUT_MS':'1500','SHUTDOWN_TIMEOUT_SECONDS':'10'}
 sample={x.split('=')[0] for x in (clone/'.env.example').read_text().splitlines() if x};assert sample==set(values)
 env=clone/'.env';env.write_text(''.join(k+'='+v+'\n' for k,v in values.items()));env.chmod(0o600)
 subprocess.run(cmd+['up','--build','-d','--wait','--wait-timeout','180'],cwd=clone,check=True,capture_output=True)
 with urllib.request.urlopen('http://127.0.0.1:'+port+'/health',timeout=5) as r:v=json.load(r);assert r.status==200 and v['code']=='OK' and v['data']['dependencies']['mysql']=='up'
 result={'status':'passed','kind':'same-machine clean clone; not independent peer review','source_commit':sha,'project':project,'port':int(port),'startup_seconds':round(time.monotonic()-started,2),'runtime_command':'docker compose -p <isolated-project> up --build -d --wait','private_values':'not recorded'}
 (root/'reports').mkdir(exist_ok=True);(root/'reports/clean-clone.json').write_text(json.dumps(result,indent=2));print(json.dumps(result))
finally:
 if clone.exists():subprocess.run(cmd+['down','-v','--remove-orphans'],cwd=clone,capture_output=True)
 shutil.rmtree(temp,ignore_errors=True)
