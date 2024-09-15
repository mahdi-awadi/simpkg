package identity

import (
	"errors"
	"sync"
	"time"

	"github.com/go-per/simpkg/format"
	"github.com/go-per/simpkg/random"
)

// Instance of the identity service
var Instance IIdentity = &Identity{
	uiKeysExpire: time.Minute * 10,
	validClients: make([]string, 0),
	clients:      make(map[string]*UiClient, 0),
	locker:       &sync.RWMutex{},
}

// IIdentity interface
type IIdentity interface {
	AddClients(clients []string)
	AddClient(client string)
	RemoveClient(client string)
	ClientsCount() int
	ClientExists(client string) bool
	GenerateUiToken(ownerId string) (string, error)
	RemoveUiToken(key string)
	IsUiTokenValid(key string) bool
	SetExpireTime(d time.Duration)
	Clients() (response map[string]*UiClient)
}

// UiClient struct
type UiClient struct {
	Owner     string    `json:"-"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	timer     *time.Timer
}

// Identity is a unique identifier for a user
type Identity struct {
	validClients []string
	clients      map[string]*UiClient
	uiKeysExpire time.Duration
	locker       *sync.RWMutex
}

// SetExpireTime set expire time
func (i *Identity) SetExpireTime(d time.Duration) {
	i.uiKeysExpire = d
}

// AddClients adds a list of clients to the identity service
func (i *Identity) AddClients(clients []string) {
	for _, user := range clients {
		i.AddClient(user)
	}
}

// AddClient adds a user to the identity service
func (i *Identity) AddClient(client string) {
	i.validClients = append(i.validClients, client)
}

// RemoveClient removes a user from the identity service
func (i *Identity) RemoveClient(client string) {
	for index, u := range i.validClients {
		if u == client {
			i.validClients = append(i.validClients[:index], i.validClients[index+1:]...)
			break
		}
	}
}

// ClientExists checks if a user exists in the identity service
func (i *Identity) ClientExists(client string) bool {
	i.locker.RLock()
	exsists := false
	for _, u := range i.validClients {
		if u == client {
			exsists = true
			break
		}
	}
	i.locker.RUnlock()
	return exsists
}

// ClientsCount returns the number of validClients in the identity service
func (i *Identity) ClientsCount() int {
	return len(i.validClients)
}

// GenerateUiToken adds a UI key to the identity service
func (i *Identity) GenerateUiToken(ownerId string) (string, error) {
	if !i.ClientExists(ownerId) {
		return "", errors.New("client does not exist")
	}

	i.locker.Lock()
	c, ok := i.clients[ownerId]
	i.locker.Unlock()
	if ok {
		c.CreatedAt = time.Now()
		c.timer.Reset(i.uiKeysExpire)
		return c.Key, nil
	}

	key := format.Format("v%v", random.String(5))
	i.locker.Lock()
	i.clients[ownerId] = &UiClient{
		Owner:     ownerId,
		Key:       key,
		CreatedAt: time.Now(),
		timer: time.AfterFunc(i.uiKeysExpire, func() {
			i.RemoveUiToken(key)
		}),
	}
	i.locker.Unlock()

	return key, nil
}

// RemoveUiToken removes a UI key from the identity service
func (i *Identity) RemoveUiToken(key string) {
	i.locker.Lock()
	for _, client := range i.clients {
		if client.Key == key {
			client.timer.Stop()
			delete(i.clients, client.Owner)
			break
		}
	}
	i.locker.Unlock()
}

// IsUiTokenValid checks if a UI key is valid
func (i *Identity) IsUiTokenValid(key string) bool {
	for _, k := range i.clients {
		if k.Key == key && k.CreatedAt.Add(i.uiKeysExpire).After(time.Now()) {
			return true
		}
	}
	return false
}

// Clients returns clients list
func (i *Identity) Clients() (response map[string]*UiClient) {
	return i.clients
}
