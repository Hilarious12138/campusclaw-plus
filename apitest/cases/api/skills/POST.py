import json
def run(a):
 t=a.login('teacher_a');s=a.login('student_a');b=a.login('teacher_b')
 before=a.sql("SELECT COUNT(*) FROM resources WHERE kind='skills';")
 body={'name':'<script>alert(1)</script>','instructions':'只引导下一步','conditions':['方程'],'tools':[],'enabled':True}
 a.request('POST','/api/skills',body,token=s,status=403)
 assert a.sql("SELECT COUNT(*) FROM resources WHERE kind='skills';")==before
 a.request('POST','/api/skills',body,status=401)
 a.request('POST','/api/skills',body,token=t,headers={'Origin':'https://evil.invalid'},status=403)
 v=a.request('POST','/api/skills',body,token=t)['data'];assert v['revision']==1
 a.request('GET','/api/skills/'+v['id'],token=b,status=404)
 a.request('GET','/api/skills/'+'0'*32,token=b,status=404)
 got=a.request('GET','/api/skills/'+v['id'],token=t)['data'];assert got['body']==body
 a.request('PUT','/api/skills/'+v['id'],body,token=t,headers={'If-Match':'99'},status=409)
 # In isolated test DB, remove audit target temporarily to prove transaction rollback.
 a.sql('RENAME TABLE audit_events TO audit_events_test_backup;')
 try:
  count=a.sql('SELECT COUNT(*) FROM resources;')
  a.request('POST','/api/skills',body,token=t,status=503)
  assert a.sql('SELECT COUNT(*) FROM resources;')==count
 finally:a.sql('RENAME TABLE audit_events_test_backup TO audit_events;')
 a.request('DELETE','/api/skills/'+v['id'],token=t)
