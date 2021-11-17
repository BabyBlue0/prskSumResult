package main

import (
	"math"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PRSKSong struct {
	No      int
	Default string
	Title   string
	Unit    string
	ComboEx int
	ComboM  int
	Time    string
}

func calcEditDistance(s1, s2 string) int {
	ns1 := strings.ReplaceAll(s1, " ", "")
	ns2 := strings.ReplaceAll(s2, " ", "")

	r1, r2 := []rune(ns1), []rune(ns2)
	//fmt.Println( r1 )
	//fmt.Println( r2 )

	ed := make([][]int, len(r2)+1)
	for y, _ := range ed {
		ed[y] = make([]int, len(r1)+1)
		for x, _ := range ed[y] {
			if y == 0 || x == 0 {
				ed[y][x] = int(math.Max(float64(x), float64(y)))
			} else {
				var minN int
				if r1[x-1] == r2[y-1] {
					minN = ed[y-1][x-1]
				} else {
					minN = ed[y-1][x-1] + 1 //replace
				}
				if ed[y-1][x]+1 < minN {
					minN = ed[y-1][x] + 1 //insert
				}
				if ed[y][x-1]+1 < minN {
					minN = ed[y][x-1] + 1 //delete
				}

				ed[y][x] = minN
			}
		}
	}
	return ed[len(r2)][len(r1)]
}

func getAllSongTitle() ([]PRSKSong, error) {
	songs := []PRSKSong{}

	//プロジェクトセカイ攻略Wikiからスクレイピングする

	//When the test, It don't Scraping
	///*
	urlPjsekaiwiki := "https://pjsekai.com/"
	docHome, err := goquery.NewDocument(urlPjsekaiwiki)
	if err != nil {
		return nil, err
	}

	selAllSongs := docHome.Find("a[title=\"収録楽曲\"]").First()
	page, _ := selAllSongs.Attr("href")
	urlAllSongs := urlPjsekaiwiki + page[2:]

	docSongs, err := goquery.NewDocument(urlAllSongs)
	//*/

	///*** ON DEV ***
	//data write
	//res,_ = docSongs.Find("body").Html()
	//ioutil.WriteFile("./prjsekaiSONGS.html", []byte(res), os.ModePerm)

	//data read
	//fiSongs, _ := ioutil.ReadFile("./prjsekaiSONGS.html")
	//srSongs := strings.NewReader(string(fiSongs))
	//docSongs, err := goquery.NewDocumentFromReader(srSongs)
	//****************/

	if err != nil {
		return nil, err
	}

	//docSongs.Find("table#sortable_table1 > tbody > tr > td:nth-child(4)").Each(func(_ int, s *goquery.Selection) {
	docSongs.Find("table#sortable_table1 > tbody > tr").Each(func(_ int, s *goquery.Selection) {
		song := PRSKSong{}
		song.No, _ = strconv.Atoi(s.Find("td:nth-child(1)").Text())
		song.Default = s.Find("td:nth-child(2)").Text()
		song.Title = s.Find("td:nth-child(4)").Text()
		song.Unit = s.Find("td:nth-child(5)").Text()
		song.ComboEx, _ = strconv.Atoi(s.Find("td:nth-child(11)").Text())
		song.ComboM, _ = strconv.Atoi(s.Find("td:nth-child(12)").Text())
		song.Time = s.Find("td:nth-child(13)").Text()

		songs = append(songs, song)
	})

	return songs, nil
}

func searchStringWithEDFromPRSKSong(inStr string, inPRSKSongs []PRSKSong) (PRSKSong, int, error) {
	minED := 1000000
	correct := PRSKSong{}
	for _, ps := range inPRSKSongs {
		if calcED := calcEditDistance(inStr, ps.Title); calcED < minED {
			minED = calcED
			correct = ps
		}
		if minED == 0 {
			break
		}
	}
	return correct, minED, nil
}

func searchStringWithED(inStr string, inStrList []string) (string, int, error) {
	minED := 1000000
	correct := ""
	for _, str := range inStrList {
		if calcED := calcEditDistance(inStr, str); calcED < minED {
			minED = calcED
			correct = str
		}
		if minED == 0 {
			break
		}
	}
	return correct, minED, nil
}
