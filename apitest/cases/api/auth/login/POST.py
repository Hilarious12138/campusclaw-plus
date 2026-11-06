import uuid
def run(a):
 token=a.login('teacher_a');assert token
 u=a.request('GET','/api/me',token=token)['data'];assert u['role']=='teacher' and u['class_id']=='a'*32
 a.request('GET','/api/me',status=401)
 a.request('GET','/api/me',token=token+'bad',status=401)
 a.request('GET','/api/me?class_id=other',token=token,status=400)
 a.request('GET','/api/me',headers={'Cookie':'access_token='+token},status=401)
 a.request('POST','/api/auth/logout',{},token=token,headers={'Origin':'https://evil.invalid'},status=403)
 a.request('POST','/api/auth/logout',{},token=token)
 a.request('GET','/api/me',token=token,status=401)
 a.tokens.pop('teacher_a')
 assert a.sql("SELECT COUNT(*) FROM users WHERE password_hash LIKE '$argon2id$%';")=='4'
 assert a.sql('SELECT COUNT(DISTINCT password_hash) FROM users;')=='4'
 assert int(a.sql('SELECT COUNT(*) FROM revoked_tokens;'))>=1

 # The account fingerprint is separate from legitimate users and never logged in plaintext.
 rate_user='rate-'+uuid.uuid4().hex
 for _ in range(5):a.request('POST','/api/auth/login',{'username':rate_user,'password':'wrong'},status=401)
 a.request('POST','/api/auth/login',{'username':rate_user,'password':'wrong'},status=429)
