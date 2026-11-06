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
   p=subprocess.run(['mysql','-N','-B','-h',env.get('DB_HOST','db'),'-u',env['DB_USER'],env['DB_NAME']],input=query,text=True,capture_output=True,env=env,timeout=20)
  assert p.returncode==0,'test SQL execution failed'
  return p.stdout.strip()
