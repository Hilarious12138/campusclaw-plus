import json,os,subprocess,urllib.request,urllib.error
from pathlib import Path
class API:
 def __init__(self):
  self.base=os.environ.get('BASE_URL','http://127.0.0.1:3011');self.passwords={name:os.environ['SEED_'+name.upper()+'_PASSWORD'] for name in ['teacher_a','student_a','teacher_b','student_b']};self.password=self.passwords['teacher_a'];self.tokens={};self.trace=[]
 def request(self,method,path,body=None,token=None,status=200,headers=None):
  h={'Content-Type':'application/json'}
  if token:h['Authorization']='Bearer '+token
  h.update(headers or {})
  raw=json.dumps(body).encode() if body is not None else None
  req=urllib.request.Request(self.base+path,data=raw,headers=h,method=method)
  try:r=urllib.request.urlopen(req,timeout=10)
  except urllib.error.HTTPError as e:r=e
  payload=r.read();
  try:preview=json.loads(payload)
  except Exception:preview={}
  self.trace.append({'method':method,'path':path.split('?')[0],'expected_http':status,'actual_http':r.status,'code':preview.get('code'),'request_id':preview.get('request_id'),'has_data':'data' in preview,'actor':next((name for name,value in self.tokens.items() if value==token),'anonymous')})
  assert r.status==status,(method,path,r.status,status)
  out=json.loads(payload);expected_code={200:'OK',202:'ACCEPTED',400:'INVALID_ARGUMENT',401:'UNAUTHENTICATED',403:'FORBIDDEN',404:'NOT_FOUND',409:'CONFLICT',429:'RATE_LIMITED',503:'DEPENDENCY_UNAVAILABLE'}.get(status);assert out['code']==expected_code,(path,out['code'],expected_code);assert out['request_id'];assert r.headers['X-Request-ID']==out['request_id']
  if status>=400:assert 'data' not in out
  assert all(v.encode() not in payload for v in self.passwords.values())
  return out
 def login(self,name):
  if name not in self.tokens:self.tokens[name]=self.request('POST','/api/auth/login',{'username':name,'password':self.passwords[name]})['data']['access_token']
  return self.tokens[name]
 def sql(self,query):
  project=os.environ.get('APITEST_DOCKER_PROJECT')
  if project:
   cmd=[os.environ.get('DOCKER_BIN','docker'),'exec','-i',project+'-db-1','sh','-c','MYSQL_PWD="$MYSQL_PASSWORD" mysql -N -B --raw -u "$MYSQL_USER" "$MYSQL_DATABASE"']
   p=subprocess.run(cmd,input=query,text=True,capture_output=True,timeout=20)
  else:
   env=os.environ.copy();env['MYSQL_PWD']=env['DB_PASSWORD']
   p=subprocess.run(['mysql','--skip-ssl','-N','-B','--raw','-h',env.get('DB_HOST','db'),'-u',env['DB_USER'],env['DB_NAME']],input=query,text=True,capture_output=True,env=env,timeout=20)
  assert p.returncode==0,'test SQL execution failed'
  return p.stdout.strip()
 def upload(self,title,text,token,key,status=202):
  import uuid
  boundary='----'+uuid.uuid4().hex
  raw=('--'+boundary+'\r\nContent-Disposition: form-data; name="file"; filename="'+title+'"\r\nContent-Type: text/plain\r\n\r\n'+text+'\r\n--'+boundary+'--\r\n').encode()
  req=urllib.request.Request(self.base+'/api/materials',data=raw,headers={'Content-Type':'multipart/form-data; boundary='+boundary,'Authorization':'Bearer '+token,'Idempotency-Key':key},method='POST')
  try:r=urllib.request.urlopen(req,timeout=10)
  except urllib.error.HTTPError as e:r=e
  out=json.loads(r.read());self.trace.append({'method':'POST','path':'/api/materials','expected_http':status,'actual_http':r.status,'code':out.get('code'),'request_id':out.get('request_id'),'has_data':'data' in out,'actor':next((name for name,value in self.tokens.items() if value==token),'anonymous')});assert r.status==status,(r.status,status,out.get('code'))
  expected_code={202:'ACCEPTED',400:'INVALID_ARGUMENT',401:'UNAUTHENTICATED',403:'FORBIDDEN',409:'IDEMPOTENCY_CONFLICT',429:'TASK_CAPACITY_EXCEEDED',503:'DEPENDENCY_UNAVAILABLE'}.get(status);assert out['code']==expected_code,(out['code'],expected_code)
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
