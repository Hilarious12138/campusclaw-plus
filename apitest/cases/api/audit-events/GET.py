def run(a):
 t=a.login('teacher_a');s=a.login('student_a');b=a.login('teacher_b')
 uid=a.request('GET','/api/me',token=s)['data']['id']
 a.request('POST','/api/skills',{'name':'denied'},token=s,status=403)
 rows=a.request('GET','/api/audit-events?actor_id='+uid+'&status=403',token=t)['data'];assert rows
 for row in rows:
  assert row['class_id']=='a'*32 and row['actor_id']==uid and row['http_status']==403
  for k in ['occurred_at','operation','object_type','object_id','outcome','skill_id','usage','duration_ms']:assert k in row
 a.request('GET','/api/audit-events',token=s,status=403)
 assert a.request('GET','/api/audit-events?actor_id='+uid,token=b)['data']==[]
 a.request('GET','/api/audit-events?class_id=other',token=t,status=400)
