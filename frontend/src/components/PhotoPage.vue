<script setup>
import {reactive} from "vue";
import {useRoute} from "vue-router";

const route = useRoute()

const data = reactive({
  photo: "",
  photoIndex: 0,
  photoCount: 0,
  remoteAddr: route.params.address
})


//
// function getPhoto(photoIndex) {
//   let getPhotoResult = null
//   GetPhoto(photoIndex).then(result => {
//     getPhotoResult = result
//   })
//   return getPhotoResult
// }
//
// function nextPhoto() {
//   if (data.photoIndex < data.photoCount - 1) {
//     ++data.photoIndex
//     data.photo = getPhoto(data.photoIndex)
//   }
// }
//
// function prevPhoto() {
//   if (data.photoIndex > 0) {
//     data.photoIndex--
//     data.photo = getPhoto(data.photoIndex)
//   }
// }

</script>

<template>
  <!--鼠标点击控制左右，浏览图片-->
  <main>
    <div id="content">
      <div>
        <h1 style="width: 500px">
          查看手机照片
          <a style="font-size: 20px; text-align: right">{{ data.remoteAddr }}</a>
          <!--          123456-->
        </h1>
      </div>
      <!--    <div class="mid">-->
      <div class="img">
        <img src="https://wails.io/zh-Hans/assets/images/wails-c83b438ac6f00d3db85089a2d1d0c887.webp" alt="图片1">
      </div>
      <!--    </div>-->
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
}

h1 {
  color: #111111;
  text-align: left;
}
</style>