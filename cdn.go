package main

import (
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func TencentCloudInit() (*cdn.Client, error) {

	// 实例化一个认证对象
	credential := common.NewCredential(
		secretID,
		secretKey,
	)

	// 实例化一个client选项
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"

	// 实例化要请求产品的client对象
	client, err := cdn.NewClient(credential, "", cpf)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewStopCdnDomainRequests(cdn_domain string) (string, error) {
	// 调用TencentCloudInit函数获取初始化后的client对象
	client, err := TencentCloudInit()
	if err != nil {
		panic(err)
	}

	// 实例化一个请求对象
	request := cdn.NewStopCdnDomainRequest()
	request.Domain = common.StringPtr(cdn_domain)

	// 发送请求并处理响应
	response, err := client.StopCdnDomain(request)
	if err != nil {
		panic(err)
	}

	return response.ToJsonString(), nil
}
