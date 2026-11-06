import tempfile,unittest
from pathlib import Path
from security_scan import scan
class SecurityScanTests(unittest.TestCase):
 def test_red_then_green(self):
  with tempfile.TemporaryDirectory() as d:
   root=Path(d);p=root/'sample.txt'
   p.write_text('sk'+'-'+('artificial_noncredential_fixture_'*2))
   self.assertEqual(scan(root),1)
   p.write_text('safe sample')
   self.assertEqual(scan(root),0)
 def test_configured_secret_is_not_printed(self):
  with tempfile.TemporaryDirectory() as d:
   root=Path(d);value='generated-fixture-only-123456789'
   (root/'.env').write_text('JWT_SECRET='+value)
   (root/'bundle.js').write_text(value)
   self.assertEqual(scan(root),1)
