package security

import (
	_ "crypto"
	"crypto/rsa"
	"crypto/x509"
	"database/sql/driver"
	"encoding/base64"
	"errors"
)

type MPublicKey rsa.PublicKey

// 写入数据库之前，对数据做类型转换
func (s MPublicKey) Value() (driver.Value, error) {
	derPkix, err := x509.MarshalPKIXPublicKey(s)
	return base64.StdEncoding.EncodeToString(derPkix), err
}

// 将数据库中取出的数据，赋值给目标类型
func (s *MPublicKey) Scan(v interface{}) (err error) {
	switch v.(type) {
	case string:
		str := v.(string)
		var bytes []byte
		var pk interface{}
		bytes, err = base64.StdEncoding.DecodeString(str)
		if err == nil {
			pk, err = x509.ParsePKIXPublicKey(bytes)
			*s = pk.(MPublicKey)
		}
	default:
		err = errors.New("类型处理错误")
	}
	return err

}
