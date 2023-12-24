package main

import (
	"bufio"
	"crypto/md5"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/google/uuid"
	"github.com/micmonay/keybd_event"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

//go:embed all:frontend/dist
var assets embed.FS

var app *App

// 保存当前已连接的设备
var connectedDevices = map[string]Device{}

// 保存每个设备的照片，键为设备主机地址，值为照片文件Base64编码
var devicePhotos = map[string][]string{
	"localhost": {},
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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
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

func loadLocalPhotos() {
	_, err := os.Stat("./local-photos")
	if err != nil {
		err = os.Mkdir("./local-photos", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	files, err := os.ReadDir("./local-photos")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		devicePhotos["localhost"] = append(devicePhotos["localhost"], loadBase64Picture("./local-photos/"+file.Name()))
	}
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
	// 加载本地照片
	loadLocalPhotos()

	// Create an instance of the app structure
	app = NewApp()

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

	// 初始化键盘事件
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
	mux.HandleFunc("/cursor-text", cursorTextHandler)

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
func greetingHandler(w http.ResponseWriter, _ *http.Request) {
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
	//goland:noinspection GoPreferNilSlice
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
		goals = append(goals, Goal{Action: "ring", Information: strconv.Itoa(ringInfo.Duration)})
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
	// 检查文件存在
	fileInfo, err := os.Stat(fileName)
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
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	readingFile, err := os.Open(fileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, "File open error.")
		return
	}
	defer func(readingFile *os.File) {
		err := readingFile.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, "File close error.")
			return
		}
	}(readingFile)
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
	host := addrToHost(r.RemoteAddr)
	fileName := r.FormValue("fileName")
	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "File name is required.")
		return
	}
	// 计算远程主机地址的MD5值
	remoteAddrMD5 := md5.Sum([]byte(host))
	// 转换为无横杠十六进制
	remoteAddrMD5String := fmt.Sprintf("%x", remoteAddrMD5)
	// 创建目录
	_, err := os.Stat("./upload")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("./upload", os.ModePerm)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintln(w, "File upload error.")
				return
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, "File upload error.")
			return
		}
	}
	_, err = os.Stat("./upload/" + remoteAddrMD5String)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("./upload/"+remoteAddrMD5String, os.ModePerm)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintln(w, "File upload error.")
				return
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, "File upload error.")
			return
		}
	}
	// 写入文件
	writingFile, err := os.OpenFile("./upload/"+remoteAddrMD5String+"/"+fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, "File save error.")
		return
	}
	buffer := make([]byte, 1024*1024)
	// 类型：application/octet-stream，大小：1GB
	for {
		n, err := r.Body.Read(buffer)
		if err == io.EOF {
			_, _ = writingFile.Write(buffer[:n])
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
		_, err = writingFile.Write(buffer[:n])
	}
	err = writingFile.Sync()
	if err != nil {
		return
	}
	// 关闭文件
	err = writingFile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, "File save error.")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println(host + "已上传文件：" + fileName)
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
	//w.WriteHeader(http.StatusOK)
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

	//goland:noinspection GoPreferNilSlice
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
	cursorText := r.FormValue("text")
	if cursorText == "" {
		_, _ = fmt.Fprintln(w, "Cursor text is required.")
		return
	}

	setCursorText(cursorText)
	fmt.Println("接收到光标文本：" + cursorText)
}

// setCursorText 在服务端机器上显示光标文本
func setCursorText(cursorText string) {
	robotgo.TypeStr(cursorText)
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
