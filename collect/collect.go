package collect

import (
	"bufio"
	"fmt"
	"go-crawler/proxy"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Fetcher interface {
	Get(url *Request) ([]byte, error)
}

type BaseFetch struct {
}

func (BaseFetch) Get(req *Request) ([]byte, error) {
	resp, err := http.Get(req.Url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code: %d", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := DeterminEncoding(bodyReader)

	// 将 HTML 文本从特定编码转换为 UTF-8 编码
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return io.ReadAll(utf8Reader)
}

type BrowserFetch struct {
	Timeout time.Duration
	Proxy   proxy.ProxyFunc
	Logger  *zap.Logger
}

// 模拟浏览器访问
func (b BrowserFetch) Get(request *Request) ([]byte, error) {
	client := &http.Client{
		Timeout: b.Timeout,
	}

	// 代理服务器
	if b.Proxy != nil {
		transport := http.DefaultTransport.(*http.Transport)
		transport.Proxy = b.Proxy
		client.Transport = transport
	}

	req, err := http.NewRequest("GET", request.Url, nil)
	if err != nil {
		return nil, fmt.Errorf("get url failed:%v", err)
	}

	if len(request.Cookie) > 0 {
		req.Header.Set("Cookie", request.Cookie)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")

	resp, err := client.Do(req)
	time.Sleep(request.WaitTime)

	if err != nil {
		b.Logger.Error("fetch failed", zap.Error(err))
		return nil, err
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := DeterminEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return io.ReadAll(utf8Reader)
}

// 检测并返回当前 HTML 文本的编码格式
func DeterminEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)

	// 如果返回的 HTML 文本小于 1024 字节，我们认为当前 HTML 文本有问题，直接返回默认的 UTF-8 编码就好了
	if err != nil {
		fmt.Printf("fetch error:%v", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
