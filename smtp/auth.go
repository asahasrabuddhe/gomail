package smtp

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
)

// loginAuth implements LOGIN Auth Mechanism for smtp.Auth interface
type loginAuth struct {
	username, password string
	host               string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if !server.TLS {
		advertised := false
		for _, mechanism := range server.Auth {
			if mechanism == "LOGIN" {
				advertised = true
				break
			}
		}
		if !advertised {
			return "", nil, errors.New("pigeon: unencrypted connection")
		}
	}
	if server.Name != a.host {
		return "", nil, errors.New("pigeon: wrong host name")
	}
	return "LOGIN", nil, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}

	switch {
	case bytes.Equal(fromServer, []byte("Username:")):
		return []byte(a.username), nil
	case bytes.Equal(fromServer, []byte("Password:")):
		return []byte(a.password), nil
	default:
		return nil, fmt.Errorf("pigeon: unexpected server challenge: %s", fromServer)
	}
}
