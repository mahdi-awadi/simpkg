package identity

import (
	"testing"
	"time"
)

func TestIdentity_AddClient(t *testing.T) {
	type fields struct {
		validClients []string
		clients      map[string]*UiClient
		uiKeysExpire time.Duration
	}
	type args struct {
		client string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"TestIdentity_AddClient", fields{[]string{}, make(map[string]*UiClient, 0), 0}, args{"client"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Identity{
				validClients: tt.fields.validClients,
				clients:      tt.fields.clients,
				uiKeysExpire: tt.fields.uiKeysExpire,
			}
			i.AddClient(tt.args.client)
		})
	}
}

func TestIdentity_AddClients(t *testing.T) {
	type fields struct {
		validClients []string
		clients      map[string]*UiClient
		uiKeysExpire time.Duration
	}
	type args struct {
		clients []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"TestIdentity_AddClients", fields{[]string{}, make(map[string]*UiClient, 0), 0}, args{[]string{"client1", "client2"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Identity{
				validClients: tt.fields.validClients,
				clients:      tt.fields.clients,
				uiKeysExpire: tt.fields.uiKeysExpire,
			}
			i.AddClients(tt.args.clients)
		})
	}
}
