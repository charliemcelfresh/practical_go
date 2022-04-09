package jwt_maker

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	log "github.com/sirupsen/logrus"
)

func Run(jwtDuration string, iss string, aud string, adminOrUserID int64) {
	dur, err := time.ParseDuration(jwtDuration)
	if err != nil {
		log.Fatalln(err)
	}
	prvKey, err := ioutil.ReadFile(os.Getenv("JWT_PRIVATE_KEY"))
	if err != nil {
		log.Fatalln(err)
	}
	pubKey, err := ioutil.ReadFile(os.Getenv("JWT_PUBLIC_KEY"))
	if err != nil {
		log.Fatalln(err)
	}

	jwtToken := NewJWT(prvKey, pubKey)

	tok, err := jwtToken.Create(dur, iss, aud, adminOrUserID)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("TOKEN:", tok)

	_, err = jwtToken.Validate(tok)
	if err != nil {
		log.Fatalln(err)
	}
}

type JWT struct {
	privateKey []byte
	publicKey  []byte
}

func NewJWT(privateKey []byte, publicKey []byte) JWT {
	return JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (j JWT) Create(ttl time.Duration, iss string, aud string, adminOrUserID int64) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["exp"] = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()          // The time at which the token was issued.
	claims["iss"] = iss                 // issuing service
	claims["aud"] = aud                 // audience, eg user, admin, service-to-service
	if aud == "user" {
		claims["user_id"] = adminOrUserID
	} else if aud == "admin" {
		claims["admin_id"] = adminOrUserID
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}
	return token, nil
}

func (j JWT) Validate(token string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	tok, err := jwt.Parse(
		token, func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
			}

			return key, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims, nil
}
