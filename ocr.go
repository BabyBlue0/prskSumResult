package main

import (
	"image"
	"strconv"
	"strings"

	"github.com/otiai10/gosseract"
)

func extractTextFromBytes(image_byte []byte) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetLanguage("jpn", "eng")
	client.Trim = true
	client.SetImageFromBytes(image_byte)

	text, err := client.Text()
	return text, err
}

func extractNumFromBytes(image_byte []byte) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetLanguage("eng")
	client.SetWhitelist("0123456789")
	client.SetImageFromBytes(image_byte)

	text, err := client.Text()
	return text, err
}

func getTextFromImageByOCR(img image.Image, rec image.Rectangle) (string, error) {
	buf, err := croppingImageToBytes(img, rec)
	if err != nil {
		return "", err
	}
	text, err := extractTextFromBytes(buf)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(text, " ", ""), nil
}

func getNumFromImageByOCR(img image.Image, rec image.Rectangle) (uint, error) {
	buf, err := croppingImageToBytes(img, rec)
	if err != nil {
		return 0, err
	}
	text, err := extractNumFromBytes(buf)
	if err != nil {
		return 0, err
	}
	num, err := strconv.Atoi(text)
	if err != nil {
		return 0, nil
	}

	return uint(num), nil
}
