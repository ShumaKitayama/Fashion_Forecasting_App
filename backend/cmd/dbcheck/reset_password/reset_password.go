package reset_password

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func Run() {
	// .env.localを読み込む
	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatalf("Error loading .env.local file: %v", err)
	}

	log.Println("管理者ユーザーのパスワードをリセットします...")

	// データベース接続情報を取得
	dbHost := "localhost"
	dbUser := "trendscout"
	dbPassword := "trendscout"
	dbName := "trendscout"

	// 接続文字列を作成
	connString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		dbUser, dbPassword, dbHost, dbName)

	// データベースに接続
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("接続文字列の解析エラー: %v", err)
	}

	pgPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}
	defer pgPool.Close()

	// 接続テスト
	if err := pgPool.Ping(ctx); err != nil {
		log.Fatalf("データベース接続確認エラー: %v", err)
	}
	log.Println("PostgreSQLに接続しました")

	// 管理者ユーザーの確認
	var userID int
	var email string
	var passwordHash string

	err = pgPool.QueryRow(ctx, "SELECT id, email, password_hash FROM users WHERE email = 'admin@example.com'").Scan(&userID, &email, &passwordHash)
	if err != nil {
		log.Fatalf("管理者ユーザー検索エラー: %v", err)
	}

	log.Printf("管理者ユーザー検出: ID=%d, Email=%s", userID, email)
	log.Printf("現在のパスワードハッシュ: %s", passwordHash)

	// 新しいパスワードの生成
	password := "password"
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("パスワードハッシュ生成エラー: %v", err)
	}

	// ハッシュの比較
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		log.Printf("現在のパスワードハッシュと 'password' は一致しません: %v", err)
		log.Println("パスワードをリセットします...")

		// パスワードのリセット
		_, err = pgPool.Exec(ctx, "UPDATE users SET password_hash = $1 WHERE id = $2", newHashedPassword, userID)
		if err != nil {
			log.Fatalf("パスワード更新エラー: %v", err)
		}
		log.Printf("パスワードを 'password' にリセットしました")
	} else {
		log.Println("現在のパスワードハッシュは 'password' と一致します。更新は不要です。")
	}

	// 更新後のパスワードハッシュを確認
	var updatedHash string
	err = pgPool.QueryRow(ctx, "SELECT password_hash FROM users WHERE id = $1", userID).Scan(&updatedHash)
	if err != nil {
		log.Fatalf("更新後のハッシュ取得エラー: %v", err)
	}
	log.Printf("更新後のパスワードハッシュ: %s", updatedHash)

	log.Println("パスワードリセット処理が完了しました。")
	log.Println("以下の認証情報でログインできるようになりました:")
	log.Println("Email: admin@example.com")
	log.Println("Password: password")
}
