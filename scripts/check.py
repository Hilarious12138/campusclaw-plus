#!/usr/bin/env python3
"""Repository invariants; intentionally runnable without external packages."""
import json,re,subprocess,sys,os
from pathlib import Path
ROOT=Path(__file__).resolve().parents[1]
def fail(s):raise SystemExit(s)
# Dependency direction: application imports allowed by layer, standard IO boundaries included.
allowed={'domain':set(),'service':{'domain'},'repository':{'domain'},'httpapi':{'service','domain'},'lifecycle':set(),'config':set()}
for p in (ROOT/'backend/internal').rglob('*.go'):
 if p.name.endswith('_test.go'):continue
 layer=p.relative_to(ROOT/'backend/internal').parts[0];s=p.read_text()
 for imp in re.findall(r'"([^"\n]+)"',s):
  if imp.startswith('campusclaw-plus/internal/') and imp.split('/')[2] not in allowed[layer]:fail(f'reverse dependency: {p.relative_to(ROOT)} -> {imp}')
  if layer in ('domain','service') and imp in ('net/http','database/sql'):fail(f'layer IO violation: {p.relative_to(ROOT)}')
for p in (ROOT/'frontend/src').rglob('*'):
 if p.is_file() and any(x in p.read_text() for x in ['dangerouslySetInnerHTML','localStorage.setItem','mysql://']):fail('unsafe frontend operation')
# Only key names in example, exact ignore boundary and Docker context exclusions.
for line in (ROOT/'.env.example').read_text().splitlines():
 if line and not re.fullmatch(r'[A-Z][A-Z0-9_]*=',line):fail('env example contains values')
for file in ['.gitignore','.dockerignore']:
 if '.env' not in (ROOT/file).read_text().splitlines():fail('env ignore missing')
# All committed historical file contents, plus tracked/current source. Compare configured secrets without printing.
secrets=[]
if (ROOT/'.env').exists():
 for line in (ROOT/'.env').read_text().splitlines():
  if '=' in line:
   k,v=line.split('=',1)
   if any(x in k for x in ['PASSWORD','SECRET','TOKEN','API_KEY']) and len(v)>=16:secrets.append(v)
secrets.extend(v for k,v in os.environ.items() if any(x in k for x in ['PASSWORD','SECRET','TOKEN','API_KEY']) and len(v)>=16)
patterns=[re.compile(r'gh[pousr]_[A-Za-z0-9]{25,}'),re.compile(r'sk-[A-Za-z0-9_-]{24,}'),re.compile(r'-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----')]
def scan(data,label):
 if any(v.encode() in data for v in secrets):fail('configured secret detected in '+label)
 if any(p.search(data.decode('utf-8',errors='ignore')) for p in patterns):fail('credential pattern detected in '+label)
tracked=subprocess.check_output(['git','ls-files','-z'],cwd=ROOT).decode().split('\0')
for name in tracked:
 p=ROOT/name
 if name and p.is_file() and '.git' not in p.parts:scan(p.read_bytes(),name)
objects=subprocess.check_output(['git','rev-list','--objects','--all'],cwd=ROOT,text=True).splitlines()
for line in objects:
 oid=line.split()[0]
 if subprocess.check_output(['git','cat-file','-t',oid],cwd=ROOT,text=True).strip()=='blob':scan(subprocess.check_output(['git','cat-file','blob',oid],cwd=ROOT),'Git history blob')
# Frontmatter validation for narrow project skill, no optional PyYAML dependency needed.
skill=(ROOT/'.agents/skills/verify-change/SKILL.md').read_text()
if not skill.startswith('---\n') or '\nname: verify-change\n' not in skill or '\ndescription:' not in skill:fail('skill metadata invalid')
print(json.dumps({'check':'architecture/security/history/config/skill','status':'passed'},ensure_ascii=False))
