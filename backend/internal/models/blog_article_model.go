package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BlogArticle represents a blog article in MongoDB
type BlogArticle struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	KeywordID   int       `bson:"keyword_id" json:"keyword_id"`
	Title       string    `bson:"title" json:"title"`
	URL         string    `bson:"url" json:"url"`
	Author      string    `bson:"author" json:"author"`
	Content     string    `bson:"content" json:"content"`
	PublishDate time.Time `bson:"publish_date" json:"publish_date"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

// articlesCollection returns the blog articles collection
func articlesCollection() *mongo.Collection {
	return MongoDB.Collection("blog_articles")
}

// CreateBlogArticle adds a new blog article to the database
func CreateBlogArticle(ctx context.Context, article *BlogArticle) error {
	if article.CreatedAt.IsZero() {
		article.CreatedAt = time.Now()
	}

	_, err := articlesCollection().InsertOne(ctx, article)
	return err
}

// GetBlogArticlesForKeyword retrieves blog articles for a specific keyword
func GetBlogArticlesForKeyword(ctx context.Context, keywordID int, limit int) ([]*BlogArticle, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}

	opts := options.Find().
		SetSort(bson.D{primitive.E{Key: "publish_date", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := articlesCollection().Find(ctx, 
		bson.M{"keyword_id": keywordID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []*BlogArticle
	if err := cursor.All(ctx, &articles); err != nil {
		return nil, err
	}

	return articles, nil
}

// GetBlogArticlesByDate retrieves blog articles published on a specific date
func GetBlogArticlesByDate(ctx context.Context, date time.Time) ([]*BlogArticle, error) {
	// 指定された日付の開始と終了を計算
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// 検索条件を設定
	filter := bson.M{
		"publish_date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	// 検索オプションを設定
	opts := options.Find().
		SetSort(bson.D{primitive.E{Key: "publish_date", Value: -1}})

	// クエリを実行
	cursor, err := articlesCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 結果を変換
	var articles []*BlogArticle
	if err := cursor.All(ctx, &articles); err != nil {
		return nil, err
	}

	return articles, nil
} 