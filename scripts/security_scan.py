#!/usr/bin/env python3
"""Scan without printing matched material. Fixture root permits safe red/green tests."""
import argparse,json,re,subprocess
from pathlib import Path

def scan(root,history=False):
 patterns=[re.compile(rb'gh[pousr]_[A-Za-z0-9]{25,}'),re.compile(rb'sk-[A-Za-z0-9_-]{24,}'),re.compile(rb'-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----')]
 values=[]
 env=root/'.env'
 if env.exists():
  for line in env.read_text().splitlines():
   if '=' in line:
    k,v=line.split('=',1)
    if any(x in k for x in ('PASSWORD','SECRET','TOKEN','API_KEY')) and len(v)>=16:values.append(v.encode())
 def bad(b):return any(p.search(b) for p in patterns) or any(v in b for v in values)
 findings=0
 if (root/'.git').exists():
  names=subprocess.check_output(['git','ls-files','--cached','--others','--exclude-standard','-z'],cwd=root).decode().split('\0')
  files=[root/n for n in names if n]
 else:files=[p for p in root.rglob('*') if p.is_file() and p.name!='.env']
 files+=list((root/'frontend/dist').rglob('*'))
 for p in set(files):
  if p.is_file() and p.name!='.env' and bad(p.read_bytes()):findings+=1
 if history:
  for line in subprocess.check_output(['git','rev-list','--objects','--all'],cwd=root,text=True).splitlines():
   oid=line.split()[0]
   if subprocess.check_output(['git','cat-file','-t',oid],cwd=root,text=True).strip()=='blob' and bad(subprocess.check_output(['git','cat-file','blob',oid],cwd=root)):findings+=1
 return findings
if __name__=='__main__':
 p=argparse.ArgumentParser();p.add_argument('--root',type=Path,default=Path(__file__).resolve().parents[1]);p.add_argument('--history',action='store_true');a=p.parse_args()
 n=scan(a.root,a.history);print(json.dumps({'check':'secret-scan','findings':n,'status':'failed' if n else 'passed'}));raise SystemExit(bool(n))
