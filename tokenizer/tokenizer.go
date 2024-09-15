package tokenizer

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-per/simpkg/client"
	"github.com/go-per/simpkg/encryption"
	"github.com/go-per/simpkg/format"
	"github.com/go-per/simpkg/random"
	"github.com/imroc/req/v3"
)

// Instance is Tokenizer instance
var Instance *Tokenizer

// Tokenizer is the main struct of the tokenizer package
type Tokenizer struct {
	Encryptor encryption.IEncryptor
	client    *req.Client
}

// initialize the tokenizer
func init() {
	Instance = &Tokenizer{
		Encryptor: encryption.New(),
		client:    client.New().SetUserAgent("-"),
	}

	// set encryptor default key
	Instance.Encryptor.SetKey(encryptionKey)
}

// GetAccessToken get access token from server
func (tokenizer *Tokenizer) GetAccessToken(serverAddr, clientID, sourceName string) (string, error) {
	body := format.Format(`{"data":["%s", "%s", "%s"]}`, random.ID(), random.String(16), random.String(32))
	resp, err := tokenizer.client.R().SetBodyJsonString(body).Post(fmt.Sprintf("%s/access/%s/%s", serverAddr, sourceName, clientID))
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("could not load data from server")
	}
	var response = struct {
		Output struct {
			Data    string `json:"data"`
			Error   bool   `json:"error"`
			Message string `json:"message"`
		} `json:"output"`
	}{}

	if err == nil {
		body, err := resp.ToBytes()
		if err == nil {
			err = json.Unmarshal(body, &response)
		}
	}
	if err != nil {
		return "", errors.New("could not unmarshal response: " + err.Error())
	}
	if response.Output.Error {
		var message = "error occurred while getting access token"
		if response.Output.Message != "" {
			message = response.Output.Message
		}
		return "", errors.New(message)
	}

	// decrypt data
	token, err := tokenizer.Encryptor.Decrypt(response.Output.Data)
	if err != nil || token == "" || len(token) < 1 {
		return "", errors.New("server response is invalid")
	}

	return token, nil
}
