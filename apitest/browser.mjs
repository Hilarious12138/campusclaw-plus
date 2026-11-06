import fs from 'node:fs';import assert from 'node:assert/strict';import {createRequire} from 'node:module';
const require=createRequire('/work/frontend/package.json');const {chromium}=require('@playwright/test');
const base=process.env.BASE_URL,dir=process.env.GRADE_RESULT_DIR||'/reports';
const browser=await chromium.launch({headless:true,executablePath:'/usr/bin/chromium-browser',args:['--no-sandbox','--disable-dev-shm-usage']});
try{
 const page=await browser.newPage({viewport:{width:1280,height:1000}});
 await page.goto(base);await page.getByLabel('账号',{exact:true}).fill('teacher_a');await page.getByLabel('密码',{exact:true}).fill(process.env.SEED_PASSWORD);await page.getByRole('button',{name:'登录',exact:true}).click();await page.getByRole('heading',{name:'班级知识库',exact:true}).waitFor();
 await page.screenshot({path:dir+'/knowledge.png',fullPage:true});
 const login=await page.request.post(base+'/api/auth/login',{data:{username:'teacher_a',password:process.env.SEED_PASSWORD}});const token=(await login.json()).data.access_token;
 const h={Authorization:'Bearer '+token};const key='xss-'+Date.now();
 const created=await page.request.post(base+'/api/skills',{headers:h,data:{name:'<script>alert(1)</script>',instructions:'普通文本',conditions:[],tools:[],enabled:false}});assert.equal(created.status(),200);
 const skill=(await created.json()).data;
 await page.getByRole('button',{name:'教学配置',exact:true}).click();await page.getByText('Skill技能与导入',{exact:true}).click();
 await page.getByLabel('编辑技能',{exact:true}).selectOption(skill.id);
 assert.equal(await page.getByLabel('技能名称',{exact:true}).inputValue(),'<script>alert(1)</script>');
 const option=page.locator('option').filter({hasText:'<script>alert(1)</script>'});assert.ok((await option.evaluate(el=>el.outerHTML)).includes('&lt;script&gt;'));assert.equal(await page.locator('script').evaluateAll(nodes=>nodes.some(n=>n.textContent?.includes('alert(1)'))),false);
 await page.screenshot({path:dir+'/xss-escaped.png',fullPage:true});
 await page.getByText('MCP工具接入',{exact:true}).click();await page.screenshot({path:dir+'/mcp.png',fullPage:true});
 await page.getByRole('button',{name:'AI学习应用',exact:true}).click();await page.getByLabel('选择应用',{exact:true}).selectOption({label:'数学助手'});await page.getByRole('button',{name:'新建会话',exact:true}).click();await page.getByLabel('你的问题',{exact:true}).fill('解方程怎样保持等式两边相等？');await page.getByRole('button',{name:'发送问题',exact:true}).click();await page.getByText('已完成',{exact:true}).waitFor({timeout:70000});
 await page.screenshot({path:dir+'/conversation.png',fullPage:true});
 await page.getByRole('button',{name:'操作审计',exact:true}).click();await page.getByRole('heading',{name:'本班操作审计',exact:true}).waitFor();await page.screenshot({path:dir+'/audit.png',fullPage:true});
 assert.equal(await page.evaluate(()=>Object.keys(localStorage).length),0);assert.equal(await page.evaluate(()=>Object.keys(sessionStorage).length),0);
 await page.request.delete(base+'/api/skills/'+skill.id,{headers:h});
 fs.writeFileSync(dir+'/browser-business.json',JSON.stringify({status:'passed',evidence:'Real Linux Chromium, escaped DOM, in-memory token, fixture model; screenshots ungenerated',screenshots:['knowledge.png','xss-escaped.png','mcp.png','conversation.png','audit.png']}));
}finally{await browser.close()}
