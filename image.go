package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"time"
	//online
)

type NumPos struct {
	Num uint
	Rec image.Rectangle
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, err
}

func outputImage(filePath string, img image.Image) error {
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		return err
	}

	return nil
}

func croppingImageToImage(img image.Image, pos image.Rectangle) image.Image {
	//cropping data from result
	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	crpimg := img.(SubImager).SubImage(pos)
	return crpimg.(image.Image)
}

func croppingImageToBytes(img image.Image, pos image.Rectangle) ([]byte, error) {
	crpimg := croppingImageToImage(img, pos)
	//image.Image to bytes.Buffer
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, crpimg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func validateImagePos(img image.Image, pos PRSKPositionOfData) error {
	//checking file extend
	//checking some positions of cropping
	//checking file size
	if img.Bounds().Dx() != int(pos.Width) ||
		img.Bounds().Dy() != int(pos.Height) {
		return fmt.Errorf("Image size is not same.\nImage: %vx%v, Pos:%vx%v", img.Bounds().Dx(), img.Bounds().Dy(), pos.Width, pos.Height)
	}

	return nil
}

func getTimestampByExif(imgPath string) (time.Time, error) {
	//file, err := os.Open(imgPath)
	//if err != nil {
	//  fmt.Println("image.go: error in os.Open")
	//  return time.Time{}, err
	//}

	//x, err := exif.Decode(file)
	//if err != nil {
	//  fmt.Println("image.go: error in exif.Decode")
	//  return time.Time{}, err
	//}

	//dt, err := x.DateTime()
	//if err != nil {
	//  fmt.Println("image.go: error in x.DateTime")
	//  return time.Time{}, err
	//}
	return time.Time{}, nil
}
