# Certbot-hook-for-dnspod  

## 这是什么？  
一个以 `DNS-01` 方式自动化Let's Encrypt认证流程的辅助程序，适用于DNSPod用户，使用腾讯云API 3.0编写。将它作为hook程序即可自动化Certbot的renew流程。  
如果不HOOK的话，需要到DNSPoD控制台手动添加TXT解析记录，非常不方便，也无法自动续订，本程序实现了HOOK程序，只需配置好ID和KEY，即可利用腾讯云API自动添加并清理TXT记录，使Certbot的认证流程自动化。  

## 构成  
- `main.go`  
一个使用腾讯云3.0 API写的小程序，可以添加或删除DNSPod中的域名解析记录  
- `/etc/dnspod/config.yml`  
配置文件，储存腾讯云API ID&KEY以及自己的域名  
- `authenticator.sh`  
hook脚本，调用上述go程序，添加Certbot生成的TXT记录，完成认证  
- `cleanup.sh`  
hook脚本，调用上述go程序，清理上面添加的TXT记录  


## 使用方法  
1. 编译并安装go程序
2. 到腾讯云控制台申请ID和KEY并配置到 `/etc/dnspod/config.yml` ，具体格式参考  `config.yml`
3. 利用hook脚本执行Certbot，可以利用Crontab实现自动renew  

## 调用示例  
**注意：实际使用需去掉 `--dry-run` 参数**
- 借助hook脚本申请证书  
```bash
sudo certbot certonly --dry-run --manual --manual-public-ip-logging-ok --preferred-challenges=dns --manual-auth-hook ./authenticator.sh --manual-cleanup-hook ./cleanup.sh -d subdomain.yourdomain.com
```

- 借助hook脚本续订所有证书  
```bash
sudo certbot renew --dry-run --manual --manual-public-ip-logging-ok --preferred-challenges=dns --manual-auth-hook ./authenticator.sh --manual-cleanup-hook ./cleanup.sh
```

- 借助hook脚本续订特定证书  
```bash
sudo certbot renew --dry-run --manual --manual-public-ip-logging-ok --preferred-challenges=dns --manual-auth-hook ./authenticator.sh --manual-cleanup-hook ./cleanup.sh --cert-name subdomain.yourdomain.com
```

## 一种自动化renew配置方案  
- 将go程序编译好安装到PATH  
- 配置好 `/etc/dnspod/config.yml` 中的 `secretId`、`secretKey`、`domain`  
- 将两个HOOK脚本放到 `/etc/dnspod/`  
- 新建 `/etc/cron.daily/certbot-renew`，并加上**可执行权限**，这样即可每天自动检查并renew所有子域名下证书，并在更新证书之后重启nginx  
    ```bash
    # !/bin/bash
    certbot renew --manual --manual-public-ip-logging-ok --preferred-challenges=dns --manual-auth-hook /etc/dnspod/authenticator.sh --manual-cleanup-hook /etc/dnspod/cleanup.sh --post-hook "systemctl reload nginx" --quiet
    ```