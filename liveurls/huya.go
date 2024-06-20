// Package liveurls
// @Time:2023/06/05 23:34
// @File:huya.go
// @SoftWare:Goland

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

	"github.com/tidwall/gjson"
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

var streamInfo = map[string]any{"flv": make(map[string]string), "hls": make(map[string]string)}

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
	TimeLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		TimeLocation = time.FixedZone("CST", 8*60*60)
	}
	now := time.Now().In(TimeLocation)
	q, _ := url.ParseQuery(antiCode)
	q.Set("t", "102")
	q.Set("ctype", "tars_mp")
	q.Set("wsTime", strconv.FormatInt(time.Now().Unix()+21600, 16))
	q.Set("ver", "1")
	q.Set("sv", now.Format("2006010215"))
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

func format(jsonStr string, uid int) map[string]any {
	cdnType := map[string]string{"HY": "hycdn", "AL": "alicdn", "TX": "txcdn", "HW": "hwcdn", "HS": "hscdn", "WS": "wscdn"}
	ojsonStr := gjson.Get(jsonStr, "roomInfo.tLiveInfo.tLiveStreamInfo.vStreamInfo").String()
	qreg := regexp.MustCompile(`(?i){"_proto"[\s\S]*?"value":([\s\S]*),"_classname"`)
	qres := qreg.FindStringSubmatch(ojsonStr)
	gjson.Parse(qres[1]).ForEach(func(_, value gjson.Result) bool {
		sFlvUrl := value.Get("sFlvUrl").String()
		sFlvUrlSuffix := value.Get("sFlvUrlSuffix").String()
		sHlsUrl := value.Get("sHlsUrl").String()
		sHlsUrlSuffix := value.Get("sHlsUrlSuffix").String()
		sStreamName := value.Get("sStreamName").String()
		sCdnType := value.Get("sCdnType").String()
		sFlvAntiCode := value.Get("sFlvAntiCode").String()
		sHlsAntiCode := value.Get("sHlsAntiCode").String()
		if sFlvUrl != "" {
			streamInfo["flv"].(map[string]string)[cdnType[sCdnType]] = sFlvUrl + "/" + sStreamName + "." + sFlvUrlSuffix + "?" + processAntiCode(sFlvAntiCode, uid, sStreamName)
		}
		if sHlsUrl != "" {
			streamInfo["hls"].(map[string]string)[cdnType[sCdnType]] = sHlsUrl + "/" + sStreamName + "." + sHlsUrlSuffix + "?" + processAntiCode(sHlsAntiCode, uid, sStreamName)
		}
		return true
	})
	return streamInfo
}

func (h *Huya) GetLiveUrl() any {
	liveurl := "https://m.huya.com/" + h.Rid
	client := &http.Client{}
	r, _ := http.NewRequest("GET", liveurl, nil)
	r.Header.Add("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1")
	r.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	resp, _ := client.Do(r)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	str := string(body)
	freg := regexp.MustCompile(`(?i)<script>[\s\S]window.HNF_GLOBAL_INIT = ([\s\S]*?) </script>`)
	res := freg.FindStringSubmatch(str)
	if len(res) > 1 && !strings.Contains(res[1], "\"exceptionType\":0") {
		jsonStr := res[1]
		liveStatus := gjson.Get(jsonStr, "roomInfo.eLiveStatus").Int()
		var mediaurl any

		if liveStatus == 2 {
			realurl := format(jsonStr, uid)
			if h.Type == "display" {
				return realurl
			}
			for k, v := range realurl {
				if k == h.Media {
					if urlarr, ok := v.(map[string]string); ok {
						for k, v := range urlarr {
							if k == h.Cdn {
								mediaurl = strings.Replace(v, "http://", "https://", 1)
							}
						}
					}
				}
			}
		} else if liveStatus == 3 {
			liveLineUrl := gjson.Get(jsonStr, "roomProfile.liveLineUrl").String()
			if liveLineUrl != "" {
				decodedLiveLineUrl, _ := base64.StdEncoding.DecodeString(liveLineUrl)
				mediaurl = "https:" + string(decodedLiveLineUrl)
			}
		} else {
			mediaurl = nil
		}
		return mediaurl
	} else if strings.Contains(res[1], "\"exceptionType\":0") {
		var h5info any
		ostr, _ := getContent("https://www.huya.com/" + h.Rid)
		nstr := string(ostr)
		lreg := regexp.MustCompile(`(?i)<script>[\s\S]*hyPlayerConfig =[\s\S]*stream: ([\s\S]*)};[\s\S]*window.TT_LIVE_TIMING`)
		lres := lreg.FindStringSubmatch(nstr)
		gjson.Get(lres[1], "data").ForEach(func(key, value gjson.Result) bool {
			if strings.Contains(value.String(), "gameStreamInfoList") {
				cdnType := map[string]string{"HY": "hycdn", "AL": "alicdn", "TX": "txcdn", "HW": "hwcdn", "HS": "hscdn", "WS": "wscdn"}
				gjson.Get(value.String(), "gameStreamInfoList").ForEach(func(_, value gjson.Result) bool {
					sFlvUrl := value.Get("sFlvUrl").String()
					sFlvUrlSuffix := value.Get("sFlvUrlSuffix").String()
					sHlsUrl := value.Get("sHlsUrl").String()
					sHlsUrlSuffix := value.Get("sHlsUrlSuffix").String()
					sStreamName := value.Get("sStreamName").String()
					sCdnType := value.Get("sCdnType").String()
					sFlvAntiCode := value.Get("sFlvAntiCode").String()
					sHlsAntiCode := value.Get("sHlsAntiCode").String()
					if sFlvUrl != "" {
						streamInfo["flv"].(map[string]string)[cdnType[sCdnType]] = sFlvUrl + "/" + sStreamName + "." + sFlvUrlSuffix + "?" + processAntiCode(sFlvAntiCode, uid, sStreamName)
					}
					if sHlsUrl != "" {
						streamInfo["hls"].(map[string]string)[cdnType[sCdnType]] = sHlsUrl + "/" + sStreamName + "." + sHlsUrlSuffix + "?" + processAntiCode(sHlsAntiCode, uid, sStreamName)
					}
					return true
				})
				if h.Type == "display" {
					h5info = streamInfo
				} else {
					for k, v := range streamInfo {
						if k == h.Media {
							if urlarr, ok := v.(map[string]string); ok {
								for k, v := range urlarr {
									if k == h.Cdn {
										h5info = strings.Replace(v, "http://", "https://", 1)
									}
								}
							}
						}
					}
				}
				return false
			}
			return true
		})
		return h5info
	}

	return nil
}
