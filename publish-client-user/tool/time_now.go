package tool

import (
	"strconv"
	"time"
)

//获取系统当前时间
func GetTimeNow() *time.Time {
	//获取系统当前时间
	currentTimeData := time.Date(time.Now().Year(), time.Now().Month(),
		time.Now().Day(), time.Now().Hour(), time.Now().Minute(),
		time.Now().Second(), time.Now().Nanosecond(), time.Local)
	return &currentTimeData
}

//处理时间格式
func SubTime(str string) string {
	return str[0 : len(str)-15]
}

//获取当前时间戳
func GetTimeUnix() string {
	timeUnix := time.Now().UnixNano()
	i := int64(timeUnix)
	fileName := strconv.FormatInt(i, 10)
	return fileName
}

//time格式转string
func TimeFormatString(time time.Time) string {
	baseFormat := "2006-01-02 15:04:05"
	strTime := time.Format(baseFormat)
	return strTime
}

//string格式转time
func StringFormatTime(times string) time.Time {
	baseFormat := "2006-01-02 15:04:05"
	parseStrTime, _ := time.Parse(baseFormat, times)
	return parseStrTime
}

//获取当前时间 string 有小尾巴
func GetTimeStringNow() string {
	return TimeFormatString(*GetTimeNow())
}

//获取当前时间 string 没有小尾巴
func GetTimeStringNowNoTail() string {
	t := TimeFormatString(*GetTimeNow())
	return t[:len(t)-9]
}

//截取时间字符串中的年份有小尾巴2019-01-04 17:42:29
func GetYearString(time string) string {
	return time[0 : len(time)-15]
}

//截取时间字符串中的月份有小尾巴2019-01-04 17:42:29
func GetMonthString(time string) string {
	return time[5 : len(time)-12]
}

//截取时间字符串中的日期有小尾巴2019-01-04 17:42:29
func GetDayString(time string) string {
	return time[8 : len(time)-8]
}

//截取时间字符串中的年份没有小尾巴2019-01-04
func GetYearStringNoTail(time string) string {
	return time[0 : len(time)-6]
}

//截取时间字符串中的月份没有小尾巴2019-01-04
func GetMonthStringNoTail(time string) string {
	return time[5 : len(time)-3]
}

//截取时间字符串中的日期没有小尾巴2019-01-04
func GetDayStringNoTail(time string) string {
	return time[8:len(time)]
}

//获取系统当前年份+1
func GetCurrentYear() (string) {
	return strconv.Itoa(time.Now().Year() + 1)
}
