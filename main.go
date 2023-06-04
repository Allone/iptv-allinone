// Package Golang
// @Time:2023/06/02 10:28
// @File:main.go
// @SoftWare:Goland
// @Author:feiyang
// @Contact:TG@feiyangdigital

package main

import (
	"Golang/list"
	"Golang/liveurls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/forgoer/openssl"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
)

func duanyan(adurl string, realurl any) string {
	var liveurl string
	if str, ok := realurl.(string); ok {
		liveurl = str
	} else {
		liveurl = adurl
	}
	return liveurl
}

func getTestVideoUrl(c *gin.Context) {
	fmt.Fprintln(c.Writer, "#EXTM3U")
	fmt.Fprintln(c.Writer, "#EXTINF:-1 tvg-name=\"4K60PSDR-H264-AAC测试\" tvg-logo=\"https://cdn.jsdelivr.net/gh/youshandefeiyang/IPTV/logo/tg.jpg\" group-title=\"4K频道\",4K60PSDR-H264-AAC测试")
	fmt.Fprintln(c.Writer, "http://159.75.85.63:5680/d/ad/h264/playad.m3u8")
	fmt.Fprintln(c.Writer, "#EXTINF:-1 tvg-name=\"4K60PHLG-HEVC-EAC3测试\" tvg-logo=\"https://cdn.jsdelivr.net/gh/youshandefeiyang/IPTV/logo/tg.jpg\" group-title=\"4K频道\",4K60PHLG-HEVC-EAC3测试")
	fmt.Fprintln(c.Writer, "http://159.75.85.63:5680/d/ad/playad.m3u8")
}

func getLivePrefix(c *gin.Context) string {
	firstUrl := c.DefaultQuery("url", "https://www.goodiptv.club")
	realUrl, _ := url.QueryUnescape(firstUrl)
	return realUrl
}

func setupRouter(adurl string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/douyin", func(c *gin.Context) {
		vrurl := c.Query("url")
		douyinobj := &liveurls.Douyin{}
		douyinobj.Shorturl = vrurl
		c.Redirect(http.StatusMovedPermanently, duanyan(adurl, douyinobj.GetRealurl()))
	})

	r.GET("/huyayqk.m3u", func(c *gin.Context) {
		yaobj := &list.HuyaYqk{}
		res, _ := yaobj.HuYaYqk("https://live.cdn.huya.com/liveHttpUI/getLiveList?iGid=2135")
		var result list.YaResponse
		json.Unmarshal(res, &result)
		pageCount := result.ITotalPage
		pageSize := result.IPageSize
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=huyayqk.m3u")
		getTestVideoUrl(c)

		for i := 1; i <= pageCount; i++ {
			apiRes, _ := yaobj.HuYaYqk(fmt.Sprintf("https://live.cdn.huya.com/liveHttpUI/getLiveList?iGid=2135&iPageNo=%d&iPageSize=%d", i, pageSize))
			var res list.YaResponse
			json.Unmarshal(apiRes, &res)
			data := res.VList
			for _, value := range data {
				fmt.Fprintf(c.Writer, "#EXTINF:-1 tvg-logo=\"%s\" group-title=\"%s\", %s\n", value.SAvatar180, value.SGameFullName, value.SNick)
				fmt.Fprintf(c.Writer, "%s/huya/%v\n", getLivePrefix(c), value.LProfileRoom)
			}
		}
	})

	r.GET("/douyuyqk.m3u", func(c *gin.Context) {
		yuobj := &list.DouYuYqk{}
		resAPI, _ := yuobj.Douyuyqk("https://www.douyu.com/gapi/rkc/directory/mixList/2_208/list")

		var result list.DouYuResponse
		json.Unmarshal(resAPI, &result)
		pageCount := result.Data.Pgcnt

		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=douyuyqk.m3u")
		getTestVideoUrl(c)

		for i := 1; i <= pageCount; i++ {
			apiRes, _ := yuobj.Douyuyqk("https://www.douyu.com/gapi/rkc/directory/mixList/2_208/" + strconv.Itoa(i))

			var res list.DouYuResponse
			json.Unmarshal(apiRes, &res)
			data := res.Data.Rl

			for _, value := range data {
				fmt.Fprintf(c.Writer, "#EXTINF:-1 tvg-logo=\"https://apic.douyucdn.cn/upload/%s_big.jpg\" group-title=\"%s\", %s\n", value.Av, value.C2name, value.Nn)
				fmt.Fprintf(c.Writer, "%s/douyu/%v\n", getLivePrefix(c), value.Rid)
			}
		}
	})

	r.GET("/yylunbo.m3u", func(c *gin.Context) {
		yylistobj := &list.Yylist{}
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=yylunbo.m3u")
		getTestVideoUrl(c)

		i := 1
		for {
			apiRes := yylistobj.Yylb(fmt.Sprintf("https://rubiks-idx.yy.com/nav/other/pnk1/448772?channel=appstore&compAppid=yymip&exposured=80&hdid=8dce117c5c963bf9e7063e7cc4382178498f8765&hostVersion=8.25.0&individualSwitch=1&ispType=2&netType=2&openCardLive=1&osVersion=16.5&page=%d&stype=2&supportSwan=0&uid=1834958700&unionVersion=0&y0=8b799811753625ef70dbc1cc001e3a1f861c7f0261d4f7712efa5ea232f4bd3ce0ab999309cac0d7869449a56b44c774&y1=8b799811753625ef70dbc1cc001e3a1f861c7f0261d4f7712efa5ea232f4bd3ce0ab999309cac0d7869449a56b44c774&y11=9c03c7008d1fdae4873436607388718b&y12=9d8393ec004d98b7e20f0c347c3a8c24&yv=1&yyVersion=8.25.0", i))
			var res list.ApiResponse
			json.Unmarshal([]byte(apiRes), &res)
			for _, value := range res.Data.Data {
				fmt.Fprintf(c.Writer, "#EXTINF:-1 tvg-logo=\"%s\" group-title=\"%s\", %s\n", value.Avatar, value.Biz, value.Desc)
				fmt.Fprintf(c.Writer, "%s/yy/%v\n", getLivePrefix(c), value.Sid)
			}
			if res.Data.IsLastPage == 1 {
				break
			}
			i++
		}
	})

	r.GET("/:path/:rid", func(c *gin.Context) {
		path := c.Param("path")
		rid := c.Param("rid")
		switch path {
		case "douyin":
			douyinobj := &liveurls.Douyin{}
			douyinobj.Rid = rid
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, douyinobj.GetDouYinUrl()))
		case "douyu":
			douyuobj := &liveurls.Douyu{}
			douyuobj.Rid = rid
			douyuobj.Stream_type = c.DefaultQuery("stream", "hls")
			douyuobj.Cdn_type = c.DefaultQuery("cdn", "akm-tct")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, douyuobj.GetRealUrl()))
		case "huya":
			huyaobj := &liveurls.Huya{}
			huyaobj.Rid = rid
			huyaobj.Cdn = c.DefaultQuery("cdn", "hwcdn")
			huyaobj.Media = c.DefaultQuery("media", "flv")
			huyaobj.Type = c.DefaultQuery("type", "nodisplay")
			if huyaobj.Type == "display" {
				c.JSON(200, huyaobj.GetLiveUrl())
			} else {
				c.Redirect(http.StatusMovedPermanently, duanyan(adurl, huyaobj.GetLiveUrl()))
			}
		case "bilibili":
			biliobj := &liveurls.BiliBili{}
			biliobj.Rid = rid
			biliobj.Platform = c.DefaultQuery("platform", "web")
			biliobj.Quality = c.DefaultQuery("quality", "10000")
			biliobj.Line = c.DefaultQuery("line", "second")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, biliobj.GetPlayUrl()))
		case "youtube":
			ytbObj := &liveurls.Youtube{}
			ytbObj.Rid = rid
			ytbObj.Quality = c.DefaultQuery("quality", "1080")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, ytbObj.GetLiveUrl()))
		case "yy":
			yyObj := &liveurls.Yy{}
			yyObj.Rid = rid
			yyObj.Quality = c.DefaultQuery("quality", "4")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, yyObj.GetLiveUrl()))
		}
	})
	return r
}

func main() {
//	key := []byte("6354127897263145")
//	defstr, _ := base64.StdEncoding.DecodeString("Mf5ZVkSUHH5xC9fH2Sao+2LgjRfydmzMgHNrVYX4AcSoI0nktkV7z1jSU6nSihf7ny+PexV73YjDoEtG7qu+Cw==")
//	defurl, _ := openssl.AesECBDecrypt(defstr, key, openssl.PKCS7_PADDING)


	defurl, _ := base64.StdEncoding.DecodeString("aHR0cDovL3R2LnhpbWl4LnVzOjg4L3Rlc3QubTN1OA==")
	r := setupRouter(string(defurl))
	r.Run(":35455")
}
