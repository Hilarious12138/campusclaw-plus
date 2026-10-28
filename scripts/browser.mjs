import {createRequire} from 'node:module';
import fs from 'node:fs';
const require=createRequire(new URL('../frontend/package.json',import.meta.url));
const {chromium}=require('@playwright/test');
const base=process.env.PLUS_BASE_URL;
if(!base)throw new Error('PLUS_BASE_URL required');
const browser=await chromium.launch({headless:true,channel:process.env.PLAYWRIGHT_CHANNEL||'chrome'});
try{
 const page=await browser.newPage();
 for(const width of [320,768,1280]){
  await page.setViewportSize({width,height:900});await page.goto(base);
  await page.getByRole('status').filter({hasText:'服务已就绪'}).waitFor();
  if(await page.evaluate(()=>document.documentElement.scrollWidth>innerWidth))throw new Error('responsive overflow');
 }
 await page.screenshot({path:'reports/homepage.png',fullPage:true});
 await page.route('**/health',route=>route.fulfill({status:503,contentType:'application/json',body:JSON.stringify({code:'DEPENDENCY_UNAVAILABLE'})}));
 await page.getByRole('button',{name:'重新检查'}).click();await page.getByRole('alert').waitFor();
 await page.screenshot({path:'reports/homepage-error.png',fullPage:true});
 await page.unroute('**/health');await page.getByRole('button',{name:'重新检查'}).focus();await page.keyboard.press('Enter');await page.getByRole('status').filter({hasText:'服务已就绪'}).waitFor();
 fs.writeFileSync('reports/browser.json',JSON.stringify({status:'passed',breakpoints:[320,768,1280],states:['ready','error','retry'],source:'real Chrome; error response injected for UI only'},null,2));
 console.log('Chrome responsive, ready/error/retry passed');
}finally{await browser.close()}
