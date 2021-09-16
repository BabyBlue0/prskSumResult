package main

import (
	"fmt"
	"image"
	"sort"

	"gocv.io/x/gocv"
)

func extractNumFromImageByTM(img image.Image, rec image.Rectangle, tmps []string, confidence float32) (uint, error) {
	//give some room
	rec = image.Rect(rec.Min.X-10, rec.Min.Y-10, rec.Max.X+10, rec.Max.Y+10)

	fullImg, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return 0, err
	}
	defer fullImg.Close()

	//cropping image
	img_orig := fullImg.Region(rec)

	img_gray := gocv.NewMat()
	gocv.CvtColor(img_orig, &img_gray, gocv.ColorBGRToGray)
	defer img_gray.Close()
	mask := gocv.NewMat()
	defer mask.Close()

	//return value, number and position
	var recs []NumPos

	for idx, tmp := range tmps {
		template := gocv.IMRead(tmp, gocv.IMReadGrayScale)
		defer template.Close()
		tmp_h, tmp_w := template.Size()[0], template.Size()[1]

		// matching template
		res := gocv.NewMat()
		defer res.Close()
		gocv.MatchTemplate(img_gray, template, &res, 3, mask)

		//get result
		//gocv.Threshold( res, &res, 0.9, 1, 3 )
		for y := 0; y < res.Rows(); y++ {
			for x := 0; x < res.Cols(); x++ {
				if confidence <= res.GetFloatAt(y, x) && res.GetFloatAt(y, x) <= 1.00 {
					recs = append(recs, NumPos{Rec: image.Rect(rec.Min.X+x, rec.Min.Y+y, rec.Min.X+x+tmp_w, rec.Min.Y+y+tmp_h), Num: uint(idx)})
				}
			}
		}
	}

	/* draw rectangle
		img_orig := gocv.IMRead(ip, gocv.IMReadColor)
		defer img_orig.Close()

		recs, _ := extractNumFromImageByTM(img, pos[2], tmps)
		for _, rec := range recs {
			gocv.Rectangle(&img_orig, rec.Rec, color.RGBA{0, 0, 0, 0}, 3)
	    fmt.Println("Num: ", rec.Num, "min: ", rec.Rec.Min, "max: ", rec.Rec.Max )
		}
		gocv.IMWrite("out_"+ip, img_orig)
	*/

	//sort recs
	sort.SliceStable(recs, func(i, j int) bool { return recs[i].Rec.Min.X < recs[j].Rec.Min.X })

	//get num from NumPos
	retnum := uint(0)
	prevX := 0
	for _, rec := range recs {
		//if pixcel is so neary, the pixcel is ignored
		if rec.Rec.Min.X < prevX+3 {
			continue
		}
		retnum = retnum*10 + rec.Num
		prevX = rec.Rec.Min.X
	}

	return retnum, nil
}

func getScore(img image.Image, rec image.Rectangle) (uint, error) {
	tmps := []string{}
	for i := 0; i <= 9; i++ {
		tmps = append(tmps, fmt.Sprintf("img/tmp/tmp_score_%d.png", i))
	}

	num, err := extractNumFromImageByTM(img, rec, tmps, 0.999)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func getCombo(img image.Image, rec image.Rectangle) (uint, error) {
	tmps := []string{}
	for i := 0; i <= 9; i++ {
		tmps = append(tmps, fmt.Sprintf("img/tmp/tmp_combo_%d.png", i))
	}

	num, err := extractNumFromImageByTM(img, rec, tmps, 0.999)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func getDetail(img image.Image, rec image.Rectangle) (uint, error) {
	tmps := []string{}
	for i := 0; i <= 9; i++ {
		tmps = append(tmps, fmt.Sprintf("img/tmp/tmp_detail_%d.png", i))
	}

	num, err := extractNumFromImageByTM(img, rec, tmps, 0.99)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func cv_getImagePaths() []string {
	return []string{
		"result1.png",
		"result2.png",
		"result3.png",
		"result4.png",
		"result5.png",
		"result6.png",
		"result_offline.png",
		"result_online.png",
	}

}

func cv_main() {

	//for image cropping

	pos := PRSKPositionOfData{
		Width:   2160,
		Height:  1620,
		Name:    image.Rect(189, 15, 1000, 63),
		Score:   image.Rect(80, 547, 783, 665),
		Combo:   image.Rect(260, 955, 760, 1200),
		Level:   image.Rect(185, 87, 390, 140),
		Perfect: image.Rect(1120, 935, 1240, 975),
		Great:   image.Rect(1120, 995, 1240, 1035),
		Good:    image.Rect(1120, 1054, 1240, 1094),
		Bad:     image.Rect(1120, 1113, 1240, 1153),
		Miss:    image.Rect(1120, 1172, 1240, 1212),
	}

	ips := cv_getImagePaths()
	for _, ip := range ips {
		img, err := getImageFromFilePath(ip)
		if err != nil {
			fmt.Println(err)
			return
		}

		score, _ := getScore(img, pos.Score)
		combo, _ := getCombo(img, pos.Combo)
		perfect, _ := getDetail(img, pos.Perfect)
		great, _ := getDetail(img, pos.Great)
		good, _ := getDetail(img, pos.Good)
		bad, _ := getDetail(img, pos.Bad)
		miss, _ := getDetail(img, pos.Miss)

		//output
		fmt.Println("Image: ", ip)
		fmt.Println("Score: ", score)
		fmt.Println("Combo: ", combo)
		fmt.Println("perfect: ", perfect)
		fmt.Println("great: ", great)
		fmt.Println("good: ", good)
		fmt.Println("bad: ", bad)
		fmt.Println("Miss: ", miss)
		fmt.Println("")
	}
}
