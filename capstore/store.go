package capstore

import (
	"errors"
	"sync"
)

// Instance store instance
var Instance *CaptchaStore
var storeNameLocker = sync.RWMutex{}
var actionNameLocker = sync.RWMutex{}

// IStore store interface
type IStore interface {
	Pool() IPool // Pool returns pool instance
}

// CaptchaStore captcha store
type CaptchaStore struct {
	stores      map[string]IStore
	activeStore string
	actionName  []string
}

// initialize
func init() {
	Instance = New()
}

// New creates new instance
func New() *CaptchaStore {
	instance := &CaptchaStore{
		stores: make(map[string]IStore, 0),
	}

	// add default store
	instance.AddStore(DefaultStoreKey, newDefaultStore())

	return instance
}

// AddStore add new store
func (store *CaptchaStore) AddStore(name string, s IStore) {
	store.stores[name] = s
}

// Use set active store name
func (store *CaptchaStore) Use(storeName string) {
	storeNameLocker.Lock()
	store.activeStore = storeName
	storeNameLocker.Unlock()
}

// WithAction set action name
func (store *CaptchaStore) WithAction(action ...string) *CaptchaStore {
	actionNameLocker.Lock()
	store.actionName = action
	actionNameLocker.Unlock()
	return store
}

// GetToken returns token
func (store *CaptchaStore) GetToken() (token *Token, err error) {
	current := store.Current()
	if current == nil {
		err = errors.New("no active store")
		return
	}

	token, err = current.Pool().Get(store.actionName...)
	store.resetAction()
	return
}

// resetAction reset action name
func (store *CaptchaStore) resetAction() {
	store.actionName = []string{}
}

// GetActiveName returns active store name
func (store *CaptchaStore) GetActiveName() string {
	return store.activeStore
}

// Current returns active store
func (store *CaptchaStore) Current() IStore {
	return store.stores[store.activeStore]
}

// Pool returns Current store pool
func (store *CaptchaStore) Pool() IPool {
	current := store.Current()
	if current == nil {
		return nil
	}

	return current.Pool()
}
