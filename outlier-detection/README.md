# Envoy Outlier Detection サンプル

Envoy の Outlier Detection（異常検知によるホスト切り離し）機能のデモ環境です。

## 概要

Outlier Detection は、実際のトラフィックを監視し、エラー率の上昇や連続した失敗を検知して、一時的にそのエンドポイントを負荷分散から外す機能です。

### 設定値

| パラメータ | 値 | 説明 |
|-----------|-----|------|
| `consecutive_5xx` | 3 | 3回連続で 5xx エラーが発生したら切り離し |
| `interval` | 5s | 5秒ごとにスイープ（チェック）を行う |
| `base_ejection_time` | 30s | 切り離す基本時間 |
| `max_ejection_percent` | 50 | 切り離すホストの最大割合 |

## 構成

```
┌─────────┐     ┌─────────┐     ┌──────────┐
│   k6    │────▶│  Envoy  │────▶│ backend1 │
│ (load)  │     │ :10000  │  ┌─▶│  :8080   │
└─────────┘     └────┬────┘  │  └──────────┘
                     │       │
                     └───────┤  ┌──────────┐
                             └─▶│ backend2 │
                                │  :8080   │
                                └──────────┘
```

## 起動方法

```bash
# ビルド & 起動
docker-compose up -d --build

# ログ確認
docker-compose logs -f
```

## アクセス先

| サービス | URL | 説明 |
|----------|-----|------|
| Envoy | http://localhost:10000 | フロントエンド（リクエスト受付） |
| Envoy Admin | http://localhost:9901 | Envoy 管理画面 |
| backend1 | http://localhost:8081 | バックエンド1（直接アクセス） |
| backend2 | http://localhost:8082 | バックエンド2（直接アクセス） |
| Grafana | http://localhost:3000 | ダッシュボード（認証なし） |
| Prometheus | http://localhost:9090 | メトリクス |

## 動作確認

### 1. 正常時の確認

リクエストが `backend1` と `backend2` に交互に振り分けられることを確認します。

```bash
for i in {1..10}; do curl -s localhost:10000; done
```

### 2. エラーを発生させる

backend1 のエラーレートを 100% に設定します。

```bash
# backend1 を 100% エラーにする（直接アクセス）
curl -X POST "http://localhost:8081/error-rate?rate=100"

# 確認
curl http://localhost:8081/error-rate
```

> **Note**: Envoy 経由（`:10000`）ではなく、backend1 のポート（`:8081`）に直接アクセスしています。

### 3. Outlier Detection の発動確認

リクエストを送り続けると、backend1 が切り離されます。

```bash
# リクエストを送信（エラー後、backend2 のみになる）
for i in {1..20}; do curl -s localhost:10000; sleep 0.5; done
```

### 4. Envoy Admin で確認

```bash
# クラスタの状態確認
curl -s localhost:9901/clusters | grep -E "(backend|ejections)"
```

`outlier_detection.ejections_active` が `1` になっていれば、切り離しが発動しています。

### 5. エラーレートを戻す

```bash
# backend1 を正常に戻す
curl -X POST "http://localhost:8081/error-rate?rate=0"
```

`base_ejection_time`（30秒）経過後、backend1 が復帰します。

## 負荷テスト

k6 を使用して負荷をかけながら Grafana で状態を確認できます。

```bash
# k6 がインストールされている場合
k6 run k6/script.js

# Docker で実行する場合
docker run --rm -i --network=host grafana/k6 run - < k6/script.js
```

## Grafana ダッシュボード

http://localhost:3000 にアクセスすると、以下のメトリクスを確認できます：

- **Active Ejections** - 現在切り離されているホスト数
- **Request Rate by Response Code** - レスポンスコード別リクエストレート
- **Ejection Events** - 切り離しイベント（タイプ別）
- **Active Connections & Requests** - アクティブな接続・リクエスト数
- **Total Ejections** - 累計切り離し回数
- **Cluster Membership** - 健全なホスト数 vs 全ホスト数

## 停止

```bash
docker-compose down
```

## 参考

- [Envoy: Outlier Detection (Architecture Overview)](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/outlier)
- [Envoy: OutlierDetection (v3 API Reference)](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/cluster/v3/outlier_detection.proto)
