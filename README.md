# TencentCDNAuthentication
简单的腾讯云鉴权服务，支持QPS IP_QPS 等限制，也可在高响应次数的场景下自动关闭CDN，适用于腾讯云的远程鉴权服务。

注意：上线生成环境前要测试CDN回源节点连通性！！！由于腾讯云CDN限制，只支持单CDN熔断。

![img](https://qcloudimg.tencent-cloud.cn/image/document/d26215edda04745d6fdee7d68ef64cc5.jpg?1)

## How To Use?

1. 从release下载适用于对应版本的二进制可执行文件。
2. 部署至服务器，根据说明配置.env文件。
3. 配置blacklist.json文件或者whitelist.json文件（如果未配置其中一个文件，则将判断为运行，比如我没有配置whitelist，那么将所有路径及referer视为白名单，blacklist同理）。具体配置请参考仓库。
4. 开放相应端口或者启用反向代理，程序将运行在`5276`端口。
5. CDN启用相关配置（请求方法使用GET，根据需求配置其他选项，为防止程序崩溃或者服务崩溃影响您的业务，建议将超时时长调低）![img](https://qcloudimg.tencent-cloud.cn/image/document/b9a476dda2f433adc8dc49d0d263d4aa.png)

Tips: 当然，你可以将超时设置调的尽量低，仅仅将此程序视为面对大量访问情况下的分析工具，触发max_qps后将会自动关停腾讯云CDN。


## .env配置说明

| 配置项     | 示例       | 是否必填 | 说明                                                         |
| ---------- | :--------- | -------- | :----------------------------------------------------------- |
| secretID   | 114514     | 否       | 腾讯云secretID，需前往[官网](https://console.cloud.tencent.com/cam/capi)获取，注意KEY安全 |
| secretKey  | 114514     | 否       | 腾讯云secretKey，需前往[官网](https://console.cloud.tencent.com/cam/capi)获取，注意KEY安全 |
| cdn_domain | 114514.com | 否       | 需要自动关闭的CDN站点域名                                    |
| qps        | 10000      | 是       | 一分钟内的QPS限制，超过限制将返回403                         |
| ip_qps     | 500        | 是       | 单IP一分钟内的QPS限制，超过限制将返回403                     |
| max_qps    | 500000     | 是       | 一分钟内的QPS限制，超过限制将关闭CDN，如果未配置secret建议将此项配置尽量高 |
