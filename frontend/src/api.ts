export type User={id:string;username:string;role:'teacher'|'student';class_id:string};
export type Resource={id:string;revision:number;body:Record<string,unknown>};
export type Material={id:string;title:string;task_id:string;status:string};
export type Citation={chunk_id:string;material_id:string;material_title:string;chunk_index:number;start_char:number;end_char:number;excerpt:string};
export async function request<T>(token:string,path:string,method='GET',body?:unknown,headers:Record<string,string>={}):Promise<T>{
 const isForm=body instanceof FormData;
 const r=await fetch(path,{method,headers:{...(!isForm&&body!==undefined?{'Content-Type':'application/json'}:{}),...(token?{Authorization:'Bearer '+token}:{}),...headers},body:body===undefined?undefined:isForm?body:JSON.stringify(body)});
 const v=await r.json();if(!r.ok)throw new Error(({401:'登录已失效，请重新登录',403:'当前账号没有操作权限',404:'资源不存在或不可访问',409:'内容已变更或请求冲突，请刷新重试',429:'请求过多，请稍后重试',503:'服务暂不可用，请稍后重试'} as Record<number,string>)[r.status]||'输入无效，请检查后重试');return v.data;
}
export async function streamEvents(response:Response,onEvent:(kind:string,data:Record<string,unknown>)=>void){
 if(!response.ok||!response.body)throw new Error('会话请求失败，请检查登录和服务状态');
 const reader=response.body.getReader(),decoder=new TextDecoder();let buffer='',done=false;
 for(;;){const v=await reader.read();buffer+=decoder.decode(v.value,{stream:!v.done}).replaceAll('\r\n','\n');let index;
  while((index=buffer.indexOf('\n\n'))>=0){const frame=buffer.slice(0,index);buffer=buffer.slice(index+2);let kind='message';const data:string[]=[];for(const line of frame.split('\n')){if(line.startsWith('event:'))kind=line.slice(6).trim();if(line.startsWith('data:'))data.push(line.slice(5).trimStart());}if(data.length){const value=JSON.parse(data.join('\n'));onEvent(kind,value);if(kind==='done')done=true;if(kind==='error')throw new Error(String(value.message||'回答未完成'));}}
  if(v.done)break;
 }
 if(!done)throw new Error('回答连接中断，未标记完成，请重新查询会话');
}
