

> 最多只支持三级域名



```shell
certbot certonly --manual --preferred-challenges=dns 
--manual-auth-hook  "/path/certbot --action=add --accessKey_id=阿里云ACCESS_KEY_ID --accessKey_secret=阿里云ACCESS_KEY_SECRET"  
--manual-cleanup-hook "/path/certbot --action=del --accessKey_id=阿里云ACCESS_KEY_ID --accessKey_secret=阿里云ACCESS_KEY_SECRET"
-d *.test.ainiok.com --register-unsafely-without-email
```

