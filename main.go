package main

import (
	"bufio"
	"crypto/md5"
	"embed"
	"encoding/base64"
	"encoding/json"
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

// RingInfo 表示响铃任务的信息。
type RingInfo struct {
	Tone     string // 铃声文件名
	Duration int    // 响铃持续秒数
}

var keyBonding keybd_event.KeyBonding

func loadBase64Picture(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	fileContent := make([]byte, fileInfo.Size())
	n, err := file.Read(fileContent)
	if err != nil {
		panic(err)
	}
	if n != int(fileInfo.Size()) {
		panic("Read file error.")
	}
	return base64.StdEncoding.EncodeToString(fileContent)
}

func addrToHost(addr string) string {
	// 找到冒号
	colonIndex := -1
	for i, c := range addr {
		if c == ':' {
			colonIndex = i
			break
		}
	}
	if colonIndex == -1 {
		return addr
	}
	return addr[:colonIndex]
}

func main() {
	//base64Picture := loadBase64Picture("C:\\Users\\ab123\\Pictures\\艾丝妲\\asta.png")
	//fmt.Println(base64Picture)
	//devicePhotos["localhost"] = append(devicePhotos["localhost"], base64Picture)
	//base64Picture = loadBase64Picture("C:\\Users\\ab123\\Pictures\\个人头像\\美乐蒂玩电脑.jpg")
	//fmt.Println(base64Picture)
	//devicePhotos["localhost"] = append(devicePhotos["localhost"], base64Picture)

	// Create an instance of the app structure
	app := NewApp()

	// Create http server
	err := httpServer()
	if err != nil {
		println("Error:", err.Error())
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "MobiConn",
		Width:  1200,
		Height: 1000,
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
	mux.HandleFunc("/heartbeat", heartbeatHandler)
	mux.HandleFunc("/file/download", fileDownloadHandler)
	mux.HandleFunc("/file/upload", fileUploadHandler)
	mux.HandleFunc("/photo/count", photoCountHandler)
	mux.HandleFunc("/photo/upload", photoUploadHandler)
	mux.HandleFunc("/ppt", powerPointHandler)
	mux.HandleFunc("/cursorText", cursorTextHandler)

	//mux.HandleFunc("/download", downloadHandler)
	//mux.HandleFunc("/upload", uploadHandler)
	//mux.HandleFunc("/transit", transitHandler)
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

	return nil
}

// greetingHandler 处理问候请求
func greetingHandler(w http.ResponseWriter, r *http.Request) {
	message := "Hello World!"
	_, _ = fmt.Fprintln(w, message)
}

type Goal struct {
	Action      string `json:"action"`
	Information string `json:"information"`
}

// heartbeatHandler 处理心跳请求
func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	device := connectedDevices[addrToHost(r.RemoteAddr)]
	device.LastConnection = time.Now()
	connectedDevices[addrToHost(r.RemoteAddr)] = device

	//goals := []Goal{{Action: "123", Information: "456"}}
	goals := []Goal{}

	// 文件下载
	for _, fileName := range deviceDownloadList[addrToHost(r.RemoteAddr)] {
		_, err := os.Stat(fileName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintln(w, "File not found.")
			continue
		}
		goals = append(goals, Goal{Action: "download", Information: fileName})
		delete(deviceDownloadList, addrToHost(r.RemoteAddr))
	}

	// 响铃
	if ringInfo, ok := deviceRingList[addrToHost(r.RemoteAddr)]; ok {
		goals = append(goals, Goal{Action: "ring", Information: ringInfo.Tone})
		delete(deviceRingList, addrToHost(r.RemoteAddr))
	}

	// 震动
	if vibrateDuration, ok := deviceVibrateList[addrToHost(r.RemoteAddr)]; ok {
		goals = append(goals, Goal{Action: "vibration", Information: strconv.Itoa(vibrateDuration)})
		delete(deviceVibrateList, addrToHost(r.RemoteAddr))
	}

	// 返回
	marshal, err := json.Marshal(goals)
	if err != nil {
		_, _ = fmt.Fprintln(w, "JSON marshal error.")
		return
	}
	_, _ = fmt.Fprintln(w, string(marshal))
}

// downloadHandler 处理文件下载请求
func fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.FormValue("fileName")
	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "File name is required.")
		return
	}
	//if !stringInSlice(fileName, deviceDownloadList[addrToHost(r.RemoteAddr)]) {
	//	w.WriteHeader(http.StatusNotFound)
	//	_, _ = fmt.Fprintln(w, "File is not in download list.")
	//	return
	//}
	// 检查文件存在
	_, err := os.Stat(fileName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintln(w, "File not found.")
		array := deviceDownloadList[addrToHost(r.RemoteAddr)]
		// remove if value = fileName
		for i, v := range array {
			if v == fileName {
				array = append(array[:i], array[i+1:]...)
				break
			}
		}
		deviceDownloadList[addrToHost(r.RemoteAddr)] = array
		return
	}
	readingFile, err := os.Open(fileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, "File open error.")
		return
	}
	defer readingFile.Close()
	reader := bufio.NewReader(readingFile)
	buffer := make([]byte, 1024*1024)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			_, _ = w.Write(buffer[:n])
			break
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, "File read error.")
			return
		}
		if n == 0 {
			break
		}
		_, err = w.Write(buffer[:n])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, "File write error.")
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// fileUploadHandler 处理文件上传请求
func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 计算远程主机地址的MD5值
	remoteAddrMD5 := md5.Sum([]byte(addrToHost(r.RemoteAddr)))
	// 转换为无横杠十六进制
	remoteAddrMD5String := fmt.Sprintf("%x", remoteAddrMD5)
	// 创建UploadFiles目录
	if err := os.Mkdir("./uploadFiles/"+remoteAddrMD5String, os.ModePerm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, "File upload error.")
		return
	}
	formFile, fileHeader, err := r.FormFile("formFile")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "File is required.")
		return
	}

	// 打开并写入文件
	writingFile, err := os.OpenFile("./uploadFiles/"+remoteAddrMD5String+"/"+fileHeader.Filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, "File save error.")
		return
	}
	defer writingFile.Close()
	writer := bufio.NewWriter(writingFile)
	buffer := make([]byte, 1024*1024)
	for {
		n, err := formFile.Read(buffer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, "File read error.")
			return
		}
		if n == 0 {
			break
		}
		_, err = writer.Write(buffer[:n])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, "File save error.")
			return
		}
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// photoCountHandler 处理照片数量设置请求
func photoCountHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "Photo count handler")
	countString := r.FormValue("count")
	if countString == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "Photo count is required.")
		return
	}
	count, err := strconv.Atoi(countString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "Photo count is invalid.")
		return
	}
	devicePhotos[addrToHost(r.RemoteAddr)] = make([]string, count)
	w.WriteHeader(http.StatusOK)
	host := addrToHost(r.RemoteAddr)
	_, _ = fmt.Fprintln(w, host+"的照片数量已设置为"+countString)
}

// photoUploadHandler 处理照片上传请求
func photoUploadHandler(w http.ResponseWriter, r *http.Request) {
	indexString := r.FormValue("index")
	if indexString == "" {
		//w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "Index is required.")
		return
	}
	index, err := strconv.Atoi(indexString)
	if err != nil {
		//w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "Index is invalid.")
		return
	}

	// 读取照片文件
	// r.Body.Read()

	photo := []byte{}
	for {
		photoFileContent := make([]byte, 1024*1024)
		n, err := r.Body.Read(photoFileContent)
		if err == io.EOF {
			photo = append(photo, photoFileContent[:n]...)
			break
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, "File read error.")
			return
		}
		photo = append(photo, photoFileContent[:n]...)
	}
	encoded := base64.StdEncoding.EncodeToString(photo)
	// base64加密
	devicePhotos[addrToHost(r.RemoteAddr)][index] = encoded
	//devicePhotos[addrToHost(r.RemoteAddr)][index] = photo
	//w.WriteHeader(http.StatusOK)
	host := addrToHost(r.RemoteAddr)
	_, _ = fmt.Fprintln(w, host+"的第"+indexString+"张照片已上传")
	_, _ = fmt.Fprintln(w, "照片内容为："+encoded)
}

// powerPointHandler 处理PPT操作请求
func powerPointHandler(w http.ResponseWriter, r *http.Request) {
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

// cursorTextHandler 处理光标文本设置请求
func cursorTextHandler(w http.ResponseWriter, r *http.Request) {
	//// 检查Token
	//if r.FormValue("token") != connectedDevices[addrToHost(r.RemoteAddr)].Token.String() {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	_, _ = fmt.Fprintln(w, "Token is not matched. Please connect again.")
	//	return
	//}

	//// 检查是否允许发送光标文本
	//if !connectedDevices[addrToHost(r.RemoteAddr)].AllowCursorText {
	//	_, _ = fmt.Fprintln(w, "Cursor text is not allowed.")
	//	return
	//}

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

//func powerPointDemoHandler(w http.ResponseWriter, r *http.Request) {
//	for {
//		switchPowerPoint(Next)
//		time.Sleep(5 * time.Second)
//	}
//	return
//}

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
