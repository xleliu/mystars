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
	c := "```"
	return fmt.Sprintf(`Usage:

1. Get a access token from [github](https://github.com/settings/tokens)
2. Execute:
	%sbash
	go run github.com/xiaoler/mystars/cmd/main.go -o README.md -t {your token}
	%s 
`, c, c)
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
