package totp

import (
	"fmt"
	"net/url"
	"strings"
)

type OTP interface {
	Issuer() string
	User() string
	Secret() string
	OriginalURI() string
}

type otp struct {
	issuer string
	user   string
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

	return &otp{
		issuer: issuer,
		user:   user,
		secret: secret[0],
		uri:    uri,
	}, nil
}

func (otp otp) Issuer() string {
	return otp.issuer
}

func (otp otp) User() string {
	return otp.user
}

func (otp otp) Secret() string {
	return otp.secret
}

func (otp otp) OriginalURI() string {
	return otp.uri.String()
}

func getUser(path string) (string, string) {
	if strings.Contains(path, ":") {
		values := strings.SplitN(path, ":", 2)
		return values[0], values[1]
	}
	return "", path
}
