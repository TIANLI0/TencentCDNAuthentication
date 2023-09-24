package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	secretID    string
	secretKey   string
	ip_QPM_int  int
	QPM_int     int
	max_QPM_int int
	cdn_domains string
	maxTraffic  float64
)

func loadconfig() {
	// 从.env文件中读取secretID和secretKey
	err := godotenv.Load()
	if err != nil {
		colorRedBold := "\033[1;31m"
		colorReset := "\033[0m"
		fmt.Println(colorRedBold + "无法加载.env文件" + colorReset)
		os.Exit(1)
		return
	}
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		colorRedBold := "\033[1;31m"
		colorReset := "\033[0m"
		fmt.Println(colorRedBold + "未找到.env文件" + colorReset)
		os.Exit(1)
		return
	}
	secretID = os.Getenv("secretID")
	secretKey = os.Getenv("secretKey")
	ip_QPM_str := os.Getenv("ip_QPM")

	ip_QPM_int, err = strconv.Atoi(ip_QPM_str)
	if err != nil {
		fmt.Println("ip_QPM转换失败，请检查.env文件")
		ip_QPM_int = 1000
	}

	QPM_str := os.Getenv("QPM")
	QPM_int, err = strconv.Atoi(QPM_str)
	if err != nil {
		fmt.Println("QPM转换失败，请检查.env文件")
		QPM_int = 10000
	}

	max_QPM_str := os.Getenv("max_QPM")
	max_QPM_int, err = strconv.Atoi(max_QPM_str)
	if err != nil {
		fmt.Println("max_QPM转换失败，请检查.env文件")
		max_QPM_int = 1000000
	}

	if secretID == "" || secretKey == "" {
		fmt.Println("secretID或secretKey为空，请检查.env文件")
	}

	cdn_domains = os.Getenv("cdn_domain")
	if cdn_domains == "" {
		fmt.Println("cdn_domain为空，请检查.env文件")
	}

	maxTraffic_str := os.Getenv("max_traffic")
	maxTraffic, err = strconv.ParseFloat(maxTraffic_str, 64)
	if err != nil {
		fmt.Println("maxTraffic转换失败，请检查.env文件")
		maxTraffic = 1000000
	}
}
