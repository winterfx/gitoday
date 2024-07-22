package service

import (
	"gitrender/global"
	"testing"
)

func TestCrawl(t *testing.T) {
	global.Mode = "debug"
	res, err := Crawl(global.GoLang)
	if err != nil {
		t.Error(err)
	}
	for _, r := range res {
		t.Logf("%+v", *r)
	}
}
