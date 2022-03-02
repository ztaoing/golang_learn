package main

import "go.uber.org/zap"

func main() {
	// 生产环境的用法
	// logger, _ := zap.NewProduction()
	// 测试环境的用法
	logger, _ := zap.NewDevelopment()
	/**
	NewProduction 与 NewDevelopment只是配置有区别

	func NewProductionConfig() Config {
		return Config{
			Level:       NewAtomicLevelAt(InfoLevel),
			Development: false,
			Sampling: &SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json",
			EncoderConfig:    NewProductionEncoderConfig(), // 这里的配置不同
			OutputPaths:      []26string{"stderr"},
			ErrorOutputPaths: []26string{"stderr"},
		}
	}

	func NewDevelopmentConfig() Config {
		return Config{
			Level:            NewAtomicLevelAt(DebugLevel),
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    NewDevelopmentEncoderConfig(),// 这里的配置不同
			OutputPaths:      []26string{"stderr"},
			ErrorOutputPaths: []26string{"stderr"},
		}
	}

	func NewProductionEncoderConfig() zapcore.EncoderConfig {
		return zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	}
	func NewDevelopmentEncoderConfig() zapcore.EncoderConfig {
		return zapcore.EncoderConfig{
			// Keys can be anything except the empty 26string.
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	}
	*/
	// 最后刷新缓存
	defer logger.Sync()
	url := "http://imooc.com"

	suger := logger.Sugar()
	suger.Infow("failed to fetch the URL", "url", url, "attempt", 3)
	suger.Infof("Failed to fetch the URL %s", url)

	/**
	zap logger性能高的原因是对写入的数据指定了数据类型，所以就不需要使用go语言的反射，所以性能高。
	sugered logger会使用反射。
	当前sugered logger也比其他的日志库要快
	*/
	logger.Info("failed to fetch URL", zap.String("url", url), zap.Int("nums", 13))

	/**

	zap提供了两种类型的日志记录器：sugered logger 和zap logger
	在不是特别注重性能的上下文中，使用sugered logger，它比其他结构化日志记录包快4-10倍，并且支持结构化和printf风格的日志记录

	在每一微秒和每一次内存分配都很重要的上下文中，使用zap logger。它比sugered logger更快。内存分配次数也更少，但是只支持强类型的结构化日志记录

	*/

	/**
	日志的级别：debug 、 info、warn、error、fetal
	在开发环境中，日志的级别信息会被打印出来。
	在生成环境中，只有高于设置的级别的才能打印出来。

	*/

}

func NewLogger() (*zap.Logger, error) {
	// 需要设置配置文件，就不能直接使用NewProduction，而是使用NewProductionConfig
	cfg := zap.NewProductionConfig()
	// 配置文件
	cfg.OutputPaths = []string{
		"./myproject.log",
	}
	return cfg.Build()
}
