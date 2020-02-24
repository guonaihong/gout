## how to use gout
## 设置http body
### [写入JSON到http body里面](./01-color-json.go)
```SetJSON```接口，写入json至http.body里

### [设置和解析xml](./11-xml.go)
使用```SetXML```写入xml格式数据，使用```BindXML```读取
### [设置和解析yaml](./12-yaml.go)
使用```SetYAML```写入yaml格式数据，使用```BindYAML```读取
### [设置formdata](./13-form-data.go)
使用```SetForm```接口接入form-data格式，该接口支持多种数据类型。
### [upload file](./14-upload-file.go)
这里使用的是```SetBody```进行上传文件

### [设置非json/xml/yaml数据到body](./04a-set-body.go)
需要使用```SetBody```接口

### [解析非json/xml/yaml数据body到变量里面](./04b-bind-body.go)
需要使用```BindBody```接口

## 设置控制字段
### [设置Query Parameters](./02-query.go)
需要使用```SetQuery```接口，该接口支持多种数据类型

### [设置http header](./03a-set-header.go)
需要使用```SetHeader```接口，丰富的数据类型让你停不下来

### [解析http header](./03b-bind-header.go)
需要```BindHeader```接口，支持多类型自动绑定


### [设置timeout](./05a-timeout.go)
需使用```SetTimeout```接口

### [取消一个正在发送的请求](./05b-cancel.go)
使用```WithContext```接口可取消正在发送的请求

### [发送x-www-form-urlencoded格式数据](./06-x-www-form-urlencoded.go)
使用```SetWWWForm```接口

### [处理多种body，一个接口既可以处理json，也可以处理html 404](./07-callback.go)
使用```Callback```接口

### [设置cookie](./08-cookie.go)
使用```SetCookies```接口，可传一个或者多个cookie

### [修改传输成为unixsocket](./09-unix.go)

### debug 模式
####  [打开debug模式or 打开debug关闭颜色高亮](./10a-debug.go)
都是使用```Debug()```接口，只是里面传递的策略函数不一样

#### [自定义debug模式](./10b-debug-custom.go)
debug接口具有强大的扩展性能，简单啪啪两下写个策略函数，就可以扩展该接口，比如设置某个环境变量才打开debug接口
#### [trace 功能，主要诊断接口各个阶段的性能](./10c-debug-trace.go)
```Debug()```里面传递 ```gout.Trace()```策略函数就可以打开这个功能

#### [保存debug信息](./15-debug-save-file.go)
自定义```Debug()```接口的策略函数

### 压测功能
#### [压测一定次数](./16a-benchmark-number.go)
```Number()```控制次数
#### [压测固定时间](./16b-benchmark-duration.go)
```Durations()```控制时间
#### [已某个固定频率压测](./16c-benchmark-rate.go)
```Rate()```控制压测频率
#### [压测功能和apache ab 的对比，远比ab性能要好](./16d-benchmark-vs-ab.go)
与apache ab的性能pk
#### [基于回调函数的自定义压测模式](./16e-customize-bench.go)
```Loop()```接口可传递回调函数
### import
#### [导入纯文本请求并发送](./17a-import-rawhttp.go)
```RawText()```接口可完成该功能
### export
#### [生成curl命令](./18a-gen-curl.go)
```Curl().Do()```可实现

### 指数回退重试
#### [重试](./19a-retry.go)
```Retry()```下面的接口
#### [使用冷备地址进行重试](./19b-retry-customize-backup.go)
```Func()```可传回调函数进行自定义设置
#### [基于某个http code进行重试，比如ES 返回209告知资源不可用](19c-retry-httpcode.go)
```Func()``` 的入参有code信息，使用```filter.ErrRetry```告知gout需要重试
#### [socks5](./20-socks5.go)
