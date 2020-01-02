package untils

import (
	"crypto/md5"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
)

type RespData struct {
	Code int8 `json:"code"`
	Msg string `json:"msg"`
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

func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func GenUserPassword(password string) (md5Pwd string){
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
		Log.Error("gen jwt token exeception", zap.String("msg", err.Error()))
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

func GetCurrentMonth() string {
	d := time.Now()
	d = d.AddDate(0, 0, 0)
	return d.Format("2006-01")
}

func GetNextMonth() string {
	d := time.Now()
	d = d.AddDate(0, 1, 0)
	return d.Format("2006-01")
}

func GetUUID() string {
	return uuid.NewV4().String()
}
