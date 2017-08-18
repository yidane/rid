package http

import (
	"io/ioutil"
	"net/http"
	"time"
)

func NewRequest(method, url string, cookies []*http.Cookie) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Host", "rid.nxin.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Origin", "http://rid.nxin.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36") // nothing
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")                                                            // must
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Referer", "http://rid.nxin.com/login.html")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")

	if cookies != nil {
		for i := 0; i < len(cookies); i++ {
			req.AddCookie(cookies[i])
		}
	}

	return req, nil
}

func GetResponseContent(req *http.Request) ([]byte, []*http.Cookie, error) {
	client := &http.Client{}
	client.Timeout = time.Hour
	response, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}

	return data, response.Cookies(), nil
}
