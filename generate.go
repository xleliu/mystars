package mystars

import (
	"fmt"
	"time"
)

func Title() string {
	date := time.Now().Format(time.RFC1123Z)
	return "# My Stars\nUpdated at: " + date + "\n\n"
}

func Desc() string {
	return `Usage:

`
}

func Category(lang string) string {
	return "\n### " + lang + "\n"
}

func Repo(abst *Abstract) string {
	template := "- [%s](%s) - %s (`stars: %d`, `license: %s`)\n"

	return fmt.Sprintf(template,
		abst.Name,
		abst.Url,
		abst.Desc,
		abst.StarCount,
		abst.License,
	)
}
