<script setup>
import Sidebar from "./components/Sidebar.vue";
</script>

<template>
  <!--  <img id="logo" alt="Wails logo" src="./assets/images/logo-universal.png"/>-->
  <!--  <HelloWorld/>-->
  <!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>手机相册</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.4/css/all.min.css">
<!--    <link rel="stylesheet" href="../styles/img_browser.css">-->
<!--    <link rel="stylesheet" href="../styles/sidebar.css">-->
<!--    改变窗口大小-->

  </head>
  <body>


  <Sidebar/>
  <router-view></router-view>
  </body>
  </html>

</template>

<style>
#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  display: flex;
  align-items: flex-start;
  height: 100vh;
  margin: 0;
  background-color: whitesmoke;
}
</style>
