package untils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/chujieyang/commonops/ops/conf"
)

func ConvertUtcTimeToLocal(utcTime string, timeLayout string) (localTime string) {
	formate := "2006-01-02 15:04:05"
	parseTime, _ := time.Parse(timeLayout, utcTime)
	local, _ := time.LoadLocation("Local")
	return parseTime.In(local).Format(formate)
}

func GetNowTime() JSONTime {
	var cstZone = time.FixedZone("CST", 8*3600)
	return JSONTime{
		Time: time.Now().In(cstZone),
	}
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func DesEncode(data string) string {
	key := []byte(conf.DesKey)
	src := []byte(data)
	block, err := des.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}
	src = PKCS5Padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(src))
	blockMode.CryptBlocks(crypted, src)
	return base64.StdEncoding.EncodeToString(crypted)
}

func DesDecode(data string) string {
	key := []byte(conf.DesKey)
	crypted, _ := base64.StdEncoding.DecodeString(data)
	block, _ := des.NewCipher(key)
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return string(origData)
}
