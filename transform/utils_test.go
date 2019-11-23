package transform

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"testing"
)

type TransformTestSuite struct {
	suite.Suite
	source image.Image
}

func (suite *TransformTestSuite) SetupSuite() {
	suite.source = loadImg("./../resources/test/source.jpeg")
}

func (suite *TransformTestSuite) TestTransformation_300x300() {
	transformed := ProcessImg(300, 300, suite.source)
	img300x300 := loadImg("./../resources/test/300x300.jpeg")
	res := FastCompare(transformed, img300x300)
	assert.True(suite.T(), res, "Images comparing failed")
}

func (suite *TransformTestSuite) TestTransformation_300x100() {
	transformed := ProcessImg(300, 100, suite.source)
	img300x100 := loadImg("./../resources/test/300x100.jpeg")
	res := FastCompare(transformed, img300x100)
	assert.True(suite.T(), res, "Images comparing failed")
}

func (suite *TransformTestSuite) TestTransformation_300x500() {
	transformed := ProcessImg(300, 500, suite.source)
	img300x500 := loadImg("./../resources/test/300x500.jpeg")
	res := FastCompare(transformed, img300x500)
	assert.True(suite.T(), res, "Images comparing failed")
}

func (suite *TransformTestSuite) TestTransformation_300x1000() {
	transformed := ProcessImg(300, 1000, suite.source)
	img300x1000 := loadImg("./../resources/test/300x1000.jpeg")
	res := FastCompare(transformed, img300x1000)
	assert.True(suite.T(), res, "Images comparing failed")
}

func (suite *TransformTestSuite) TestTransformation_500x300() {
	transformed := ProcessImg(500, 300, suite.source)
	img500x300 := loadImg("./../resources/test/500x300.jpeg")
	res := FastCompare(transformed, img500x300)
	assert.True(suite.T(), res, "Images comparing failed")
}

func TestTransformTestSuite(t *testing.T) {
	suite.Run(t, new(TransformTestSuite))
}

func FastCompare(img1, img2 image.Image) bool {
	b := img1.Bounds()
	if !b.Eq(img2.Bounds()) {
		log.Println("different image sizes")
		return false
	}

	var sum int64
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			sum += diff(r1, r2)
			sum += diff(g1, g2)
			sum += diff(b1, b2)
		}
	}

	nPixels := (b.Max.X - b.Min.X) * (b.Max.Y - b.Min.Y)
	dif := float64(sum*100) / (float64(nPixels) * 0xffff * 3)

	if dif < 1 {
		return true
	} else {
		return false
	}
}

func diff(a, b uint32) int64 {
	if a > b {
		return int64(a - b)
	}
	return int64(b - a)
}

func loadImg(path string) image.Image {
	file, _ := os.Open(path)
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("Failed to load source img " + err.Error())
	}
	return img
}
