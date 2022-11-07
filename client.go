package mystars

import (
	"context"
	"log"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

type Github struct {
	client   *github.Client
	NextPage int
	LastPage int
}

type Abstract struct {
	Name      string
	Url       string
	Desc      string
	StarCount int
	License   string
}

type Following struct {
	Url  string
	Name string
	Desc string
}

func GetLicense(l *github.License) string {
	if l == nil {
		return ""
	}
	return *l.Key
}

func GetString(ptr *string, placeholder string, limit int) string {
	if ptr == nil {
		return placeholder
	}
	// 限制长度
	rs := []rune(*ptr)
	if limit > 0 && len(rs) > limit {
		return string(rs[:limit]) + "……"
	}
	return *ptr
}

func NewGithub(ctx context.Context, accessToken string) *Github {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &Github{
		client:   client,
		NextPage: 1,
	}
}

func (g *Github) ResetPage() {
	g.NextPage, g.LastPage = 1, 0
}

func (g *Github) HasNext() bool {
	log.Printf("WIP… Next Page: %d, Last Page: %d\n", g.NextPage, g.LastPage)
	if g.NextPage == 0 {
		return false
	}
	return g.NextPage == 1 || g.NextPage <= g.LastPage
}

func (g *Github) MyStars(ctx context.Context) ([]*github.StarredRepository, error) {
	// list all repositories for the authenticated user
	opts := &github.ActivityListStarredOptions{
		// Sort:      "full_name",
		// Direction: "desc",
		ListOptions: github.ListOptions{
			Page: g.NextPage,
			// PerPage: 30,
		},
	}
	repos, resp, err := g.client.Activity.ListStarred(ctx, "", opts)
	g.NextPage, g.LastPage = resp.NextPage, resp.LastPage
	return repos, err
}

func (g *Github) MyFollowing(ctx context.Context) ([]*github.User, error) {
	opts := &github.ListOptions{
		Page: g.NextPage,
		// PerPage: 30,
	}
	users, resp, err := g.client.Users.ListFollowing(ctx, "", opts)
	g.NextPage, g.LastPage = resp.NextPage, resp.LastPage
	return users, err
}
