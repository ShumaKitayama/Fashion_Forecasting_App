package visualizer

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/trendscout/backend/internal/models"
	"github.com/trendscout/backend/internal/trend"
)

func main() {
	// コマンドライン引数からキーワードIDを取得
	if len(os.Args) < 2 {
		log.Fatal("使用法: trend_visualizer <keyword_id>")
	}
	
	keywordID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("無効なキーワードID: %v", err)
	}
	
	// キーワード情報の取得
	ctx := context.Background()
	keyword, err := models.GetKeywordByID(ctx, keywordID)
	if err != nil || keyword == nil {
		log.Fatalf("キーワード取得エラー: %v", err)
	}
	
	// 過去データの取得
	records, err := models.GetLatestTrendRecords(ctx, keywordID, 30)
	if err != nil {
		log.Fatalf("トレンドデータ取得エラー: %v", err)
	}
	
	if len(records) == 0 {
		log.Fatalf("キーワード '%s' のトレンドデータがありません", keyword.Keyword)
	}
	
	// 予測の実行
	trendService := trend.NewService()
	predictions, err := trendService.PredictTrend(ctx, keywordID, 14)
	if err != nil {
		log.Printf("予測エラー: %v", err)
		log.Println("予測データなしでグラフを生成します")
		predictions = nil
	}
	
	// テキスト形式でデータを出力
	fmt.Printf("「%s」のトレンド分析\n", keyword.Keyword)
	fmt.Println("========================")
	fmt.Println("日付 | ボリューム")
	fmt.Println("------------------------")
	
	// 実績データ
	for _, r := range records {
		fmt.Printf("%s | %d\n", r.RecordDate.Format("2006-01-02"), r.Volume)
	}
	
	// 予測データ
	if predictions != nil && len(predictions) > 0 {
		fmt.Println("------------------------")
		fmt.Println("予測データ:")
		for _, p := range predictions {
			fmt.Printf("%s | %d\n", p.Date.Format("2006-01-02"), p.Volume)
		}
	}
	
	fmt.Println("========================")
	fmt.Println("注: HTML形式のグラフを生成するには go-echarts パッケージをインストールしてください")
	fmt.Println("go get github.com/go-echarts/go-echarts/v2")
} 