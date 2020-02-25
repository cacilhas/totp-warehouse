package totp

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hgfischer/go-otp"
)

// OTP interface for OTP pogo
type OTP interface {
	// Issuer returns the service issuer
	Issuer() string
	// User returns the service user
	User() string
	// Length returns the secret length in digits
	Length() uint8
	// Secret the secret's B32
	Secret() string
	// Token returns the current TOTP token
	Token() string
	// Key returns the storage key
	Key() string
	// String the URI used to supply the OTP data
	String() string
}

type pogo struct {
	issuer string
	user   string
	length uint8
	secret string
	uri    *url.URL
}

// Import imports OTP data from URI
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

func (pogo pogo) Token() string {
	totp := otp.TOTP{
		Secret:         strings.ToUpper(pogo.secret),
		Length:         pogo.length,
		IsBase32Secret: true,
	}
	return totp.Now().Get()
}

func (pogo pogo) Key() string {
	return fmt.Sprintf(
		"%v@%v",
		pogo.user,
		strings.ReplaceAll(pogo.issuer, " ", "+"),
	)
}

func (pogo pogo) String() string {
	return pogo.uri.String()
}

func (pogo pogo) Bytes() []byte {
	return []byte(pogo.String())
}

func getUser(path string) (string, string) {
	if strings.Contains(path, ":") {
		values := strings.SplitN(path, ":", 2)
		return values[0], values[1]
	}
	return "", path
}
