package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chujieyang/commonops/ops/opslog"

	"github.com/chujieyang/commonops/ops/conf"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RespData struct {
	Code int8        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type JSONTime struct {
	time.Time
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func GenUserPassword(password string) (md5Pwd string) {
	m5 := md5.New()
	m5.Write([]byte(password))
	m5.Write([]byte(conf.SecretSalt))
	bMd5 := m5.Sum(nil)
	md5Pwd = hex.EncodeToString(bMd5)
	return
}

func GenJWT(userInfo map[string]interface{}) (tokenString string) {
	jtoken := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["userInfo"] = userInfo
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(2)).Unix()
	claims["iat"] = time.Now().Unix()
	jtoken.Claims = claims
	tokenString, err := jtoken.SignedString([]byte(conf.SecretSalt))
	if err != nil {
		opslog.Error().Printf("gen jwt token exeception: %s \n", err)
	}
	return
}

func GetCurrentUserId(c *gin.Context) uint {
	userId, _ := c.Get("userId")
	return uint(userId.(float64))
}

func GetCurrentUsername(c *gin.Context) string {
	username, _ := c.Get("username")
	return username.(string)
}

func GetCurrentTime() string {
	d := time.Now()
	return d.Format("2006-01-02 15:04:05")
}

func GetUUID() string {
	return uuid.New().String()
}

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
	fmt.Println(conf.DesKey)
	key := []byte(conf.DesKey)
	src := []byte(data)
	block, err := des.NewCipher(key)
	if err != nil {
		opslog.Error().Println(err)
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

func ExtractUriPath(uri string) string {
	uriPathRegexp := regexp.MustCompile(`^(.*)\?|.*`)
	result := uriPathRegexp.FindStringSubmatch(uri)
	return strings.Replace(result[0], "?", "", -1)
}
