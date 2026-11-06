import uuid
def run(a):
 t=a.login('teacher_a');s=a.login('student_a');b=a.login('student_b');key=uuid.uuid4().hex
 title='数学验证-'+key+'.md';text='# 数学\n解方程应保持等式两边相等，先确认未知数，再两边减去同一个数。'
 a.upload(title,text,s,key,status=403)
 a.upload('bad.pdf',text,t,key,status=400)
 a.upload('bad.md','x'*(2*1024*1024+1),t,key,status=400)
 v=a.upload(title,text,t,key)['data'];assert 'result' not in v
 same=a.upload(title,text,t,key)['data'];assert same['task_id']==v['task_id']
 a.upload(title,text+'changed',t,key,status=409)
 a.request('GET','/api/tasks/'+v['task_id'],token=b,status=404)
 a.request('GET','/api/materials/'+v['material_id'],token=b,status=404)
 a.request('GET','/api/materials/'+v['material_id']+'/file',token=b,status=404)
 done=a.wait_task(v['task_id'],t);assert done['status']=='succeeded',done.get('error_code')
 assert a.sql("SELECT COUNT(*) FROM task_results WHERE task_id='"+v['task_id']+"';")=='1'
 assert int(a.sql("SELECT COUNT(*) FROM knowledge_chunks WHERE material_id='"+v['material_id']+"';"))>=1
