package logging

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const serviceName = "lab8-consumer"

func NewJSONLogger(output io.Writer, module string) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(output),
		zap.InfoLevel,
	))

	return logger.With(
		zap.String("service", serviceName),
		zap.String("module", module),
		zap.Int("pid", os.Getpid()),
	)
}
