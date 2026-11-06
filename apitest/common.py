import json,os,subprocess,urllib.request,urllib.error
from pathlib import Path
class API:
 def __init__(self):
  self.base=os.environ.get('BASE_URL','http://127.0.0.1:3011');self.password=os.environ['SEED_PASSWORD'];self.tokens={}
 def request(self,method,path,body=None,token=None,status=200,headers=None):
  h={'Content-Type':'application/json'}
  if token:h['Authorization']='Bearer '+token
  h.update(headers or {})
  raw=json.dumps(body).encode() if body is not None else None
  req=urllib.request.Request(self.base+path,data=raw,headers=h,method=method)
  try:r=urllib.request.urlopen(req,timeout=10)
  except urllib.error.HTTPError as e:r=e
  payload=r.read();assert r.status==status,(method,path,r.status,status)
  out=json.loads(payload);assert out['request_id'];assert r.headers['X-Request-ID']==out['request_id']
  if status>=400:assert 'data' not in out
  assert self.password.encode() not in payload
  return out
 def login(self,name):
  if name not in self.tokens:self.tokens[name]=self.request('POST','/api/auth/login',{'username':name,'password':self.password})['data']['access_token']
  return self.tokens[name]
 def sql(self,query):
  project=os.environ.get('APITEST_DOCKER_PROJECT')
  if project:
   cmd=[os.environ.get('DOCKER_BIN','docker'),'exec','-i',project+'-db-1','sh','-c','MYSQL_PWD="$MYSQL_PASSWORD" mysql -N -B -u "$MYSQL_USER" "$MYSQL_DATABASE"']
   p=subprocess.run(cmd,input=query,text=True,capture_output=True,timeout=20)
  else:
   env=os.environ.copy();env['MYSQL_PWD']=env['DB_PASSWORD']
   p=subprocess.run(['mysql','--skip-ssl','-N','-B','-h',env.get('DB_HOST','db'),'-u',env['DB_USER'],env['DB_NAME']],input=query,text=True,capture_output=True,env=env,timeout=20)
  assert p.returncode==0,'test SQL execution failed'
  return p.stdout.strip()
 def upload(self,title,text,token,key,status=202):
  import uuid
  boundary='----'+uuid.uuid4().hex
  raw=('--'+boundary+'\r\nContent-Disposition: form-data; name="file"; filename="'+title+'"\r\nContent-Type: text/plain\r\n\r\n'+text+'\r\n--'+boundary+'--\r\n').encode()
  req=urllib.request.Request(self.base+'/api/materials',data=raw,headers={'Content-Type':'multipart/form-data; boundary='+boundary,'Authorization':'Bearer '+token,'Idempotency-Key':key},method='POST')
  try:r=urllib.request.urlopen(req,timeout=10)
  except urllib.error.HTTPError as e:r=e
  out=json.loads(r.read());assert r.status==status,(r.status,status,out.get('code'))
  if status>=400:assert 'data' not in out
  return out
 def wait_task(self,tid,token,timeout=90):
  import time
  end=time.monotonic()+timeout
  while time.monotonic()<end:
   t=self.request('GET','/api/tasks/'+tid,token=token)['data']
   if t['status'] in ['succeeded','failed']:return t
   time.sleep(1)
  raise AssertionError('task did not reach terminal state')
 def fixture(self,mode='ok',delay=0):
  req=urllib.request.Request(os.environ.get('FIXTURE_URL','http://127.0.0.1:3012')+'/control',data=json.dumps({'mode':mode,'delay':delay}).encode(),headers={'Content-Type':'application/json'},method='POST')
  with urllib.request.urlopen(req,timeout=5) as r:return json.load(r)
