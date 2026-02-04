# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

miniflux-syncは、Miniflux RSSフィードリーダーの設定をYAMLファイルで管理・同期するGo製CLIツール。バージョン管理による宣言的なフィード管理を実現する。

## 開発環境

Nix環境を使用。すべてのGoコマンドは `nix develop -c` を経由して実行する。

```bash
# テスト実行（カバレッジ付き）
nix develop -c go test -cover ./...

# 特定パッケージのテスト
nix develop -c go test -cover ./diff

# 単一テスト実行
nix develop -c go test -run TestCalculateDiff ./diff

# リント（GITHUB_TOKEN必要）
curl -sSL https://github.com/revett/dotfiles/raw/main/.golangci.yml -o .golangci.yml
nix develop -c golangci-lint run --config .golangci.yml

# リリース（GITHUB_TOKEN必要）
./scripts/release.sh
```

## アーキテクチャ

### 同期フロー
1. **ローカル状態**（YAML）→ `diff.State` にパース
2. **リモート状態**（Miniflux API）→ `diff.State` にフェッチ
3. **差分計算** → `[]diff.Action` を生成（フィード・カテゴリの作成/削除）
4. **アクション適用** → リモートMinifluxインスタンスを更新

### 主要データ構造

```go
// diff.State - 状態の中心的な表現
type State struct {
    FeedURLsByCategoryTitle map[string][]string  // カテゴリ → フィードURL群
}

// diff.Action - 同期操作
type Action struct {
    Type          ActionType  // CreateCategory | CreateFeed | DeleteCategory | DeleteFeed
    CategoryTitle string
    FeedURL       string
}
```

### パッケージ構成

- `cmd/` - CLIコマンド実装（sync, dump）
- `api/` - Miniflux APIクライアント、フェッチ、更新、状態変換
- `diff/` - 状態比較とアクション計算（コアアルゴリズム）
- `parse/` - YAMLパース
- `config/` - CLIフラグ解析
- `log/` - ロギングユーティリティ（zerologラッパー）

### 差分アルゴリズム（diff/diff.go）
1. リモートにあってローカルにないフィード → DeleteFeed
2. リモートにあってローカルにないカテゴリ → DeleteCategory
3. ローカルにあってリモートにないカテゴリ → CreateCategory
4. ローカルにあってリモートにないフィード → CreateFeed
5. アクションはソートされる: CreateCategory → CreateFeed → DeleteFeed → DeleteCategory

## YAMLフォーマット

```yaml
カテゴリ名:
  - https://example.com/feed1.xml
  - https://example.com/feed2.xml
```

## テストパターン

- テーブル駆動テスト with サブテスト
- testify/require を使用
- `t.Parallel()` で並行実行
- `diff/diff_test.go` に包括的なテストシナリオあり
