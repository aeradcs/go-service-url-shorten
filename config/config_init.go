package config

import "flag"

var (
	Port    = flag.String("a", "8080", "port")
	BaseUrl = flag.String("b", "http://localhost:8080/", "base url used for shortened url")
)
