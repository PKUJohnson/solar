package dateutil

import (
	"fmt"
	"testing"
	"time"
)

func TestGetFirstOfMonth(t *testing.T) {
	fmt.Println(GetFirstOfMonth())
}

func TestGetNMinuteAgo(t *testing.T) {
	fmt.Println(GetNMinuteAgo(10))
	fmt.Println(GetNMinuteAgo(5))
}

func TestWeekweek(t *testing.T) {
	tt := time.Now()
	t2 := GetNMonthTime(-1)
	fmt.Println(tt.ISOWeek())
	fmt.Println(t2.ISOWeek())
	fmt.Println(t2)
}

func TestWeek2(t *testing.T) {
	for i := 0; i < 12; i++ {
		t1 := GetNDayAgoTime(i * 7)
		fmt.Println(t1.ISOWeek())
	}

	fmt.Println(GetNDayAgoTime(10))
}

func TestYear(t *testing.T) {
	tt := GetDateTime("20060102", "20160101")
	fmt.Println(tt)
	fmt.Println(tt.Unix())

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}

func TestGetNDayAgoTime(t *testing.T) {
	fmt.Println(GetYesterdayTime())
}

func TestGetDayMiddleTime(t *testing.T) {
	tt := time.Unix(1486078800, 0)
	fmt.Println(tt)
	fmt.Println(GetDayMiddleTime(tt, 1))
	fmt.Println(GetDayBeginTime(tt))
}
