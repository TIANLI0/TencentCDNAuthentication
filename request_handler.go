package main

import (
	"bufio"
	"fmt"
	"io"
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
	mutex         sync.RWMutex
	ipAccessMap   map[string]AccessRecord
	ipAccessSlice []string
	logMutex      sync.Mutex
	istrigger     bool
	dbFile        = "CDN.DB"
	traffic       float64
)

func init() {
	ipAccessMap = make(map[string]AccessRecord)
	ipAccessSlice = make([]string, 0)
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
		go func() {
			if Traffic(c.Request.URL.Path) {
				fmt.Println("流量超过阈值，已停止CDN服务")
				result, _ := NewStopCdnDomainRequests(cdn_domains)
				fmt.Println("请检查是否有人在攻击你的CDN！", result)
				// 限制NewStopCdnDomainRequests函数的调用频率
				time.Sleep(60 * time.Second)
			}
		}()
		c.JSON(http.StatusOK, gin.H{"message": "OK.Tianli's CDN is working."})
		fmt.Println("QPS:", getQPS())
	}
}

func ip_QPS(ip string) bool {
	if ip == "" {
		return true
	}

	mutex.Lock()
	defer mutex.Unlock()

	currentTime := time.Now()

	// 删除一分钟之前的访问记录
	for _, accessTime := range ipAccessSlice {
		record := ipAccessMap[accessTime]
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

func QPS() bool {
	mutex.Lock()
	defer mutex.Unlock()

	currentTime := time.Now()

	// 删除一分钟之前的访问记录
	for _, accessTime := range ipAccessSlice {
		record := ipAccessMap[accessTime]
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
			result, _ := NewStopCdnDomainRequests(cdn_domains)
			fmt.Println("请检查是否有人在攻击你的CDN！", result)
			// 限制NewStopCdnDomainRequests函数的调用频率
			time.Sleep(60 * time.Second)
		}
	}()

	return record.Count <= qps_int
}

func getQPS() int {
	mutex.RLock()
	defer mutex.RUnlock()

	record, exists := ipAccessMap["QPS"]
	if !exists {
		return 0
	}
	return record.Count
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

	logMutex.Lock()
	defer logMutex.Unlock()

	os.Mkdir("./log", 0777)
	logPath := fmt.Sprintf("./log/%s.log", time.Now().Format("2006-01-02"))
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("创建日志文件失败：", err)
		return
	}
	defer logFile.Close()
	logFile.WriteString(fmt.Sprintf("%s\n%s %s %s %s\n\n", time.Now(), c.Request.RemoteAddr, c.Request.Method, c.Request.URL, c.Request.Header))
}

func Traffic(path string) bool {
	go func() {
		if !istrigger {
			time.Sleep(24 * time.Hour)
			istrigger = true
			// 清零流量统计
			traffic = 0.0
		}
	}()

	traffic_0 := readDB(path)
	if traffic_0 == 0.0 {
		bodySize := getBodySize(path)
		recordDB(path, bodySize)
		traffic += float64(bodySize) / 1024.0 / 1024.0
	} else {
		traffic += traffic_0
	}

	fmt.Printf("估计流量：%.2fMB 设置阈值:%.2fMB ", traffic, maxTraffic)
	return traffic >= maxTraffic
}

// 读取CDN.DB文件的数据
func readDB(path string) float64 {
	file, err := os.OpenFile(dbFile, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("打开DB失败:", err)
		return 0
	}
	defer file.Close()
	var bodySize int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		linePath, err := fmt.Sscanf(line, "%s %d\n", &path, &bodySize)
		if err != nil {
			fmt.Println("读取失败:", err)
			return 0.0
		}
		if linePath != 2 {
			fmt.Println("读取失败: 数据格式错误")
			return 0.0
		}

		// 找到匹配的路径后返回大小
		if path == "desired_path" {
			return float64(bodySize) / 1024 / 1024
		}
	}

	if scanner.Err() != nil {
		fmt.Println("读取失败:", scanner.Err())
		return 0.0
	}

	return 0.0
}

// 记录相关数据到CDN.DB
func recordDB(path string, bodySize int) {
	file, err := os.OpenFile(dbFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("打开DB失败:", err)
		return
	}
	defer file.Close()

	record := fmt.Sprintf("%s %d\n", path, bodySize)
	_, err = file.WriteString(record)
	if err != nil {
		fmt.Println("写入DB失败:", err)
		return
	}
}

// 获取请求地址的body大小
func getBodySize(path string) int {
	// 组装请求地址
	url := fmt.Sprintf("http://%s%s", cdn_domains, path)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return 0
	}
	defer resp.Body.Close()

	// 读取response的body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("请求失败:", err)
		return 0
	}

	return len(bodyBytes)
}
