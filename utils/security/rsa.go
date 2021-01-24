package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/wenzhenxi/gorsa"
	"hita/config"
)

func GenerateRSAKeysStr() (publicKey string, privateKey string, err error) {
	prK, er := rsa.GenerateKey(rand.Reader, 1024)
	if er != nil {
		err = er
		return
	}
	prKB, er2 := x509.MarshalPKCS8PrivateKey(prK)
	if er2 != nil || prK == nil {
		err = er2
		return
	}
	puK := &prK.PublicKey
	puKB, er := x509.MarshalPKIXPublicKey(puK)
	if er != nil {
		err = er
		return
	}
	blockPu := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: puKB,
	}
	blockPr := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: prKB,
	}
	publicKey = string(pem.EncodeToMemory(blockPu))
	privateKey = string(pem.EncodeToMemory(blockPr))
	return
}

func EncryptWithServerKey(content string) string {
	res, _ := gorsa.PublicEncrypt(content, config.ServerPublicKey)
	return res
}

func DecryptWithServerKey(content string) string {
	res, _ := gorsa.PriKeyDecrypt(content, config.ServerPrivateKey)
	return res
}

func EncryptWithPublicKey(content string, publicKey string) string {
	str, _ := gorsa.PublicEncrypt(content, publicKey)
	return str
}

func EncryptWithPrivateKey(content string, privateKey string) string {
	str, _ := gorsa.PriKeyEncrypt(content, privateKey)
	return str
}
func DecryptWithPrivateKey(content string, privateKey string) string {
	str, _ := gorsa.PriKeyDecrypt(content, privateKey)
	return str
}

func DecryptWithPublicKey(content string, publicKey string) string {
	str, _ := gorsa.PublicDecrypt(content, publicKey)
	return str
}
