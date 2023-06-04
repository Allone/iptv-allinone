// Package liveurls
// @Time:2023/02/03 01:59
// @File:douyin.go
// @SoftWare:Goland
// @Author:feiyang
// @Contact:TG@feiyangdigital

package liveurls

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

type Douyin struct {
	Shorturl string
	Rid      string
}

func GetRoomId(url string) any {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	r, _ := http.NewRequest("GET", url, nil)
	r.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	r.Header.Add("authority", "v.douyin.com")
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	reurl := resp.Header.Get("Location")
	reg := regexp.MustCompile(`\d{19}`)
	res := reg.FindAllStringSubmatch(reurl, -1)
	if res == nil {
		return nil
	}
	return res[0][0]
}

func (d *Douyin) GetRealurl() any {
	var mediamap map[string]map[string]map[string]map[string]any
	var roomstatus map[string]map[string]map[string]int
	var roomid string
	if str, ok := GetRoomId(d.Shorturl).(string); ok {
		roomid = str
	} else {
		return nil
	}
	client := &http.Client{}
	r, _ := http.NewRequest("GET", "https://webcast.amemv.com/webcast/room/reflow/info/?type_id=0&live_id=1&room_id="+roomid+"&sec_user_id=&app_id=1128", nil)
	r.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	r.Header.Add("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1")
	resp, _ := client.Do(r)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	str, _ := url.QueryUnescape(string(body))
	json.Unmarshal([]byte(str), &mediamap)
	json.Unmarshal([]byte(str), &roomstatus)
	var realurl any
	if mediaslice, ok := mediamap["data"]["room"]["stream_url"]["hls_pull_url_map"].(map[string]any); ok {
		for k, v := range mediaslice {
			if k == "FULL_HD1" {
				realurl = fmt.Sprintf("%v", v)
			}
			if roomstatus["data"]["room"]["status"] != 2 {
				realurl = nil
			}
		}
	}
	return realurl
}

func (d *Douyin) GetDouYinUrl() any {
	liveurl := "https://live.douyin.com/" + d.Rid
	client := &http.Client{}
	var mediamap map[string]map[string]any
	re, _ := http.NewRequest("GET", liveurl, nil)
	re.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	re.Header.Add("upgrade-insecure-requests", "1")
	oresp, _ := client.Do(re)
	defer oresp.Body.Close()
	oreg := regexp.MustCompile(`(?i)__ac_nonce=(.*?);`)
	ores := oreg.FindStringSubmatch(oresp.Header["Set-Cookie"][0])
	r, _ := http.NewRequest("GET", liveurl, nil)
	cookie1 := &http.Cookie{Name: "__ac_nonce", Value: ores[1]}
	r.AddCookie(cookie1)
	r.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	r.Header.Add("upgrade-insecure-requests", "1")
	resp, _ := client.Do(r)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	str, _ := url.QueryUnescape(string(body))
	reg := regexp.MustCompile(`(?i)\"roomid\"\:\"([0-9]+)\"`)
	res := reg.FindAllStringSubmatch(str, -1)
	if res == nil {
		return nil
	}
	nnreg := regexp.MustCompile(`(?i)\"id_str\":\"` + res[0][1] + `(?i)\"[\s\S]*?\"hls_pull_url\"`)
	nnres := nnreg.FindAllStringSubmatch(str, -1)
	nnnreg := regexp.MustCompile(`(?i)\"hls_pull_url_map\"[\s\S]*?}`)
	nnnres := nnnreg.FindAllStringSubmatch(nnres[0][0], -1)
	if nnnres == nil {
		return nil
	}
	json.Unmarshal([]byte(`{`+nnnres[0][0]+`}`), &mediamap)
	return mediamap["hls_pull_url_map"]["FULL_HD1"]
}
