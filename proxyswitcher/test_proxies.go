package proxyswitcher

import (
	"strings"
	"time"

	"github.com/go-per/simpkg/client"
	"github.com/go-per/simpkg/format"
	"github.com/go-per/simpkg/parse"
)

// testResponse struct
type testResponse struct {
	Ip string `json:"ip"`
}

// TestAll tests all proxies in the proxy list
func (ps *Switcher) TestAll(callback func(result string, err error)) {
	if ps.proxies == nil || len(ps.proxies) == 0 {
		callback("", format.Error("Proxy list is empty"))
		return
	}

	// get machine ip
	callback(getOutgoingIp(1, nil))

	// test all proxies
	for i, proxy := range ps.proxies {
		callback(getOutgoingIp(i+2, proxy))
	}
}

// getOutgoingIp returns outgoing ip address
func getOutgoingIp(index int, proxy IProxy) (string, error) {
	var response *testResponse
	c := client.New().SetTimeout(time.Second * 10)
	if proxy != nil {
		c.SetProxyURL(proxy.RawUrl())
	}

	resp, err := c.R().Get("https://antcpt.com/score_detector/getMyIp.php")
	if err != nil {
		return "", err
	}

	if resp == nil || resp.GetStatusCode() != 200 {
		return "", format.Error("Status code is not valid or response is empty")
	}

	body, err := resp.ToBytes()
	if err != nil {
		return "", err
	}
	if err := parse.Decode(body, &response); err != nil {
		return "", err
	}
	if response == nil {
		return "", format.Error("Response is empty")
	}

	if proxy == nil {
		return format.Format("#%v %v\t\t~>\t%v\t[Local]", index, "Local Ip", response.Ip), nil
	}

	if !strings.Contains(proxy.RawUrl(), response.Ip) {
		url := proxy.Url()
		return format.Format("#%v %v\t~>\t%v\t[Not Match]", index, url.Hostname(), response.Ip), nil
	}

	return format.Format("#%v %v \t~>\t%v\t[OK]", index, proxy.Url().Host, response.Ip), nil
}
