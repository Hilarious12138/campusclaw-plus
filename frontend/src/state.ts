export type Health = {code:string;data?:{version:string;dependencies:{mysql:string}}};
export type State = {kind:'loading'|'ready'|'error'|'empty';version?:string;mysql?:string};
export function stateFromResponse(ok:boolean,value:Health|null):State {
 if(!ok||value?.code!=='OK') return {kind:'error'};
 if(!value.data?.version||value.data.dependencies?.mysql!=='up') return {kind:'empty'};
 return {kind:'ready',version:value.data.version,mysql:value.data.dependencies.mysql};
}
export const loading:State={kind:'loading'};
