import json,os,uuid
from pathlib import Path
from common import API
out=Path('/reports/broker-recovery.json');a=API();t=a.login('teacher_a')
if os.environ.get('BROKER_PHASE')=='down':
 health=a.request('GET','/health');assert health['data']['status']=='degraded' and health['data']['dependencies']['broker']=='down'
 key=uuid.uuid4().hex;v=a.upload('broker-'+key+'.md','方程队列恢复 '+key,t,key)['data'];assert a.request('GET','/api/tasks/'+v['task_id'],token=t)['data']['status']=='pending';out.write_text(json.dumps(v))
else:
 v=json.loads(out.read_text());done=a.wait_task(v['task_id'],t,90);assert done['status']=='succeeded';v['status']='passed';out.write_text(json.dumps(v))
print('broker recovery phase passed (real broker outage)')
