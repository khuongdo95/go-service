package token

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"testing"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/stretchr/testify/assert"
)

const _name = "smm"

// TestGenerateKeys generates and saves ed25519 keys to disk after
// encoding into PEM format
func TestGenerateKeys(t *testing.T) {
	//t.Skip("Manual generate new ed25519 key pair")

	var (
		err   error
		b     []byte
		block *pem.Block
		pub   ed25519.PublicKey
		priv  ed25519.PrivateKey
	)

	pub, priv, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Printf("Generation error : %s", err)
		os.Exit(1)
	}

	b, err = x509.MarshalPKCS8PrivateKey(priv)
	assert.NoError(t, err)

	block = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	}

	err = os.WriteFile(_name, pem.EncodeToMemory(block), 0600)
	assert.NoError(t, err)

	// public key
	b, err = x509.MarshalPKIXPublicKey(pub)
	assert.NoError(t, err)

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}

	fileName := _name + ".pub"
	err = os.WriteFile(fileName, pem.EncodeToMemory(block), 0644)
	assert.NoError(t, err)

	// You can create jwk.Key from a raw key, too
	fromRawKey, err := jwk.New(priv)
	assert.NoError(t, err)

	// Keys can be serialized back to JSON
	jsonBuf, err := json.Marshal(fromRawKey)
	assert.NoError(t, err)

	fileNameJWK := _name + ".jwk"
	err = os.WriteFile(fileNameJWK, jsonBuf, 0644)
	assert.NoError(t, err)
}
