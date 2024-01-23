// Package liveurls
// @Time:2024/01/12 22:00
// @File:huya.go
// @SoftWare:Goland


package liveurls

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Huya struct {
	Rid     string
	Cdn     string
	CdnType string
}

func MD5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

func parseAntiCode(antiCode, streamName string) string {
	qr, _ := url.ParseQuery(antiCode)
	t := "0"
	f := strconv.FormatInt(time.Now().UnixNano()/100, 10)
	wsTime := qr.Get("wsTime")

	decodeString, _ := base64.StdEncoding.DecodeString(qr.Get("fm"))
	fm := string(decodeString)
	fm = strings.ReplaceAll(fm, "$0", t)
	fm = strings.ReplaceAll(fm, "$1", streamName)
	fm = strings.ReplaceAll(fm, "$2", f)
	fm = strings.ReplaceAll(fm, "$3", wsTime)

	return fmt.Sprintf("wsSecret=%s&wsTime=%s&u=%s&seqid=%s&txyp=%s&fs=%s&sphdcdn=%s&sphdDC=%s&sphd=%s&u=0&t=100&ratio=0",
		MD5([]byte(fm)), wsTime, t, f, qr.Get("txyp"), qr.Get("fs"), qr.Get("sphdcdn"), qr.Get("sphdDC"), qr.Get("sphd"))
}

func (h *Huya) GetLiveUrl() any {
	liveurl := "https://m.huya.com/" + h.Rid
	client := &http.Client{}
	r, _ := http.NewRequest("GET", liveurl, nil)
	r.Header.Add("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1")
	r.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	resp, _ := client.Do(r)
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	reg := regexp.MustCompile("<script> window.HNF_GLOBAL_INIT = (.*)</script>")
	matches := reg.FindStringSubmatch(string(result))
	if matches == nil || len(matches) < 2 {
		return nil
	}
	return h.extractInfo(matches[1])
}

func (h *Huya) extractInfo(content string) any {
	parse := gjson.Parse(content)
	streamInfo := parse.Get("roomInfo.tLiveInfo.tLiveStreamInfo.vStreamInfo.value")
	if len(streamInfo.Array()) == 0 {
		return nil
	}
	var cdnSlice []string
	var finalurl string
	streamInfo.ForEach(func(key, value gjson.Result) bool {
		var cdnType = gjson.Get(value.String(), "sCdnType").String()
		cdnSlice = append(cdnSlice, cdnType)
		if cdnType == h.Cdn {
			urlStr := fmt.Sprintf("%s/%s.%s?%s",
				value.Get("sFlvUrl").String(),
				value.Get("sStreamName").String(),
				value.Get("sFlvUrlSuffix").String(),
				parseAntiCode(value.Get("sFlvAntiCode").String(), value.Get("sStreamName").String()))
			finalurl = strings.Replace(urlStr, "http://", "https://", 1)
		}
		return true
	})
	if h.CdnType == "display" {
		return cdnSlice
	}
	return finalurl
}
