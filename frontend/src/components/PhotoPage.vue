<script setup>
import {reactive} from "vue";
import {useRoute} from "vue-router";
import {GetPhoto, GetPhotoCount} from "../../wailsjs/go/main/App.js";

const route = useRoute()

const data = reactive({
  photo: "",
  photoIndex: 0,
  photoCount: 0,
  remoteAddr: route.params.address
})

function getPhoto() {
  GetPhoto(data.remoteAddr, data.photoIndex).then(result => {
    data.photo = result
  })
}

function getPhotoCount() {
  GetPhotoCount(data.remoteAddr).then(result => {
    data.photoCount = result
    if (data.photoCount > 0) {
      getPhoto(data.remoteAddr, data.photoIndex)
    }
  })
}

function nextPhoto() {
  if (data.photoIndex < data.photoCount - 1) {
    ++data.photoIndex
    getPhoto(data.remoteAddr, data.photoIndex)
  }
}

function prevPhoto() {
  if (data.photoIndex > 0) {
    --data.photoIndex
    getPhoto(data.remoteAddr, data.photoIndex)
  }
}

getPhotoCount()

</script>

<template>
  <!--鼠标点击控制左右，浏览图片-->
  <main>
    <div id="content">
      <div class="mid">
        <div>
          <h1 style="width: 500px">
            查看手机照片
            <a style="font-size: 20px; text-align: right">{{ data.remoteAddr }}</a>
            <!--          123456-->
          </h1>
        </div>
        <div class="button-container">
          <button class="button" @click="prevPhoto()">上一张</button>
          <textarea style="margin-left: 70px; width: 220px; height:43px; font-size: 16px; resize: none"
                    disabled="disabled"
                    :value="data.photoIndex + 1"></textarea>
          <p style="margin-left: 10px">/{{ data.photoCount }}</p>
          <button class="button" @click="nextPhoto()" style="margin-left: 70px">下一张</button>
        </div>
        <!--      <div>-->
        <!--        <h2>{{ data.photo }}</h2>-->
        <!--      </div>-->
        <div class="img">
          <!--                    <img id="photo" src="https://wails.io/zh-Hans/assets/images/wails-c83b438ac6f00d3db85089a2d1d0c887.webp"-->
          <!--                         style="width: 800px" alt="aaa">-->
          <img id="photo" style="max-width: 800px; max-height: 600px; width: auto; height: auto" :src="'data:image;base64, ' + data.photo"
               :alt="data.photoCount === 0 ? '请先选择设备' : data.photo">
        </div>

      </div>
      <!--    <div class="right">-->
      <!--      <div class="img">-->
      <!--        <img src="https://wails.io/zh-Hans/assets/images/wails-mac-ebdd63462decbbae1c0e3d4112f7b468.webp" alt="图片2">-->
      <!--      </div>-->
      <!--    </div>-->
      <!--    <div class="image-container">-->
      <!--      <img :src="getPhoto(0)" alt="Image"/>-->
      <!--    </div>-->
      <!--      <p>{{ getPhoto(0) }}}</p>-->
    </div>

    <!--    <div id="addr" hidden="hidden" onload="alert(this.innerText)">{{$route.params.address}}</div>-->
  </main>
</template>

<style scoped>
#content {
  margin-left: 300px;
  width: 50%;
  color: #111111;
}

h1 {
  color: #111111;
  text-align: left;
}

.button-container {
  display: flex;
  justify-content: space-between;
}

.button {
  width: 450px;
  height: 50px;
}

.img {
  margin-top: 20px;
  width: 800px;
}

</style>