package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/go-toast/toast"
	"github.com/google/uuid"
	"github.com/micmonay/keybd_event"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// 保存当前已连接的设备
var connectedDevices = map[string]Device{}

// 维护每个设备的下载列表
var deviceDownloadList = map[string][]string{}

// 维护每个设备的上传列表
var deviceUploadList = map[string][]UploadInfo{}

// 记录将要响铃的设备主机和铃声及响铃持续秒数
var deviceRingList = map[string]RingInfo{}

// 记录每个设备需要读取的的相册路径
var deviceAlbumFetchPaths = map[string][]string{}

// 记录每个设备需要读取的的相册文件
var deviceAlbumFetchFiles = map[string][]string{}

// Device 表示一个设备的信息。
type Device struct {
	Hostname        string
	MacAddress      string
	LastConnection  time.Time
	AllowCursorText bool
	AllowPowerPoint bool
	Token           uuid.UUID
}

// DownloadInfo 表示下载请求的信息。
type DownloadInfo struct {
	FileName string
	Length   int64
}

// AcceptStatus 表示上传请求是否被接受或被拒绝。
type AcceptStatus int

const (
	Accept AcceptStatus = iota
	Reject
	Pending
)

// UploadInfo 表示上传请求的信息。
type UploadInfo struct {
	FileName string       // 文件名
	Length   int64        // 文件长度
	Accepted AcceptStatus // 上传请求是否被接受或被拒绝
}

// RingInfo 表示响铃任务的信息。
type RingInfo struct {
	Tone     string // 铃声文件名
	Duration int    // 响铃持续秒数
}

var keyBonding keybd_event.KeyBonding

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create download directory if not exists
	_, err := os.Stat("./download")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("./download", os.ModePerm)
			if err != nil {
				fmt.Println("Error:", err.Error())
			}
		}
	}

	// Create http server
	err = httpServer()
	if err != nil {
		println("Error:", err.Error())
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "MobiConn",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			&Device{},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

	if keyBonding, err = keybd_event.NewKeyBonding(); err != nil {
		panic(err)
	}
}

func httpServer() error {
	// 创建ServeMux实例
	mux := http.NewServeMux()

	// 绑定处理函数到路由路径
	mux.HandleFunc("/greeting", greetingHandler)
	mux.HandleFunc("/connect", connectHandler)
	mux.HandleFunc("/heartbeat", heartbeatHandler)
	mux.HandleFunc("/download", downloadHandler)
	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/transit", transitHandler)
	mux.HandleFunc("/cursorText", cursorTextHandler)
	mux.HandleFunc("/powerPoint", powerPointHandler)
	mux.HandleFunc("/photo", photoHandler)

	mux.HandleFunc("/powerPointDemo", powerPointDemoHandler)

	// 启动HTTP服务器
	go func() {
		err := http.ListenAndServe(":25236", mux)
		if err != nil {
			fmt.Println("HTTP server error:", err)
		}
	}()

	return nil
}

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	message := "Hello World!"
	_, _ = fmt.Fprintln(w, message)
}

func readIP(req *http.Request) (ip string, port string, err error) {
	//fmt.Printf("<h1>static file server</h1><p><a href='./static'>folder</p></a>")

	ip, port, err = net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		fmt.Printf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
		fmt.Printf("userip: %q is not IP:port", req.RemoteAddr)
		return
	}

	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr
	// Header.Get is case-insensitive
	forward := req.Header.Get("X-Forwarded-For")

	fmt.Printf("<p>IP: %s</p>", ip)
	fmt.Printf("<p>Port: %s</p>", port)
	fmt.Printf("<p>Forwarded for: %s</p>", forward)
	return ip, port, err
}

func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	// 检查客户端是否已连接过
	if _, ok := connectedDevices[r.RemoteAddr]; !ok {
		_, _ = fmt.Fprintln(w, "You have never connected to this server yet. Please connect first.")
		return
	}

	// 检查Token是否匹配
	if r.FormValue("token") != connectedDevices[r.RemoteAddr].Token.String() {
		_, _ = fmt.Fprintln(w, "Token is not matched. Please connect again.")
		return
	}

	// 更新最后连接时间
	device := connectedDevices[r.RemoteAddr]
	device.LastConnection = time.Now()
	connectedDevices[r.RemoteAddr] = device

	var goals []any
	// 检查是否有文件需要下载
	for _, s := range deviceDownloadList[r.RemoteAddr] {
		fileInfo, err := os.Stat("./download/" + s)
		if err != nil {
			// 文件不存在，输出控制台警告
			fmt.Println("Error:", err.Error())
			continue
		}
		goal := map[string]any{
			"Action": "download",
			"Information": DownloadInfo{
				FileName: s,
				Length:   fileInfo.Size(),
			},
		}
		goals = append(goals, goal)
	}
	// 检查响铃任务信息
	if ringInfo, ok := deviceRingList[r.RemoteAddr]; ok {
		goal := map[string]any{
			"Action":      "ring",
			"Information": ringInfo,
		}
		goals = append(goals, goal)
	}

	// 检查是否有文件需要上传
	for _, uploadInfo := range deviceUploadList[r.RemoteAddr] {
		if uploadInfo.Accepted == Pending {
			continue
		} else if uploadInfo.Accepted == Reject {
			goal := map[string]any{
				"Action": "upload",
				"Information": map[string]any{
					"FileName": uploadInfo.FileName,
					"Status":   "rejected",
				},
			}
			goals = append(goals, goal)
			continue
		} else if uploadInfo.Accepted == Accept {
			goal := map[string]any{
				"Action": "upload",
				"Information": map[string]any{
					"FileName": uploadInfo.FileName,
					"Status":   "accepted",
				},
			}
			goals = append(goals, goal)
		}

		// 从上传列表中删除该文件
		for i, info := range deviceUploadList[r.RemoteAddr] {
			if info.FileName == uploadInfo.FileName {
				deviceUploadList[r.RemoteAddr] = append(deviceUploadList[r.RemoteAddr][:i], deviceUploadList[r.RemoteAddr][i+1:]...)
				break
			}
		}
		continue
	}

	// 检查是否有相册需要读取
	for _, path := range deviceAlbumFetchPaths[r.RemoteAddr] {
		goal := map[string]any{
			"Action": "album",
			"Path":   path,
		}
		goals = append(goals, goal)
	}

	// 检查是否有相册文件需要读取
	for _, file := range deviceAlbumFetchFiles[r.RemoteAddr] {
		goal := map[string]any{
			"Action": "photo",
			"File":   file,
		}
		goals = append(goals, goal)
	}

	// 将goals加入响应
	response := map[string]any{
		"Status": "success",
	}
	if len(goals) != 0 {
		response["goals"] = goals
	}

	// 将响应转换为JSON
	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")
	// Send response
	_, _ = fmt.Fprintln(w, string(responseJson))
	//}
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	// 连接请求必须带有MAC地址
	macAddress := r.FormValue("macAddress")
	if macAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "MAC address is required.")
		return
	}

	device := Device{
		Hostname:        r.RemoteAddr,
		MacAddress:      macAddress,
		LastConnection:  time.Now(),
		AllowCursorText: true,
		AllowPowerPoint: true,
	}

	// 询问服务端是否允许连接
	if !connectionAllowed(&device) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintln(w, "Connection is not allowed.")
		return
	}

	// 生成Token
	token := uuid.New()
	device.Token = token
	connectedDevices[r.RemoteAddr] = device

	// 返回Token和超时时间
	response := map[string]any{
		"Status":  "success",
		"Token":   token,
		"Timeout": 5,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	// 发送响应
	_, _ = fmt.Fprintln(w, string(responseJson))

	// 发送文件
	sendFileTo(r.RemoteAddr, "amber.png")

	// 发送响铃任务
	SendRingTo(r.RemoteAddr, "ringtone.mp3", 5)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// 检查Token
	if r.FormValue("token") != connectedDevices[r.RemoteAddr].Token.String() {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintln(w, "Token is not matched. Please connect again.")
		return
	}

	// 从Query中获取文件名
	fileName := r.FormValue("fileName")
	if fileName == "" {
		_, _ = fmt.Fprintln(w, "File name is required.")
		return
	}

	// 检查是否为文件下载完成的回调
	if r.FormValue("status") == "success" {
		// 删除下载列表中的文件
		for i, s := range deviceDownloadList[r.RemoteAddr] {
			if s == fileName {
				deviceDownloadList[r.RemoteAddr] = append(deviceDownloadList[r.RemoteAddr][:i], deviceDownloadList[r.RemoteAddr][i+1:]...)
				break
			}
		}
		// 返回成功信息
		_, _ = fmt.Fprintln(w, "File download success.")
		return
	}

	// 检查文件是否在下载清单内
	if !stringInSlice(fileName, deviceDownloadList[r.RemoteAddr]) {
		_, _ = fmt.Fprintln(w, "File is not in download list.")
		return
	}

	// 将文件内容发送给客户端
	// 设置响应头
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	file, err := os.Open("./download/" + fileName)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}(file)

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// 检查Token
	if r.FormValue("token") != connectedDevices[r.RemoteAddr].Token.String() {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintln(w, "Token is not matched. Please connect again.")
		return
	}

	// 从Query中获取文件名
	fileName, lengthString := r.FormValue("fileName"), r.FormValue("length")
	if fileName == "" {
		_, _ = fmt.Fprintln(w, "File name is required.")
		return
	}
	if lengthString == "" {
		_, _ = fmt.Fprintln(w, "File length is required.")
		return
	}
	// 将length转换为int64
	length, err := strconv.ParseInt(lengthString, 10, 64)
	if err != nil {
		_, _ = fmt.Fprintln(w, "File length is invalid.")
		return
	}

	// 将上传请求加入上传列表
	deviceUploadList[r.RemoteAddr] = append(deviceUploadList[r.RemoteAddr], UploadInfo{
		FileName: fileName,
		Length:   length,
		Accepted: Pending,
	})

	// 返回成功加入上传列表信息
	_, _ = fmt.Fprintln(w, "File upload request is pending.")
	return
}

func transitHandler(w http.ResponseWriter, r *http.Request) {
	// 检查Token
	if r.FormValue("token") != connectedDevices[r.RemoteAddr].Token.String() {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintln(w, "Token is not matched. Please connect again.")
		return
	}

	// 接收文件上传请求
	file, handler, err := r.FormFile("file")
	if err != nil {
		_, _ = fmt.Fprintln(w, "File is required.")
		return
	}

	// 检查文件大小
	if handler.Size > 1024*1024*1024 {
		_, _ = fmt.Fprintln(w, "File size is too large.")
		return
	}

	// 检查文件是否在上传列表内
	var uploadInfo UploadInfo
	for _, info := range deviceUploadList[r.RemoteAddr] {
		if info.FileName == handler.Filename {
			uploadInfo = info
			break
		}
	}
	if uploadInfo.FileName == "" {
		_, _ = fmt.Fprintln(w, "File is not in upload list.")
		return
	}

	// 检查文件长度是否匹配
	if uploadInfo.Length != handler.Size {
		_, _ = fmt.Fprintln(w, "File length is not matched.")
		return
	}

	// 检查是否允许上传
	if uploadInfo.Accepted != Accept {
		if uploadInfo.Accepted == Pending {
			_, _ = fmt.Fprintln(w, "File upload request is not accepted yet.")
		} else {
			_, _ = fmt.Fprintln(w, "File upload request is rejected.")
		}
		return
	}

	// 读取文件，并保存到 ./upload/{远程主机名}/{文件名} 中
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		_, _ = fmt.Fprintln(w, "File read error.")
		return
	}
	// 计算端口号，冒号以后
	port := r.RemoteAddr[strings.Index(r.RemoteAddr, ":")+1:]
	err = os.WriteFile("./upload/"+port+"/"+handler.Filename, fileBytes, 0644)
	if err != nil {
		_, _ = fmt.Fprintln(w, "File write error.")
		return
	}
	return
}

// powerPointHandler 处理PPT操作请求
func powerPointHandler(w http.ResponseWriter, r *http.Request) {
	// 检查Token
	if r.FormValue("token") != connectedDevices[r.RemoteAddr].Token.String() {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintln(w, "Token is not matched. Please connect again.")
		return
	}

	// 检查是否允许操作PPT
	if !connectedDevices[r.RemoteAddr].AllowPowerPoint {
		_, _ = fmt.Fprintln(w, "PowerPoint is not allowed.")
		return
	}

	// 检查是否为切换PPT请求
	if r.FormValue("direction") == "next" {
		switchPowerPoint(Next)
	} else if r.FormValue("direction") == "previous" {
		switchPowerPoint(Previous)
	} else {
		_, _ = fmt.Fprintln(w, "Direction must be \"next\" or \"previous\".")
		return
	}
	return
}

type SwitchDirection int

const (
	Next SwitchDirection = iota
	Previous
)

// switchPowerPoint 切换PPT
func switchPowerPoint(direction SwitchDirection) {
	if direction == Next {
		keyBonding.SetKeys(keybd_event.VK_RIGHT)
	} else {
		keyBonding.SetKeys(keybd_event.VK_LEFT)
	}
	if err := keyBonding.Launching(); err != nil {
		panic(err)
	}
	return
}

func photoHandler(w http.ResponseWriter, r *http.Request) {
	// 检查Token
	if r.FormValue("token") != connectedDevices[r.RemoteAddr].Token.String() {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintln(w, "Token is not matched. Please connect again.")
		return
	}

	if r.FormValue("action") == "path" {
		if r.FormValue("path") == "" {
			_, _ = fmt.Fprintln(w, "Path is required.")
			return
		}
		// 创建新路径
		if err := os.Mkdir("./album/"+r.RemoteAddr+"/"+r.FormValue("path"), os.ModePerm); err != nil {
			_, _ = fmt.Fprintln(w, "OK.")
			return
		}
		_, _ = fmt.Fprintln(w, "OK.")
		return
	}

	if r.FormValue("action") == "file" {
		content, handler, err := r.FormFile("file")
		if err != nil {
			_, _ = fmt.Fprintln(w, "File is required.")
			return
		}
		// 保存文件
		fileBytes, err := io.ReadAll(content)
		if err != nil {
			_, _ = fmt.Fprintln(w, "File read error.")
			return
		}
		err = os.WriteFile("./album/"+r.RemoteAddr+"/"+r.FormValue("path")+"/"+handler.Filename, fileBytes, 0644)
	}
}

func connectionAllowed(device *Device) bool {
	// Pass to frontend to check
	return true
}

func sendFileTo(remoteHostname string, fileName string) {
	// 检查fileName是否已在下载列表中
	if stringInSlice(fileName, deviceDownloadList[remoteHostname]) {
		return
	}
	deviceDownloadList[remoteHostname] = append(deviceDownloadList[remoteHostname], fileName)
}

func SendRingTo(remoteHostname string, tone string, duration int) {
	deviceRingList[remoteHostname] = RingInfo{
		Tone:     tone,
		Duration: duration,
	}
}

func isMobilePhoneConnected(remoteAddr string, timeout float64) (isConnected bool) {
	lastConnection := connectedDevices[remoteAddr].LastConnection
	// 检查是否超时
	return time.Now().Sub(lastConnection).Seconds() < timeout
}

func refreshUploadList() {
	// 让前端刷新上传列表
	return
}

func isPortAvailable(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("Error:", err.Error())
		}
	}(listener)
	return true
}

func findAvailablePort(begin int, end int) int {
	// 获取随机未被占用的端口
	if begin > end {
		return -1
	}
	for port := begin; port <= end; port++ {
		if isPortAvailable(port) {
			return port
		}
	}
	return -1
}

func cursorTextHandler(w http.ResponseWriter, r *http.Request) {
	// 检查Token
	if r.FormValue("token") != connectedDevices[r.RemoteAddr].Token.String() {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintln(w, "Token is not matched. Please connect again.")
		return
	}

	// 检查是否允许发送光标文本
	if !connectedDevices[r.RemoteAddr].AllowCursorText {
		_, _ = fmt.Fprintln(w, "Cursor text is not allowed.")
		return
	}

	// 从Query中获取光标文本
	cursorText := r.FormValue("cursorText")
	if cursorText == "" {
		_, _ = fmt.Fprintln(w, "Cursor text is required.")
		return
	}

	// 在服务端机器上显示光标文本
	setCursorText(cursorText, 5*time.Second)
	_, _ = fmt.Fprintln(w, "Successfully set cursor text: "+cursorText)
	return
}

// setCursorText 在服务端机器上显示光标文本
func setCursorText(cursorText string, duration time.Duration) {
	// 在服务端机器上显示光标文本
	notification := toast.Notification{
		AppID:   "MobiConn",
		Title:   "MobiConn",
		Message: cursorText,
	}
	err := notification.Push()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	return
}

// stringInSlice 检查字符串是否在字符串切片中
func stringInSlice(str string, list []string) bool {
	for _, s := range list {
		if str == s {
			return true
		}
	}
	return false
}

func powerPointDemoHandler(w http.ResponseWriter, r *http.Request) {
	for {
		switchPowerPoint(Next)
		time.Sleep(5 * time.Second)
	}
	return
}
