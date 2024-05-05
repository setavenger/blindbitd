package lib

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os/user"
	"strings"
)

func ConvertIntToThousandString(num int) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", num)
}

func ConvertFloatToThousandString(num float64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%f", num)
}

func ResolvePath(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	return strings.Replace(path, "~", dir, 1)
}
