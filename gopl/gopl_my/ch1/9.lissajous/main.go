// Lissajous generates GIF animations of random Lissajous figures.
// Run with "web" command-line argument for web server.
// See page 13.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// const声明
const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

// 复合声明
var palette = []color.Color{color.White, color.Black}

func main() {
	// The sequence of images is deterministic unless we seed
	// the pseudo-random number generator using the current time.
	// Thanks to Randall McPherson for pointing out the omission.
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) > 1 && os.Args[1] == "web" {
		// 定义了一个 HTTP 请求处理的匿名函数 handler
		handler := func(w http.ResponseWriter, r *http.Request) {
			lissajous(w)
		}
		// 注册一个 HTTP 处理函数，当请求路径为 "/" 时，调用 handler 函数来处理请求
		http.HandleFunc("/", handler)
		fmt.Println("http://localhost:8000/")
		// 启动一个 HTTP 服务器，监听在本地主机的 8000 端口上。log.Fatal 会在服务器启动失败时(如端口已被占用)输出错误信息并结束程序。
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
		return
	}
	// 默认输出到标准输出
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // x轴旋转5圈
		res     = 0.001 // 角度分辨率，即每一步的角度增量0.01
		size    = 100   // 图像画布的大小， [-100, +100]
		nframes = 64    // 动画的帧数 64帧
		delay   = 8     // 帧之间的延迟8ms
	)
	freq := rand.Float64() * 3.0        // y轴振荡器的相对频率 [0.0, 3.0)
	anim := gif.GIF{LoopCount: nframes} // 声明一个gif.GIF类型的结构体anim，设置其 LoopCount 属性为循环64帧，其他属性为默认的零值
	phase := 0.0                        // 两个振荡器之间的相位差
	for i := 0; i < nframes; i++ {      // 处理动画的每一帧
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)   // 创建一个矩形区域 rect，其大小为 2*size+1
		img := image.NewPaletted(rect, palette)        // 使用这个矩形创建一个新的调色板图像 img，使用常量声明的调色板 palette
		for t := 0.0; t < cycles*2*math.Pi; t += res { // 处理图像中的每一个像素，计算对应的x和y坐标
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				blackIndex) // 将x、y坐标映射到图像的像素位置上，设置为黑色
		}
		phase += 0.1                           // 每帧结束后，相位差会增加 0.1，从而产生动画效果
		anim.Delay = append(anim.Delay, delay) // 将延迟时间添加到 anim 的 Delay 字段中
		anim.Image = append(anim.Image, img)   // 将当前帧的图像添加到 anim 的 Imag切片字段中
	}
	// NOTE: ignoring encoding errors
	gif.EncodeAll(out, &anim) // 将 anim结构体编码为 GIF 格式，并通过传入的 out 参数将数据写入到输出流中。
}
