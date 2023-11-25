// Package liveurls
// @Time:2023/02/10 01:03
// @File:bilibili.go
// @SoftWare:Goland
// @Author:feiyang
// @Contact:TG@feiyangdigital

package liveurls

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

type BiliBili struct {
	Rid      string
	Line     string
	Quality  string
	Platform string
}

func (b *BiliBili) GetRealRoomID() any {
	var firstmap = make(map[string]any)
	var realroomid string
	apiurl := "https://api.live.bilibili.com/room/v1/Room/room_init?id=" + b.Rid
	client := &http.Client{}
	r, _ := http.NewRequest("GET", apiurl, nil)
	r.Header.Add("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1")
	resp, _ := client.Do(r)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &firstmap)
	if firstmap["msg"] == "直播间不存在" {
		return nil
	}
	if newmap, ok := firstmap["data"].(map[string]any); ok {
		if newmap["live_status"] != float64(1) {
			return nil
		} else {
			if flt, ok := newmap["room_id"].(float64); ok {
				realroomid = fmt.Sprintf("%v", int(flt))
			}
		}

	}
	return realroomid
}

func (b *BiliBili) GetPlayUrl() any {
	var roomid string
	var realurl string
	if str, ok := b.GetRealRoomID().(string); ok {
		roomid = str
	} else {
		return nil
	}
	client := &http.Client{}
	params := map[string]string{
		"room_id":  roomid,
		"protocol": "0,1",
		"format":   "0,1,2",
		"codec":    "0,1",
		"qn":       b.Quality,
		"platform": b.Platform,
		"ptype":    "8",
	}
	r, _ := http.NewRequest("GET", "https://api.live.bilibili.com/xlive/web-room/v2/index/getRoomPlayInfo", nil)
	q := r.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	r.URL.RawQuery = q.Encode()
	resp, _ := client.Do(r)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var json = string(body)
	value := gjson.Get(json, "data.playurl_info.playurl.stream")
	value.ForEach(func(key, value gjson.Result) bool {
		newvalue := gjson.Get(value.String(), "format.0.format_name")
		if newvalue.String() == "ts" {
			nnvalue := gjson.Get(value.String(), "format.#")
			valuelast := fmt.Sprintf("%v", nnvalue.Int()-1)
			codeclen := gjson.Get(value.String(), "format."+valuelast+".codec.#")
			codeclast := fmt.Sprintf("%v", codeclen.Int()-1)
			base_url := gjson.Get(value.String(), "format."+valuelast+".codec."+codeclast+".base_url")
			url_info := gjson.Get(value.String(), "format."+valuelast+".codec."+codeclast+".url_info")
			url_info.ForEach(func(key, value gjson.Result) bool {
				keyval := fmt.Sprintf("%v", key)
				switch b.Line {
				case "first":
					if keyval == "0" {
						host := gjson.Get(value.String(), "host")
						extra := gjson.Get(value.String(), "extra")
						realurl = fmt.Sprintf("%v%v%v", host, base_url, extra)
					}

				case "second":
					if keyval == "1" {
						host := gjson.Get(value.String(), "host")
						extra := gjson.Get(value.String(), "extra")
						realurl = fmt.Sprintf("%v%v%v", host, base_url, extra)
					}
				case "third":
					if keyval == "2" {
						host := gjson.Get(value.String(), "host")
						extra := gjson.Get(value.String(), "extra")
						realurl = fmt.Sprintf("%v%v%v", host, base_url, extra)
					}
				}

				return true
			})
		}
		return true
	})
	return realurl
}
