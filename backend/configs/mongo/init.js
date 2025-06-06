// imagesコレクションの初期化
db = db.getSiblingDB("trendscout");

// 既存のコレクションを削除
db.images.drop();

// imagesコレクションの作成
db.createCollection("images");

// インデックスの作成
db.images.createIndex({ keyword_id: 1 });
db.images.createIndex({ fetched_at: 1 });

// サンプルデータの挿入
db.images.insertMany([
  {
    keyword_id: 1,
    image_url: "https://example.com/sample1.jpg",
    caption: "サンプル画像1",
    tags: ["fashion", "trend", "summer"],
    fetched_at: new Date(),
  },
  {
    keyword_id: 2,
    image_url: "https://example.com/sample2.jpg",
    caption: "サンプル画像2",
    tags: ["style", "casual", "spring"],
    fetched_at: new Date(),
  },
]);
