package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"sort"

	"github.com/xiaoler/mystars"
)

func main() {
	var accessToken, mdFile string

	flag.StringVar(&accessToken, "t", os.Getenv("GITHUB_TOKEN"), "The github personal access token")
	flag.StringVar(&mdFile, "o", "README.md", "The markdoan file output to……")
	flag.Parse()

	if accessToken == "" {
		log.Fatalln("Please set a github access token……")
	}

	ctx := context.Background()
	github := mystars.NewGithub(ctx, accessToken)

	var collect = make(map[string][]*mystars.Abstract)
	for github.HasNext() {
		repos, err := github.MyStars(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		for _, repo := range repos {
			abst := &mystars.Abstract{
				Name:      *repo.Repository.FullName,
				Url:       *repo.Repository.HTMLURL,
				StarCount: *repo.Repository.StargazersCount,
				Desc:      mystars.GetString(repo.Repository.Description, "", 200),
				License:   mystars.GetLicense(repo.Repository.License),
			}
			lang := mystars.GetString(repo.Repository.Language, "Others", 0)
			collect[lang] = append(collect[lang], abst)
		}
	}

	var langs = []string{}
	var others = false
	for l := range collect {
		if l == "Others" {
			others = true
			continue
		}
		langs = append(langs, l)
	}
	sort.Strings(langs)
	if others {
		langs = append(langs, "Others")
	}

	file, err := os.OpenFile(mdFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(mystars.Title())
	writer.WriteString(mystars.Desc())
	// 按语言顺序写入
	for _, lang := range langs {
		writer.WriteString(mystars.Category(lang))
		for _, abst := range collect[lang] {
			writer.WriteString(mystars.Repo(abst))
		}
	}
	writer.Flush()

	log.Println("DONE!")
}
