#!/usr/bin/env python3
"""Real MySQL/RabbitMQ orchestration phases. Embedding latency is an explicit fixture."""
import json,os,time,uuid,math,urllib.request
from pathlib import Path
from common import API
out=Path(os.environ.get('GRADE_RESULT_DIR','/reports'));out.mkdir(parents=True,exist_ok=True)
a=API();token=a.login('teacher_a');phase=os.environ.get('ASYNC_PHASE','prepare');delay=int(os.environ.get('ASYNC_DELAY_MS','10000'));file=out/('async-'+str(delay)+'.json')
if phase=='prepare':
 a.fixture(delay=delay);times=[];ids=[]
 for i in range(20):
  key=uuid.uuid4().hex;start=time.monotonic();v=a.upload('异步验证-'+key+'.md','方程等式验证 '+key,token,key)['data'];times.append((time.monotonic()-start)*1000);ids.append(v)
  assert a.request('GET','/api/tasks/'+v['task_id'],token=token)['data']['status']=='pending'
 p95=sorted(times)[math.ceil(len(times)*.95)-1];assert p95<1000 and p95<delay*.1
 file.write_text(json.dumps({'fixture':True,'delay_ms':delay,'samples_ms':times,'p95_ms':p95,'tasks':ids,'status':'pending-confirmed'},ensure_ascii=False,indent=2))
elif phase=='complete':
 v=json.loads(file.read_text());first=v['tasks'][0];done=a.wait_task(first['task_id'],token,90);assert done['status']=='succeeded',done.get('error_code')
 # Complete-task durability is observed before acknowledgement; shell can kill the worker here.
 v['status']='completed';v['attempts']=done['attempts'];v['deliveries_before_replay']=done['delivery_count'];v['first_task']=first['task_id'];v['worker_duration_ms']=a.sql("SELECT ROUND(TIMESTAMPDIFF(MICROSECOND,started_at,finished_at)/1000) FROM tasks WHERE id='"+first['task_id']+"';")
 v['result_count']=a.sql("SELECT COUNT(*) FROM task_results WHERE task_id='"+first['task_id']+"';");assert v['result_count']=='1'
 file.write_text(json.dumps(v,ensure_ascii=False,indent=2))
 # Remove only other synthetic perf tasks in isolated test DB, avoiding a 20*30s queue drain.
 for x in v['tasks'][1:]:
  tid=x['task_id'];mid=x['material_id'];a.sql("DELETE FROM task_requests WHERE task_id='"+tid+"';DELETE FROM task_outbox WHERE task_id='"+tid+"';DELETE FROM task_attempts WHERE task_id='"+tid+"';DELETE FROM tasks WHERE id='"+tid+"';DELETE FROM materials WHERE id='"+mid+"';")
 a.fixture()
elif phase=='replayed':
 v=json.loads(file.read_text());tid=v['first_task'];end=time.monotonic()+30
 while time.monotonic()<end:
  t=a.request('GET','/api/tasks/'+tid,token=token)['data']
  if t['delivery_count']>v['deliveries_before_replay']:break
  time.sleep(1)
 assert t['delivery_count']>v['deliveries_before_replay'] and t['attempts']==v['attempts'] and t['status']=='succeeded'
 assert a.sql("SELECT COUNT(*) FROM task_results WHERE task_id='"+tid+"';")=='1'
 v['duplicate_redelivery']=t['delivery_count'];v['status']='passed';file.write_text(json.dumps(v,ensure_ascii=False,indent=2))
else:raise SystemExit('unknown async phase')
print('async',phase,delay,'passed (real DB/broker; fixture embedding)')
