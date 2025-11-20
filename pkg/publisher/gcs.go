package publisher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/shouni/gemini-reviewer-core/pkg/adapters"

	"github.com/shouni/go-remote-io/pkg/factory"
	"github.com/shouni/go-remote-io/pkg/remoteio"
)

// ReviewMetadata はレポートのヘッダー生成に必要な最小限のメタデータです。
// CLIのConfigやWebのRequestから、この構造体に詰め替えて渡します。
type ReviewMetadata struct {
	RepoURL       string
	BaseBranch    string
	FeatureBranch string
}

// ReviewPublisher インターフェース
type ReviewPublisher interface {
	Publish(ctx context.Context, reviewMarkdown string, meta ReviewMetadata) error
}

// GCSPublisher 実装
type GCSPublisher struct {
	ioFactory   factory.Factory
	targetURI   string
	contentType string
}

// NewGCSPublisher コンストラクタ
func NewGCSPublisher(ioFactory factory.Factory, uri string, contentType string) *GCSPublisher {
	if contentType == "" {
		contentType = "text/html; charset=utf-8"
	}
	return &GCSPublisher{
		ioFactory:   ioFactory,
		targetURI:   uri,
		contentType: contentType,
	}
}

// Publish メイン処理
func (p *GCSPublisher) Publish(ctx context.Context, reviewMarkdown string, meta ReviewMetadata) error {
	bucketName, objectPath, err := remoteio.ParseGCSURI(p.targetURI)
	if err != nil {
		return fmt.Errorf("GCS URI解析失敗: %w", err)
	}

	// meta を渡す
	htmlReader, err := p.convertMarkdownToHTML(ctx, reviewMarkdown, meta)
	if err != nil {
		return fmt.Errorf("HTML変換失敗: %w", err)
	}

	slog.Info("GCSへアップロード開始", "uri", p.targetURI)

	// Upload処理
	writer, err := p.ioFactory.NewOutputWriter()
	if err != nil {
		return err
	}

	// エラーハンドリングを追加（任意ですが、あると親切です）
	if err := writer.WriteToGCS(ctx, bucketName, objectPath, htmlReader, p.contentType); err != nil {
		return fmt.Errorf("GCSへの書き込み失敗: %w", err)
	}

	return nil
}

// convertMarkdownToHTML 内部ヘルパーメソッド
func (p *GCSPublisher) convertMarkdownToHTML(ctx context.Context, reviewMarkdown string, meta ReviewMetadata) (io.Reader, error) {
	// タイトル定義
	const ReviewTitle = "AIコードレビュー結果"

	markdownRunner, err := adapters.NewMarkdownToHtmlRunner(ctx)
	if err != nil {
		return nil, err
	}

	// meta の情報を使ってヘッダーを作成
	summaryMarkdown := fmt.Sprintf(
		"レビュー対象リポジトリ: `%s`\n\nブランチ差分: `%s` ← `%s`\n\n",
		meta.RepoURL,
		meta.BaseBranch,
		meta.FeatureBranch,
	)

	var buffer bytes.Buffer
	buffer.WriteString("## " + ReviewTitle + "\n\n")
	buffer.WriteString(summaryMarkdown + "\n\n")
	buffer.WriteString(reviewMarkdown)

	return markdownRunner.Run(ctx, buffer.Bytes())
}
