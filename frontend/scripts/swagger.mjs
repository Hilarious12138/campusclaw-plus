import fs from 'node:fs';
fs.mkdirSync('dist/swagger',{recursive:true});
for(const name of ['swagger-ui-bundle.js','swagger-ui.css'])fs.copyFileSync('node_modules/swagger-ui-dist/'+name,'dist/swagger/'+name);
fs.writeFileSync('dist/swagger/index.html','<!doctype html><html lang="zh-CN"><meta charset="utf-8"><title>CampusClaw Plus API</title><link rel="stylesheet" href="./swagger-ui.css"><div id="swagger-ui"></div><script src="./swagger-ui-bundle.js"></script><script>SwaggerUIBundle({url:"/openapi.json",dom_id:"#swagger-ui",persistAuthorization:false,validatorUrl:null});</script></html>');
