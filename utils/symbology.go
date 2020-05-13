package utils

import (
	"fmt"
	"strings"
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
