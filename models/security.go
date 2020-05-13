package models

import (
	"gitlab.com/alphaticks/alphac/enum"
	"strings"
)

func (s *Security) Format(format string) string {
	switch s.SecurityType {
	case enum.SecurityType_CRYPTO_SPOT:
		res := format
		res = strings.Replace(res, "%B", strings.ToUpper(s.Underlying.Symbol), -1)
		res = strings.Replace(res, "%Q", strings.ToUpper(s.QuoteCurrency.Symbol), -1)
		res = strings.Replace(res, "%b", strings.ToLower(s.Underlying.Symbol), -1)
		res = strings.Replace(res, "%q", strings.ToLower(s.QuoteCurrency.Symbol), -1)
		return res
	default:
		return s.Symbol
	}
}
