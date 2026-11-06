// Explicit test double. Not included in production Compose or used as real-model evidence.
import http from 'node:http';import crypto from 'node:crypto';
let delay=0,mode='ok',calls=0,chatCalls=0;let lastChat=null;
const vocab=[['方程','等式','未知数','equation'],['面积','正方形','square'],['语文','阅读','poem'],['单位','长度','meter'],['光合','植物'],['SQL','database'],['课程','教学'],['噪声','zzzz']];
function embedding(s){const out=Array(8).fill(0);for(let i=0;i<vocab.length;i++)for(const k of vocab[i])if(s.toLowerCase().includes(k.toLowerCase()))out[i]+=1;if(out.every(x=>x===0)){out[7]=1;}return out;}
http.createServer(async(req,res)=>{let raw='';for await(const b of req){raw+=b;if(raw.length>4e6){res.writeHead(413).end();return;}}let b={};try{if(raw)b=JSON.parse(raw);}catch{res.writeHead(400).end();return;}
 const json=(s,v)=>{res.writeHead(s,{'Content-Type':'application/json'});res.end(JSON.stringify(v));};
 if(req.url==='/health'){json(200,{fixture:true});return;}
 if(req.url==='/control'&&req.method==='POST'){delay=Number(b.delay||0);mode=b.mode||'ok';json(200,{fixture:true});return;}
 if(req.url==='/metrics'){json(200,{fixture:true,calls,chatCalls,lastChat});return;}
 if(req.headers.authorization!=='Bearer '+process.env.FIXTURE_API_KEY){json(401,{error:'unauthorized fixture'});return;}
 if(req.url==='/v1/embeddings'){calls++;await new Promise(r=>setTimeout(r,delay));if(mode==='fail'){json(429,{error:'controlled fixture failure'});return;}json(200,{data:b.input.map((s,index)=>({index,embedding:embedding(s)}))});return;}
 if(req.url==='/v1/chat/completions'){chatCalls++;lastChat=b;res.writeHead(200,{'Content-Type':'text/event-stream'});for(const s of ['请先确认已知条件，','再检查下一步的依据。','[1]']){res.write('data: '+JSON.stringify({choices:[{delta:{content:s}}]})+'\n\n');await new Promise(r=>setTimeout(r,30));}res.end('data: [DONE]\n\n');return;}
 json(404,{error:'fixture path not found'});
}).listen(8091,'0.0.0.0');
