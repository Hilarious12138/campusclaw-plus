import {test} from 'node:test';
import assert from 'node:assert/strict';
import {loading,stateFromResponse} from './state.ts';
test('normal ready S13',()=>assert.deepEqual(stateFromResponse(true,{code:'OK',data:{version:'test',dependencies:{mysql:'up'}}}),{kind:'ready',version:'test',mysql:'up'}));
test('failure HTTP or business rejection S14',()=>{assert.equal(stateFromResponse(false,{code:'DEPENDENCY_UNAVAILABLE'}).kind,'error');assert.equal(stateFromResponse(true,{code:'NOT_FOUND'}).kind,'error')});
test('boundary pending and empty S14',()=>{assert.equal(loading.kind,'loading');assert.equal(stateFromResponse(true,{code:'OK'}).kind,'empty')});
