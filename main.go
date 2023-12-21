package main

import (
	"embed"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-toast/toast"
	"github.com/micmonay/keybd_event"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// 保存当前已连接的设备
var connectedDevices = map[string]Device{}

// 保存每个设备的照片，键为设备主机地址，值为照片文件Base64编码
var devicePhotos = map[string][]string{
	"localhost": {"(Base64 encoded photo 1)", "(Base64 encoded photo 2)", "(Base64 encoded photo 3)"},
}

// 维护每个设备的下载列表
var deviceDownloadList = map[string][]string{}

//// 维护每个设备的上传列表
//var deviceUploadList = map[string][]UploadInfo{}

// 记录将要响铃的设备主机和铃声及响铃持续秒数
var deviceRingList = map[string]RingInfo{}

// 记录每个设备需要震动的时长
var deviceVibrateList = map[string]int{}

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

	k = 3

	// 绑定处理函数到路由路径
	mux.HandleFunc("/greeting", greetingHandler)
	//mux.HandleFunc("/heartbeat", heartbeatHandler)
	mux.HandleFunc("/photo/upload", uploadPhotosHandler)
	mux.HandleFunc("/ppt", powerPointHandler)

	//mux.HandleFunc("/download", downloadHandler)
	//mux.HandleFunc("/upload", uploadHandler)
	//mux.HandleFunc("/transit", transitHandler)
	//mux.HandleFunc("/cursorText", cursorTextHandler)
	//mux.HandleFunc("/photo", photoHandler)
	//
	//mux.HandleFunc("/powerPointDemo", powerPointDemoHandler)

	// 启动HTTP服务器
	go func() {
		err := http.ListenAndServe(":25236", mux)
		if err != nil {
			fmt.Println("HTTP server error:", err)
		}
	}()

	//// 创建ServeMux实例，以实现与Vue前端的交互
	//frontendMux := http.NewServeMux()
	//
	//// 绑定处理函数到路由路径
	//frontendMux.HandleFunc("/getAlbumCount", frontendGetAlbumCountHandler)
	//
	//// 启动HTTP服务器
	//go func() {
	//	err := http.ListenAndServe(":25237", frontendMux)
	//	if err != nil {
	//		fmt.Println("HTTP server error:", err)
	//	}
	//}()

	return nil
}

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	message := "Hello World!"
	_, _ = fmt.Fprintln(w, message)
}

//func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
//	// 检查客户端是否已连接过
//	if _, ok := connectedDevices[r.RemoteAddr]; !ok {
//		_, _ = fmt.Fprintln(w, "You have never connected to this server yet. Please connect first.")
//		return
//	}
//
//	// 检查Token是否匹配
//	if r.FormValue("token") != connectedDevices[r.RemoteAddr].Token.String() {
//		_, _ = fmt.Fprintln(w, "Token is not matched. Please connect again.")
//		return
//	}
//
//	// 更新最后连接时间
//	device := connectedDevices[r.RemoteAddr]
//	device.LastConnection = time.Now()
//	connectedDevices[r.RemoteAddr] = device
//
//	var goals []any
//
//	for count := albumCount[r.RemoteAddr]; count != -1; {
//		goal := map[string]any{
//			"Action": "albumCount",
//		}
//		goals = append(goals, goal)
//	}
//
//	// 检查是否有文件需要下载
//	for _, s := range deviceDownloadList[r.RemoteAddr] {
//		fileInfo, err := os.Stat("./download/" + s)
//		if err != nil {
//			// 文件不存在，输出控制台警告
//			fmt.Println("Error:", err.Error())
//			continue
//		}
//		goal := map[string]any{
//			"Action": "download",
//			"Information": DownloadInfo{
//				FileName: s,
//				Length:   fileInfo.Size(),
//			},
//		}
//		goals = append(goals, goal)
//	}
//	// 检查响铃任务信息
//	if ringInfo, ok := deviceRingList[r.RemoteAddr]; ok {
//		goal := map[string]any{
//			"Action":      "ring",
//			"Information": ringInfo,
//		}
//		goals = append(goals, goal)
//	}
//
//	// 检查是否有文件需要上传
//	for _, uploadInfo := range deviceUploadList[r.RemoteAddr] {
//		if uploadInfo.Accepted == Pending {
//			continue
//		} else if uploadInfo.Accepted == Reject {
//			goal := map[string]any{
//				"Action": "upload",
//				"Information": map[string]any{
//					"FileName": uploadInfo.FileName,
//					"Status":   "rejected",
//				},
//			}
//			goals = append(goals, goal)
//			continue
//		} else if uploadInfo.Accepted == Accept {
//			goal := map[string]any{
//				"Action": "upload",
//				"Information": map[string]any{
//					"FileName": uploadInfo.FileName,
//					"Status":   "accepted",
//				},
//			}
//			goals = append(goals, goal)
//		}
//
//		// 从上传列表中删除该文件
//		for i, info := range deviceUploadList[r.RemoteAddr] {
//			if info.FileName == uploadInfo.FileName {
//				deviceUploadList[r.RemoteAddr] = append(deviceUploadList[r.RemoteAddr][:i], deviceUploadList[r.RemoteAddr][i+1:]...)
//				break
//			}
//		}
//		continue
//	}
//
//	// 检查是否有相册需要读取
//	for _, path := range deviceAlbumFetchPaths[r.RemoteAddr] {
//		goal := map[string]any{
//			"Action": "album",
//			"Path":   path,
//		}
//		goals = append(goals, goal)
//	}
//
//	// 检查是否有相册文件需要读取
//	for _, file := range deviceAlbumFetchFiles[r.RemoteAddr] {
//		goal := map[string]any{
//			"Action": "photo",
//			"File":   file,
//		}
//		goals = append(goals, goal)
//	}
//
//	// 将goals加入响应
//	response := map[string]any{
//		"Status": "success",
//	}
//	if len(goals) != 0 {
//		response["Goals"] = goals
//	}
//
//	// 将响应转换为JSON
//	responseJson, err := json.Marshal(response)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// Set response header
//	w.Header().Set("Content-Type", "application/json")
//	// Send response
//	_, _ = fmt.Fprintln(w, string(responseJson))
//	//}
//}

func uploadPhotosHandler(w http.ResponseWriter, r *http.Request) {
	indexString := r.FormValue("index")
	if indexString == "" {
		_, _ = fmt.Fprintln(w, "Photo index is required.")
		return
	}
	index, err := strconv.Atoi(indexString)
	if err != nil {
		_, _ = fmt.Fprintln(w, "Photo index is invalid.")
		return
	}
	photoFile, photoHeader, err := r.FormFile("photo")
	if err != nil {
		_, _ = fmt.Fprintln(w, "Photo is required.")
		return
	}
	photoFileContent := make([]byte, photoHeader.Size)
	n, err := photoFile.Read(photoFileContent)
	if err != nil {
		_, _ = fmt.Fprintln(w, "Photo read error.")
		return
	}
	if n != int(photoHeader.Size) {
		_, _ = fmt.Fprintln(w, "Photo read error.")
		return
	}
	devicePhotos[r.RemoteAddr][index] = string(photoFileContent)
	_, _ = fmt.Fprintln(w, "OK.")
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

//func frontendGetAlbumCountHandler(w http.ResponseWriter, r *http.Request) {
//	remoteAddr := r.FormValue("RemoteAddress")
//	if remoteAddr == "" {
//		_, _ = fmt.Fprintln(w, "Remote address is required.")
//		return
//	}
//	albumCount[remoteAddr] = -1
//	for albumCount[remoteAddr] == -1 {
//		time.Sleep(100 * time.Millisecond)
//	}
//	_, _ = fmt.Fprintln(w, albumCount[remoteAddr])
//}

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
