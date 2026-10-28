import os,subprocess,tempfile,unittest
from pathlib import Path
ROOT=Path(__file__).resolve().parents[1]
class DatabaseConfig(unittest.TestCase):
 def test_normal_and_failure_root_secrets(self):
  with tempfile.TemporaryDirectory(prefix='campusclaw-plus-db-check-') as temp:
   fake=Path(temp)/'docker-entrypoint.sh';fake.write_text('#!/bin/sh\nexit 0\n');fake.chmod(0o700)
   for root,user,good in [('a'*24,'b'*24,True),('a'*23,'b'*24,False),('a'*24,'a'*24,False),('a'*24,'b'*23,False)]:
    with self.subTest(length=len(root),same=root==user):
     r=subprocess.run(['sh',str(ROOT/'deploy/db-entrypoint.sh'),'mysqld'],env=dict(os.environ,PATH=temp+os.pathsep+os.environ['PATH'],MYSQL_ROOT_PASSWORD=root,MYSQL_PASSWORD=user),capture_output=True,text=True)
     self.assertEqual(r.returncode==0,good)
     self.assertNotIn(root,r.stdout+r.stderr);self.assertNotIn(user,r.stdout+r.stderr)
