#!/usr/bin/env python3
import importlib.util,json,time,os,traceback,html
from pathlib import Path
from common import API

def run():
 a=API();rows=[]
 paths=list((Path(__file__).parent/'cases').rglob('*.py'))
 order=['auth','materials','retrieval','conversations','mcp','skills','audit-events']
 def priority(p):
  parts=p.parts
  return next((i for i,k in enumerate(order) if k in parts),99),str(p)
 for p in sorted(paths,key=priority):
  start=time.monotonic();row={'case':str(p.relative_to(Path(__file__).parent))}
  try:
   spec=importlib.util.spec_from_file_location('test_case',p);m=importlib.util.module_from_spec(spec);spec.loader.exec_module(m);m.run(a);row['status']='passed'
  except Exception as e:row.update(status='failed',detail=type(e).__name__+': '+str(e))
  row['elapsed_ms']=round((time.monotonic()-start)*1000);rows.append(row);print(row['case'],row['status'])
 out=Path(os.environ.get('APITEST_REPORT_DIR',str(Path(__file__).parents[1]/'reports/api')));out.mkdir(parents=True,exist_ok=True)
 raw=json.dumps(rows,ensure_ascii=False,indent=2)
 # Never include raw request/response bodies or tokens in reports.
 for v in [*a.passwords.values(),*a.tokens.values()]:raw=raw.replace(v,'[REDACTED]')
 (out/'report.json').write_text(raw);(out/'report.html').write_text('<!doctype html><meta charset="utf-8"><h1>CampusClaw Plus API verification</h1><pre>'+html.escape(raw)+'</pre>')
 return bool(not rows or any(r['status']!='passed' for r in rows))
if __name__=='__main__':raise SystemExit(run())
