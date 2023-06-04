// Package liveurls
// @Time:2023/02/05 23:34
// @File:huya.go
// @SoftWare:Goland
// @Author:feiyang
// @Contact:TG@feiyangdigital

package liveurls

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Huya struct {
	Rid   string
	Media string
	Type  string
	Cdn   string
}

type Data struct {
}

type Payload struct {
	AppId   int    `json:"appId"`
	ByPass  int    `json:"byPass"`
	Context string `json:"context"`
	Version string `json:"version"`
	Data    Data   `json:"data"`
}

type ResponseData struct {
	Data struct {
		Uid string `json:"uid"`
	} `json:"data"`
}

func getContent(apiUrl string) ([]byte, error) {
	payload := Payload{
		AppId:   5002,
		ByPass:  3,
		Context: "",
		Version: "2.4",
		Data:    Data{},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(jsonPayload)))
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func getUid() string {
	content, _ := getContent("https://udblgn.huya.com/web/anonymousLogin")
	var responseData ResponseData
	json.Unmarshal(content, &responseData)
	uid := responseData.Data.Uid
	return uid
}

var uid, _ = strconv.Atoi(getUid())

func getUUID() int64 {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	randNum := rand.Intn(1000)
	return ((now % 10000000000 * 1000) + int64(randNum)) % 4294967295
}

func processAntiCode(antiCode string, uid int, streamName string) string {
	q, _ := url.ParseQuery(antiCode)
	q.Set("ver", "1")
	q.Set("sv", "2110211124")
	seqId := strconv.Itoa(uid + int(time.Now().UnixNano()/int64(time.Millisecond)))
	q.Set("seqid", seqId)
	q.Set("uid", strconv.Itoa(uid))
	q.Set("uuid", strconv.FormatInt(getUUID(), 10))
	h := md5.New()
	h.Write([]byte(seqId + "|" + q.Get("ctype") + "|" + q.Get("t")))
	ss := hex.EncodeToString(h.Sum(nil))
	fm, _ := base64.StdEncoding.DecodeString(q.Get("fm"))
	q.Set("fm", strings.Replace(strings.Replace(strings.Replace(strings.Replace(string(fm), "$0", q.Get("uid"), -1), "$1", streamName, -1), "$2", ss, -1), "$3", q.Get("wsTime"), -1))
	h.Reset()
	h.Write([]byte(q.Get("fm")))
	q.Set("wsSecret", hex.EncodeToString(h.Sum(nil)))
	q.Del("fm")
	if _, ok := q["txyp"]; ok {
		q.Del("txyp")
	}
	return q.Encode()
}

func format(liveArr map[string]any, uid int) map[string]any {
	streamInfo := map[string]any{"flv": make(map[string]string), "hls": make(map[string]string)}
	cdnType := map[string]string{"HY": "hycdn", "AL": "alicdn", "TX": "txcdn", "HW": "hwcdn", "HS": "hscdn", "WS": "wscdn"}
	for _, s := range liveArr["roomInfo"].(map[string]any)["tLiveInfo"].(map[string]any)["tLiveStreamInfo"].(map[string]any)["vStreamInfo"].(map[string]any)["value"].([]any) {
		sMap := s.(map[string]any)
		if sMap["sFlvUrl"] != nil {
			streamInfo["flv"].(map[string]string)[cdnType[sMap["sCdnType"].(string)]] = sMap["sFlvUrl"].(string) + "/" + sMap["sStreamName"].(string) + "." + sMap["sFlvUrlSuffix"].(string) + "?" + processAntiCode(sMap["sFlvAntiCode"].(string), uid, sMap["sStreamName"].(string))
		}
		if sMap["sHlsUrl"] != nil {
			streamInfo["hls"].(map[string]string)[cdnType[sMap["sCdnType"].(string)]] = sMap["sHlsUrl"].(string) + "/" + sMap["sStreamName"].(string) + "." + sMap["sHlsUrlSuffix"].(string) + "?" + processAntiCode(sMap["sHlsAntiCode"].(string), uid, sMap["sStreamName"].(string))
		}
	}
	return streamInfo
}

func (h *Huya) GetLiveUrl() any {
	var liveArr map[string]any
	liveurl := "https://m.huya.com/" + h.Rid
	client := &http.Client{}
	r, _ := http.NewRequest("GET", liveurl, nil)
	r.Header.Add("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1")
	r.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	resp, _ := client.Do(r)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	str := string(body)
	freg := regexp.MustCompile(`(?i)<script> window.HNF_GLOBAL_INIT = (.*) </script>`)
	res := freg.FindStringSubmatch(str)

	json.Unmarshal([]byte(res[1]), &liveArr)
	var mediaurl any
	if roomInfo, ok := liveArr["roomInfo"].(map[string]any); ok {
		if liveStatus, ok := roomInfo["eLiveStatus"].(float64); ok && liveStatus == 2 {
			realurl := format(liveArr, uid)
			if h.Type == "display" {
				return realurl
			}
			for k, v := range realurl {
				switch k {
				case h.Media:
					if urlarr, ok := v.(map[string]string); ok {
						for k, v := range urlarr {
							switch k {
							case h.Cdn:
								mediaurl = v
							}
						}
					}
				}
			}

		} else if liveStatus, ok := roomInfo["eLiveStatus"].(float64); ok && liveStatus == 3 {
			if roomProfile, ok := liveArr["roomProfile"].(map[string]any); ok {
				if liveLineUrl, ok := roomProfile["liveLineUrl"].(string); ok {
					decodedLiveLineUrl, _ := base64.StdEncoding.DecodeString(liveLineUrl)
					mediaurl = "http:" + string(decodedLiveLineUrl)
				}
			}
		} else {
			mediaurl = nil
		}
	}
	return mediaurl

}
