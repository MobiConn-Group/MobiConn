// 模拟接收后端信息并改变信号图标颜色的示例代码
setTimeout(function () {
    document.querySelector('.signal-icon').classList.add('connected'); // 当接收到后端信息后，添加 connected 类来改变颜色为绿色
}, 1000); // 3秒后模拟接收到后端信息

const showMessageModal = (message) => {
    const modal = new bootstrap.Modal(document.getElementById('messageModal'));
    document.getElementById('messageContent').innerText = message;
    modal.show();
};