import {defineConfig} from 'vite';
export default defineConfig({server:{proxy:{'/health':'http://127.0.0.1:3001','/live':'http://127.0.0.1:3001','/api':'http://127.0.0.1:3001'}},build:{target:'es2022'}});
