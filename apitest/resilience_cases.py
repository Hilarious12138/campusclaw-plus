"""Retry classification and real durable dead-letter publishing; no real AI service."""
import json,os,time,uuid
from pathlib import Path
from common import API
out=Path(os.environ.get('GRADE_RESULT_DIR','/reports'));a=API();t=a.login('teacher_a');results=[]
try:
 for mode,expected in [('fail',3),('permanent',1)]:
  a.fixture(mode=mode);key=uuid.uuid4().hex;v=a.upload('retry-'+key+'.md','方程重试边界 '+key,t,key)['data'];task=a.wait_task(v['task_id'],t,90)
  assert task['status']=='failed' and task['attempts']==expected,(task['status'],task['attempts'],task.get('error_code'))
  until=time.monotonic()+15
  while time.monotonic()<until:
   count=a.sql("SELECT COUNT(*) FROM task_outbox WHERE task_id='"+v['task_id']+"' AND kind='dead' AND published_at IS NOT NULL;")
   if count=='1':break
   time.sleep(1)
  assert count=='1','dead-letter not confirmed'
  assert a.sql("SELECT COUNT(*) FROM task_results WHERE task_id='"+v['task_id']+"';")=='0'
  results.append({'fixture_mode':mode,'attempts':task['attempts'],'status':task['status'],'error_code':task['error_code'],'dead_letter_confirmed':True})
finally:a.fixture()
(out/'retry-dead-letter.json').write_text(json.dumps({'status':'passed','cases':results},ensure_ascii=False,indent=2))
print('retry classification, bounded attempts and dead-letter confirmed with real broker/DB')
