package tool

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

func init() {
	file, _ := os.OpenFile("sportseckill.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	log.SetOutput(file)
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

}

var traceId = uuid.New().String()

func Info(args ...interface{}) {
	log.Info(traceId, args)
}

// Timestamp 获取timestamp格式时间
func Timestamp(time time.Time) int64 {
	return time.UnixNano() / 1e6
}

// UTCTime 获取UTC格式时间，2021-05-15T18:02:48.284Z
func UTCTime() string {
	utcTime := time.Now().UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}

func FormatTime(d time.Duration) string {
	utcTime := time.Now().Add(d).UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10])
	return iso
}

func FormatTimeByTime(date time.Time, d time.Duration) string {
	utcTime := date.Add(d).UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10])
	return iso
}

func StringToTime(str string) time.Time {
	local, _ := time.LoadLocation("Local")
	time, _ := time.ParseInLocation("2006-01-02 15:04:05", str, local)
	return time
}

// ParseRequestParams 转化请求参数为json
func ParseRequestParams(params interface{}) (string, *bytes.Reader, error) {
	if params == nil {
		return "", nil, errors.New("illegal parameter")
	}
	data, err := json.Marshal(params)
	if err != nil {
		return "", nil, errors.New("json convert string error")
	}
	jsonBody := string(data)
	byteBody := bytes.NewReader(data)
	return jsonBody, byteBody, nil
}

// GetSign 生成okex的签名
func GetSign(timestamp string, method string, requestUri string, body string, secretKey string) (string, error) {
	content := timestamp + strings.ToUpper(method) + requestUri + body
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(content))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}
