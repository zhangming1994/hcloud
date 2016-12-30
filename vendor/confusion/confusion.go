package confusion

import (
	// "fmt"
	// "github.com/456vv/verifycode"
	"github.com/astaxie/beego"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	// "strconv"
	"github.com/nfnt/resize"
	"image/color"
	"math/rand"
	"strings"
	"time"
)

//生成工号图片
// func createVerify(uname string, width int, height int, validatepath string) (fileName string) {
// 	//验证码颜色
// 	c := []string{"#ff8080FF", "#00ff0000", "#8080c0FD"}
// 	colors, err := verifycode.NewColor(c)
// 	if err != nil {
// 		fmt.Println("NewColor: %v", err)
// 		os.Exit(-1)
// 	}
// 	//验证码背景
// 	b := []string{"#804040FF"}
// 	backgrounds, err := verifycode.NewColor(b)
// 	if err != nil {
// 		fmt.Println("NewColor: %v", err)
// 		os.Exit(-1)
// 	}
// 	//字体
// 	f := []string{"./static/fonts/simsunb.ttf"}
// 	fonts, err := verifycode.NewFont(f)
// 	if err != nil {
// 		fmt.Println("NewFont: %v", err)
// 		os.Exit(-1)
// 	}
// 	verifyCode := verifycode.NewVerifyCode()
// 	verifyCode.SetDPI(72) //也可以不用设置这个
// 	verifyCode.SetColor(colors)
// 	verifyCode.SetBackground(backgrounds)
// 	verifyCode.SetFont(fonts)
// 	verifyCode.SetWidthWithHeight(width, height) // 宽，高
// 	verifyCode.SetFontSize(float64(height / 2))
// 	verifyCode.SetHinting(false) //也可以不用设置这个
// 	verifyCode.SetKerning(-5, 5) //随机字距，最小-100，最大100

// 	fileName = validatepath + uname + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".png"

// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		fmt.Println("创建文件出错 %v", err)
// 		os.Exit(-1)
// 	}
// 	err = verifyCode.PNG(uname, file)
// 	if err != nil {
// 		fmt.Println("生成验证码出错 %v", err)
// 		os.Exit(-1)
// 	}
// 	return fileName
// }

// //混淆图片 srcImg原始图片 dstImg新生成发图片 uname工号 validatepath验证码图片
// func MixupImg(srcImg string, dstImg string, uname string, validatepath string) {
// 	//原始图片
// 	imgb, err := os.Open(srcImg)
// 	fmt.Println(err)
// 	var (
// 		img image.Image
// 	)
// 	lastIndex := strings.LastIndex(srcImg, ".")
// 	ext := srcImg[lastIndex+1:]

// 	if ext == "png" {
// 		img, err = png.Decode(imgb)
// 		if err != nil {
// 			fmt.Println("the .png ext file is error")
// 		}
// 	} else if ext == "jpg" || ext == "jpeg" {
// 		img, _ = jpeg.Decode(imgb)
// 	}
// 	defer imgb.Close()

// 	//工号图片
// 	verifyImgName := createVerify(uname, img.Bounds().Dx(), img.Bounds().Dy(), validatepath)

// 	wmb, _ := os.Open(verifyImgName)
// 	watermark, _ := png.Decode(wmb)
// 	defer wmb.Close()

// 	b := img.Bounds()
// 	m := image.NewNRGBA(b)

// 	draw.Draw(m, b, watermark, image.ZP, draw.Src)
// 	draw.Draw(m, b, img, image.ZP, draw.Over)

// 	//生成新图片new.jpg，并设置图片质量..
// 	imgw, _ := os.Create(dstImg)
// 	jpeg.Encode(imgw, m, &jpeg.Options{100})

// 	defer imgw.Close()

// 	fmt.Println("create the image success" + dstImg)
// }

func MixUp(srcImg string, dstImg string) error {
	//原始图片
	imgb, err := os.Open(srcImg)
	if err != nil {
		beego.Error(err)
		return err
	}
	var (
		img image.Image
	)
	lastIndex := strings.LastIndex(srcImg, ".")
	ext := srcImg[lastIndex+1:]

	if ext == "png" {
		img, err = png.Decode(imgb)
		if err != nil {
			beego.Error("the .png ext file is error:" + err.Error())
			return err
		}
	} else if ext == "jpg" || ext == "jpeg" {
		img, err = jpeg.Decode(imgb)
		if err != nil {
			beego.Error(err)
			return err
		}
	}
	defer imgb.Close()

	b := img.Bounds()
	m := image.NewNRGBA(b)

	// DisturbBitmap(m)

	draw.Draw(m, b, img, image.ZP, draw.Src)

	//生成新图片new.jpg，并设置图片质量..
	imgw, err := os.Create(dstImg)
	if err != nil {
		beego.Error(err)
		return err
	}
	jpeg.Encode(imgw, m, &jpeg.Options{100})

	defer imgw.Close()
	return nil
}

//绘制干扰背景
func DisturbBitmap(img *image.NRGBA) {
	r := rand.New(rand.NewSource(int64(time.Now().Second())))
	imagexx := img.Rect.Max.X
	imageyy := img.Rect.Max.Y

	for i := 0; i < imagexx; i += 50 {
		for j := 0; j < imageyy; j += 50 {
			a := r.Intn(imagexx)
			b := r.Intn(imageyy)
			c := color.NRGBA{uint8(i), uint8(j), 0, 255}
			img.Set(a, b, c)
		}
	}
}

// 执行操作
func Imagepro(path string, newpath string) error {
	lastIndex := strings.LastIndex(path, ".")
	ext := path[lastIndex+1:]
	if ext == "png" {
		err := Imag_thumbpng(path, 172, 129, newpath)
		return err
	} else if ext == "jpg" || ext == "jpeg" {
		err := Imag_thumbjpg(path, 172, 129, newpath)
		return err
	}
	return nil
}

//jpge小图
func Imag_thumbjpg(file string, width uint, height uint, to string) error {
	// 打开图片并解码
	file_origin, _ := os.Open(file)
	origin, _ := jpeg.Decode(file_origin)
	defer file_origin.Close()
	canvas := resize.Resize(width, height, origin, resize.Lanczos3)
	file_out, err := os.Create(to)
	if err != nil {
		beego.Error(err)
		return err
	}
	defer file_out.Close()
	err = jpeg.Encode(file_out, canvas, &jpeg.Options{80})
	return err
}

// png小图
func Imag_thumbpng(file string, width uint, height uint, to string) error {
	// 打开图片并解码
	file_origin, _ := os.Open(file)
	origin, _ := png.Decode(file_origin)
	defer file_origin.Close()
	canvas := resize.Resize(width, height, origin, resize.Lanczos3)
	file_out, err := os.Create(to)
	if err != nil {
		beego.Error(err)
		return err
	}
	defer file_out.Close()
	err = png.Encode(file_out, canvas)
	return err
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) (bool, error) {
	_, err := os.Stat(filename)
	boo := (err == nil || os.IsExist(err))
	return boo, err
}

// //测试
// func main() {
// 	mixupImg("123.jpg", "mix.jpg", "9999")
// }

// func main() {
// 	os.Mkdir("uploadfile", os.ModePerm)
// 	Imagepro("upload/", "uploadfile/")
// 	beego.Run()
// }
