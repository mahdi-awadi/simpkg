package capstore

import (
	"fmt"
	"hash/crc32"
	"sync"
	"time"

	"github.com/go-per/simpkg/i18n"
)

// DefaultKey default key
const DefaultKey = "_default"

type interfaceMap map[string]any
type tokenMap map[string]map[string]Token

// Token struct
type Token struct {
	Value      string    `json:"value"`
	Data       any       `json:"data"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiryTime time.Time `json:"expiry_time"`
}

// IPool interface
type IPool interface {
	Tokens() tokenMap
	SetTokenLifeTime(t time.Duration)
	SetMinToken(min int, action ...string)
	SubscribeOnAdd(handler func())
	SubscribeOnRemove(handler func())
	Push(token string, data any, action ...string) *Token
	Get(action ...string) (*Token, error)
	Len() interfaceMap
}

// Pool instance.
type Pool struct {
	tokens           tokenMap
	lk               sync.RWMutex
	tokenLifeTime    time.Duration
	minTokens        map[string]int
	tokensChecksum   interfaceMap
	onAddHandlers    []func()
	onRemoveHandlers []func()
}

// NewPool create New pool instance.
func NewPool() IPool {
	m := &Pool{
		tokenLifeTime:    time.Second * 60,
		lk:               sync.RWMutex{},
		minTokens:        make(map[string]int),
		onAddHandlers:    make([]func(), 0),
		onRemoveHandlers: make([]func(), 0),
	}
	m.reset()

	return m
}

// Tokens returns all tokens
func (pool *Pool) Tokens() tokenMap {
	return pool.tokens
}

// SetTokenLifeTime set token lifetime
func (pool *Pool) SetTokenLifeTime(t time.Duration) {
	pool.tokenLifeTime = t
}

// SetMinToken set minimum required tokens count
func (pool *Pool) SetMinToken(min int, action ...string) {
	if len(action) == 0 || action[0] == "" {
		action = []string{DefaultKey}
	}
	pool.minTokens[action[0]] = min
}

// SubscribeOnAdd subscribe on add event
func (pool *Pool) SubscribeOnAdd(fn func()) {
	pool.onAddHandlers = append(pool.onAddHandlers, fn)
}

// SubscribeOnRemove subscribe on remove event
func (pool *Pool) SubscribeOnRemove(fn func()) {
	pool.onRemoveHandlers = append(pool.onRemoveHandlers, fn)
}

// Push append Token to list
func (pool *Pool) Push(token string, data any, action ...string) *Token {
	if token == "" {
		return nil
	}
	if len(action) == 0 || action[0] == "" {
		action = []string{DefaultKey}
	}

	pool.lk.Lock()
	defer pool.lk.Unlock()

	// create token checksum
	checksum := pool.makeChecksum(token)
	_, exists := pool.tokensChecksum[checksum]
	if exists {
		return nil
	}

	actionName := action[0]
	if _, ok := pool.tokens[actionName]; !ok {
		pool.tokens[actionName] = make(map[string]Token)
	}

	// insert Text
	t := Token{
		Value:      token,
		Data:       data,
		CreatedAt:  time.Now(),
		ExpiryTime: time.Now().Add(pool.tokenLifeTime),
	}

	pool.tokensChecksum[checksum] = ""
	pool.tokens[actionName][checksum] = t
	if pool.tokenLifeTime > 0 {
		go func(actionName, checksum string) {
			time.AfterFunc(pool.tokenLifeTime, func() {
				pool.remove(actionName, checksum)
			})

			if pool.onAddHandlers != nil {
				for _, fn := range pool.onAddHandlers {
					go fn()
				}
			}
		}(actionName, checksum)
	}

	return &t
}

// Get returns first Token item and remove it from list
func (pool *Pool) Get(action ...string) (*Token, error) {
	if len(action) == 0 || action[0] == "" {
		action = []string{DefaultKey}
	}
	tokens, ok := pool.tokens[action[0]]
	if !ok || len(tokens) == 0 {
		return nil, i18n.TranslateAsError("no_captcha_exists")
	}

	// get first item
	for _, token := range tokens {
		pool.remove(action[0], pool.makeChecksum(token.Value))
		return &token, nil
	}

	return nil, i18n.TranslateAsError("no_captcha_exists")
}

// Len returns tokens length
func (pool *Pool) Len() interfaceMap {
	tokens := interfaceMap{}
	for action, tokenMap := range pool.Tokens() {
		tokens[action] = len(tokenMap)
	}
	return tokens
}

// makeChecksum makes token checksum
func (pool *Pool) makeChecksum(token string) string {
	table := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum([]byte(token), table)
	return fmt.Sprintf("_%x", checksum)
}

// remove delete an item by checksum
func (pool *Pool) remove(action, checksum string) {
	pool.lk.Lock()
	defer pool.lk.Unlock()
	delete(pool.tokens[action], checksum)
	delete(pool.tokensChecksum, checksum)

	// on remove callback
	if pool.onRemoveHandlers != nil {
		for _, fn := range pool.onRemoveHandlers {
			go fn()
		}
	}
}

// reset clear the tokens.
func (pool *Pool) reset() *Pool {
	pool.tokens = make(tokenMap)
	pool.tokensChecksum = make(interfaceMap, 0)
	return pool
}
