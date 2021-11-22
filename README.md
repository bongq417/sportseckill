# sportseckill

抢文体中心场地

liunx构建如下命令
> GOOS=linux GOARCH=amd64 go build
>
参数：
> ./sportseckill -showId 752 -offsetDay 2ns -start '2021-11-22 09:00:00' -hall '16:00,17:00' -runtime 30ns
> 
注意：
> 记得替换client.Request里的token值