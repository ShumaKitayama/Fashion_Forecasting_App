package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Image represents an image document in MongoDB
type Image struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	KeywordID int                `bson:"keyword_id" json:"keyword_id"`
	ImageURL  string             `bson:"image_url" json:"image_url"`
	Caption   string             `bson:"caption" json:"caption"`
	Tags      []string           `bson:"tags" json:"tags"`
	FetchedAt time.Time          `bson:"fetched_at" json:"fetched_at"`
}

// SocialMediaPost represents a social media post document in MongoDB
type SocialMediaPost struct {
	ID          string    `bson:"_id,omitempty"`
	KeywordID   int       `bson:"keyword_id"`
	Platform    string    `bson:"platform"`
	PostID      string    `bson:"post_id"`
	Username    string    `bson:"username"`
	Caption     string    `bson:"caption"`
	ImageURL    string    `bson:"image_url"`
	LikeCount   int       `bson:"like_count"`
	CommentCount int      `bson:"comment_count"`
	PostDate    time.Time `bson:"post_date"`
	CreatedAt   time.Time `bson:"created_at"`
}

// imagesCollection returns the images collection
func imagesCollection() *mongo.Collection {
	return MongoDB.Collection("images")
}

// CreateImage adds a new image to the database
func CreateImage(ctx context.Context, image *Image) error {
	if image.ID.IsZero() {
		image.ID = primitive.NewObjectID()
	}
	
	if image.FetchedAt.IsZero() {
		image.FetchedAt = time.Now()
	}

	_, err := imagesCollection().InsertOne(ctx, image)
	return err
}

// GetImagesForKeyword retrieves images for a specific keyword
func GetImagesForKeyword(ctx context.Context, keywordID int, limit int) ([]*Image, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}

	opts := options.Find().
		SetSort(bson.D{primitive.E{Key: "fetched_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := imagesCollection().Find(ctx, 
		bson.M{"keyword_id": keywordID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var images []*Image
	if err := cursor.All(ctx, &images); err != nil {
		return nil, err
	}

	return images, nil
}

// GetImagesByKeywordAndDate retrieves images for a specific keyword on a specific date
func GetImagesByKeywordAndDate(ctx context.Context, keywordID int, date time.Time) ([]*Image, error) {
	// 指定された日付の開始と終了を計算
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// 検索条件を設定
	filter := bson.M{
		"keyword_id": keywordID,
		"fetched_at": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	// 検索オプションを設定（最新のものから）
	opts := options.Find().
		SetSort(bson.D{primitive.E{Key: "fetched_at", Value: -1}})

	// クエリを実行
	cursor, err := imagesCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 結果を変換
	var images []*Image
	if err := cursor.All(ctx, &images); err != nil {
		return nil, err
	}

	return images, nil
}

// GetRecentImages retrieves recent images across all keywords
func GetRecentImages(ctx context.Context, limit int) ([]*Image, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}

	opts := options.Find().
		SetSort(bson.D{primitive.E{Key: "fetched_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := imagesCollection().Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var images []*Image
	if err := cursor.All(ctx, &images); err != nil {
		return nil, err
	}

	return images, nil
}

// GetImagesByTags retrieves images that contain specific tags
func GetImagesByTags(ctx context.Context, tags []string, limit int) ([]*Image, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}

	opts := options.Find().
		SetSort(bson.D{primitive.E{Key: "fetched_at", Value: -1}}).
		SetLimit(int64(limit))

	filter := bson.M{"tags": bson.M{"$in": tags}}
	
	cursor, err := imagesCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var images []*Image
	if err := cursor.All(ctx, &images); err != nil {
		return nil, err
	}

	return images, nil
}

// DeleteImage removes an image from the database
func DeleteImage(ctx context.Context, id primitive.ObjectID) error {
	result, err := imagesCollection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// GetImagesByKeywordAndDateRange retrieves images for a keyword within a date range
func GetImagesByKeywordAndDateRange(ctx context.Context, keywordID int, startDate, endDate time.Time) ([]Image, error) {
	filter := bson.M{
		"keyword_id": keywordID,
		"fetched_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	opts := options.Find().
		SetSort(bson.D{primitive.E{Key: "fetched_at", Value: -1}})

	cursor, err := imagesCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var images []Image
	for cursor.Next(ctx) {
		var image Image
		if err := cursor.Decode(&image); err != nil {
			continue // Skip invalid documents
		}
		images = append(images, image)
	}

	return images, cursor.Err()
} 