package main

import (
	"image"
	"strconv"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

func extractTextFromBytes(byteImage *[]byte) (string, error) {
	return extractTextFromBytesSpecifyingWhitelist(byteImage, "")
}
func extractTextFromBytesSpecifyingWhitelist(byteImage *[]byte, whiteList string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetLanguage("jpn", "eng")
	client.Trim = true
	if whiteList != "" {
		client.SetWhitelist(whiteList)
	}
	client.SetImageFromBytes(*byteImage)

	text, err := client.Text()
	return text, err
}

func extractNumFromBytes(byteImage *[]byte) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetLanguage("eng")
	client.SetWhitelist("0123456789")
	client.SetImageFromBytes(*byteImage)

	text, err := client.Text()
	return text, err
}

func getTextFromImageByOCR(imImg *image.Image, rec image.Rectangle) (string, error) {
	return getTextFromImageByOCRSpecifyngWhitelist(imImg, rec, "")
}

func getTextFromImageByOCRSpecifyngWhitelist(imImg *image.Image, rec image.Rectangle, whiteList string) (string, error) {
	buf, err := croppingImageToBytes(imImg, rec)
	if err != nil {
		return "", err
	}
	text, err := extractTextFromBytesSpecifyingWhitelist(&buf, whiteList)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(text, " ", ""), nil
}

func getNumFromImageByOCR(imImg *image.Image, rec image.Rectangle) (uint, error) {
	buf, err := croppingImageToBytes(imImg, rec)
	if err != nil {
		return 0, err
	}
	text, err := extractNumFromBytes(&buf)
	if err != nil {
		return 0, err
	}
	num, err := strconv.Atoi(text)
	if err != nil {
		return 0, nil
	}

	return uint(num), nil
}
