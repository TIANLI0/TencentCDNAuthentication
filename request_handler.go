package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessRecord struct {
	Count          int
	LastAccessTime time.Time
}

var (
	mutex       sync.Mutex
	ipAccessMap map[string]AccessRecord
)

func init() {
	ipAccessMap = make(map[string]AccessRecord)
}

func handleRequest(c *gin.Context) {

	ip := c.Request.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.ClientIP()
	}
	if !isPathWhitelisted(c.Request.URL.Path) || !isRefererWhitelisted(c.Request.Referer()) || isPathBlacklisted(c.Request.URL.Path) || isRefererBlacklisted(c.Request.Referer()) || !ip_QPS(ip) || !QPS() {
		c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden."})
	} else {

		c.JSON(http.StatusOK, gin.H{"message": "OK.Tianli's CDN is working."})
		fmt.Println("QPS:", ipAccessMap["QPS"].Count)
	}
}

// IP QPS限制函数

func ip_QPS(ip string) bool {
	if ip == "" {
		return true
	}
	mutex.Lock()
	defer mutex.Unlock()

	currentTime := time.Now()

	// 删除一分钟之前的访问记录
	for accessTime, record := range ipAccessMap {
		if currentTime.Sub(record.LastAccessTime) > time.Minute {
			delete(ipAccessMap, accessTime)
		}
	}

	// 记录当前的访问次数和时间
	record, exists := ipAccessMap[ip]
	if !exists {
		record = AccessRecord{}
	}
	record.Count++
	record.LastAccessTime = currentTime
	ipAccessMap[ip] = record

	return record.Count <= ip_qps_int
}

// QPS限制函数，只需要记录单位时间内的访问次数是否超过阈值即可
func QPS() bool {
	mutex.Lock()
	defer mutex.Unlock()

	currentTime := time.Now()

	// 删除一分钟之前的访问记录
	for accessTime, record := range ipAccessMap {
		if currentTime.Sub(record.LastAccessTime) > time.Minute {
			delete(ipAccessMap, accessTime)
		}
	}

	// 记录当前的访问次数和时间
	record, exists := ipAccessMap["QPS"]
	if !exists {
		record = AccessRecord{}
	}
	record.Count++
	record.LastAccessTime = currentTime
	ipAccessMap["QPS"] = record
	go func() {
		if record.Count > max_qps_int {
			fmt.Println("QPS超过阈值，当前QPS为：", record.Count)
			// 限制NewStopCdnDomainRequests函数的调用频率
			time.Sleep(60 * time.Second)
			result, _ := NewStopCdnDomainRequests(cdn_domains)
			fmt.Println("请检查是否有人在攻击你的CDN！", result)
		}
	}()

	return record.Count <= qps_int
}
