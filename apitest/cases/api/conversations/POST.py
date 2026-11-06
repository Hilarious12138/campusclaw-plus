import json,urllib.request,uuid
def run(a):
 t=a.login('teacher_a');s=a.login('student_a');b=a.login('student_b')
 agents=a.request('GET','/api/agents',token=t)['data'];assert len(agents)>=3
 agent=next(v for v in agents if v['body']['name']=='数学助手')
 mats=a.request('GET','/api/materials',token=t)['data'];assert mats
 body=agent['body'];body['material_ids']=[mats[0]['id']]
 a.request('PUT','/api/agents/'+agent['id'],body,token=t,headers={'If-Match':str(agent['revision'])})
 c=a.request('POST','/api/conversations',{'agent_id':agent['id']},token=s)['data']
 a.request('GET','/api/conversations/'+c['id'],token=b,status=404)
 req=urllib.request.Request(a.base+'/api/conversations/'+c['id']+'/messages',data=json.dumps({'content':'解方程怎样保持等式两边相等？'}).encode(),headers={'Authorization':'Bearer '+s,'Content-Type':'application/json'},method='POST')
 with urllib.request.urlopen(req,timeout=70) as res:
  assert res.headers['Content-Type']=='text/event-stream';raw=res.read().decode();assert 'event: meta' in raw and 'event: delta' in raw and 'event: done' in raw,raw
 turns=a.request('GET','/api/conversations/'+c['id'],token=s)['data']['turns'];assert len(turns)==1 and turns[0]['status']=='succeeded';assert turns[0]['skill_name']=='数学解题引导' and turns[0]['citations']
 assert int(a.sql("SELECT JSON_LENGTH(snapshot) FROM conversation_turns WHERE id='"+turns[0]['id']+"';"))>=2

 snapshot=json.loads(a.sql("SELECT snapshot FROM conversation_turns WHERE id='"+turns[0]['id']+"';"))
 snapshot='\n'.join(v['content'] for v in snapshot)
 assert '数学解题引导' in snapshot and '只引导' in snapshot and '<evidence>' in snapshot
 assert '语文阅读引导' not in snapshot
