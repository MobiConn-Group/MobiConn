import * as wails from '@wailsapp/runtime';

wails.Events.On('img_browser', (imageData) => {
    const imageDisplay = document.getElementById('imageDisplay');
    imageDisplay.src = 'data:image/png;base64,' + imageData;
  });
  
  function sendNumberToBackend(number) {
    wails.Events.Emit('sendNumberToBackend', number);
  }

// 模拟接收后端信息并改变信号图标颜色的示例代码
setTimeout(function () {
    document.querySelector('.signal-icon').classList.add('connected'); // 当接收到后端信息后，添加 connected 类来改变颜色为绿色
}, 1000); // 3秒后模拟接收到后端信息

var maxValue = 100; // 初始化最大值为100
var minValue = 1

// 发送 AJAX 请求以获取最大值
var xhr = new XMLHttpRequest();
xhr.onreadystatechange = function () {
    if (xhr.readyState === XMLHttpRequest.DONE) {
        if (xhr.status === 200) {
            var response = JSON.parse(xhr.responseText);
            maxValue = response.maxValue;

            // 在获得最大值后更新页面上的最大值元素
            var maxInput = document.getElementById('maxValue');
            maxInput.value = maxValue.toString();
        } else {
            console.error('Failed to get max value.');
        }
    }
};
xhr.open('GET', '/album'); // 向根路径发送 GET 请求
xhr.send();

function nextPage() {
    var number = document.getElementById('photo_id');
    var currentValue = parseInt(number.textContent);
    if (currentValue < maxValue) {
        number.textContent = currentValue + 1;
    } else {
        number.textContent = minValue.toString();
    }
}

function prevPage() {
    var number = document.getElementById('photo_id');
    var currentValue = parseInt(number.textContent);
    if (currentValue > minValue) {
        number.textContent = currentValue - 1;
    } else {
        number.textContent = maxValue.toString();
    }
}

function requestTotalImageCount() {
	const url = '/image/count';

	fetch(url)
		.then(response => response.json())
		.then(result => {
			// 在这里处理返回的图片总数量
            maxValue = result.total;
		})
		.catch(error => {
			// 处理错误
			console.error(error);
		});
}

// 调用函数获取图片总数量
requestTotalImageCount();