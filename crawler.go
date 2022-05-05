package crawler

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Crawler struct {
	Client          *http.Client
	URL             string
	Timeout         int
	Proxy           string
	TLSClientConfig *tls.Config
	Headers         http.Header
}

// 新建一个Crawler实例
func NewCrawler(url string) *Crawler {

	return &Crawler{
		Client: &http.Client{
			Timeout: time.Duration(10) * time.Second,
		},
		URL:     url,
		Timeout: 10,
		Headers: http.Header{},
	}
}

// 配置证书
func (c *Crawler) SetCert(caCertPath, certFile, keyFile string) *Crawler {
	var pool *x509.CertPool
	var caCrt []byte // 根证书
	caCrt, err := ioutil.ReadFile(caCertPath)
	pool = x509.NewCertPool()
	if err != nil {
		panic(err)
	}
	pool.AppendCertsFromPEM(caCrt)
	var cliCrt tls.Certificate // 具体的证书加载对象
	cliCrt, err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	c.TLSClientConfig = &tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{cliCrt},
	}
	c.loadTransport()
	return c
}

// SkipCert 跳过证书验证
func (c *Crawler) SkipCert() *Crawler {
	c.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	c.loadTransport()
	return c
}

// 设置代理
func (c *Crawler) SetProxy(proxy string) *Crawler {
	c.Proxy = proxy
	c.loadTransport()
	return c
}

// 更新 Crawler 实例的 配置
func (c *Crawler) loadTransport() {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(c.Proxy)
	}
	c.Client.Transport = &http.Transport{
		TLSClientConfig: c.TLSClientConfig,
		Proxy:           proxy,
	}
}

// 设置超时时间
func (c *Crawler) SetTimeout(timeout int) *Crawler {
	c.Timeout = timeout
	c.Client.Timeout = time.Duration(c.Timeout) * time.Second
	return c
}

// 设置全局请求头
func (c *Crawler) AddHeader(field, value string) *Crawler {
	c.Headers.Add(field, value)
	return c
}

func (c *Crawler) LoadHeader(req *http.Request) {
	req.Header = c.Headers
}

func (c *Crawler) Execute(method, path string, data url.Values) (string, error) {
	// start := time.Now()
	fullPath := fmt.Sprintf("%s%s", c.URL, path)
	Url, _ := url.Parse(fullPath)
	var req *http.Request
	var err error
	var urlPath string
	switch strings.ToUpper(method) {
	case "GET":
		Url.RawQuery = data.Encode()
		urlPath = Url.String()
		req, err = http.NewRequest("GET", urlPath, nil)
	case "POST":
		urlPath = Url.String()
		req, err = http.NewRequest("POST", urlPath, bytes.NewBuffer([]byte(data.Encode())))
	default:
		return "", fmt.Errorf("[%s]方法不存在", method)
	}
	if err != nil {
		return "", err
	}
	c.LoadHeader(req)
	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("[%s](%s)失败：%s", method, urlPath, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("[%s](%s)解析失败：%s", method, urlPath, err)
	}
	// log.Printf("[%s](%s)\nresponse(%s): %s", method, urlPath, ElapsedTime(start), string(body))
	return string(body), err
}

func ElapsedTime(start time.Time) time.Duration {
	return time.Since(start)
}
