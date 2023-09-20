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
	ip_qps_int  int
	qps_int     int
	max_qps_int int
	cdn_domains string
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
	ip_qps_str := os.Getenv("ip_qps")

	ip_qps_int, err = strconv.Atoi(ip_qps_str)
	if err != nil {
		fmt.Println("ip_qps转换失败，请检查.env文件")
		ip_qps_int = 1000
	}

	qps_str := os.Getenv("qps")
	qps_int, err = strconv.Atoi(qps_str)
	if err != nil {
		fmt.Println("qps转换失败，请检查.env文件")
		qps_int = 10000
	}

	max_qps_str := os.Getenv("max_qps")
	max_qps_int, err = strconv.Atoi(max_qps_str)
	if err != nil {
		fmt.Println("max_qps转换失败，请检查.env文件")
		max_qps_int = 1000000
	}

	if secretID == "" || secretKey == "" {
		fmt.Println("secretID或secretKey为空，请检查.env文件")
	}

	cdn_domains = os.Getenv("cdn_domain")
	if cdn_domains == "" {
		fmt.Println("cdn_domain为空，请检查.env文件")
	}
}
