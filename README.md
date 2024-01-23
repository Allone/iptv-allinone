# **使用说明：**  
**如果你没有软路由或者服务器，那么推荐白嫖Vercel使用，[点击查看部署方法](https://github.com/papagaye744/iptv-go)！**
## 一、推荐使用Docker一键运行，并配置watchtower监听Docker镜像更新，直接一劳永逸：
### 1，使用Docker一键配置allinone
```
docker run -d --restart unless-stopped --privileged=true -p 35455:35455 --name allinone youshandefeiyang/allinone
```
### 2，一键配置watchtower每天凌晨两点自动监听allinone镜像更新，同步GitHub仓库：
```
docker run -d --name watchtower --restart unless-stopped -v /var/run/docker.sock:/var/run/docker.sock  containrrr/watchtower -c  --schedule "0 0 2 * * *"
```
## 二、直接运行：
首先去action中下载对应平台二进制执行文件，然后解压并直接执行
```
chmod 777 allinone && ./allinone
```
建议搭配进程守护工具进行使用，windows直接双击运行！
## 三、详细使用方法
## **虎牙、斗鱼、YY实时M3U获取：**
### 虎牙一起看：
```
http://你的IP:35455/huyayqk.m3u
```
### 斗鱼一起看：
```
http://你的IP:35455/douyuyqk.m3u
```
### YY轮播：
```
http://你的IP:35455/yylunbo.m3u
```
### 如果使需要自定义M3U文件中的前缀域名，可以传入url参数（需要注意的是，当域名中含有特殊字符时，需要对链接进行urlencode处理）：
```
http://你的IP:35455/xxxyqk.m3u?url=http://192.168.10.1:35455
```
## **抖音：**
### 默认最高画质，浏览器打开并复制`(live.douyin.com/)xxxxxx`，只需要复制后面的xxxxx即可（可选flv和hls两种种流媒体传输方式，默认flv）：
```
http://你的IP:35455/douyin/xxxxx(?stream=hls)
```
## **斗鱼：**
### 1，可选m3u8和flv以及xs三种流媒体传输方式【`(www.douyu.com/)xxxxxx 或 (www.douyu.com/xx/xx?rid=)xxxxxx`，默认flv】：
```
http://你的IP:35455/douyu/xxxxx(?stream=flv)
```
## **BiliBili`(live.bilibili.com/)xxxxxx`：**
### 1，平台platform参数选择（默认web，如果有问题，可以切换h5平台）：
```
"web"   => "桌面端"
"h5"    => "h5端"
```
### 2，线路line参数选择（默认线路二，如果卡顿/看不了，请切换线路一或者三，一般直播间只会提供两条线路，所以建议线路一/二之间切换）：
```
"first"  => "线路一"
"second" => "线路二"
"third"  => "线路三"
```
### 3，画质quality参数选择（默认原画，可以看什么画质去直播间看看，能选什么画质就能加什么参数，参数错误一定不能播放）：
```
"30000" => "杜比"
"20000" => "4K"
"10000" => "原画"
"400"   => "蓝光"
"250"   => "超清"
"150"   => "高清"
"80"    => "流畅"
```
### 4，最后的代理链接示例：
```
http://你的IP:35455/bilibili/xxxxxx(?platform=h5&line=first&quality=10000)
```
## **虎牙`(huya.com/)xxxxxx`：**  
### 1，查看CDN支持
```
http://你的IP:35455/huya/xxxxx?cdntype=display
```
### 2，切换CDN（默认HW，可选HY、AL、TX、HS、WS，具体支持哪些cdn，可以参考“1”）：
```
http://你的IP:35455/huya/xxxxx?cdn=AL
```
### 3，最后的代理链接示例：
```
http://你的IP:35455/huya/xxxxx(?cdn=xxx)
```
## **YouTube:**
```
https://www.youtube.com/watch?v=cK4LemjoFd0
Rid: cK4LemjoFd0
http://你的IP:35455/youtube/cK4LemjoFd0(?quality=1080/720...)
```
## **YY（默认最高画质，参数为4）:**
```
https://www.yy.com/xxxx
http://你的IP:35455/yy/xxxx(?quality=1/2/3/4...)
```
## 更多平台后续会酌情添加
