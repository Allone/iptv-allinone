// Package liveurls
// @Time:2023/02/05 06:36
// @File:douyu.go
// @SoftWare:Goland
// @Author:feiyang
// @Contact:TG@feiyangdigital

package liveurls

import (
	"Golang/utils"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Douyu struct {
	Rid         string
	Stream_type string
}

func md5V3(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}

func (d *Douyu) GetRoomId() any {
	liveurl := "https://m.douyu.com/" + d.Rid
	client := &http.Client{}
	r, _ := http.NewRequest("GET", liveurl, nil)
	r.Header.Add("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1")
	r.Header.Add("upgrade-insecure-requests", "1")
	resp, _ := client.Do(r)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	roomidreg := regexp.MustCompile(`(?i)rid":(\d{1,8}),"vipId`)
	roomidres := roomidreg.FindStringSubmatch(string(body))
	if roomidres == nil {
		return nil
	}
	realroomid := roomidres[1]
	return realroomid
}

func (d *Douyu) GetRealUrl() any {
	var jsUtil = &utils.JsUtil{}
	did := "10000000000000000000000000001501"
	var timestamp = time.Now().Unix()
	var realroomid string
	rid := d.GetRoomId()
	if str, ok := rid.(string); ok {
		realroomid = str
	} else {
		return nil
	}
	liveurl := "https://www.douyu.com/" + realroomid
	client := &http.Client{}
	r, _ := http.NewRequest("GET", liveurl, nil)
	r.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	r.Header.Add("upgrade-insecure-requests", "1")
	resp, _ := client.Do(r)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	reg := regexp.MustCompile(`(?i)(vdwdae325w_64we[\s\S]*function ub98484234[\s\S]*?)function`)
	res := reg.FindStringSubmatch(string(body))
	nreg := regexp.MustCompile(`(?i)eval.*?;}`)
	strfn := nreg.ReplaceAllString(res[1], "strc;}")
	var funcContent1 []string
	funcContent1 = append(append(funcContent1, strfn), "ub98484234")
	result := jsUtil.JsRun(funcContent1, "ub98484234")
	nres := fmt.Sprintf("%s", result)
	nnreg := regexp.MustCompile(`(?i)v=(\d+)`)
	nnres := nnreg.FindStringSubmatch(nres)
	unrb := fmt.Sprintf("%v%v%v%v", realroomid, did, timestamp, nnres[1])
	rb := md5V3(unrb)
	nnnreg := regexp.MustCompile(`(?i)return rt;}\);?`)
	strfn2 := nnnreg.ReplaceAllString(nres, "return rt;}")
	strfn3 := strings.Replace(strfn2, `(function (`, `function sign(`, -1)
	strfn4 := strings.Replace(strfn3, `CryptoJS.MD5(cb).toString()`, `"`+rb+`"`, -1)
	var funcContent2 []string
	funcContent2 = append(append(funcContent2, strfn4), "sign")
	result2 := jsUtil.JsRun(funcContent2, realroomid, did, timestamp)
	param := fmt.Sprintf("%s", result2)
	realparam := param + "&rate=0"
	r1, n4err := http.Post("https://www.douyu.com/lapi/live/getH5Play/"+realroomid, "application/x-www-form-urlencoded", strings.NewReader(realparam))
	if n4err != nil {
		return nil
	}
	defer r1.Body.Close()
	body1, _ := io.ReadAll(r1.Body)
	var s1 map[string]any
	json.Unmarshal(body1, &s1)
	var flv_url string
	var rtmp_url string
	var rtmp_live string
	for k, v := range s1 {
		if k == "error" {
			if s1[k] != float64(0) {
				return nil
			}
		}
		if v, ok := v.(map[string]any); ok {
			for k, v := range v {
				if k == "rtmp_url" {
					if urlstr, ok := v.(string); ok {
						rtmp_url = urlstr
					}
				} else if k == "rtmp_live" {
					if urlstr, ok := v.(string); ok {
						rtmp_live = urlstr
					}
				}
			}
		}
	}
	flv_url = rtmp_url + "/" + rtmp_live
	n4reg := regexp.MustCompile(`(?i)(\d{1,8}[0-9a-zA-Z]+)_?\d{0,4}(.flv|/playlist)`)
	houzhui := n4reg.FindStringSubmatch(flv_url)
	var real_url string
	switch d.Stream_type {
	case "hls":
		real_url = strings.Replace(flv_url, houzhui[1]+".flv", houzhui[1]+".m3u8", -1)
	case "flv":
		real_url = flv_url
	case "xs":
		real_url = strings.Replace(flv_url, houzhui[1]+".flv", houzhui[1]+".xs", -1)
	}
	return real_url
}
