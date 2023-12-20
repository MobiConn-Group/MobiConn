
// 模拟接收后端信息并改变信号图标颜色的示例代码
setTimeout(function () {
    document.querySelector('.signal-icon').classList.add('connected'); // 当接收到后端信息后，添加 connected 类来改变颜色为绿色
}, 1000); // 3秒后模拟接收到后端信息

const socket = new WebSocket('ws://localhost:8080/ws');
var info = document.getElementById('transfer_info');
info.value += "\n"
socket.onmessage = function (event) {
    const data = JSON.parse(event.data);
    const messageInput = document.getElementById('transfer_info');
    info += data.message; // 实时更新文本框中的消息字符串
    info += "你好吗";
};

var info = document.getElementById('transfer_info');
info.value += "你好";


document.getElementById('upload_btn').addEventListener('click', function () {
    var fileInput = document.getElementById('fileInput');
    var file = fileInput.files[0];
    alert("Message！")
    if (file) {
        var formData = new FormData();
        formData.append('file', file);

        fetch('/upload', {
            method: 'POST',
            body: formData
        })
            .then(response => {
                if (response.ok) {
                    alert('File uploaded successfully');
                } else {
                    alert('Failed to upload file');
                }
            })
            .catch(error => {
                console.error('Error:', error);
            });
    } else {
        alert('Please select a file to upload');
    }
});