package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var sampleSecret []byte

func init() {
	assertAvailablePRNG()
	bytez, err := generateRandomBytes(256)
	if err != nil {
		println(err)
	}
	sampleSecret = bytez
	//privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	//if err != nil {
	//	fmt.Printf("Cannot generate RSA key\n")
	//	os.Exit(1)
	//}
	//hmacSampleSecret = &privatekey.PublicKey
}

func assertAvailablePRNG() {
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

const iss = "Zoom Clone"
const expirationTime int64 = 60 * 2

const UserTokenKey = "userToken"
const UsernameKey = "username"

type JwtLoginClaims struct {
	*jwt.StandardClaims
	Username  string
	UserToken string
}

func CreateToken(username string, userToken string) (string, error) {
	claims := &JwtLoginClaims{
		&jwt.StandardClaims{
			Issuer:    iss,
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		username,
		userToken,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(sampleSecret)
}

func GetTokenData(tokenStr string) (*JwtLoginClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtLoginClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return sampleSecret, nil
	})

	if claims, ok := token.Claims.(*JwtLoginClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func Authenticate(authHeader string) (*JwtLoginClaims, bool) {
	if authHeader != "" {
		authHeaderSplit := strings.Split(authHeader, " ")
		if len(authHeaderSplit) == 2 {
			data, err := GetTokenData(authHeaderSplit[1])
			if err == nil {
				return data, true
			}
		}
	}
	return nil, false
}
