package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const targetEnvName = "GO_ENV"

// envLoad env load
func EnvLoad() {
	// 環境未設定の場合はローカルを設定
	if "" == os.Getenv(targetEnvName) {
		_ = os.Setenv(targetEnvName, "local")
	}
	err := godotenv.Load(fmt.Sprintf("./common/env/%s.env", os.Getenv(targetEnvName)))
	if err != nil {
		log.Fatalf("Error loading env target env is %s", os.Getenv(targetEnvName))
	}
}
