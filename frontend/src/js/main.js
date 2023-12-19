document.getElementById('upload-btn').addEventListener('click', function () {
    var fileInput = document.getElementById('file-input');
    var file = fileInput.files[0];

    if (file) {
        var formData = new FormData();
        formData.append('file', file);

        fetch('http://example.com/upload', {
            method: 'POST',
            body: formData
        })
            .then(response => response.json())
            .then(data => {
                console.log('后端返回的数据：', data);
                var leftContent = document.querySelector('.left-content');
                leftContent.innerHTML = '<h2>后端返回的文本信息</h2><p>' + data.text + '</p>';
            })
            .catch(error => {
                console.error('发生错误：', error);
            });
    } else {
        console.log('请选择一个文件');
    }
});

// 模拟接收后端信息并改变信号图标颜色的示例代码
setTimeout(function () {
    document.querySelector('.signal-icon').classList.add('connected'); // 当接收到后端信息后，添加 connected 类来改变颜色为绿色
}, 1000); // 3秒后模拟接收到后端信息


// 假设后端返回的相册图片数据为一个包含图片URL的数组
const photoData = ["/path/to/photo1.jpg", "/path/to/photo2.jpg", "/path/to/photo3.jpg"];

// 选择图片容器
const photoGallery = document.getElementById('photoGallery');

// 遍历图片数据，并创建 img 元素并添加到图片容器中
photoData.forEach(photoUrl => {
    const img = document.createElement('img');
    img.src = photoUrl;
    img.alt = '手机相册照片';
    photoGallery.appendChild(img);
});