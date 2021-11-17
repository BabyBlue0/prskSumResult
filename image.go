package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"time"
)

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	imImg, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return imImg, err
}

func outputImage(filePath string, imImg image.Image) error {
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	if err := png.Encode(out, imImg); err != nil {
		return err
	}

	return nil
}

func croppingImageToImage(imImg *image.Image, pos image.Rectangle) image.Image {
	//cropping data from result
	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	crpimg := (*imImg).(SubImager).SubImage(pos)
	return crpimg.(image.Image)
}

func croppingImageToBytes(imImg *image.Image, pos image.Rectangle) ([]byte, error) {
	imCrpImg := croppingImageToImage(imImg, pos)
	//image.Image to bytes.Buffer
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, imCrpImg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func validateImagePos(imImg image.Image, pos PRSKPositionOfData) error {
	//checking file extend
	//checking some positions of cropping
	//checking file size
	if imImg.Bounds().Dx() != int(pos.Width) || imImg.Bounds().Dy() != int(pos.Height) {
		return fmt.Errorf("Image size is not same.\nImage: %vx%v, Pos:%vx%v", imImg.Bounds().Dx(), imImg.Bounds().Dy(), pos.Width, pos.Height)
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
