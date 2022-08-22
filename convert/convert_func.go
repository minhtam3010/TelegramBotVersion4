package convert

import (
	"time"
)

// func ConvertMysqlTimeUnixTime(mysqlTime string) string {
// 	if mysqlTime == "0" {
// 		return "0"
// 	}
// 	res1 := strings.Replace(mysqlTime, "T", " ", 1)
// 	res2 := res1[11:19]
// 	layout := "15:04:05"
// 	t, err := time.Parse(layout, res2)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return t.String()[11:19]
// }

func ConvertToday(today time.Time) string {
	return today.Format("2006-01-02")
}
