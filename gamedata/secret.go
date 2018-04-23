package gamedata

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"crypto/md5"
	"encoding/hex"
)

func GenRsaKey() ([]byte, []byte, error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	var priBuf bytes.Buffer
	err = pem.Encode(&priBuf, block)
	if err != nil {
		return nil, nil, err
	}

	// 生成公钥
	publicKey := &privateKey.PublicKey
	derPkix, _ := x509.MarshalPKIXPublicKey(publicKey)

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	var pubBuf bytes.Buffer
	err = pem.Encode(&pubBuf, block)
	if err != nil {
		return nil, nil, err
	}
	return priBuf.Bytes(), pubBuf.Bytes(), nil
}

func RsaEncrept(originData, publicKey []byte) (string, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)
	data, err := rsa.EncryptPKCS1v15(rand.Reader, pub, originData)
	return base64.StdEncoding.EncodeToString(data), err
}

func RsaDecrypt(ciphertext, privateKey []byte) (string, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	cipherStr := h.Sum(nil)
	cipher := hex.EncodeToString(cipherStr)
	return cipher
}
