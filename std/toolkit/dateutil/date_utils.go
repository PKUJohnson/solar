package dateutil

import (
	"strconv"
	"time"
)

const (
	TimeDivide = 1000 * 1000
    TimeFormat = "2006/1/2 15:04:05"
    DataFormat = "20060102"
	FrontendTimeFormat = "2006-01-02 15:04:05"
)

var (
	ShangHaiTimeZone, _ = time.LoadLocation("Asia/Shanghai")
)

func GetNowTimestamp() int64 {
	now := time.Now().UnixNano()

	return now / TimeDivide
}

func GetDateTime(format, dateStr string) time.Time {
	singleTime, _ := time.Parse(format, dateStr)
	return singleTime
}

func GetUnixNano(format, dateStr string) int64 {
	singleTime, _ := time.Parse(format, dateStr)

	return singleTime.UnixNano()
}

func GetNextDayMiddleTime() time.Time {
	nowTime := time.Now()
	nextDay := nowTime.AddDate(0, 0, 1)
	middleTime, _ := time.ParseInLocation("20060102", nextDay.Format("20060102"), time.Local)

	return middleTime
}

func GetDayMiddleTime(t time.Time, cnt int) time.Time {
	nextDay := t.AddDate(0, 0, cnt)
	middleTime, _ := time.ParseInLocation("20060102", nextDay.Format("20060102"), time.Local)

	return middleTime
}

func GetNextDayMiddle() int64 {
	nowTime := time.Now()
	nextDay := nowTime.AddDate(0, 0, 1)
	middleTime, _ := time.ParseInLocation("20060102", nextDay.Format("20060102"), time.Local)
	return middleTime.Unix()
}

func GetNMinuteAgo(cnt int) time.Time {
	nowTime := time.Now()
	newTime := nowTime.Add(-10 * time.Minute)

	return newTime
}

func GetNDayAgo(cnt int, format string, srcDate string) string {
	singleTime, _ := time.Parse(format, srcDate)
	nextDay := singleTime.AddDate(0, 0, -cnt)
	res := nextDay.Format(format)
	return res
}

func GetNDayAgoTime(cnt int) time.Time {
	nowTime := time.Now()
	nextDay := nowTime.AddDate(0, 0, -cnt)
	middleTime, _ := time.ParseInLocation("20060102", nextDay.Format("20060102"), time.Local)

	return middleTime
}

func GetNMonthTime(cnt int) time.Time {
	nowTime := time.Now()
	newTime := nowTime.AddDate(0, cnt, 0)
	middleTime, _ := time.ParseInLocation("20060102", newTime.Format("20060102"), time.Local)

	return middleTime
}

func FormatDateTime(datestr string, srcFormat string, toFormat string) string {
	singleTime, _ := time.Parse(srcFormat, datestr)
	result := singleTime.Format(toFormat)
	return result
}

func GetCalendarToday() string {
	return time.Now().Format("2006-01-02")
}

func GetCalendarDayDiff(day_count int64) string {
	timestamp := time.Now().Unix() + day_count*24*3600
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02")
}

func GetStandarTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func GetTodayBeginTime() time.Time {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	return t
}

func GetDayBeginTime(day time.Time) time.Time {
	timeStr := day.Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	return t
}

/*
 * 获取当月第一天
 */
func GetFirstOfMonth() time.Time {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	return time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
}

func GetStartDayOfWeek() time.Time {
	now := time.Now()
	weekday := time.Duration(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	year, month, day := now.Date()
	currentZeroDay := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return currentZeroDay.Add(-1 * (weekday - 1) * 24 * time.Hour)
}

func CheckTimeIsToday(t time.Time) bool {
	todayBeginTime := GetTodayBeginTime()
	if t.Unix() < todayBeginTime.Unix() {
		return false
	}
	return true
}

func GetNWeeks(cnt int) []string {
	arr := []string{}
	for i := 0; i < cnt; i++ {
		t := GetNDayAgoTime(i * 7)
		year, week := t.ISOWeek()
		weekName := strconv.Itoa(year)
		if week < 10 {
			weekName = weekName + "0" + strconv.Itoa(week)
		} else {
			weekName = weekName + strconv.Itoa(week)
		}
		arr = append(arr, weekName)
	}
	return arr
}

func GetYesterdayTime() time.Time {
	passTwentyFourHours := time.Now().In(ShangHaiTimeZone).AddDate(0, 0, -1)
	return passTwentyFourHours
}

func GetDatePeriodString(dayCount int64) (string, string) {
	local, _ := time.LoadLocation("Asia/Shanghai")
	endDate := time.Now().In(local)
	startDate := endDate.AddDate(0, 0, -int(dayCount))
	return startDate.Format("20060102"), endDate.Format("20060102")
}

func GetDateTimePeriodString(dayCount int64) (string, string) {
	local, _ := time.LoadLocation("Asia/Shanghai")
	endDate := time.Now().In(local)
	startDate := endDate.AddDate(0, 0, -int(dayCount))
	return startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05")
}

func ParseShanghaiTime(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, ShangHaiTimeZone)
}

func GetLastDayOfMonth(date string) string {
	t, _ := time.Parse("20060102", date)
	y,m,_ := t.Date()
	lastday:= time.Date(y,m+1,0,0,0,0,0,time.UTC)
	return lastday.Format("20060102")
}

func GetFridayOfWeek(date string) string {
	t, _ := time.Parse("20060102", date)
	y,m,d := t.Date()
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d2 := d + 5 - weekday
	lastday:= time.Date(y,m,d2,0,0,0,0,time.UTC)
	return lastday.Format("20060102")
}