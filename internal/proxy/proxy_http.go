package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

func (p *Proxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	// 创建Transport
	//transport := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}

	director := func(req *http.Request) {
		if r.URL.Scheme == "" {
			req.URL.Scheme = "http"
		}

		//请求拦截器
		if p.reqHandler != nil {
			p.reqHandler(w, r)
		}

		req.URL.Host = r.URL.Host
		req.Host = r.Host
	}

	modifyResponse := func(resp *http.Response) error {
		if resp.StatusCode == http.StatusOK {
			contentType := resp.Header.Get("Content-Type")
			if strings.Contains(contentType, "text/html") {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				resp.Body.Close()

				//相应拦截器
				modifiedBody := body
				if p.respHandler != nil {
					p.respHandler(resp, body)
				}

				resp.Body = io.NopCloser(bytes.NewReader(modifiedBody))
				resp.ContentLength = int64(len(modifiedBody))
				resp.Header.Set("Content-Length", fmt.Sprint(len(modifiedBody)))
			}
		}
		return nil
	}

	// 创建反向代理
	proxy := &httputil.ReverseProxy{
		Director:       director,
		Transport:      transport,
		ModifyResponse: modifyResponse,
	}

	proxy.ServeHTTP(w, r)
}
