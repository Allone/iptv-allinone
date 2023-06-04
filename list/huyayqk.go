// Package list
// @Time:2023/06/02 10:00
// @File:mian.go
// @SoftWare:Goland
// @Author:feiyang
// @Contact:TG@feiyangdigital

package list

import (
	"io"
	"net/http"
)

type HuyaYqk struct {
}

type YaResponse struct {
	ITotalPage int `json:"iTotalPage"`
	IPageSize  int `json:"iPageSize"`
	VList      []struct {
		SAvatar180    string `json:"sAvatar180"`
		SGameFullName string `json:"sGameFullName"`
		SNick         string `json:"sNick"`
		LProfileRoom  int    `json:"lProfileRoom"`
	} `json:"vList"`
}

func (hy *HuyaYqk) HuYaYqk(requestURL string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
