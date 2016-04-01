package requestclient

import (
	"fmt"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"time"
)

type ProxyType int

const (
	NO_PROXY     ProxyType = iota
	ENV_PROXY              = iota
	MANUAL_PROXY           = iota
	SOCKS5_PROXY           = iota
)

type ClientProxy struct {
	ProxyType ProxyType
	URL       string
}

func GetClient(clientProxy *ClientProxy) (client *http.Client, err error) {
	var transport *http.Transport

	if clientProxy == nil || clientProxy.ProxyType == NO_PROXY {
		//TODO
	} else if clientProxy.ProxyType == ENV_PROXY {
		//TODO
	} else if clientProxy.ProxyType == MANUAL_PROXY {
		fmt.Println("manual proxy")
		url, err := url.Parse(clientProxy.URL)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(url),
		}
	} else if clientProxy.ProxyType == SOCKS5_PROXY {
		dialer, err := proxy.SOCKS5("tcp", clientProxy.URL, nil,
			&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			},
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		transport = &http.Transport{
			Proxy: nil,
			Dial:  dialer.Dial,
		}
	}

	client = &http.Client{Transport: transport}

	return client, err
}
