package main

const (
	NO_PROXY     = iota
	ENV_PROXY    = iota
	MANUAL_PROXY = iota
	SOCKS5_PROXY = iota
)

const USE_PROXY = SOCKS5_PROXY
const HTTP_PROXY_URL = ""
const SOCKS5_PROXY_URL = "127.0.0.1:1080"
