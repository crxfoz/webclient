package webclient

import (
	"net/http/cookiejar"
	"net/http"
	"crypto/tls"
	"time"
	"net"
)

type Config struct {
	Timeout time.Duration
	UseKeepAlive bool
	FollowRedirect bool
}

func (c Config) New() *Webclient {
	options := cookiejar.Options{}
	jar, _ := cookiejar.New(&options)

	newWebClient := &Webclient{
		client:    &http.Client{Jar: jar},
		transport: &http.Transport{
			Proxy:              nil,
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: false,
			DisableKeepAlives:  !c.UseKeepAlive,
		},
	}

	if c.Timeout > 0 {
		newWebClient.client.Timeout = c.Timeout * time.Second
		newWebClient.transport.TLSHandshakeTimeout = c.Timeout * time.Second
		newWebClient.transport.ResponseHeaderTimeout = c.Timeout * time.Second
		newWebClient.transport.DialContext = (&net.Dialer{
			Timeout:   c.Timeout * time.Second,
			KeepAlive: 120 * time.Second,
		}).DialContext
	}

	if !c.FollowRedirect {
		newWebClient.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		newWebClient.client.CheckRedirect = nil
	}

	return newWebClient

}