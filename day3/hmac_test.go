package main

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

/*
	type Token struct {
		Raw       string // 原始Token字符串，当开始解析时填充此字段
		Method    SigningMethod  // 签名使用的方法
		Header    map[string]interface{}// JWT的header部分
		Claims    Claims   // JWT的payload部分
		Signature string   // JWT的签名部分，当开始解析时填充此字段
		Valid     bool   // JWT是否合法有效
	}
*/

func TestHmac(t *testing.T) {
	secret := []byte("niganma")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   123456,
		"name": "goeer",
	})
	fmt.Printf("%+v\n", *token)
	signedString, err := token.SignedString(secret)
	fmt.Println(signedString, err)
}
func TestPreClaims(t *testing.T) {
	mySigningKey := []byte("niganma")
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Unix(
			time.Now().Add(7*24*time.Hour).Unix(), 0)),
		Issuer: "test",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Printf("%v %v\n", ss, err)
}

type MyClaims struct {
	User string `json:"user"`
	jwt.RegisteredClaims
}

func TestCustomClaims(t *testing.T) {
	secret := []byte("niganma")

	claims := MyClaims{
		User: "114514",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "goeer",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(secret)
	fmt.Println(signedString, err)
}

func TestParse(t *testing.T) {
	secret := []byte("niganma")
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTIzNDU2LCJuYW1lIjoiZ29lZXIifQ.90w0iizcS47-3ENmbD6VVo1s3IaJ7GgTG9f-fK672dE"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不匹配的签名算法 [%s]\n", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		fmt.Println(token, err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	} else {
		fmt.Println(err)
	}
}
func TestProcess(t *testing.T) {
	secret := []byte("niganma")
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiMTE0NTE0IiwiaXNzIjoiZ29lZXIiLCJleHAiOjE3MTQyMDk5NDIsIm5iZiI6MTcxNDIwNjM0MiwiaWF0IjoxNzE0MjA2MzQyfQ.SeL7W2Mglj7lHKZ9n3SvpNl_nS_g1rbAyJtw5RVSHSE"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不匹配的签名算法 [%s]\n", token.Header["alg"])
		}
		return secret, nil
	})
	if token.Valid {
		fmt.Println("token合法")
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		fmt.Println("传入的字符串甚至连一个token都不是...")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		fmt.Println("token已经过期或者还没有生效")
	} else {
		fmt.Println("token处理异常...")
	}
}
func TestCustomClaimsParsee(t *testing.T) {
	secret := []byte("niganma")
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiMTE0NTE0IiwiaXNzIjoiZ29lZXIiLCJleHAiOjE3MTQyMTA2NDIsIm5iZiI6MTcxNDIwNzA0MiwiaWF0IjoxNzE0MjA3MDQyfQ.BtNbZ3KOZc4RjLCQ-63XRZebJXCHR1meTCgvmtrmww4"
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		fmt.Println(token, err)
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		fmt.Println(claims)
	} else {
		fmt.Println(err)
	}
}
func TestRsa(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	publicKey:=&privateKey.PublicKey
	if err!=nil{
		fmt.Println(err)
		return 
	}
	claims:=MyClaims{
		User:"114514",
		RegisteredClaims:jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer: "goeer",
		},
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodRS256,claims)
	signedString,err:=token.SignedString(privateKey)
	fmt.Println(signedString,err)
	token,err=jwt.ParseWithClaims(signedString,&MyClaims{},func(token *jwt.Token)(interface{},error){
		return publicKey,nil
	})
	if err!=nil{
		fmt.Println(err)
	}else if claims,ok:=token.Claims.(*MyClaims);ok&&token.Valid{
		fmt.Println(claims)
	}
}
