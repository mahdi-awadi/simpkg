package client

import (
	"github.com/go-per/simpkg/useragent"
	"github.com/imroc/req/v3"
)

// New creates new http client
func New() *req.Client {
	client := req.NewClient().
		SetUserAgent(useragent.Random()).
		DisableAutoDecode().
		EnableInsecureSkipVerify().
		EnableKeepAlives()

	return client
}
