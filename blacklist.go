package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type Blacklist struct {
	PathList          []PathItem_b  `json:"pathlist"`
	ReferList         []ReferItem_b `json:"referlist"`
	AllowEmptyReferer bool          `json:"allowEmptyReferer,omitempty"`
}

var blacklist Blacklist

type PathItem_b struct {
	Paths []string `json:"paths"`
}

type ReferItem_b struct {
	Refer string `json:"refer"`
}

func isPathBlacklisted(path string) bool {
	if len(whitelist.ReferList) == 0 {
		return false
	}

	for _, item := range blacklist.PathList {
		for _, p := range item.Paths {
			match, err := regexp.MatchString(p, path)
			if err != nil {
				fmt.Printf("正则匹配错误：%s", err)
				continue
			}
			if match {
				return false
			}
		}
	}
	return true
}

func isRefererBlacklisted(referer string) bool {
	if len(whitelist.ReferList) == 0 {
		return false
	}
	for _, item := range blacklist.ReferList {
		if strings.Contains(referer, item.Refer) {
			return false
		}
	}
	return true
}

func loadBlacklist() {
	// 从JSON文件加载黑名单数据
	data, err := os.ReadFile("blacklist.json")
	if err != nil {
		fmt.Println("无法加载黑名单数据:", err)
		return
	}

	if err := json.Unmarshal(data, &blacklist); err != nil {
		fmt.Println("无法解析黑名单数据:", err)
		return
	}

	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			// 在定时任务中重新加载黑名单数据
			data, err := os.ReadFile("blacklist.json")
			if err != nil {
				fmt.Println("无法加载黑名单数据:", err)
				continue
			}

			if err := json.Unmarshal(data, &blacklist); err != nil {
				fmt.Println("无法解析黑名单数据:", err)
				continue
			}
		}
	}()
}
