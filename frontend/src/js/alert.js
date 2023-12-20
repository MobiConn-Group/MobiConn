// 获取模态
var modal = document.getElementById("myModal");

// 获取打开模态框的按钮
var btn = document.getElementById("myBtn");

// 获取关闭模态的 <span> 元素
var span = document.getElementsByClassName("close")[0];

// 当用户点击按钮时，打开模态
btn.onclick = function() {
  modal.style.display = "block";
}

// 当用户点击 <span> (x) 时，关闭模态
span.onclick = function() {
  modal.style.display = "none";
}

// 当用户点击模态之外的任何地方时，关闭它
window.onclick = function(event) {
  if (event.target == modal) {
    modal.style.display = "none";
  }
}