<script setup>
import {Exit} from "../../wailsjs/go/main/App.js";

function exit() {
  Exit()
}
</script>

<template>
  <div id="sidebar">
    <ul>
      <li>
        <router-link to="/device">
          <i class="fas fa-mobile"></i>
          移动设备
        </router-link>
      </li>
      <li>
        <a @click="exit()" href="#">
          <i class="fas fa-door-open"></i>
          退出
        </a>
      </li>
<!--      <li>-->
<!--        <router-link to="/photo">-->
<!--          <i class="fas fa-images"></i>-->
<!--          查看手机照片-->
<!--        </router-link>-->
<!--      </li>-->
<!--      <li><a href="./settings.html"><i class="fas fa-cogs"></i>设置</a></li>-->
<!--      <li>-->
<!--        <router-link to="/hello">-->
<!--          <i class="fas fa-file-upload"></i>-->
<!--          你好-->
<!--        </router-link>-->
<!--      </li>-->
    </ul>
<!--    <div class="signal-container">-->
<!--      <div class="signal-icon"><i class="fas fa-signal"></i></div>-->
<!--    </div>-->
  </div>
</template>

<style scoped>
/* 设置基本字体和页面背景 */
body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  display: flex;
  height: 100vh;
  margin: 0;
  background-color: #f4f4f4;
}

/* 设置侧边栏样式 */
#sidebar {
  width: 250px;
  background-color: #333;
  color: #fff;
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  transition: width 0.5s;
  /* Smooth transition for sidebar width */
  z-index: 1000;
  /* Ensure sidebar is on top of other content */

}

/* 为了在hover时增加侧边栏的宽度 */
#sidebar:hover {
  width: 300px;
}

/* 调整点击效果 */
#sidebar a:active {
  background-color: #111;
}

/* 设置侧边栏内文字和图标间距以及效果 */
#sidebar ul {
  list-style-type: none;
  padding: 0;
  /* Reset default padding */
}

#sidebar li {
  margin-top: 50px;
}


/* 设置链接颜色和padding */
#sidebar a {
  display: flex;
  align-items: center;
  color: #ddd;
  /* Slightly muted white for links */
  text-decoration: none;
  padding: 15px 20px;
  transition: all 0.3s, width 0.5s;
  /* Smooth transition for background color and padding */
}

/* 鼠标悬停时背景色变化 */
#sidebar a:hover {
  background-color: #555;
  color: #fff;
  /* Brighter white on hover */
  /* transform: scale(1.1); */
}

/* 设置图标和文字间距 */
#sidebar a i {
  min-width: 20px;
  /* Ensure icons align nicely */
  margin-right: 10px;
  transition: transform 0.3s;
  /* Smooth transition for icon rotation */
}

/* 旋转图标 */
#sidebar a:hover i {
  transform: rotate(20deg);
}


/* 设置内容区域 */
#content {
  flex: 1;
  padding: 20px;
  padding-left: 270px;
  /* Adjust padding to match sidebar width */
  transition: padding-left 0.5s;
  /* Smooth transition for padding */
}

/* 在hover侧边栏时调整内容区域的padding */
#sidebar:hover + #content {
  padding-left: 320px;
}

/* 新增样式用于控制信号图标颜色 */
.signal-icon {
  color: red;
  /* 初始状态为红色 */
  transition: color 0.5s;
  /* 添加颜色的过渡效果 */
}

.signal-icon.connected {
  color: green;
  /* 当连接后后台时，改变为绿色 */
}

/* 新增样式用于控制容器位置 */
.signal-container {
  display: flex;
  justify-content: flex-end;
  /* 水平方向居右 */
  align-items: flex-end;
  /* 垂直方向居下 */
  position: absolute;
  right: 20px;
  bottom: 20px;
  /*占满父容器的高度 */
  padding: 10px;
  /* 添加一些 padding 以保持图标与边界的距离 */
}
</style>