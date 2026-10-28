import {test} from 'node:test';
import assert from 'node:assert/strict';
import {spawnSync} from 'node:child_process';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
const phases={checks:['C01','C02','C03','C04','C05','C06','C07','C08'],healthy:['R01','R02','R03'],down:['R04'],recovered:['R05'],stopped:['R06','R07']};
function run(setup,extra={}){const dir=fs.mkdtempSync(path.join(os.tmpdir(),'campusclaw-plus-report-unit-'));try{setup(dir);const p=spawnSync(process.execPath,['scripts/grader.mjs','finalize'],{env:{...process.env,GRADE_RESULT_DIR:dir,...extra},encoding:'utf8'});return {exit:p.status,data:JSON.parse(fs.readFileSync(path.join(dir,'report.json'),'utf8')),html:fs.readFileSync(path.join(dir,'report.html'),'utf8')};}finally{fs.rmSync(dir,{recursive:true,force:true});}}
function fixtures(dir,bad=false){for(const [phase,ids] of Object.entries(phases))fs.writeFileSync(path.join(dir,phase+'.json'),JSON.stringify(ids.map((id,i)=>({id,name:'<script>unsafe</script>',criterion:'unit fixture only',status:bad&&id==='C01'?'failed':'passed',elapsed_ms:1,detail:'synthetic proof'}))));}
test('normal report aggregation uses all fifteen required cases',()=>{const r=run(fixtures);assert.equal(r.exit,0);assert.equal(r.data.counts.passed,15);assert.equal(r.data.status,'passed');assert.equal(r.html.includes('<script>unsafe</script>'),false);assert.ok(r.html.includes('&lt;script&gt;unsafe'));});
test('failed assertion makes report and exit fail',()=>{const r=run(d=>fixtures(d,true));assert.equal(r.exit,1);assert.equal(r.data.counts.failed,1);assert.equal(r.data.status,'failed');});
test('empty or missing phases cannot pass',()=>{const r=run(()=>{});assert.equal(r.exit,1);assert.equal(r.data.counts.blocked,15);});
test('host cleanup failure cannot leave a success report',()=>{const r=run(fixtures,{HOST_FAILURE:'1',HOST_FAILURE_REASON:'cleanup'});assert.equal(r.exit,1);assert.equal(r.data.counts.failed,1);assert.equal(r.data.status,'failed');});

test('missing env still produces a failure report without host developer runtimes',()=>{const dir=fs.mkdtempSync(path.join(os.tmpdir(),'campusclaw-plus-shell-unit-'));try{fs.mkdirSync(path.join(dir,'scripts'));fs.copyFileSync('scripts/grade.sh',path.join(dir,'scripts/grade.sh'));const p=spawnSync('sh',[path.join(dir,'scripts/grade.sh')],{encoding:'utf8'});assert.equal(p.status,1);const reportDir=path.join(dir,'reports',fs.readdirSync(path.join(dir,'reports'))[0]);assert.equal(JSON.parse(fs.readFileSync(path.join(reportDir,'report.json'),'utf8')).status,'failed');assert.ok(fs.statSync(path.join(reportDir,'report.html')).size>100);}finally{fs.rmSync(dir,{recursive:true,force:true});}});
