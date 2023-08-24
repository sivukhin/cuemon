package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	authJsonSuffix   = "_AUTH_JSON"
	cookieEnv        = "COOKIE"
	authorizationEnv = "AUTHORIZATION"
	privateKeyIdEnv  = "PRIVATE_KEY_ID"
	clientEmailEnv   = "CLIENT_EMAIL"
	privateKeyEnv    = "PRIVATE_KEY"
	iapClientIdEnv   = "IAP_CLIENT_ID"
	oauthTokenUriEnv = "OAUTH_TOKEN_URI"
)

type AuthorizationMethods struct {
	Cookie                   *string
	AuthorizationHeader      *string
	ProxyAuthorizationHeader *string
}

func ParseEnvVars(envs []string) (map[string]string, error) {
	vars := make(map[string]string)
	for _, env := range envs {
		tokens := strings.SplitN(env, "=", 2)
		var value string
		if len(tokens) > 1 {
			if suffix, ok := strings.CutPrefix(tokens[1], "file://"); ok {
				content, err := os.ReadFile(suffix)
				if err != nil {
					return nil, fmt.Errorf("unable to read file %v for key %v", suffix, tokens[0])
				}
				value = string(content)
			} else {
				value = tokens[1]
			}
		}
		if prefix, ok := strings.CutSuffix(tokens[0], authJsonSuffix); ok {
			var content map[string]string
			err := json.Unmarshal([]byte(value), &content)
			if err != nil {
				return nil, fmt.Errorf("unable to deserialize JSON from env key %v", env)
			}
			for jsonKey, jsonValue := range content {
				vars[prefix+"_"+strings.ToUpper(jsonKey)] = string(jsonValue)
			}
		} else {
			vars[tokens[0]] = value
		}
	}
	return vars, nil
}

func AnalyzeSubjectAuthorization(envs []string) (map[string]AuthorizationMethods, error) {
	envVars, err := ParseEnvVars(envs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse env vars: %w", err)
	}
	subjects := make(map[string]AuthorizationMethods)
	for keyLoop, valueLoop := range envVars {
		key, value := keyLoop, valueLoop
		keyTokens := strings.SplitN(key, "_", 2)
		if len(keyTokens) != 2 {
			continue
		}
		keySubject, keyType := keyTokens[0], keyTokens[1]
		methods := subjects[keySubject]
		if keyType == cookieEnv {
			methods.Cookie = &value
		} else if keyType == authorizationEnv {
			methods.AuthorizationHeader = &value
		} else if keyType == iapClientIdEnv {
			oauthTokenUri, ok := envVars[keySubject+"_"+oauthTokenUriEnv]
			if !ok {
				return nil, fmt.Errorf("%v_%v var expected for IAP authorization", keySubject, oauthTokenUriEnv)
			}
			privateKeyId, ok := envVars[keySubject+"_"+privateKeyIdEnv]
			if !ok {
				return nil, fmt.Errorf("%v_%v var expected for IAP authorization", keySubject, privateKeyIdEnv)
			}
			clientEmail, ok := envVars[keySubject+"_"+clientEmailEnv]
			if !ok {
				return nil, fmt.Errorf("%v_%v var expected for IAP authorization", keySubject, clientEmailEnv)
			}
			privateKey, ok := envVars[keySubject+"_"+privateKeyEnv]
			if !ok {
				return nil, fmt.Errorf("%v_%v var expected for IAP authorization", keySubject, privateKeyEnv)
			}
			header, err := iapAuthorizationToken(oauthTokenUri, value, clientEmail, privateKeyId, privateKey)
			if err != nil {
				return nil, fmt.Errorf("unable to get IAP authorization header for subject %v: %w", keySubject, err)
			}
			bearer := "Bearer " + header
			methods.ProxyAuthorizationHeader = &bearer
		}
		subjects[keySubject] = methods
	}
	return subjects, nil
}

func iapAuthorizationToken(oauthTokenUri string, clientId, clientEmail, privateKeyId, privateKey string) (string, error) {
	type (
		iapHeader struct {
			Alg string `json:"alg"`
			Typ string `json:"typ"`
			Kid string `json:"kid"`
		}
		iapPayload struct {
			Iss            string `json:"iss"`
			Aud            string `json:"aud"`
			Sub            string `json:"sub"`
			TargetAudience string `json:"target_audience"`
			Exp            int64  `json:"exp"`
			Iat            int64  `json:"iat"`
		}
	)
	header := iapHeader{
		Alg: "RS256",
		Typ: "JWT",
		Kid: privateKeyId,
	}
	payload := iapPayload{
		Iss:            clientEmail,
		Aud:            oauthTokenUri,
		Sub:            clientEmail,
		TargetAudience: clientId,
		Iat:            time.Now().Unix(),
		Exp:            time.Now().Add(10 * time.Minute).Unix(),
	}
	headerBase64, err := jsonBase64(header)
	if err != nil {
		return "", fmt.Errorf("unable to serialize header: %w", err)
	}
	payloadBase64, err := jsonBase64(payload)
	if err != nil {
		return "", fmt.Errorf("unable to serialize payload: %w", err)
	}
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", fmt.Errorf("unable to decode private key")
	}
	rsaPrivateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("unable to parse RSA private key: %w", err)
	}
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(fmt.Sprintf("%v.%v", headerBase64, payloadBase64)))
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey.(*rsa.PrivateKey), crypto.SHA256, sha256Hash.Sum(nil))
	if err != nil {
		return "", fmt.Errorf("unable to sign request: %w", err)
	}
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)
	query := url.Values{}
	query.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	query.Add("assertion", fmt.Sprintf("%v.%v.%v", headerBase64, payloadBase64, signatureBase64))
	encodedData := query.Encode()
	request, err := http.NewRequest("POST", oauthTokenUri, strings.NewReader(encodedData))
	if err != nil {
		return "", fmt.Errorf("unable to create request: %w", err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("unable to get response: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode/100 != 2 {
		message, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("unexpected http status code: %v (%v)", response.StatusCode, string(message))
	}
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body: %w", err)
	}
	type oauthResponse struct {
		Token string `json:"id_token"`
	}
	var token oauthResponse
	err = json.Unmarshal(responseBytes, &token)
	if err != nil {
		return "", fmt.Errorf("unable to deserialize response body: %w", err)
	}
	if token.Token == "" {
		return "", fmt.Errorf("get empty token")
	}
	return token.Token, nil
}

func jsonBase64(v any) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("unable to serialize value: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
