package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

var k = 0

func (a *App) GetAndAddOne() int {
	k++
	return k
}

type DeviceView struct {
	Address   string `json:"address"`
	Connected bool   `json:"connected"`
}

//func NewMyStruct() *DeviceView {
//	return &DeviceView{
//		people: People{
//			List: []Person{
//				{Name: "John", Age: 25},
//				{Name: "Jane", Age: 30},
//			},
//		},
//	}
//}

const timeout = 5000

func (a *App) GetDevices() map[string]DeviceView {
	result := map[string]DeviceView{}
	result["localhost"] = DeviceView{
		Address:   "localhost",
		Connected: true,
	}
	for k, v := range connectedDevices {
		result[k] = DeviceView{
			Address:   k,
			Connected: time.Now().Sub(v.LastConnection) < timeout*time.Millisecond,
		}
	}
	return result
}

func (a *App) OpenAndSendFile(remoteAddr string) {
	dialogOptions := runtime.OpenDialogOptions{
		Title: "选择文件",
		Filters: []runtime.FileFilter{
			//{DisplayName: "文本文件", Pattern: "*.txt"},
			//{DisplayName: "图片文件", Pattern: "*.jpg;*.png"},
			{DisplayName: "所有文件", Pattern: "*.*"},
		},
	}

	selectedFile, err := runtime.OpenFileDialog(a.ctx, dialogOptions)
	if err != nil {
		// 处理错误
		return
	}
	if selectedFile == "" {
		// 用户取消了选择
		return
	}

	// 处理所选文件
	fmt.Println("向", remoteAddr, "发送文件", selectedFile)

	if stringInSlice(selectedFile, deviceDownloadList[remoteAddr]) {
		return
	}
	deviceDownloadList[remoteAddr] = append(deviceDownloadList[remoteAddr], selectedFile)
}

func (a *App) SendRing(remoteAddr string, tone string, duration int) {
	deviceRingList[remoteAddr] = RingInfo{
		Tone:     tone,
		Duration: duration,
	}
	fmt.Println("向", remoteAddr, "发送铃声", tone, duration, "ms")
}

func (a *App) SendVibrate(remoteAddr string, duration int) {
	deviceVibrateList[remoteAddr] = duration
	fmt.Println("向", remoteAddr, "发送震动", duration, "ms")
}

func (a *App) GetPhotoCount(remoteAddr string) int {
	return len(devicePhotos[remoteAddr])
}

func (a *App) GetPhoto(remoteAddr string, photoIndex int) string {
	return devicePhotos[remoteAddr][photoIndex]
}
