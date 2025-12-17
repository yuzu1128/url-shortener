# URL Shortener

Go で構築したURL短縮サービスのREST API。

## 技術スタック

- **Go** 1.25
- **標準ライブラリのみ**（外部依存なし）

## 機能

### エンドポイント

- `POST /shorten` - URLを短縮
- `GET /{shortCode}` - 短縮URLからリダイレクト
- `GET /stats/{shortCode}` - 統計情報取得

## 使用例

### URL短縮

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/very/long/url"}'
```

レスポンス:
```json
{
  "short_url": "http://localhost:8080/abc12345",
  "short_code": "abc12345",
  "original_url": "https://example.com/very/long/url"
}
```

### リダイレクト

ブラウザで `http://localhost:8080/abc12345` にアクセスすると、元のURLにリダイレクトされます。

### 統計情報

```bash
curl http://localhost:8080/stats/abc12345
```

レスポンス:
```json
{
  "short_code": "abc12345",
  "original_url": "https://example.com/very/long/url",
  "clicks": 5,
  "created_at": "2025-12-17T15:00:00Z"
}
```

## セットアップ

```bash
# ビルド
go build -o url-shortener

# 実行
./url-shortener

# または直接実行
go run main.go
```

## 開発情報

- **開発者**: YuZu
- **開発期間**: 2025年12月
- **目的**: ポートフォリオ用プロジェクト
