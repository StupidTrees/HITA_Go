package security

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql/driver"
	"encoding/base64"
	"errors"
)

type MPrivateKey rsa.PrivateKey

// 写入数据库之前，对数据做类型转换
func (s MPrivateKey) Value() (driver.Value, error) {
	derPkix := x509.MarshalPKCS1PrivateKey((*rsa.PrivateKey)(&s))
	return base64.StdEncoding.EncodeToString(derPkix), nil
}

// 将数据库中取出的数据，赋值给目标类型
func (s *MPrivateKey) Scan(v interface{}) (err error) {
	switch v.(type) {
	case string:
		str := v.(string)
		var bytes []byte
		var pk interface{}
		bytes, err = base64.StdEncoding.DecodeString(str)
		if err == nil {
			pk, err = x509.ParsePKCS1PrivateKey(bytes)
			*s = pk.(MPrivateKey)
		}
	default:
		err = errors.New("类型处理错误")
	}
	return err

}
