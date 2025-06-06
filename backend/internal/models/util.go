package models

import (
	"context"
	"strings"
)

// CreateTestUser はテスト用のユーザーを作成または検索する補助関数です
func CreateTestUser(ctx context.Context) (*User, error) {
	// テスト用ユーザーを検索
	email := "test_scraper@example.com"
	password := "test_password"
	
	user, err := GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	
	if user != nil {
		return user, nil
	}
	
	// 存在しない場合は新規作成
	return CreateUser(ctx, email, password)
}

// SanitizeString は文字列を安全に処理するユーティリティ関数です
func SanitizeString(s string) string {
	// 基本的なサニタイズ処理（必要に応じて拡張）
	// 先頭と末尾の空白を削除
	return strings.TrimSpace(s)
} 