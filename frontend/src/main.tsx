import {useEffect,useState} from 'react';
import {createRoot} from 'react-dom/client';
import {loading,stateFromResponse,type State} from './state';
import './styles.css';
import {Workspace} from './Workspace';
function App(){
 const [state,setState]=useState<State>(loading);const [retry,setRetry]=useState(0);
 useEffect(()=>{let disposed=false;const c=new AbortController();const timeout=setTimeout(()=>c.abort(),5000);setState(loading);
  fetch('/health',{signal:c.signal,cache:'no-store'}).then(async r=>{const body=await r.json();if(!disposed)setState(stateFromResponse(r.ok,body))}).catch(()=>{if(!disposed)setState({kind:'error'})}).finally(()=>clearTimeout(timeout));
  return()=>{disposed=true;clearTimeout(timeout);c.abort()};
 },[retry]);
 const label={loading:'正在检查服务…',ready:'服务已就绪',error:'服务暂不可用，请稍后重试',empty:'未返回可用状态，请重试'}[state.kind];
 return <main><header><span className="badge">工程基线</span><h1>CampusClaw Plus</h1><p>面向班级教学的智能学习平台</p></header><section aria-label="服务状态"><h2>运行状态</h2><p role={state.kind==='error'?'alert':'status'} className={state.kind}>{label}</p>{state.kind==='ready'&&<dl><dt>版本</dt><dd>{state.version}</dd><dt>MySQL</dt><dd>{state.mysql}</dd></dl>}<button disabled={state.kind==='loading'} onClick={()=>setRetry(v=>v+1)}>重新检查</button></section><Workspace/><footer>CampusClaw Plus · 班级知识库与智能学习</footer></main>;
}
createRoot(document.getElementById('root')!).render(<App/>);
