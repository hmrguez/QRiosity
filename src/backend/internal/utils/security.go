package utils

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/big"
	"net/http"
)

var (
	jwksURL       = "https://cognito-idp.us-east-2.amazonaws.com/us-east-2_XSBPyWz8o/.well-known/jwks.json"
	rsaPublicKeys map[string]*rsa.PublicKey
)

// Define the JWKS response structure
type Jwks struct {
	Keys []JwkKey `json:"keys"`
}

type JwkKey struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// Function to load the JWKS and extract the RSA keys
func loadRSAPublicKeys(jwksURL string) error {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var jwks Jwks
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	rsaPublicKeys = make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		// Decode the modulus (n) and exponent (e)
		nDecoded, err := base64.RawURLEncoding.DecodeString(key.N)
		if err != nil {
			return err
		}

		eDecoded, err := base64.RawURLEncoding.DecodeString(key.E)
		if err != nil {
			return err
		}

		// Convert exponent from bytes to int
		var eInt int
		if len(eDecoded) == 3 {
			eInt = int(binary.BigEndian.Uint32(append([]byte{0x00}, eDecoded...)))
		} else if len(eDecoded) == 2 {
			eInt = int(binary.BigEndian.Uint16(eDecoded))
		} else {
			eInt = int(eDecoded[0])
		}

		// Create the RSA public key
		rsaPublicKeys[key.Kid] = &rsa.PublicKey{
			N: new(big.Int).SetBytes(nDecoded),
			E: eInt,
		}
	}

	return nil
}

// Function to retrieve an RSA public key by KID
func getRSAPublicKey(kid string) (*rsa.PublicKey, error) {
	if rsaPublicKeys == nil {
		return nil, errors.New("RSA keys not loaded")
	}
	key, found := rsaPublicKeys[kid]
	if !found {
		return nil, errors.New("key not found for the given kid")
	}
	return key, nil
}

// Step 2: Validate JWT token using Cognito's public keys
func validateToken(tokenString string) (string, error) {

	if rsaPublicKeys == nil {
		err := loadRSAPublicKeys(jwksURL)
		if err != nil {
			return "", err
		}
	}

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
