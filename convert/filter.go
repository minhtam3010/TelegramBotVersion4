package convert

import (
	"regexp"
	"strings"
)

func FilterString(myStr string) (res []interface{}, dateTime []string) {
	replacer := strings.NewReplacer("[", "", "{", "", "}", "", "]", "")
	filter := replacer.Replace(myStr)
	split_data := strings.Split(filter, " ")

	for _, str := range split_data {
		re := regexp.MustCompile("2022-[0-9][0-9]-[0-9][0-9]")
		if len(str) == 25 && len(re.FindStringSubmatch(str)) >= 1 {
			re = regexp.MustCompile("[0-9][0-9]:[0-9][0-9]:[0-9][0-9]")
			dateTime = append(dateTime, re.FindStringSubmatch(str)[0])
			res = append(res, re.FindStringSubmatch(str)[0])
		} else {
			res = append(res, str)
		}

	}
	return res, dateTime
}
