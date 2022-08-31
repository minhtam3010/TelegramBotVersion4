package convert

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func FilterString(myStr string) (res []interface{}, dateTime []string) {
	replacer := strings.NewReplacer("[", "", "{", "", "}", "", "]", "")
	filter := replacer.Replace(myStr)
	split_data := strings.Split(filter, " ")

	for _, str := range split_data {
		re := regexp.MustCompile("2022-[0-9][0-9]-[0-9][0-9]")
		if len(str) == 10 && len(re.FindStringSubmatch(str)) >= 1 {
			// re = regexp.MustCompile("[0-9][0-9]:[0-9][0-9]:[0-9][0-9]") if len(str) == 25
			dateTime = append(dateTime, re.FindStringSubmatch(str)[0])

			res = append(res, re.FindStringSubmatch(str)[0])
		} else {
			res = append(res, str)
		}
	}
	return res, dateTime
}

func Sum(arr []int) (total int) {
	for i := 0; i < len(arr); i++ {
		total += arr[i]
	}
	return
}

func SumTotal(arr2d [][]string) (total int) {
	for i := 0; i < len(arr2d); i++ {
		getInt, err := strconv.Atoi(arr2d[i][4])
		if err != nil {
			fmt.Errorf("Error while parsing: %", err)
			return
		}
		total += getInt
	}
	return
}
