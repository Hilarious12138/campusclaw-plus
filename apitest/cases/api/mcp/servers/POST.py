import json,uuid,urllib.request
def run(a):
 t=a.login('teacher_a');s=a.login('student_a');b=a.login('teacher_b')
 a.request('POST','/api/mcp/servers',{'name':'denied','base_url':'http://mcp-demo:8090'},token=s,status=403)
 a.request('POST','/api/mcp/servers',{'name':'bad','base_url':'http://169.254.169.254'},token=t,status=400)
 server=a.request('POST','/api/mcp/servers',{'name':'tools-'+uuid.uuid4().hex[:8],'base_url':'http://mcp-demo:8090'},token=t)['data']
 assert 'token' not in json.dumps(server)
 a.request('POST','/api/mcp/servers/'+server['id']+'/connect',{},token=b,status=404)
 tools=a.request('POST','/api/mcp/servers/'+server['id']+'/connect',{},token=t)['data'];assert {v['name'] for v in tools}=={'formula_hint','unit_convert'}
 ref='formula_hint@'+server['name']
 skill=a.request('POST','/api/skills',{'name':'工具验证','conditions':['正方形'],'instructions':'先确认边长和单位，不直接给最终答案','tools':[ref],'enabled':True},token=t)['data']
 agent=a.request('POST','/api/agents',{'name':'工具验证应用','prompt':'只引导下一步','mode':'guidance','skill_ids':[skill['id']],'material_ids':[],'enabled':True},token=t)['data']
 c=a.request('POST','/api/conversations',{'agent_id':agent['id']},token=s)['data']
 def ask():
  req=urllib.request.Request(a.base+'/api/conversations/'+c['id']+'/messages',data=json.dumps({'content':'正方形面积公式是什么？'}).encode(),headers={'Authorization':'Bearer '+s,'Content-Type':'application/json'},method='POST')
  with urllib.request.urlopen(req,timeout=70) as r:return r.read().decode()
 raw=ask();assert '"status":"ok"' in raw and ref in raw,raw
 a.request('POST','/api/mcp/servers/'+server['id']+'/disconnect',{},token=t)
 raw=ask();assert '"status":"denied"' in raw and ref in raw,raw
 a.request('PUT','/api/skills/'+skill['id'],skill['body'],token=t,headers={'If-Match':str(skill['revision'])},status=403)
 events=a.request('GET','/api/audit-events?operation=mcp.tools%2Fcall&limit=100',token=t)['data'];assert any(v['object_id']==ref and v['outcome']=='ok' for v in events) and any(v['object_id']==ref and v['outcome']=='denied' for v in events)
