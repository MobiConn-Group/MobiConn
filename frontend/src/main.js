import App from './App.vue'
import DevicePage from "./components/DevicePage.vue";
import HelloWorld from "./components/HelloWorld.vue";
import PhotoPage from "./components/PhotoPage.vue";
import './style.css';
import {createRouter, createWebHashHistory} from "vue-router";
import {createApp} from "vue";

const router = createRouter({
    history: createWebHashHistory(), routes: [{
        path: '/', component: App
    }, {
        path: '/device', component: DevicePage
    }, {
        path: '/photo/:address', component: PhotoPage, name: 'photo'
    }, {
        path: '/hello', component: HelloWorld
    }],
});
// export default {
//     components: {
//         DeviceTable,
//     },
//     data() {
//         return {
//             devices: [], // 从Wails应用程序中获取的设备数据
//         };
//     },
//     mounted() {
//         // 从Wails应用程序中获取设备数据
//         this.devices = window.backend.GetDevices();
//     },
// };

const app = createApp(App);
app.use(router);
app.mount('#app');
