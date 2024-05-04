package lib

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func ConvertIntToThousandString(num int) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", num)
}

func ConvertFloatToThousandString(num float64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%f", num)
}
