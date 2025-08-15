package server

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/gozstd"
	"io"
	"net/http"
	"net/url"
	resources "proxy-dev/assets"
	"proxy-dev/internal/config"
	"strings"
)

func LoadCert() tls.Certificate {
	certFile, _ := resources.ReadByte("server.crt")
	keyFile, _ := resources.ReadByte("server.key")
	cert, _ := tls.X509KeyPair(certFile, keyFile)
	return cert
}

func ShuntHandler(w http.ResponseWriter, r *http.Request) bool {
	for _, s := range config.Conf.Rule {
		if !s.Enable {
			continue
		}

		url, _ := url.Parse(s.Surl)
		prefixAddr := fmt.Sprintf("%s", url.Host)
		reqUrl := fmt.Sprintf("%s", r.URL.Host)

		if strings.Contains(reqUrl, prefixAddr) || strings.Contains(prefixAddr, reqUrl) {
			return true
		}
	}
	return false
}

func ReqHandler(w http.ResponseWriter, r *http.Request) bool {
	log.Infof("req: %s %s", r.Method, r.URL.String())
	LogFilter.Log(r.Method, r.URL.String())

	for _, s := range config.Conf.Rule {
		if !s.Enable {
			continue
		}

		if s.Type == config.PROXY_TYPE_REDIRECT {
			prefixAddr := EnsurePort(s.Surl)
			replaceAddr := EnsurePort(s.Turl)

			reqUrl := EnsurePort(r.URL.String())

			if strings.HasPrefix(reqUrl, prefixAddr) {
				newUrl := strings.ReplaceAll(reqUrl, prefixAddr, replaceAddr)
				newURL, _ := url.Parse(newUrl)

				// 修改请求URL
				r.URL = newURL
				r.Host = newURL.Host

				log.Infof("redirect: %s %s", r.Method, r.URL.String())
				return true
			}
		}
	}
	return false
}

// EnsurePort 确保URL有端口，如果没有则添加默认端口
func EnsurePort(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	if u.Port() != "" {
		return rawURL
	}

	switch strings.ToLower(u.Scheme) {
	case "http":
		u.Host = fmt.Sprintf("%s:80", u.Hostname())
	case "https":
		u.Host = fmt.Sprintf("%s:443", u.Hostname())
	}
	return u.String()
}

func ResHandler(w *http.Response, body []byte) []byte {
	for _, s := range config.Conf.Rule {
		if !s.Enable {
			continue
		}
		if s.Type != config.PROXY_TYPE_RESMOD {
			continue
		}

		prefixAddr := EnsurePort(s.Surl)
		reqUrl := EnsurePort(w.Request.URL.String())

		if strings.HasPrefix(reqUrl, prefixAddr) {

			header, ok := w.Header["Content-Encoding"]
			if ok {
				body = Decompress(header[0], body)
			}

			modified := bytes.ReplaceAll(body, []byte(s.Sdata), []byte(s.Tdata))
			if ok {
				modified = Compress(header[0], modified)
			}
			return modified
		}
	}
	return body
}

func Decompress(encoding string, body []byte) []byte {
	dst := body

	switch encoding {
	case "gzip":
		r, _ := gzip.NewReader(bytes.NewReader(body))
		defer r.Close()
		dst, _ = io.ReadAll(r)
	case "zstd":
		dst, _ = gozstd.Decompress(nil, body)
	}
	return dst
}

func Compress(encoding string, body []byte) []byte {
	dst := body

	switch encoding {
	case "gzip":
		var b bytes.Buffer
		w := gzip.NewWriter(&b)
		_, _ = w.Write(body)
		defer w.Close()
		dst = b.Bytes()
	case "zstd":
		dst = gozstd.Compress(nil, body)
	}
	return dst
}
