// Package liveurls
// @Time:2023/02/05 06:36
// @File:douyu.go
// @SoftWare:Goland
// @Author:feiyang
// @Contact:TG@feiyangdigital

package liveurls

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	js "github.com/dop251/goja"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Douyu struct {
	Rid         string
	Stream_type string
	Cdn_type    string
}

func md5V3(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}

func getDid() string {
	client := &http.Client{}
	timeStamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	url := "https://passport.douyu.com/lapi/did/api/get?client_id=25&_=" + timeStamp + "&callback=axiosJsonpCallback1"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1")
	req.Header.Set("referer", "https://m.douyu.com/")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	re := regexp.MustCompile(`axiosJsonpCallback1\((.*)\)`)
	match := re.FindStringSubmatch(string(body))
	var result map[string]map[string]string
	json.Unmarshal([]byte(match[1]), &result)
	return result["data"]["did"]
}

func (d *Douyu) GetRealUrl() any {
	did := getDid()
	var timestamp = time.Now().Unix()
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
	reg := regexp.MustCompile(`(?i)(function ub98484234.*)\s(var.*)`)
	res := reg.FindStringSubmatch(string(body))
	nreg := regexp.MustCompile(`(?i)eval.*;}`)
	strfn := nreg.ReplaceAllString(res[0], "strc;}")
	vm := js.New()
	_, err := vm.RunString(strfn)
	if err != nil {
		panic(err)
	}
	jsfn, ok := js.AssertFunction(vm.Get("ub98484234"))
	if !ok {
		panic("这不是一个函数")
	}
	result, nerr := jsfn(
		js.Undefined(),
		vm.ToValue("ub98484234"),
	)
	if nerr != nil {
		panic(nerr)
	}
	nres := fmt.Sprintf("%s", result)
	nnreg := regexp.MustCompile(`(?i)v=(\d+)`)
	nnres := nnreg.FindStringSubmatch(nres)
	unrb := fmt.Sprintf("%v%v%v%v", realroomid, did, timestamp, nnres[1])
	rb := md5V3(unrb)
	nnnreg := regexp.MustCompile(`(?i)return rt;}\);?`)
	strfn2 := nnnreg.ReplaceAllString(nres, "return rt;}")
	strfn3 := strings.Replace(strfn2, `(function (`, `function sign(`, -1)
	strfn4 := strings.Replace(strfn3, `CryptoJS.MD5(cb).toString()`, `"`+rb+`"`, -1)
	vm2 := js.New()
	_, nnerr := vm2.RunString(strfn4)
	if nnerr != nil {
		panic(nnerr)
	}
	jsfn2, nok := js.AssertFunction(vm2.Get("sign"))
	if !nok {
		panic("这不是一个函数")
	}
	result2, n3err := jsfn2(
		js.Undefined(),
		vm2.ToValue(realroomid),
		vm2.ToValue(did),
		vm2.ToValue(timestamp),
	)
	if n3err != nil {
		panic(n3err)
	}
	param := fmt.Sprintf("%s", result2)
	realparam := param + "&ver=22107261&rid=" + realroomid + "&rate=-1"
	r1, n4err := http.Post("https://m.douyu.com/api/room/ratestream", "application/x-www-form-urlencoded", strings.NewReader(realparam))
	if n4err != nil {
		panic(n4err)
	}
	defer r1.Body.Close()
	body1, _ := io.ReadAll(r1.Body)
	var s1 map[string]any
	json.Unmarshal(body1, &s1)
	var hls_url string
	for k, v := range s1 {
		if k == "code" {
			if s1[k] != float64(0) {
				return nil
			}
		}
		if v, ok := v.(map[string]any); ok {
			for k, v := range v {
				if k == "url" {
					if urlstr, ok := v.(string); ok {
						hls_url = urlstr
					}
				}
			}
		}
	}
	n4reg := regexp.MustCompile(`(?i)(\d{1,8}[0-9a-zA-Z]+)_?\d{0,4}(.m3u8|/playlist)`)
	houzhui := n4reg.FindStringSubmatch(hls_url)
	var real_url string
	flv_url := "http://" + d.Cdn_type + ".douyucdn.cn/live/" + houzhui[1] + ".flv?uuid="
	switch d.Stream_type {
	case "hls":
		real_url = hls_url
	case "flv":
		real_url = flv_url
	}
	return real_url
}
