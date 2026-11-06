def run(a):
 t=a.login('teacher_a');b=a.login('student_b')
 for mode in ['keyword','vector','hybrid']:
  out=a.request('POST','/api/retrieval/search',{'query':'解方程保持等式两边相等','mode':mode},token=t)['data'];assert out['hits'],mode
  for hit in out['hits']:
   assert hit['material_title'] and hit['excerpt'] and hit['end_char']>hit['start_char']
   a.request('GET','/api/retrieval/chunks/'+hit['chunk_id'],token=b,status=404)
   assert a.sql("SELECT COUNT(*) FROM knowledge_chunks WHERE id='"+hit['chunk_id']+"' AND class_id='"+'a'*32+"';")=='1'
 a.request('POST','/api/retrieval/search',{'query':'x','class_id':'b'},token=t,status=400)
 a.request('POST','/api/retrieval/search',{'query':'x','mode':'invalid'},token=t,status=400)
 out=a.request('POST','/api/retrieval/search',{'query':"' OR '1'='1",'mode':'keyword'},token=t)['data'];assert not out['hits']
 assert not a.request('POST','/api/retrieval/search',{'query':'方程','mode':'keyword'},token=b)['data']['hits']
