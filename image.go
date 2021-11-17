package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"github.com/dsoprea/go-exif"
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
	rawExif, err := exif.SearchFileAndExtractExif(imgPath)
	if err != nil {
		return time.Time{}, err
	}

	im := exif.NewIfdMapping()
	err = exif.LoadStandardIfds(im)
	if err != nil {
		return time.Time{}, err
	}
	ti := exif.NewTagIndex()

	_, index, err := exif.Collect(im, ti, rawExif)
	if err != nil {
		return time.Time{}, err
	}
	//index.RootIfd.PrintTagTree(true)

	rootIfd := index.RootIfd
	exifIfd, err := exif.FindIfdFromRootIfd(rootIfd, "IFD/Exif")
	if err != nil {
		return time.Time{}, err
	}

	tagName := "DateTimeOriginal"
	results, err := exifIfd.FindTagWithName(tagName)
	if err != nil {
		return time.Time{}, err
	}

	if len(results) != 1 {
		return time.Time{}, fmt.Errorf("There wasn't exactly one result")
	}

	ite := results[0]

	valueRaw, err := index.RootIfd.TagValue(ite)
	if err != nil {
		return time.Time{}, err
	}

	value := valueRaw.(string)
	dt, err := exif.ParseExifFullTimestamp(value)
	if err != nil {
		return time.Time{}, err
	}
	return dt, nil
}
