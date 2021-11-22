package badminton

import (
	"encoding/json"
	"errors"
	"github.com/bongq417/sportseckill/tool"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: http.DefaultClient,
	}
}

func (client *Client) DoPost(uri string, data url.Values, result interface{}) (response *http.Response, err error) {
	return client.Request(POST, uri, data, result)
}

func (client *Client) Request(method string, uri string, data url.Values, result interface{}) (response *http.Response, err error) {
	// 创建请求对象
	request, err := http.NewRequest(method, uri, strings.NewReader(data.Encode()))
	if err != nil {
		tool.Info("[Request]NewRequest", response, err)
		return response, err
	}

	// 请求加签与header设置
	token := "CNZZDATA1274723626=1285222766-1629965061-%7C1630820660; latitude=30.28851890563965; longitude=120.03010559082031; WECHAT_OPENID=oAKYc01778qboqyrik4rnwZAmg1M; UM_distinctid=17bb480406b40-06c1baa9a2f02a8-5f7b1a2a-3d10d-17bb480406c73c; PHPSESSID=gn7af6rftonfnp8ctfm1pcpp72"
	SetHeaders(request, token)
	// 开始发送请求
	response, err = client.httpClient.Do(request)
	if err != nil {
		tool.Info("[Request]HttpClient.Do", response, err)
		return response, err
	}
	defer response.Body.Close()

	// 解析返回结果
	status := response.StatusCode
	message := response.Status
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response, err
	}
	responseBody := string(body)
	tool.Info("[Request]ResponseBody", responseBody)
	// 解析返回结果
	if status >= 200 && status < 300 {
		if body != nil && result != nil {
			err := json.Unmarshal(body, result)
			if err != nil {
				tool.Info("[Request]Json.Unmarshal", response, err)
				return response, err
			}
		}
		return response, nil
	} else {
		return response, errors.New(message)
	}
}

// SetHeaders 设置请求头
func SetHeaders(request *http.Request, token string) {
	request.Header.Add(ContentType, Application)
	request.Header.Set(UserAgent, ChromeBrowser)
	request.Header.Set(Charset, Utf8)
	request.Header.Set(Cookie, token)
}

func (client *Client) WaitStart(start time.Time) {
	now := time.Now()
	nowTime := tool.Timestamp(now)
	startTime := tool.Timestamp(start)
	// 如果等待时间大于1s，则开启等待，提早1s开启
	if (nowTime + 1000) < startTime {
		waitTime := startTime - nowTime - 1000
		tool.Info("[WaitStart]waitTime,start", waitTime, now)
		t := time.After(time.Duration(waitTime) * time.Millisecond)
		tool.Info("[WaitStart]waitTime,end", waitTime, <-t)
	} else {
		tool.Info("[WaitStart]no need waitTime", nowTime, startTime, (startTime-nowTime)/1000)
	}
}
