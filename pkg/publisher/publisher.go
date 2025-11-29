package publisher

import (
	"context"
)

// --- 定数定義 ---
const (
	contentTypeHTML = "text/html; charset=utf-8"
	reviewTitle     = "AIコードレビュー結果"
)

// ReviewData はレポート生成に必要なすべての情報をまとめた構造体です。
type ReviewData struct {
	RepoURL        string
	BaseBranch     string
	FeatureBranch  string
	ReviewMarkdown string
}

// =================================================================
// 共通インターフェース
// =================================================================

// Publisher はレビュー結果を指定されたURIに公開する最上位の抽象インターフェースです。
// このインターフェースにより、呼び出し元はGCSやS3の実装を意識する必要がなくなります。
type Publisher interface {
	Publish(ctx context.Context, uri string, data ReviewData) error
}
