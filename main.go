package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"encoding/json"
	"go.uber.org/zap/zapcore"
)

// Logger and Bot global variables
var (
	sugar *zap.SugaredLogger
	bot   *tgbotapi.BotAPI
)

// Structure for extending
type Update struct {
	*tgbotapi.Update
}

// Command handler for logging purposes
func (u Update) Handle(c func(update Update) (*tgbotapi.Message, error)) {
	res, err := c(u)
	if err != nil {
		sugar.Error(err)
	}

	if res != nil {
		b, err := json.MarshalIndent(res, "", "\t")
		if err != nil {
			sugar.Warn(err)
		}

		sugar.Debugf("api response: %s", string(b))
	}
}

func main() {
	// Init and read config file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		sugar.Fatalf("fatal error config file: %v", err)
	}

	// Configuration defaults
	// Log level: INFO (-1 for DEBUG)
	viper.SetDefault("log.level", 0)
	// Log type: "production" or "development"
	viper.SetDefault("log.type", "production")

	// Init logger
	var loggerConfig zap.Config
	if viper.GetString("log.type") == "production" {
		loggerConfig = zap.NewProductionConfig()
	}
	if viper.GetString("log.type") == "development" {
		loggerConfig = zap.NewDevelopmentConfig()
	}
	loggerConfig.Level.SetLevel(zapcore.Level(viper.GetInt("log.level")))

	logger, _ := loggerConfig.Build()
	defer logger.Sync()

	sugar = logger.Sugar()

	// Check token
	if !viper.IsSet("token") {
		sugar.Fatalf("fatal error token null")
	}

	bot, err = tgbotapi.NewBotAPI(viper.GetString("token"))
	if err != nil {
		sugar.Fatalf("fatal error bot api init: %v", err)
	}

	sugar.Infof("authorized on @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for u := range updates {
		// Extend update
		update := Update{
			&u,
		}

		// Skip images, stickers, etc
		if update.Message == nil {
			continue
		}

		// Commands handler
		switch update.Message.Command() {
		case "start":
			go update.Handle(Start)
		case "help":
			go update.Handle(Help)
		default:
			go update.Handle(Default)
		}
	}
}
