# Certbot-hook-for-dnspod  

## 这是什么？  
一个以 `DNS-01` 方式自动化Let's Encrypt认证流程的辅助程序，适用于DNSPod用户，使用腾讯云API3.0编写。将它作为hook程序即可自动化Certbot的renew流程。  

## 构成  
- main.go: 一个使用腾讯云3.0API写的小程序，可以添加或删除DNSPod中的域名解析记录  
- /etc/dnspod/config.yml: 配置文件，储存腾讯云API ID&KEY以及自己的域名  
- authenticator.sh: hook脚本，调用上述go程序，添加Certbot生成的TXT记录，完成认证  
- cleanup.sh: 调用上述go程序，清理上面添加的TXT记录  


## 使用方法  
1. 编译并安装go程序
2. 利用hook脚本执行Certbot，可以利用Crontab实现自动renew  

## 调用示例  

```
sudo certbot certonly --dry-run --manual --manual-public-ip-logging-ok --preferred-challenges=dns --manual-auth-hook ./authenticator.sh --manual-cleanup-hook ./cleanup.sh -d subdomain.yourdomain.com
```