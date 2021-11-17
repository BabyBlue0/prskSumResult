package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"sync"

	"path"
	"time"
)

//Global variable
var globalAllSongs []PRSKSong
var globalRecords []PRSKOutputFormatToCSV

func getPathOfImages(dir string) ([]string, error) {
	poi := []string{}
	///*
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
	//*/
	//poi = append(poi, "img/source/result_offline.png")
	//poi = append(poi, "img/source/result_online.png")
	return poi, nil
}

func getPRSKScores(imagePath string, pos PRSKPositionOfData) (PRSKScore, error) {
	imImg, err := getImageFromFilePath(imagePath)
	if err != nil {
		return PRSKScore{}, err
	}
	fmt.Printf("%v: LOAD Image\n", imagePath)

	//cheking consistency of image and pos
	//if err := validateImagePos(imImg, pos); err != nil {
	//	return PRSKScore{}, err
	//}

	//Todo getXxxxをgorutineによって並列実行
	score := PRSKScore{}
	score.Name, _ = getTextFromImageByOCR(&imImg, pos.Name)
	//fmt.Printf("%v: get name\n", imagePath)
	score.Level, _ = getLevel(&imImg, pos.Level)
	//fmt.Printf("%v: get level\n", imagePath)
	score.Score, _ = getScore(&imImg, pos.Score)
	//fmt.Printf("%v: get score\n", imagePath)
	score.Combo, _ = getCombo(&imImg, pos.Combo)
	//fmt.Printf("%v: get combo\n", imagePath)
	score.Perfect, _ = getDetail(&imImg, pos.Perfect)
	//fmt.Printf("%v: get perfect\n", imagePath)
	score.Great, _ = getDetail(&imImg, pos.Great)
	//fmt.Printf("%v: get great\n", imagePath)
	score.Good, _ = getDetail(&imImg, pos.Good)
	//fmt.Printf("%v: get good\n", imagePath)
	score.Bad, _ = getDetail(&imImg, pos.Bad)
	//fmt.Printf("%v: get bad\n", imagePath)
	score.Miss, _ = getDetail(&imImg, pos.Miss)
	//fmt.Printf("%v: get miss\n", imagePath)

	//calc edit distance and decided title by ed
	var allsongs []string
	for _, s := range globalAllSongs {
		allsongs = append(allsongs, s.Title)
	}
	title, ed, _ := searchStringWithED(score.Name, allsongs)
	if ed > 5 {
		return PRSKScore{}, fmt.Errorf("Too high edit distance!!!\nsocre.Name: %v,\tED: %v\n", score.Name, ed)
	}
	score.Name = title
	//fmt.Printf("%v: get title\n", imagePath)
	return score, nil
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

	var failedFiles []string

	//Get all songs by web scraping
	fmt.Println("Getting all title of songs...")
	allSongs, err := getAllSongTitle()
	if err != nil {
		fmt.Println(err)
		return
	}
	globalAllSongs = allSongs

	//Get All image path in "img/source" directory
	fmt.Println("Searching image files...")
	imagePaths, err := getPathOfImages("img/source")
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	mutex := sync.Mutex{}
	limit := 7
	slots := make(chan struct{}, limit)
	wg.Add(len(imagePaths))
	for _, ip := range imagePaths {
		//to avoid exec multi
		ip := ip
		pos := pos

		slots <- struct{}{}
		go func(ip string, pos PRSKPositionOfData) {
			defer func() { <-slots }()
			defer wg.Done()
			fmt.Printf("%v: Start process...\n", ip)

			score, err := getPRSKScores(ip, pos)
			if err != nil {
				fmt.Printf("%v: %v\n", ip, err)
				failedFiles = append(failedFiles, ip)
				return
			}

			//extract timestamp
			imageInfo, err := os.Stat(ip)
			if err != nil {
				fmt.Printf("%v: %v\n", ip, err)
				return
			}

			ocsv := PRSKOutputFormatToCSV{Score: score, Timestamp: imageInfo.ModTime()}
			if dt, err := getTimestampByExif(ip); err != nil {
				ocsv.Timestamp = time.Time{}
			} else {
				ocsv.Timestamp = dt
			}
			ocsv.FName = path.Base(ip)

			//append records of csv
			mutex.Lock()
			globalRecords = append(globalRecords, ocsv)
			mutex.Unlock()

			fmt.Printf("%v: DONE\n", ip)
			fmt.Printf("%v: %v\n", ip, ocsv)
		}(ip, pos)
	}
	wg.Wait()

	writeCSV("output.csv", globalRecords)
	fmt.Printf("Write in csv file.\n")

	if failedFiles != nil {
		fmt.Println("\n\nSome image was failed...")
		for _, ip := range failedFiles {
			fmt.Println(ip)
		}
	} else {
		fmt.Println("\n\nALL SUCCESSFUL!")
	}

}
