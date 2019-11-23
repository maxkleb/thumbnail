package transform

import (
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
)

// Crop and add padding if needed
func ProcessImg(x int, y int, img image.Image) image.Image {
	requiredAspectRatio := float64(y)/float64(x)
	imgX := img.Bounds().Dx()
	imgY := img.Bounds().Dy()
	imageAspectRatio := float64(imgY)/float64(imgX)
	if requiredAspectRatio == imageAspectRatio && x == imgX {
		return img
	} else if requiredAspectRatio == imageAspectRatio && x < imgX{
		return resize.Resize(uint(x), uint(y), img, resize.Lanczos3)
	} else  {
		if x >= imgX && y >= imgY {
			return addPadding(x, y, img)
		}
		ratioX := float64(x)/float64(imgX)
		ratioY := float64(y)/float64(imgY)
		if ratioX < ratioY {
			img = resize.Resize(uint(x), uint(float64(x) * imageAspectRatio), img, resize.Lanczos3)
			return addPadding(x, y, img)
		}
		newX := uint(float64(y) / imageAspectRatio)
		img = resize.Resize(newX, uint(y), img, resize.Lanczos3)
		return addPadding(x, y, img)
	}
	return img
}

// Add padding to image
func addPadding(x int, y int, img image.Image) image.Image {
	rectangle := image.NewRGBA(image.Rect(0, 0, x, y))
	black := color.RGBA{}
	draw.Draw(rectangle, rectangle.Bounds(), &image.Uniform{C: black}, image.ZP, draw.Src)
	xPad := (x - img.Bounds().Dx()) / 2
	yPad := (y - img.Bounds().Dy()) / 2
	pt := image.Point{X: xPad, Y: yPad}
	imgRect := image.Rectangle{Min: pt, Max: pt.Add(img.Bounds().Size())}
	img.Bounds().At(xPad, yPad)
	draw.Draw(rectangle, imgRect, img, image.ZP, draw.Src)
	return rectangle
}