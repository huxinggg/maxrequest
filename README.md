# maxrequest

#### 介绍
轻量实用的golang请求库，不附带任何第三方库，完全使用标准库实现，让retry变得简单，让并发变得简单。

#### 安装教程

go.mod文件加入以下代码，然后执行go mod tidy
```
github.com/HuXingGG/maxrequest latest
```


#### 发起GET请求
```azure
package main

import (
	"fmt"
	"gitee.com/justin0218/maxrequest"
)

type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func main() {
	var ret Result
	url := "http://xxx.xx"
	resp, body, err := maxrequest.New().Get(url).Result(&ret)
	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
	fmt.Println(err)
	fmt.Printf("%+v\n", ret)
}
```

#### 发起POST请求
```azure

package main

import (
	"fmt"
	"gitee.com/justin0218/maxrequest"
)

type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type PostData struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func main() {
	var req PostData
	req.UserName = "测试"
	req.Password = "xxsssss"
	var ret Result
	url := "http://xxx.xx"
	resp, body, err := maxrequest.New().Post(url, req).Result(&ret)
	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
	fmt.Println(err)
	fmt.Printf("%+v\n", ret)
}


```

#### form-data请求

只需要把"Post()"改成"PostForm()"

```azure
func main() {
	var req PostData
	req.UserName = "测试"
	req.Password = "xxsssss"
	var ret Result
	url := "http://xxx.xx"
	resp, body, err := maxrequest.New().SetTimeout(10).SetRetry(3,3, nil).PostForm(url, req).Result(&ret)
	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
	fmt.Println(err)
	fmt.Printf("%+v\n", ret)
}
```

#### 设置超时

设置了10秒超时
```azure
resp, body, err := maxrequest.New().SetTimeout(10).Post(url, req).Result(&ret)
```
#### 设置重试
SetRetry第一个参数为重试次数,第二个参数为重试周期，单位：秒
```azure

package main

    import (
	"fmt"
	"gitee.com/justin0218/maxrequest"
)

type Result struct {
	Code int         `json:"code" maxrequestRetry:"404"` //代表返回code=404时，将会进行重试，多个字段设置maxrequestRetry,其中一个满足就会进行重试
	Data interface{} `json:"data"`
}

type PostData struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func main() {
	var req PostData
	req.UserName = "测试"
	req.Password = "xxsssss"
	var ret Result
	url := "http://xxx.xx"
	resp, body, err := maxrequest.New().SetTimeout(10).SetRetry(3, 3, func(retryNotice maxrequest.RetryNotice) bool {
		fmt.Printf("重试次数：%d", retryNotice.Num)
		return false //return true 可以随时终止重试
	}).Post(url, req).Result(&ret)
	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
	fmt.Println(err)
	fmt.Printf("%+v\n", ret)
}

```

#### 并发请求
同时发起多个请求，一次性拿到所有结果
```azure

func main() {
	var req PostData
	req.UserName = "测试"
	req.Password = "xxsssss"
	url := "http://xxx.xx"
	r1 := maxrequest.New().SetTimeout(10).SetRetry(3, 3, nil).Post(url, req)
	r2 := maxrequest.New().SetTimeout(10).SetRetry(3, 3, nil).Post(url, req)
	r3 := maxrequest.New().SetTimeout(10).SetRetry(3, 3, nil).Post(url, req)
	r4 := maxrequest.New().SetTimeout(10).SetRetry(3, 3, nil).Post(url, req)
	r5 := maxrequest.New().SetTimeout(10).SetRetry(3, 3, nil).Post(url, req)
	results := maxrequest.Go(r1, r2, r3, r4, r5)
	fmt.Println("所有请求结果：%v", results)
}

```
#### 注意
1. 只有在Result()方法内传入带有maxrequestRetry的tag才会触发重试机制；
2. Result()传入的结构体，只会解析json tag，如需解析其他数据类型，例如：xml，Result内请传入nil，然后自己根据返回的body自行解析；
3. 如需发送非json请求，例如xml，可把Post方法的第二个参数直接传入bytes；
4. 并发的结果和发送的请求顺序相同，雷同于JavaScript的Promise.all。

