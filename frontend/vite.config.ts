import {defineConfig,loadEnv} from 'vite';
export default defineConfig(({command,mode})=>{
 if(command==='build')return {build:{target:'es2022'}};
 const env=loadEnv(mode,'..','WEB_PORT');
 const port=Number(env.WEB_PORT);
 if(!Number.isInteger(port)||port<1||port>65535)throw new Error('Required valid config: WEB_PORT');
 const target=`http://127.0.0.1:${port}`;
 return {server:{proxy:{'/health':target,'/live':target,'/api':target}},build:{target:'es2022'}};
});
