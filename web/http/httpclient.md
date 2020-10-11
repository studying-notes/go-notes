# 标准库 net/http

## HTTP 客户端

### 基本的 HTTP/HTTPS 请求

**GET**

```go
func main() {
	var resp *http.Response
	var body []byte
	resp, _ = http.Get("http://www.baidu.com")
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

带参数的请求必须通过 net/url 标准库处理：

```go
func main() {
	var resp *http.Response
	var body []byte
	apiUrl := "http://www.baidu.com"
	data := url.Values{}
	data.Set("key", "value")
	u, _:=url.ParseRequestURI(apiUrl)
	u.RawQuery = data.Encode()
	resp, _ = http.Get(u.String())
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

**POST**

```go
func main() {
	var resp *http.Response
	var body []byte
	apiUrl := "http://www.baidu.com"
	contentType := "application/json"
	data := `{"key":"value","age":18}`
	resp, _ = http.Post(apiUrl, contentType, strings.NewReader(data))
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

**POSTFORM**

```go
func main() {
	var resp *http.Response
	var body []byte
	resp, _ = http.PostForm("http://www.baidu.com", url.Values{"key": {"value"}})
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

### 自定义 Client

管理 HTTP 客户端的 Headers、重定向策略等，必须创建一个 Client。

```go
func main() {
	var resp *http.Response
	var body []byte
	apiUrl := "http://www.baidu.com"
	client := &http.Client{CheckRedirect: nil}
	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; "+
		"Ubuntu; Linux i686 on x86_64) AppleWebKit/537.36 (KHTML, "+
		"like Gecko) Chrome/64.0.3282.140 Safari/537.36 OPR/50.0.2762.67")
	resp, _ = client.Do(req)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

### 自定义 Transport

管理代理、TLS 配置、keep-alive、压缩等，必须创建一个 Transport。

```go
func main() {
	var resp *http.Response
	var body []byte
	apiUrl := "http://www.baidu.com"
	transport := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: nil},
		DisableCompression: true,
	}
	client := &http.Client{Transport: transport}
	resp, _ = client.Get(apiUrl)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

```go

```

