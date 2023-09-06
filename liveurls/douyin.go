// Package liveurls
// @Time:2023/02/03 01:59
// @File:douyin.go
// @SoftWare:Goland
// @Author:feiyang
// @Contact:TG@feiyangdigital

package liveurls

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

type Douyin struct {
	Stream string
	Rid    string
}

func (d *Douyin) GetDouYinUrl() any {
	liveurl := "https://live.douyin.com/" + d.Rid
	client := &http.Client{}
	re, _ := http.NewRequest("GET", liveurl, nil)
	re.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	re.Header.Add("upgrade-insecure-requests", "1")
	oresp, _ := client.Do(re)
	defer oresp.Body.Close()
	oreg := regexp.MustCompile(`(?i)__ac_nonce=(.*?);`)
	ores := oreg.FindStringSubmatch(oresp.Header["Set-Cookie"][0])
	r, _ := http.NewRequest("GET", liveurl, nil)
	cookie := &http.Cookie{Name: "__ac_nonce", Value: ores[1]}
	r.AddCookie(cookie)
	r.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	r.Header.Add("upgrade-insecure-requests", "1")
	resp, _ := client.Do(r)
	defer resp.Body.Close()
	reg := regexp.MustCompile(`(?i)ttwid=.*?;`)
	res := reg.FindStringSubmatch(fmt.Sprintf("%s", resp.Cookies()))[0]
	url := "https://live.douyin.com/webcast/room/web/enter/?aid=6383&app_name=douyin_web&live_id=1&device_platform=web&language=zh-CN&enter_from=web_live&cookie_enabled=true&screen_width=1728&screen_height=1117&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Chrome&browser_version=116.0.0.0&web_rid=" + d.Rid
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", res)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "live.douyin.com")
	req.Header.Add("Connection", "keep-alive")
	ress, _ := client.Do(req)
	defer ress.Body.Close()
	body, _ := io.ReadAll(ress.Body)
	var json = string(body)
	status, _ := strconv.Atoi(fmt.Sprintf("%s", gjson.Get(json, "data.data.0.status")))
	if status != 2 {
		return nil
	}
	var realurl string
	value := gjson.Get(json, "data.data.0.stream_url.live_core_sdk_data.pull_data.stream_data")
	value.ForEach(func(key, value gjson.Result) bool {
		switch d.Stream {
		case "flv":
			realurl = fmt.Sprintf("%s", gjson.Get(value.String(), "data.origin.main.flv"))
		case "hls":
			realurl = fmt.Sprintf("%s", gjson.Get(value.String(), "data.origin.main.hls"))
		}
		return true
	})
	return realurl
}
