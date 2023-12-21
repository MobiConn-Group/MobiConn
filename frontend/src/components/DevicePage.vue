<script setup>
import {reactive} from "vue";
import {GetDevices, OpenAndSendFile, SendRing} from "../../wailsjs/go/main/App.js";

const data = reactive({
  devices: [
    {address: "qwerty", connected: false}
    // 示例数据
    // {address: '设备1地址', connected: false},
    // {address: '设备2地址', connected: false},
    // 更多设备...
  ]
})

setInterval(() => {
  GetDevices().then(result => {
    data.devices = []
    for (const dev in result) {
      data.devices.push({address: result[dev].address, connected: result[dev].connected})
    }
  })
}, 100)

function sendFile(device) {
// 发送文件的逻辑
  OpenAndSendFile(device.address)
}

function sendRing(device) {
  // 发送响铃的逻辑
  SendRing(device.address, 'ring.mp3', 5000)
}
</script>

<template>

  <main>
    <div id="content">
      <div>
        <h1>
          已连接的设备
        </h1>
      </div>
      <table>
        <thead>
        <tr>
          <th style="color: whitesmoke">设备地址</th>
          <th style="color: whitesmoke">已连接</th>
          <th style="color: whitesmoke">查看照片</th>
          <th style="color: whitesmoke">发送文件</th>
          <th style="color: whitesmoke">发送响铃</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="device in data.devices" :key="device.address">
          <td>{{ device.address }}</td>
          <td><input type="checkbox" v-model="device.connected"></td>
          <td>
            <!--  向路由 /photo 传参 address 传值 hello -->
            <router-link :to="{path: '/photo/' + device.address}">
              <button>查看照片</button>
            </router-link>
          </td>
          <td>
            <button @click="sendFile(device)">发送文件</button>
          </td>
          <td>
            <button @click="sendRing(device)">发送响铃</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
  </main>
</template>

<style scoped>

#content {
  margin-left: 300px;
  width: 50%;
}

table {
  border-collapse: collapse;
  width: 700px;
}

th, td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
  color: dimgray;
}

th {
  background-color: gray;
}

h1 {
  color: #111111;
  text-align: left;
}
</style>