package crawler_test

import (
	"net/url"
	"testing"

	crawler "gitee.com/ntshibin/crawler-go"
)

func TestMain(t *testing.T) {
	t.Run("TestGet", GetTest)
	t.Run("TestProxy", ProxyTest)
}

func ProxyTest(t *testing.T) {
	crawler := crawler.NewCrawler("https://www.google.com.hk/").SetTimeout(5).SetProxy("http://localhost:7890")
	body := url.Values{
		"q": {"golang"},
	}
	_, err := crawler.Execute("Get", "/search", body)
	if err != nil {
		t.Errorf("使用 Get 访问 https://www.google.com.hk/search?q=golang %s", err)
	} else {
		t.Log("使用代理 Get 访问 https://www.google.com.hk/search?q=golang ok")
	}
}

func GetTest(t *testing.T) {
	crawler := crawler.NewCrawler("https://www.baidu.com").SetTimeout(5)
	body := url.Values{
		"wd": {"golang"},
	}
	_, err := crawler.Execute("Get", "/s", body)
	if err != nil {
		t.Errorf("使用 Get 访问 www.baidu.com/s?wd=golang %s", err)
	} else {
		t.Log("使用 Get 访问 www.baidu.com/s?wd=golang ok")
	}
}
