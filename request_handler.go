package main

import (
	"fmt"
	"net/http"
	"os"
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

	go func() {
		logRequest(c)
	}()

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

func logRequest(c *gin.Context) {

	go func() {
		// 检查并删除七天之前的日志文件
		files, err := os.ReadDir("./log")
		if err != nil {
			fmt.Println("读取日志文件失败：", err)
			return
		}
		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				fmt.Printf("获取文件信息失败：%s：%s\n", file.Name(), err)
				continue
			}
			if time.Since(fileInfo.ModTime()) > 7*24*time.Hour {
				err = os.Remove(fmt.Sprintf("./log/%s", file.Name()))
				if err != nil {
					fmt.Printf("删除文件失败：%s：%s\n", file.Name(), err)
				}
			}
		}
	}()

	os.Mkdir("./log", 0777)
	logPath := fmt.Sprintf("./log/%s.log", time.Now().Format("2006-01-02"))
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("创建日志文件失败：", err)
	}
	defer logFile.Close()
	logFile.WriteString(fmt.Sprintf("%s\n%s %s %s %s\n\n", time.Now(), c.Request.RemoteAddr, c.Request.Method, c.Request.URL, c.Request.Header))

}
