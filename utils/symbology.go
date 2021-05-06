package utils

import (
	"fmt"
	"strings"
	"time"
)

type CCYSymbol struct {
	Base  string
	Quote string
}

func NewCCYSymbol(symbol string) (CCYSymbol, error) {
	splits := strings.Split(symbol, "/")
	if len(splits) != 2 {
		return CCYSymbol{}, fmt.Errorf("was expecting two symbol")
	}
	return CCYSymbol{
		Base:  splits[0],
		Quote: splits[1],
	}, nil
}

func (s CCYSymbol) Format(format string) string {
	res := format
	res = strings.Replace(res, "%B", strings.ToUpper(s.Base), -1)
	res = strings.Replace(res, "%Q", strings.ToUpper(s.Quote), -1)
	res = strings.Replace(res, "%b", strings.ToLower(s.Base), -1)
	res = strings.Replace(res, "%q", strings.ToLower(s.Quote), -1)
	return res
}

func (s CCYSymbol) DefaultFormat() string {
	return s.Base + "/" + s.Quote
}

func MonthToTimeMonth(m string) (time.Month, error) {
	switch strings.ToLower(m) {
	case "f":
		return time.January, nil
	case "g":
		return time.February, nil
	case "h":
		return time.March, nil
	case "j":
		return time.April, nil
	case "k":
		return time.May, nil
	case "m":
		return time.June, nil
	case "n":
		return time.July, nil
	case "q":
		return time.August, nil
	case "u":
		return time.September, nil
	case "v":
		return time.October, nil
	case "x":
		return time.November, nil
	case "z":
		return time.December, nil
	default:
		return 0, fmt.Errorf("unknown month code: %s", m)
	}
}
