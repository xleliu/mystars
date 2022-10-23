package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"

	"github.com/xiaoler/mystars"
)

var (
	token  = flag.String("t", "", "The github personal access token")
	mdFile = flag.String("o", "README.md", "The markdoan file output to……")
)

func main() {
	var accessToken string
	if *token != "" {
		accessToken = *token
	} else {
		accessToken = os.Getenv("GITHUB_TOKEN")
	}
	if accessToken == "" {
		log.Fatalln("Please set a github access token……")
	}
	ctx := context.Background()
	github := mystars.NewGithub(ctx, accessToken)
	var collect = make(map[string][]*mystars.Abstract)

	for github.HasNext() {
		repos, err := github.MyStars(ctx)
		if err != nil {
			log.Println(err)
			return
		}
		for _, repo := range repos {
			abst := &mystars.Abstract{
				Name:      *repo.Repository.FullName,
				Url:       *repo.Repository.HTMLURL,
				StarCount: *repo.Repository.StargazersCount,
				Desc:      mystars.GetString(repo.Repository.Description, ""),
				License:   mystars.GetLicense(repo.Repository.License),
			}
			lang := mystars.GetString(repo.Repository.Language, "Others")
			collect[lang] = append(collect[lang], abst)
		}
	}

	file, err := os.OpenFile(*mdFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(mystars.Title())
	writer.WriteString(mystars.Desc())

	for lang, items := range collect {
		if lang == "Others" {
			continue
		}
		writer.WriteString(mystars.Category(lang))
		for _, abst := range items {
			writer.WriteString(mystars.Repo(abst))
		}
	}
	// 最后写入其他
	writer.WriteString(mystars.Category("Others"))
	for _, abst := range collect["Others"] {
		writer.WriteString(mystars.Repo(abst))
	}
	writer.Flush()
}

func getString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
