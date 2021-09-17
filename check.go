package main

import (
	"image"
	"image/color"
	"image/draw"

	"gocv.io/x/gocv"
)

func drawRectangle(ifImg interface{}, sx, sy, ex, ey int, c color.RGBA, px int) interface{} {
	switch img := ifImg.(type) {
	case image.Image:
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
		return dimg.(image.Image)

	case gocv.Mat:
		gocv.Rectangle(&img, image.Rect(sx, sy, ex, ey), c, px)
		return img
	default:
		return nil
	}
}

func drawRectangleFromRect(img interface{}, rect image.Rectangle, c color.RGBA, px int) interface{} {
	return drawRectangle(img, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y, c, px)
}
