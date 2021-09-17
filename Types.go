package main

import (
	"image"
	"strconv"
	"time"
)

type PRSKScore struct {
	//music name
	Name string
	//dificulty
	Level string
	//score
	Score uint
	//max combo
	Combo uint
	//detail score
	Perfect uint
	Great   uint
	Good    uint
	Bad     uint
	Miss    uint
}

type PRSKOutputFormatToCSV struct {
	Score     PRSKScore
	Timestamp time.Time
	FName     string
}

func (p PRSKOutputFormatToCSV) ToMap() map[string]string {
	return map[string]string{
		"Name":      p.Score.Name,
		"Level":     p.Score.Level,
		"Score":     strconv.FormatUint(uint64(p.Score.Score), 10),
		"Combo":     strconv.FormatUint(uint64(p.Score.Combo), 10),
		"Perfect":   strconv.FormatUint(uint64(p.Score.Perfect), 10),
		"Great":     strconv.FormatUint(uint64(p.Score.Great), 10),
		"Good":      strconv.FormatUint(uint64(p.Score.Good), 10),
		"Bad":       strconv.FormatUint(uint64(p.Score.Bad), 10),
		"Miss":      strconv.FormatUint(uint64(p.Score.Miss), 10),
		"Timestamp": p.Timestamp.Format("2006/01/02 03:04:05"),
		"FName":     p.FName}
}

func (p PRSKOutputFormatToCSV) Titles() []string {
	return []string{
		"Name",
		"Level",
		"Score",
		"Combo",
		"Perfect",
		"Great",
		"Good",
		"Bad",
		"Miss",
		"Timestamp",
		"FName",
	}
}

type PRSKPositionOfData struct {
	//image size
	Width  uint
	Height uint
	//music name
	Name image.Rectangle
	//dificulty
	Level image.Rectangle
	//score
	Score image.Rectangle
	//max combo
	Combo image.Rectangle
	//detail score
	Perfect image.Rectangle
	Great   image.Rectangle
	Good    image.Rectangle
	Bad     image.Rectangle
	Miss    image.Rectangle
}

type NumPos struct {
	Num uint
	Rec image.Rectangle
}
