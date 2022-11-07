package main

import (
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

	following := getMyFollowingGroup(ctx, github)
	stars, langs := getMyStars(ctx, github)

	file, err := os.OpenFile(mdFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	g := mystars.NewGenerater(file)
	g.Title()
	g.Desc()
	g.Separator()
	// 写入group
	g.Category("Following Organizations")
	for _, f := range following {
		g.Org(f.Name, f.Url, f.Desc)
	}
	g.Separator()
	// 按语言顺序写入
	for _, lang := range langs {
		g.Category(lang)
		for _, abst := range stars[lang] {
			g.Repo(abst)
		}
	}
	g.Flush()

	log.Println("DONE!")
}

func getMyStars(ctx context.Context, github *mystars.Github) (map[string][]*mystars.Abstract, []string) {
	log.Println("Get Starred Projects……")

	var (
		collect = make(map[string][]*mystars.Abstract)
		langs   = []string{}
	)
	defer github.ResetPage()

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

	return collect, langs
}

func getMyFollowingGroup(ctx context.Context, github *mystars.Github) []*mystars.Following {
	log.Println("Get Following Organizations……")

	var following = []*mystars.Following{}
	defer github.ResetPage()

	for github.HasNext() {
		users, err := github.MyFollowing(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		for _, user := range users {
			if *user.Type == "User" {
				continue
			}
			following = append(following, &mystars.Following{
				Url:  *user.HTMLURL,
				Name: *user.Login,
				Desc: *user.Bio,
			})
		}
	}

	return following
}
