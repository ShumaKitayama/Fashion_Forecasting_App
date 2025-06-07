package dbcheck

import (
	"fmt"
	"os"
	"strings"

	"github.com/trendscout/backend/cmd/dbcheck/check_keywords"
	"github.com/trendscout/backend/cmd/dbcheck/check_user"
	"github.com/trendscout/backend/cmd/dbcheck/reset_password"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch strings.ToLower(command) {
	case "users":
		check_user.Run()
	case "keywords":
		check_keywords.Run()
	case "reset":
		reset_password.Run()
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("使用法: dbcheck <コマンド>")
	fmt.Println("利用可能なコマンド:")
	fmt.Println("  users     - ユーザーテーブルを確認")
	fmt.Println("  keywords  - キーワードテーブルを確認")
	fmt.Println("  reset     - 管理者パスワードをリセット")
} 