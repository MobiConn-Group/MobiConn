import {createApp} from 'vue'
import App from './App.vue'
import './style.css';


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

createApp(App).mount('#app')
