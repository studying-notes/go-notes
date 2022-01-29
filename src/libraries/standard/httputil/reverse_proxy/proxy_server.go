/*
 * @Date: 2022.01.29 19:21
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.29 19:21
 */

package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewProxy 拿到 targetHost 后，创建一个反向代理
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	// 返回一个单主机代理对象
	return httputil.NewSingleHostReverseProxy(url), nil
}

// ProxyRequestHandler 使用 proxy 处理请求
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	// 返回一个代理方法
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	// 初始化反向代理并传入真正后端服务的地址（被代理的服务器）
	proxy, err := NewProxy("http://127.0.0.1:7070")
	if err != nil {
		panic(err)
	}

	// 使用 proxy 处理所有请求到你的服务
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
