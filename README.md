# Go Todo API

外部ライブラリを使わず、Go標準パッケージ(`net/http`)のみで構成したシンプルなTODO REST APIです。データストアはインメモリで、プロセスを再起動するとデータは失われます。

## 特徴

- 外部ライブラリ不使用(`net/http`, `encoding/json`, `log/slog`など標準パッケージのみ)
- インメモリストア(`sync.RWMutex`で排他制御)
- レイヤードアーキテクチャ(`Repository` → `Service` → `Handler`)
- `context.Context`を全層に伝播(将来DBなど別実装への差し替えを見越した設計)
- グレースフルシャットダウン対応
- PATCHによる部分更新(送られたフィールドだけを更新)

## ディレクトリ構成

```
.
├── cmd/server/main.go          # エントリーポイント、DI、サーバー起動
├── internal/
│   ├── todo/
│   │   ├── model.go            # Todoモデル
│   │   ├── repository.go       # インメモリストア(CRUD)
│   │   ├── service.go          # バリデーション・ビジネスロジック
│   │   ├── handler.go          # HTTPハンドラ
│   │   └── *_test.go
│   ├── server/
│   │   ├── router.go           # ルーティング定義
│   │   └── router_test.go
│   └── respond/
│       └── response.go         # JSONレスポンス共通ヘルパー
├── Dockerfile
├── docker-compose.yaml
└── .air.toml                   # ホットリロード設定(air)
```

## エンドポイント一覧

| メソッド | パス | 説明 |
|---|---|---|
| GET | `/health` | ヘルスチェック |
| GET | `/todos` | Todo一覧を取得 |
| POST | `/todos` | Todoを新規作成 |
| GET | `/todos/{id}` | Todoを1件取得 |
| PATCH | `/todos/{id}` | Todoを部分更新(送ったフィールドのみ反映) |
| DELETE | `/todos/{id}` | Todoを削除 |

## リクエスト/レスポンス例

### ヘルスチェック

```
GET /health
```

```json
{
  "code": "OK",
  "message": "ok"
}
```

### Todoの作成

```
POST /todos
Content-Type: application/json

{"description": "牛乳を買う"}
```

`201 Created`

```json
{
  "id": 1,
  "description": "牛乳を買う",
  "completed": false,
  "created_at": "2026-06-30T18:02:06.116519818+09:00",
  "updated_at": "2026-06-30T18:02:06.116519889+09:00"
}
```

### Todo一覧の取得

```
GET /todos
```

`200 OK`

```json
[
  {
    "id": 1,
    "description": "牛乳を買う",
    "completed": false,
    "created_at": "2026-06-30T18:02:06.116519818+09:00",
    "updated_at": "2026-06-30T18:02:06.116519889+09:00"
  }
]
```

### Todoの部分更新(PATCH)

`completed`だけを送れば、`description`は変更されません。

```
PATCH /todos/1
Content-Type: application/json

{"completed": true}
```

`200 OK`

```json
{
  "id": 1,
  "description": "牛乳を買う",
  "completed": true,
  "created_at": "2026-06-30T18:02:06.116519818+09:00",
  "updated_at": "2026-06-30T18:02:06.996130269+09:00"
}
```

### Todoの削除

```
DELETE /todos/1
```

`204 No Content`(ボディなし)

### エラーレスポンス

存在しないIDへのアクセスや、バリデーションエラー時は以下の形式で返ります。

```json
{
  "code": "Not Found",
  "message": "todo: not found"
}
```

| ステータス | 発生条件 |
|---|---|
| `400 Bad Request` | リクエストボディが不正、`description`が空文字、`id`が整数でない |
| `404 Not Found` | 指定した`id`のTodoが存在しない |
| `500 Internal Server Error` | サーバー内部の予期しないエラー |

## 開発環境の立ち上げ方

### Docker Composeを使う場合(推奨)

[air](https://github.com/air-verse/air)によるホットリロード付きで起動します。

```sh
docker compose up
```

`http://localhost:8080` でアクセスできます。コード変更は自動で反映されます。

### ローカルのGoで直接起動する場合

```sh
go run ./cmd/server
```

デフォルトでは`PORT`環境変数を見て、未設定の場合は`8080`番ポートで起動します。

```sh
PORT=3000 go run ./cmd/server
```

## テストの実行

```sh
# 全テスト実行
go test ./...

# レースディテクタ付き
go test ./... -race

# カバレッジ計測(パッケージをまたいだ実行も含める)
go test ./... -coverpkg=./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## 動作確認(curl)

```sh
# 作成
curl -X POST http://localhost:8080/todos -d '{"description":"牛乳を買う"}'

# 一覧取得
curl http://localhost:8080/todos

# 部分更新
curl -X PATCH http://localhost:8080/todos/1 -d '{"completed":true}'

# 削除
curl -X DELETE http://localhost:8080/todos/1
```
