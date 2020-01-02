package untils

import (
	"time"
)

func ConvertUtcTimeToLocal(utcTime string, timeLayout string) (localTime string) {
	formate := "2006-01-02 15:04:05"
	parseTime, _ := time.Parse(timeLayout, utcTime)
	local, _ := time.LoadLocation("Local")
	return parseTime.In(local).Format(formate)
}

func GetNowTime() JSONTime {
	var cstZone = time.FixedZone("CST", 8*3600)
	return JSONTime{
		Time: time.Now().In(cstZone),
	}
}
