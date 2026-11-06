// Retry only unavailable advisory transport; a vulnerability result is never retried away.
import {spawnSync} from 'node:child_process';
for(let attempt=1;attempt<=3;attempt++){
 const p=spawnSync('npm',['audit','--audit-level=high','--json','--fetch-retries=2','--fetch-retry-mintimeout=1000','--fetch-retry-maxtimeout=5000'],{cwd:process.cwd(),encoding:'utf8',timeout:45000,maxBuffer:4*1024*1024});
 let v;try{v=JSON.parse(p.stdout)}catch{}
 if(v?.metadata?.vulnerabilities){const counts=v.metadata.vulnerabilities;console.log(JSON.stringify({attempt,vulnerabilities:counts,status:p.status===0?'passed':'failed'}));process.exit(p.status===0?0:1);}
 console.log(JSON.stringify({attempt,status:'unverified',reason:'advisory transport unavailable'}));
 if(attempt<3)await new Promise(r=>setTimeout(r,1000*attempt));
}
process.exit(1);
