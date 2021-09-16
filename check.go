package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

func drawRectangle(img image.Image, sx, sy, ex, ey int, c color.RGBA, px int) (image.Image, error) {
	dimg := img.(draw.Image)

	//横線
	for i := sx; i < ex; i++ {
		for j := 0; j < px; j++ {
			if 0 <= sy-j && sy-j < img.Bounds().Max.Y {
				dimg.Set(i, sy-j, c)
			}
			if 0 <= ey+j && ey+j < img.Bounds().Max.Y {
				dimg.Set(i, ey+j, c)
			}
		}
	}

	//縦線
	for i := sy; i < ey; i++ {
		for j := 0; j < px; j++ {
			if 0 <= sx-j && sx-j < img.Bounds().Max.Y {
				dimg.Set(sx-j, i, c)
			}
			if 0 <= ex+j && ex+j < img.Bounds().Max.Y {
				dimg.Set(ex+j, i, c)
			}
		}
	}

	return dimg.(image.Image), nil
}

func drawRectangleFromRect(img image.Image, rect image.Rectangle, c color.RGBA, px int) (image.Image, error) {
	return drawRectangle(img, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y, c, px)
}

func check_main() {
	img, err := getImageFromFilePath("img/result_offline.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	// draw a rectangle around the score
	rec := image.Rect(80, 547, 783, 665)
	col := color.RGBA{255, 0, 0, 255}
	pix := 10
	img, err = drawRectangleFromRect(img, rec, col, pix)
	if err != nil {
		fmt.Println(err)
		return
	}

	// draw a rectangle around the music
	rec = image.Rect(187, 15, 1000, 63)
	col = color.RGBA{0, 255, 255, 255}
	pix = 5
	img, err = drawRectangleFromRect(img, rec, col, pix)
	if err != nil {
		fmt.Println(err)
		return
	}

	// draw a rectangle around the combo
	rec = image.Rect(290, 980, 722, 1174)
	col = color.RGBA{0, 0, 255, 255}
	pix = 10
	img, err = drawRectangleFromRect(img, rec, col, pix)
	if err != nil {
		fmt.Println(err)
		return
	}

	// draw a rectangle around the dificulty
	rec = image.Rect(185, 87, 390, 140)
	col = color.RGBA{255, 255, 0, 255}
	pix = 5
	img, err = drawRectangleFromRect(img, rec, col, pix)
	if err != nil {
		fmt.Println(err)
		return
	}

	// draw a rectangle around the detail socore
	col = color.RGBA{0, 255, 0, 255}
	pix = 5
	//perfect
	{
		rec := image.Rect(1120, 935, 1240, 975)
		img, err = drawRectangleFromRect(img, rec, col, pix)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//great
	{
		rec := image.Rect(1120, 995, 1240, 1035)
		img, err = drawRectangleFromRect(img, rec, col, pix)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//good
	{
		rec := image.Rect(1120, 1054, 1240, 1094)
		img, err = drawRectangleFromRect(img, rec, col, pix)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//bad
	{
		rec := image.Rect(1120, 1113, 1240, 1153)
		img, err = drawRectangleFromRect(img, rec, col, pix)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	//miss
	{
		rec := image.Rect(1120, 1172, 1240, 1212)
		img, err = drawRectangleFromRect(img, rec, col, pix)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if err := outputImage("img/output.png", img); err != nil {
		fmt.Println(err)
		return
	}
}
