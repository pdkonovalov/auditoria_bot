package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

const (
	LogLevelDebug = "debug"
	LogLevelProd  = "prod"
)

const (
	EventIDCharsetLetters = "letters"
	EventIDCharsetNumbers = "numbers"
)

type Admin struct {
	UserID    int64  `yaml:"UserID"`
	FirstName string `yaml:"FirstName"`
	LastName  string `yaml:"LastName"`
}

type AdminList []Admin

func (f *AdminList) SetValue(v string) error {
	admin_list := make(AdminList, 0)
	for _, admin_str := range strings.Split(v, ", ") {
		admin_str_parts := strings.Split(admin_str, " ")
		if len(admin_str_parts) != 3 ||
			len(admin_str_parts[0]) == 0 ||
			len(admin_str_parts[1]) == 0 ||
			len(admin_str_parts[2]) == 0 {
			return fmt.Errorf("Not valid telegram admin list, must be: '<firstname> <lastname> <user_id>, <firstname> ... '")
		}
		admin := Admin{
			FirstName: admin_str_parts[0],
			LastName:  admin_str_parts[1],
		}
		user_id, err := strconv.ParseInt(admin_str_parts[2], 10, 64)
		if err != nil {
			return fmt.Errorf("Not valid user id, must be int64")
		}
		admin.UserID = user_id
		admin_list = append(admin_list, admin)
	}
	*f = admin_list
	return nil
}

type PaymentDetails struct {
	Account   string `yaml:"Account"`
	FirstName string `yaml:"FirstName"`
	LastName  string `yaml:"LastName"`
}

type PaymentDetailsList []PaymentDetails

func (f *PaymentDetailsList) SetValue(v string) error {
	payment_details_list := make(PaymentDetailsList, 0)
	for _, payment_details_str := range strings.Split(v, ", ") {
		parts := strings.Split(payment_details_str, " ")
		if len(parts) != 3 ||
			len(parts[0]) == 0 ||
			len(parts[1]) == 0 ||
			len(parts[2]) == 0 {
			return fmt.Errorf("Not valid payment details, must be: '<firstname> <lastname> <account>")
		}
		payment_details_list = append(payment_details_list,
			PaymentDetails{
				Account:   parts[2],
				FirstName: parts[0],
				LastName:  parts[1],
			},
		)
	}
	*f = payment_details_list
	return nil
}

type Config struct {
	LogLevel string `yaml:"LogLevel" env:"LOG_LEVEL"`

	TelegramBotUsername              string             `yaml:"TelegramBotUsername" env:"TELEGRAM_BOT_USERNAME"`
	TelegramBotToken                 string             `yaml:"TelegramBotToken" env:"TELEGRAM_BOT_TOKEN"`
	TelegramBotAdminList             AdminList          `yaml:"TelegramBotAdminList" env:"TELEGRAM_BOT_ADMIN_LIST"`
	TelegramBotTimezone              string             `yaml:"TelegramBotTimezone" env:"TELEGRAM_BOT_TIMEZONE"`
	TelegramBotDefaultPaymentDetails PaymentDetailsList `yaml:"TelegramBotDefaultPaymentDetails" env:"TELEGRAM_BOT_DEFAULT_PAYMENT_DETAILS"`

	PostgresUser     string `yaml:"PostgresUser" env:"POSTGRES_USER"`
	PostgresPassword string `yaml:"PostgresPassword" env:"POSTGRES_PASSWORD"`
	PostgresHost     string `yaml:"PostgresHost" env:"POSTGRES_HOST"`
	PostgresPort     string `yaml:"PostgresPort" env:"POSTGRES_PORT"`
	PostgresDatabase string `yaml:"PostgresDatabase" env:"POSTGRES_DB"`
	PostgresSslMode  string `yaml:"PostgresSslMode" env:"POSTGRES_SSL_MODE"`

	EventIDCharset string `yaml:"EventIDCharset" env:"EVENT_ID_CHARSET"`
	EventIDLen     int    `yaml:"EventIDLen" env:"EVENT_ID_LEN"`
}

func New() (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, err
		}
	}

	if cfg.LogLevel != LogLevelDebug && cfg.LogLevel != LogLevelProd {
		return nil, fmt.Errorf("Invalid LogLevel config variable value: '%s', must be %s or %s", cfg.LogLevel, LogLevelDebug, LogLevelProd)
	}

	return &cfg, nil
}

func (cfg *Config) StringSecureMasked() (string, error) {
	cfg_masked := new(Config)
	*cfg_masked = *cfg

	cfg_masked.TelegramBotToken = strings.Repeat("*", len(cfg_masked.TelegramBotToken))
	cfg_masked.PostgresPassword = strings.Repeat("*", len(cfg_masked.PostgresPassword))

	cfg_masked_yml, err := yaml.Marshal(cfg_masked)
	if err != nil {
		return "", fmt.Errorf("Error marshal config to yml: %s", err)
	}

	return string(cfg_masked_yml), nil
}
