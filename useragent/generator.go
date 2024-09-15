package useragent

import "github.com/go-per/simpkg/random"

// Random returns random user agent
func Random() string {
	return userAgents[random.IntInRange(0, len(userAgents)-1)]
}
