// Mandelbrot emits a PNG image of the Mandelbrot fractal.
// See page 61.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "web" {
		handler := func(w http.ResponseWriter, r *http.Request) {
			// 设置Content-Type头部为"image/svg+xml"，表示响应的内容类型为SVG图像
			genImg(w)
			w.Header().Set("Content-Type", "image/png")
		}
		fmt.Println("start web server http://localhost:8000...")
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
		return

	}
	genImg(os.Stdout)

}

func genImg(out io.Writer) {
	// 定义图像的边界和尺寸
	const (
		// 图像 x 轴的最小值
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		// 图像的宽度和高度（像素）
		width, height = 1024, 1024
	)

	// 创建一个新的 RGBA 图像
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 遍历图像的每个像素
	for py := 0; py < height; py++ {
		// 将像素的 y 坐标转换为对应的复数的实部
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			// 将像素的 x 坐标转换为对应的复数的虚部
			x := float64(px)/width*(xmax-xmin) + xmin
			// 构建复数
			z := complex(x, y)
			// 将像素 (px, py) 设置为对应复数 z 的 Mandelbrot 颜色
			img.Set(px, py, mandelbrot(z))
		}
	}

	// 将图像编码为 PNG 格式并输出到标准输出流
	png.Encode(out, img) // NOTE: ignoring errors

}

// mandelbrot 函数用于计算并返回 Mandelbrot 集合中对应于复数 z 的颜色
func mandelbrot(z complex128) color.Color {
	// 定义迭代次数和对比度
	const iterations = 200
	const contrast = 15

	// 初始化变量 v 为 0
	var v complex128
	// 迭代计算 v 的值
	for n := uint8(0); n < iterations; n++ {
		// 更新 v 的值，v 的平方加上输入的复数 z
		v = v*v + z
		// 如果 v 的模长大于 2，则认为该点逃逸出 Mandelbrot 集合
		if cmplx.Abs(v) > 2 {
			// 返回一个灰色值，其亮度与迭代次数成反比
			return color.Gray{255 - contrast*n}
		}
	}
	// 如果迭代结束后 v 的模长仍未超过 2，则认为该点属于 Mandelbrot 集合，返回黑色
	return color.Black
}

// Some other interesting functions:

// acos 函数计算复数 z 的反余弦值，并将结果映射到颜色空间
func acos(z complex128) color.Color {
	// 计算 z 的反余弦值
	v := cmplx.Acos(z)
	// 将反余弦值的实部映射到蓝色通道
	blue := uint8(real(v)*128) + 127
	// 将反余弦值的虚部映射到红色通道
	red := uint8(imag(v)*128) + 127
	// 返回一个 YCbCr 颜色，其中 Y 分量设置为 192，Cb 和 Cr 分量分别由蓝色和红色值确定
	return color.YCbCr{192, blue, red}
}

// sqrt 函数计算复数 z 的平方根，并将结果映射到颜色空间
func sqrt(z complex128) color.Color {
	// 计算 z 的平方根
	v := cmplx.Sqrt(z)
	// 将平方根值的实部映射到蓝色通道
	blue := uint8(real(v)*128) + 127
	// 将平方根值的虚部映射到红色通道
	red := uint8(imag(v)*128) + 127
	// 返回一个 YCbCr 颜色，其中 Y 分量设置为 128，Cb 和 Cr 分量分别由蓝色和红色值确定
	return color.YCbCr{128, blue, red}
}

// f(x) = x^4 - 1
//
// z' = z - f(z)/f'(z)
//
//	= z - (z^4 - 1) / (4 * z^3)
//	= z - (z - 1/z^3) / 4
//
// newton 函数使用牛顿迭代法来逼近方程 f(x) = x^4 - 1 的根，并将迭代次数映射到颜色空间
func newton(z complex128) color.Color {
	// 定义迭代次数和对比度
	const iterations = 37
	const contrast = 7
	// 迭代更新 z 的值
	for i := uint8(0); i < iterations; i++ {
		// 使用牛顿迭代公式更新 z 的值
		z -= (z - 1/(z*z*z)) / 4
		// 如果 z 的四次方与 1 的差的绝对值小于 1e-6，则认为找到了一个根
		if cmplx.Abs(z*z*z*z-1) < 1e-6 {
			// 返回一个灰色值，其亮度与迭代次数成反比
			return color.Gray{255 - contrast*i}
		}
	}
	// 如果没有找到根，则返回黑色
	return color.Black
}
