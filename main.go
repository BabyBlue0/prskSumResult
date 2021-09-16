package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"

	//"io/ioutil"
	"path"
	"time"
	//online
)

func getPathOfImages(dir string) ([]string, error) {
	poi := []string{}
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, fi := range fis {
		if !fi.IsDir() {
			p := path.Join(dir, fi.Name())
			poi = append(poi, p)
		}
	}

	return poi, nil
}

func main() {
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

	//Get all songs by web scraping
	allSongs, err := getAllSongTitle()
	if err != nil {
		fmt.Println(err)
		return
	}

	//Get All image path in "img/source" directory
	impaths, err := getPathOfImages("img/source")
	if err != nil {
		fmt.Println(err)
		return
	}

	records := []PRSKOutputFormatToCSV{}
	for _, ip := range impaths {
		img, err := getImageFromFilePath(ip)
		if err != nil {
			fmt.Println(err)
			return
		}

		//cheking consistency of image and pos
		if err := validateImagePos(img, pos); err != nil {
			fmt.Println(err)
			return
		}

		//Todo getXxxxをgorutineによって並列実行
		score := PRSKScore{}
		score.Name, _ = getTextFromImageByOCR(img, pos.Name)
		score.Level, _ = getTextFromImageByOCR(img, pos.Level)
		score.Score, _ = getScore(img, pos.Score)
		score.Combo, _ = getCombo(img, pos.Combo)
		score.Perfect, _ = getDetail(img, pos.Perfect)
		score.Great, _ = getDetail(img, pos.Great)
		score.Good, _ = getDetail(img, pos.Good)
		score.Bad, _ = getDetail(img, pos.Bad)
		score.Miss, _ = getDetail(img, pos.Miss)

		//Get timestamp
		imageInfo, err := os.Stat(ip)
		if err != nil {
			fmt.Println(err)
			return
		}

		//calc edit distance and decided title by ed
		fmt.Printf("IMG: %v\n", ip)
		title, ed, _ := searchStringWithED(score.Name, allSongs)
		if ed > 5 {
			fmt.Println("Too high edit distance!!!")
			fmt.Printf("socre.Name: %v,\tED: %v\n", score.Name, ed)
			continue
		}
		score.Name = title

		//extract timestamp
		ocsv := PRSKOutputFormatToCSV{Score: score, Timestamp: imageInfo.ModTime()}
		if dt, err := getTimestampByExif(ip); err != nil {
			fmt.Println(err)
			ocsv.Timestamp = time.Time{}
		} else {
			ocsv.Timestamp = dt
		}
		ocsv.FName = path.Base(ip)

		//append records of csv
		records = append(records, ocsv)
	}

	writeCSV("output.csv", records)
}
