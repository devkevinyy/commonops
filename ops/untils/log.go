package untils

import (
	"github.com/labstack/gommon/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var Log *zap.Logger

const baseLogPath = "./logs/"

func init() {
	writer, err := rotatelogs.New(baseLogPath+"ops_log.%Y%m%d",
		rotatelogs.WithMaxAge(time.Hour*24*7),   // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24), 	// 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	w := zapcore.AddSync(writer)
	var level  = zap.InfoLevel

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		level)
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
	Log.Info("DefaultLogger init success")
}