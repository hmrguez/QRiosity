package utils

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
)

var (
	jwksURL       = "https://cognito-idp.us-east-2.amazonaws.com/us-east-2_XSBPyWz8o/.well-known/jwks.json"
	rsaPublicKeys map[string]*rsa.PublicKey
)

// Step 1: Load JWKS keys from Cognito
func loadRSAPublicKeys() error {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var jwks struct {
		Keys []struct {
			Kid string `json:"kid"`
			Kty string `json:"kty"`
			Alg string `json:"alg"`
			Use string `json:"use"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}

	if err := json.Unmarshal(body, &jwks); err != nil {
		return err
	}

	rsaPublicKeys = make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		if key.Kty == "RSA" {
			rsaKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key.N))
			if err != nil {
				return err
			}
			rsaPublicKeys[key.Kid] = rsaKey
		}
	}

	return nil
}

// Step 2: Validate JWT token using Cognito's public keys
func validateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Extract the 'kid' from the header
		kid := token.Header["kid"].(string)

		// Get the public key for the 'kid'
		rsaKey, exists := rsaPublicKeys[kid]
		if !exists {
			return nil, fmt.Errorf("unknown kid %s", kid)
		}

		return rsaKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check additional claims like 'aud' (audience), 'iss' (issuer)
		if claims["iss"] != "https://cognito-idp.us-east-2.amazonaws.com/us-east-2_XSBPyWz8o" {
			return "", errors.New("invalid token issuer")
		}

		if claims["aud"] != "5cpkviq3bb7e1dvdhcuqidju1g" {
			return "", errors.New("invalid token audience")
		}

		return "Authorized", nil
	}

	return "", errors.New("invalid token")
}

func CheckAuthorization(ctx context.Context, event AppSyncEvent) error {
	// Extract the token from headers
	tokenString, exists := event.Headers["authorization"]
	if !exists {
		return errors.New("authorization header missing")
	}

	// Validate the token using the existing validateToken function
	_, err := validateToken(tokenString)
	return err
}

func CheckAuthorizationTokenOnly(ctx context.Context, token string) error {
	// Validate the token using the existing validateToken function
	_, err := validateToken(token)
	return err
}

func SecureResolver(ctx context.Context, event AppSyncEvent, resolver func(context.Context, json.RawMessage) (json.RawMessage, error)) (json.RawMessage, error) {
	err := CheckAuthorization(ctx, event)
	if err != nil {
		return nil, err
	}

	// Call the next function
	response, err := resolver(ctx, event.Arguments)
	return response, err
}
