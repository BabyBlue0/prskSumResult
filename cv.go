package main

import (
	"fmt"
	"image"
	"sort"

	"gocv.io/x/gocv"
)

func extractNumFromImageByTM(imImg image.Image, rec image.Rectangle, templatePaths []string, confidence float32) (uint, error) {
	//give some room
	rec = image.Rect(rec.Min.X-10, rec.Min.Y-10, rec.Max.X+10, rec.Max.Y+10)

	cvFullImg, err := gocv.ImageToMatRGB(imImg)
	if err != nil {
		return 0, err
	}
	defer cvFullImg.Close()

	//cropping image
	cvImgOrig := cvFullImg.Region(rec)

	cvImgGray := gocv.NewMat()
	gocv.CvtColor(cvImgOrig, &cvImgGray, gocv.ColorBGRToGray)
	defer cvImgGray.Close()
	mask := gocv.NewMat()
	defer mask.Close()

	var recs []NumPos

	for idx, tp := range templatePaths {
		cvTmp := gocv.IMRead(tp, gocv.IMReadGrayScale)
		defer cvTmp.Close()
		tmpH, tmpW := cvTmp.Size()[0], cvTmp.Size()[1]

		// matching template
		res := gocv.NewMat()
		defer res.Close()
		gocv.MatchTemplate(cvImgGray, cvTmp, &res, 3, mask)

		//get result
		//gocv.Threshold( res, &res, 0.9, 1, 3 )
		for y := 0; y < res.Rows(); y++ {
			for x := 0; x < res.Cols(); x++ {
				if confidence <= res.GetFloatAt(y, x) && res.GetFloatAt(y, x) <= 1.00 {
					recs = append(recs, NumPos{Rec: image.Rect(rec.Min.X+x, rec.Min.Y+y, rec.Min.X+x+tmpW, rec.Min.Y+y+tmpH), Num: uint(idx)})
				}
			}
		}
	}

	/* draw rectangle
		for _, rec := range recs {
	    cvFullImg = drawRectangleFromRect( &cvFullImg, rec.Rec, color.RGBA{0,0,0,0}, 3).(gocv.Mat)
		  fmt.Println("Num: ", rec.Num, "min: ", rec.Rec.Min, "max: ", rec.Rec.Max )
		}
		gocv.IMWrite("out_check.png", cvImgOrig)
	*/

	//sort recs
	sort.SliceStable(recs, func(i, j int) bool { return recs[i].Rec.Min.X < recs[j].Rec.Min.X })

	//get num from NumPos
	retnum := uint(0)
	prevX := 0
	for _, rec := range recs {
		//if point is so neary, the point is ignored
		if rec.Rec.Min.X < prevX+3 {
			continue
		}
		retnum = retnum*10 + rec.Num
		prevX = rec.Rec.Min.X
	}

	return retnum, nil
}

func getScore(imImg image.Image, rec image.Rectangle) (uint, error) {
	tmps := []string{}
	for i := 0; i <= 9; i++ {
		tmps = append(tmps, fmt.Sprintf("img/tmp/tmp_score_%d.png", i))
	}

	num, err := extractNumFromImageByTM(imImg, rec, tmps, 0.999)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func getCombo(imImg image.Image, rec image.Rectangle) (uint, error) {
	tmps := []string{}
	for i := 0; i <= 9; i++ {
		tmps = append(tmps, fmt.Sprintf("img/tmp/tmp_combo_%d.png", i))
	}

	num, err := extractNumFromImageByTM(imImg, rec, tmps, 0.999)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func getDetail(imImg image.Image, rec image.Rectangle) (uint, error) {
	tmps := []string{}
	for i := 0; i <= 9; i++ {
		tmps = append(tmps, fmt.Sprintf("img/tmp/tmp_detail_%d.png", i))
	}

	num, err := extractNumFromImageByTM(imImg, rec, tmps, 0.99)
	if err != nil {
		return 0, err
	}

	return num, nil
}
