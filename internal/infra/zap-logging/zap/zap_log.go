package zap

import (
	"log"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitZap init zap logger
func InitZap(logConf LoggerConfig) *zap.Logger {
	cores := make([]zapcore.Core, 0)
	log.Println("Initializing Zap logger...")

	fileCores := createFileZapCore(logConf.name, logConf.path, logConf.maxAge, logConf.rotationTime, logConf.callerFullPath)
	if len(fileCores) > 0 {
		log.Println("File cores created successfully")
		cores = append(cores, fileCores...)
	} else {
		log.Println("No file cores were created")
	}

	if logConf.stdout {
		log.Println("Adding stdout logging")
		cores = append(cores, createStdCore(logConf.callerFullPath))
	}
	core := zapcore.NewTee(cores...)
	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(2)
	logger := zap.New(core, caller, callerSkip, zap.Development())

	zap.ReplaceGlobals(logger)
	if _, err := zap.RedirectStdLogAt(logger, zapcore.ErrorLevel); err != nil {
		log.Printf("Error redirecting std log: %v", err)
		panic(err)
	}
	log.Println("Zap logger initialized successfully")
	return logger
}

// createStdCore create stdout core
func createStdCore(callerFullPath bool) zapcore.Core {
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeTime = timeEncoder
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if callerFullPath {
		consoleEncoderConfig.EncodeCaller = customCallerEncoder
	}
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	return zapcore.NewCore(consoleEncoder, consoleDebugging, zapcore.DebugLevel)
}

// createFileZapCore info file => contain all log; error file => only contain error log
func createFileZapCore(name, logPath string, maxAge, rotationTime time.Duration, callerFullPath bool) (cores []zapcore.Core) {
	log.Printf("Creating file cores with logPath: %s", logPath)
	if len(logPath) == 0 {
		log.Println("No log path provided")
		return
	}
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		log.Printf("Log path does not exist, creating directory: %s", logPath)
		if err = os.MkdirAll(logPath, os.ModePerm); err != nil {
			log.Fatalf("Failed to create log directory: %v", err)
		}
	}
	logPath = path.Join(logPath, name)

	errWriter, err := rotatelogs.New(
		logPath+"_err_%Y-%m-%d.log",
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotationTime),
	)
	if err != nil {
		log.Fatalf("Failed to create error log file: %v", err)
	}

	infoWriter, err := rotatelogs.New(
		logPath+"_info_%Y-%m-%d.log",
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotationTime),
	)
	if err != nil {
		log.Fatalf("Failed to create info log file: %v", err)
	}

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	errorCore := zapcore.AddSync(errWriter)
	infocore := zapcore.AddSync(infoWriter)
	fileEncodeConfig := zap.NewProductionEncoderConfig()
	fileEncodeConfig.EncodeTime = timeEncoder
	if callerFullPath {
		fileEncodeConfig.EncodeCaller = customCallerEncoder
	}
	fileEncoder := zapcore.NewConsoleEncoder(fileEncodeConfig)

	cores = make([]zapcore.Core, 0)
	cores = append(cores, zapcore.NewCore(fileEncoder, errorCore, highPriority))
	cores = append(cores, zapcore.NewCore(fileEncoder, infocore, lowPriority))
	log.Println("File cores created and appended successfully")
	return cores
}

// customCallerEncoder set caller fullpath
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(caller.FullPath())
}

// timeEncoder format time
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
