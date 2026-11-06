import json,urllib.request
def run(a):
 with urllib.request.urlopen(a.base+'/openapi.json') as r:spec=json.load(r)
 count=0
 for path,methods in spec['paths'].items():
  for method,op in methods.items():
   if op.get('security'):
    a.request(method.upper(),path.replace('{id}','0'*32),body={} if method.upper() in ['POST','PUT'] else None,status=401);count+=1
 assert count>=25
 with urllib.request.urlopen(a.base+'/swagger/') as r:page=r.read().decode();assert 'SwaggerUIBundle' in page and 'persistAuthorization:false' in page and 'cdn.' not in page
