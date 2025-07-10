package auto_env

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mengri/utils-store/store"
	store_mysql "github.com/mengri/utils-store/store/mysql"
	"github.com/mengri/utils/autowire-v2"
	"log"
	"os"
	"strings"
)

func init() {
	autowire.Auto(createDbFromEnv)
}

type config struct {
	DatabaseURL string `mapstructure:"DATABASE_URL" validate:"required"`
}

func createDbFromEnv() store.IDB {
	databaseUrl, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	return store_mysql.CreateDb(databaseUrl)
}
func loadConfig() (string, error) {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		// .env 文件不存在不是错误
		fmt.Println("Warning: .env file not found, using environment variables only")
	}

	dataBaseUrl := os.Getenv("DATABASE_URL")

	if dataBaseUrl == "" {
		return "", fmt.Errorf("DATABASE_URL is required")
	}
	
	return strings.TrimSuffix(dataBaseUrl, "mysql://"), nil
}
