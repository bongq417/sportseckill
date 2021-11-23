package main

import (
	"flag"
	"fmt"
	"github.com/bongq417/sportseckill/badminton"
	"github.com/bongq417/sportseckill/tool"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var offsetDay = flag.Duration("offsetDay", 4, "偏移天数0-4,默认4")
var start = flag.String("start", "", "开始时间，如：08:59:00")
var runtime = flag.Duration("runtime", 60, "运行时长，默认60s")
var hallTime = flag.String("hall", "18:00,20:00", "场次范围【13:00,14:00】、【20:00,22:00】")
var showId = flag.String("showId", "753", "三楼竞训馆=753，一楼综合馆=752")
var username = flag.String("username", "", "姓名")
var mobile = flag.String("mobile", "", "手机号")
var idCard = flag.String("card", "", "身份证")

func init() {
	flag.Parse()
}

func main() {
	// 参数校验
	if *start == "" || *username == "" || *mobile == "" || *idCard == "" {
		fmt.Println("开始时间、姓名、手机号、身份证都不能为空")
		return
	}
	client := badminton.NewClient()
	startTime := tool.StringToTime(tool.FormatTime(0) + " " + *start)
	if startTime.Before(time.Now()) {
		startTime = startTime.Add(time.Hour * 24)
	}
	endTime := startTime.Add(time.Second * (*runtime))
	var buyInfo badminton.BuyInfo
	buyInfo.Username = *username
	buyInfo.Mobile = *mobile
	buyInfo.IdCard = *idCard
	buyInfo.ShowId = *showId
	buyInfo.HallTime = *hallTime
	buyInfo.OffsetDay = *offsetDay
	tool.Info("[程序]开始运行", startTime, endTime, buyInfo.HallTime, tool.FormatTime(buyInfo.OffsetDay*24*time.Hour))
	// 等待启动
	client.WaitStart(startTime)
	// 提前dns缓存
	dnsCache("https://xihuwenti.juyancn.cn")
	success := false
	count := 0

	for hallId := 1; hallId < 6; {
		if success {
			break
		}
		GetBadminton(client, buyInfo, hallId, &success)
		if hallId < 6 {
			hallId++
		}
		if hallId == 6 {
			hallId = 1
			count++
			time.Sleep(time.Millisecond * 50)
		}
		if endTime.Before(time.Now()) {
			tool.Info("[程序]运行超时，结束程序")
			break
		}
	}
}

func GetBadminton(client *badminton.Client, buyInfo badminton.BuyInfo, hallId int, success *bool) {
	data := url.Values{}
	// https://xihuwenti.juyancn.cn/wechat/product/details?id=752&time=1631548800
	data.Set("show_id", buyInfo.ShowId)
	data.Set("date", tool.FormatTime(buyInfo.OffsetDay*24*time.Hour))
	// 一楼综合馆 752 86-91 表示为1号场地~6号场地
	// 三楼竞训馆 753 80-85 表示为1号场地~6号场地
	hallId += 80
	if buyInfo.ShowId == "752" {
		hallId += 5
	}
	//data.Set("data[]", sport+",20:00,22:00")
	data.Set("data[]", strconv.Itoa(hallId)+","+buyInfo.HallTime)
	var saveResult badminton.SaveResult
	client.DoPost("https://xihuwenti.juyancn.cn/wechat/product/save", data, &saveResult)
	tool.Info("[Save]", data, saveResult)
	if saveResult.Code != 0 {
		return
	}
	data2 := url.Values{}
	data2.Set("show_id", buyInfo.ShowId)
	data2.Set("username", buyInfo.Username)
	data2.Set("mobile", buyInfo.Mobile)
	data2.Set("smscode", "")
	data2.Set("id_card", buyInfo.IdCard)
	data2.Set("certType", "10001")
	data2.Set("param", saveResult.Msg)
	data2.Set("activityid", "0")
	data2.Set("couponId", "0")
	var addResult badminton.AddResult
	client.DoPost("https://xihuwenti.juyancn.cn/wechat/order/add", data2, &addResult)
	tool.Info("[Add]", addResult, data2)
	if addResult.OrderNum != "" {
		*success = true
	}
}

func dnsCache(url string) {
	http.Get("https://xihuwenti.juyancn.cn")
}
