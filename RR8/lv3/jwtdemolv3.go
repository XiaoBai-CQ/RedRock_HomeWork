package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}
type Payload struct {
	Sub  string `json:"sub"`
	Iat  int64  `json:"iat"`
	Exp  int64  `json:"exp"`
	Name string `json:"name"`
}

func main() {
	secret := "by__kq"

	payload := Payload{
		Sub:  "00000",
		Iat:  time.Now().Unix(),
		Exp:  time.Now().Add(time.Hour).Unix(),
		Name: "kq",
	}

	token := generateJwt(payload, secret)
	fmt.Println(token)

	payloadnow, exists := verifyJwt(token, secret)
	if exists {
		fmt.Println(payloadnow)
	} else {
		fmt.Println("Not exists")
	}
}

func sign(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func generateJwt(payload Payload, secret string) string {
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}
	headerJson, _ := json.Marshal(header)
	encodeHeader := base64.RawURLEncoding.EncodeToString(headerJson)

	payloadJson, _ := json.Marshal(payload)
	encodePayload := base64.RawURLEncoding.EncodeToString(payloadJson)

	signature := sign(encodeHeader+"."+encodePayload, secret)

	return fmt.Sprintf("%s.%s.%s", encodeHeader, encodePayload, signature)
}

func verifyJwt(jwt string, secret string) (Payload, bool) {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return Payload{}, false
	}

	encodeheader, encodepayload, signature := parts[0], parts[1], parts[2]
	fmt.Println(encodeheader, encodepayload, signature)
	claims := encodeheader + "." + encodepayload
	fmt.Println(claims)
	rightSign := sign(claims, secret)
	fmt.Println(rightSign)

	if !hmac.Equal([]byte(signature), []byte(rightSign)) {
		return Payload{}, false
	}

	var payload Payload
	payloadjsonbytes, err1 := base64.RawURLEncoding.DecodeString(encodepayload)
	if err1 != nil {
		return Payload{}, false
	}
	err := json.Unmarshal(payloadjsonbytes, &payload)
	fmt.Println(payload)
	if err != nil {
		return Payload{}, false
	}

	currentTime := time.Now().Unix()
	if currentTime > payload.Exp {
		return Payload{}, false
	}

	return payload, true
}
