package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetTimeStampByExif(t *testing.T) {
	path := "img/source/result_online.png"
	ts, err := getTimestampByExif(path)
	if err != nil {
		t.Fatal(err)
	}
	jp := time.FixedZone("Asia/Tokyo", 9*60*60)
	fmt.Println(ts.In(jp).Format("2006/01/02 03:04:05"))
}
