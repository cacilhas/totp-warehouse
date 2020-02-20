package totp

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hgfischer/go-otp"
)

type OTP interface {
	Issuer() string
	User() string
	Length() uint8
	Secret() string
	Token() string
	OriginalURI() string
}

type pogo struct {
	issuer string
	user   string
	length uint8
	secret string
	uri    *url.URL
}

func Import(code string) (OTP, error) {
	var uri *url.URL
	var err error
	if uri, err = url.Parse(code); err != nil {
		return nil, err
	}
	if uri.Scheme != "otpauth" {
		return nil, fmt.Errorf("unsupported scheme %v", uri.Scheme)
	}

	issuer, user := getUser(uri.Path[1:])
	query := uri.Query()

	maybeIssuer := query["issuer"]
	if len(maybeIssuer) > 0 {
		issuer = maybeIssuer[0]
	}

	secret := query["secret"]
	if len(secret) == 0 {
		return nil, fmt.Errorf("secret not supplied")
	}

	maybeLength := query["length"]
	if len(maybeLength) == 0 {
		maybeLength = []string{"6"}
	}
	var length int
	if length, err = strconv.Atoi(maybeLength[0]); err != nil {
		length = 6
	}

	return &pogo{
		issuer: issuer,
		user:   user,
		length: uint8(length),
		secret: secret[0],
		uri:    uri,
	}, nil
}

func (pogo pogo) Issuer() string {
	return pogo.issuer
}

func (pogo pogo) User() string {
	return pogo.user
}

func (pogo pogo) Length() uint8 {
	return pogo.length
}

func (pogo pogo) Secret() string {
	return pogo.secret
}

func (pogo pogo) OriginalURI() string {
	return pogo.uri.String()
}

func (pogo pogo) Token() string {
	totp := otp.TOTP{
		Secret:         pogo.secret,
		Length:         pogo.length,
		IsBase32Secret: true,
	}
	return totp.Now().Get()
}

func getUser(path string) (string, string) {
	if strings.Contains(path, ":") {
		values := strings.SplitN(path, ":", 2)
		return values[0], values[1]
	}
	return "", path
}
