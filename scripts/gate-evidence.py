#!/usr/bin/env python3
import json,subprocess,tempfile
from pathlib import Path
root=Path(__file__).resolve().parents[1];out=root/'reports';out.mkdir(exist_ok=True)
sha=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip()
with tempfile.TemporaryDirectory() as d:
 p=Path(d)/'fixture.txt';p.write_text('sk'+'-'+('artificial_noncredential_fixture_'*2))
 red=subprocess.run(['python3',str(root/'scripts/security_scan.py'),'--root',d],capture_output=True,text=True)
 p.unlink();green=subprocess.run(['python3',str(root/'scripts/security_scan.py'),'--root',d],capture_output=True,text=True)
 assert red.returncode!=0 and green.returncode==0
 (out/'gate-red-green.md').write_text('# 安全扫描门禁红绿证明\n\n提交：'+sha+'\n\n隔离人工标识fixture，不是真实凭证，也不加入Git历史。相同扫描器，相同期望：命中必须拒绝，删除后扫描通过。此记录只证明秘密扫描子门禁，完整依赖/安全门禁见当日基线。\n\n红退出码：'+str(red.returncode)+'\n```json\n'+red.stdout+'```\n绿退出码：'+str(green.returncode)+'\n```json\n'+green.stdout+'```\n')
print(json.dumps({'red_exit':red.returncode,'green_exit':green.returncode,'fixture':'artificial only'}))
