package service

import (
	"bytes"
	"fmt"
	"gitoday/global"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

var path = "https://github.com/trending"

type Repo struct {
	Name      string
	Url       string
	Desc      string
	Lang      string
	Star      string
	Fork      string
	TodayStar string
}

func Crawl(lang global.Language) ([]*Repo, error) {
	url := path
	if lang != global.All {
		url = fmt.Sprintf("%s/%s?since=daily", path, lang)
	}
	body, err := fetch(url)
	if err != nil {
		err := errors.Wrap(err, "fetch error")
		return nil, err
	}
	res, err := parse(body)
	if err != nil {
		err := errors.Wrap(err, "parse error")
		return nil, err
	}
	return res, nil
}

func parse(body []byte) ([]*Repo, error) {
	buf := bytes.NewBuffer(body)
	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return nil, err
	}
	repoList := make([]*Repo, 0)
	doc.Find("article.Box-row").Each(func(i int, s *goquery.Selection) {
		repo := &Repo{}
		href, ok := s.Find("h2.h3.lh-condensed").Find("a").Attr("href")
		if !ok {
			return
		}
		name := strings.SplitN(href, "/", 2)[1]
		repo.Name = name
		repo.Url = "https://www.github.com" + href
		repo.Desc = s.Find("p.col-9.color-fg-muted.my-1.pr-4").Text()
		repo.Lang = s.Find("div.f6.color-fg-muted.mt-2").Find("span.d-inline-block.ml-0.mr-3").Find("span").Text()
		repo.Star = s.Find(fmt.Sprintf("a[href=\"/%s/stargazers\"]", name)).Text()
		repo.Fork = s.Find(fmt.Sprintf("a[href=\"/%s/forks\"]", name)).Text()
		repo.TodayStar = s.Find("span.d-inline-block.float-sm-right").Text()
		repo.format()
		repoList = append(repoList, repo)
	})
	return repoList, nil
}

func fetch(url string) ([]byte, error) {
	var read io.Reader
	if global.IsPreviewMode() {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		//read from a file path
		filePath := filepath.Join(wd, "service", "debug.html")
		r, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		read = r
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status code is not 200")
		}
		read = resp.Body
	}

	body, err := io.ReadAll(read)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func (r *Repo) format() {
	r.Desc = strings.TrimSpace(strings.ReplaceAll(r.Desc, "\n", ""))
	r.Star = strings.TrimSpace(strings.ReplaceAll(r.Star, "\n", ""))
	r.Fork = strings.TrimSpace(strings.ReplaceAll(r.Fork, "\n", ""))
	r.TodayStar = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(r.TodayStar, "stars today", ""), "\n", ""))
}
