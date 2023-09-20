package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type Whitelist struct {
	PathList          []PathItem  `json:"pathlist"`
	ReferList         []ReferItem `json:"referlist"`
	AllowEmptyReferer bool        `json:"allowEmptyReferer,omitempty"`
}

var whitelist Whitelist

type PathItem struct {
	Paths []string `json:"paths"`
}

type ReferItem struct {
	Refer string `json:"refer"`
}

func isPathWhitelisted(path string) bool {

	if len(whitelist.ReferList) == 0 {
		return true
	}

	for _, item := range whitelist.PathList {
		for _, p := range item.Paths {
			match, err := regexp.MatchString(p, path)
			if err != nil {
				fmt.Printf("正则匹配错误：%s", err)
				continue
			}
			if match {
				return true
			}
		}
	}
	return false
}

func isRefererWhitelisted(referer string) bool {

	if len(whitelist.ReferList) == 0 {
		return true
	}

	if whitelist.AllowEmptyReferer && referer == "" {
		return true
	}

	if len(whitelist.ReferList) == 0 {
		return true
	}

	for _, item := range whitelist.ReferList {
		if strings.Contains(referer, item.Refer) {
			return true
		}
	}
	return false
}

func loadWhitelist() {
	// 从JSON文件加载白名单数据
	data, err := os.ReadFile("whitelist.json")
	if err != nil {
		fmt.Println("无法加载白名单数据:", err)
		return
	}

	if err := json.Unmarshal(data, &whitelist); err != nil {
		fmt.Println("无法解析白名单数据:", err)
		return
	}

	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			// 在定时任务中重新加载白名单数据
			data, err := os.ReadFile("whitelist.json")
			if err != nil {
				fmt.Println("无法加载白名单数据:", err)
				continue
			}

			if err := json.Unmarshal(data, &whitelist); err != nil {
				fmt.Println("无法解析白名单数据:", err)
				continue
			}
		}
	}()
}
