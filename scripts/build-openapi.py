#!/usr/bin/env python3
import json
from pathlib import Path
S={'type':'string'}
def obj(props,required=()):return {'type':'object','properties':props,'required':list(required),'additionalProperties':False}
paths={}
def add(path,method,summary,schema=None,public=False,multipart=False,params=None,accepted=False):
 p={'summary':summary,'tags':[path.split('/')[2] if path.startswith('/api/') else 'system'],'security':[] if public else [{'BearerAuth':[]}], 'responses':{str(n):{'description':v,'content':{'application/json':{'schema':{'$ref':'#/components/schemas/Envelope'}}}} for n,v in [(202 if accepted else 200,'Accepted' if accepted else 'OK'),(400,'INVALID_ARGUMENT'),(401,'UNAUTHENTICATED'),(403,'FORBIDDEN'),(404,'NOT_FOUND'),(409,'CONFLICT'),(429,'RATE_LIMITED'),(503,'DEPENDENCY_UNAVAILABLE')]}}
 pp=[]
 for name in ['id']:
  if '{'+name+'}' in path:pp.append({'name':name,'in':'path','required':True,'schema':S})
 if params:pp+=params
 if pp:p['parameters']=pp
 if schema:p['requestBody']={'required':True,'content':{'multipart/form-data' if multipart else 'application/json':{'schema':schema}}}
 if path.endswith('/messages'):
  p['responses']['200']={'description':'SSE meta / delta / done；失败error，未出现done不算完成','content':{'text/event-stream':{'schema':S}}}
 paths.setdefault(path,{})[method.lower()]=p
add('/health','GET','依赖就绪',public=True);add('/live','GET','进程存活',public=True)
add('/api/auth/login','POST','账号密码登录',obj({'username':S,'password':{'type':'string','format':'password'}},['username','password']),public=True)
add('/api/auth/logout','POST','撤销当前访问令牌');add('/api/me','GET','当前用户与认证班级')
for kind in ['agents','skills','prompts']:
 fields={'name':S}
 if kind=='agents':fields.update(prompt=S,mode={'type':'string','enum':['knowledge','guidance']},skill_ids={'type':'array','items':S},material_ids={'type':'array','items':S},enabled={'type':'boolean'})
 if kind=='skills':fields.update(instructions=S,material_hint=S,conditions={'type':'array','items':S},tools={'type':'array','items':S,'maxItems':2},enabled={'type':'boolean'})
 if kind=='prompts':fields['content']=S
 add('/api/'+kind,'GET','本班列表');add('/api/'+kind,'POST','教师创建，原子审计',obj(fields,['name']))
 add('/api/'+kind+'/{id}','GET','本班详情');add('/api/'+kind+'/{id}','PUT','教师更新，If-Match为当前revision',obj(fields,['name']),params=[{'name':'If-Match','in':'header','required':True,'schema':{'type':'integer','minimum':1}}]);add('/api/'+kind+'/{id}','DELETE','教师删除')
file=obj({'file':{'type':'string','format':'binary'}},['file'])
add('/api/skills/import','POST','教师预览技能文本，确认后通过创建技能接口生效',file,multipart=True)
add('/api/materials','GET','本班材料列表',params=[{'name':x,'in':'query','schema':{'type':'integer','minimum':0 if x=='offset' else 1,'maximum':100000 if x=='offset' else 100}} for x in ['limit','offset']])
add('/api/materials','POST','教师接收txt/md≤2MiB并创建可靠任务',file,multipart=True,accepted=True,params=[{'name':'Idempotency-Key','in':'header','required':True,'schema':{'type':'string','minLength':1,'maxLength':128}}])
for path,summary in [('/api/materials/{id}','材料详情'),('/api/materials/{id}/file','鉴权读取材料正文'),('/api/tasks/{id}','只读任务状态'),('/api/retrieval/chunks/{id}','鉴权回查切片')]:add(path,'GET',summary)
add('/api/retrieval/search','POST','本班可溯源检索',obj({'query':{'type':'string','minLength':1,'maxLength':1000},'mode':{'type':'string','enum':['keyword','vector','hybrid'],'default':'hybrid'}},['query']))
add('/api/conversations','GET','本人会话');add('/api/conversations','POST','新建本人会话',obj({'agent_id':S},['agent_id']));add('/api/conversations/{id}','GET','本人会话和消息');add('/api/conversations/{id}/messages','POST','检索与技能/工具后流式回答',obj({'content':{'type':'string','minLength':2,'maxLength':2000}},['content']))
add('/api/mcp/servers','GET','教师查询本班MCP服务器');add('/api/mcp/servers','POST','教师登记白名单MCP地址；不接收密钥',obj({'name':S,'base_url':S},['name','base_url']))
for op in ['connect','disconnect']:add('/api/mcp/servers/{id}/'+op,'POST','教师'+('连接并发现工具' if op=='connect' else '断开并清缓存'))
add('/api/mcp/servers/{id}/tools','GET','教师查询已发现工具')
add('/api/audit-events','GET','教师查询本班操作审计',params=[{'name':k,'in':'query','schema':S} for k in ['actor_id','operation','status','from','until','limit','offset']])
spec={'openapi':'3.0.3','info':{'title':'CampusClaw Plus API','version':'0.1.0','description':'同源/api；Bearer仅内存，班级来自认证身份。跨班与不存在统一404。'},'servers':[{'url':'/'}],'paths':paths,'components':{'securitySchemes':{'BearerAuth':{'type':'http','scheme':'bearer','bearerFormat':'JWT'}},'schemas':{'Envelope':obj({'code':S,'request_id':S,'data':{},'error':obj({'message':S})},['code','request_id'])}}}
root=Path(__file__).resolve().parents[1]
(root/'frontend/public').mkdir(exist_ok=True);(root/'frontend/public/openapi.json').write_text(json.dumps(spec,ensure_ascii=False,indent=2)+'\n')
