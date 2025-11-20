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

// 定数定義
const (
	contentTypeHTML = "text/html; charset=utf-8"
	reviewTitle     = "AIコードレビュー結果"
)

// ReviewMetadata はレポート生成および保存に必要なメタデータです。
type ReviewMetadata struct {
	RepoURL        string
	BaseBranch     string
	FeatureBranch  string
	DestinationURI string
}

// OutputWriter は go-remote-io の Writer が満たすべきインターフェースです。
type OutputWriter interface {
	WriteToGCS(ctx context.Context, bucket, object string, reader io.Reader, contentType string) error
}

// GCSPublisher はレビュー結果をGCSに公開するための構造体です。
type GCSPublisher struct {
	writer OutputWriter // 生成済みのWriter
}

// NewGCSPublisher は GCSPublisher のコンストラクタです。
func NewGCSPublisher(ioFactory factory.Factory) (*GCSPublisher, error) {
	// 1. Writerの生成 (Fail Fast)
	// アプリケーション起動時にGCSクライアント等の初期化エラーを検知します。
	writer, err := ioFactory.NewOutputWriter()
	if err != nil {
		return nil, fmt.Errorf("OutputWriterの生成に失敗しました: %w", err)
	}

	// インターフェースへのキャストチェック
	w, ok := writer.(OutputWriter)
	if !ok {
		return nil, fmt.Errorf("writer が OutputWriter インターフェースを実装していません")
	}

	return &GCSPublisher{
		writer: w,
	}, nil
}

// Publish はメインの公開処理を行います。
// 指定された DestinationURI を解析し、MarkdownをHTMLに変換してアップロードします。
func (p *GCSPublisher) Publish(ctx context.Context, reviewMarkdown string, meta ReviewMetadata) error {
	// 1. 保存先URIの解析 (Runtime)
	// ヘルパーメソッドを呼び出します（ここで空チェックも行われます）
	bucketName, objectPath, err := p.parseDestination(meta.DestinationURI)
	if err != nil {
		return err // parseDestination 内でエラーメッセージは整形済み
	}

	// 2. Markdown -> HTML 変換
	htmlReader, err := p.convertMarkdownToHTML(ctx, reviewMarkdown, meta)
	if err != nil {
		return fmt.Errorf("HTML変換に失敗しました: %w", err)
	}

	slog.Info("GCSへアップロード開始", "bucket", bucketName, "path", objectPath)

	// 3. Upload処理
	if err := p.writer.WriteToGCS(ctx, bucketName, objectPath, htmlReader, contentTypeHTML); err != nil {
		return fmt.Errorf("GCSへの書き込みに失敗しました: %w", err)
	}

	return nil
}

// convertMarkdownToHTML はMarkdownをHTMLに変換し、ヘッダー情報を付与する内部ヘルパーメソッドです。
func (p *GCSPublisher) convertMarkdownToHTML(ctx context.Context, reviewMarkdown string, meta ReviewMetadata) (io.Reader, error) {
	markdownRunner, err := adapters.NewMarkdownToHtmlRunner(ctx)
	if err != nil {
		return nil, err
	}

	// meta の情報を使ってレポートのヘッダーを作成
	summaryMarkdown := fmt.Sprintf(
		"レビュー対象リポジトリ: `%s`\n\nブランチ差分: `%s` ← `%s`\n\n",
		meta.RepoURL,
		meta.BaseBranch,
		meta.FeatureBranch,
	)

	var buffer bytes.Buffer
	buffer.WriteString("# " + reviewTitle + "\n\n")
	buffer.WriteString(summaryMarkdown + "\n\n")
	buffer.WriteString(reviewMarkdown)

	return markdownRunner.Run(ctx, buffer.Bytes())
}

// parseDestination はURIを解析し、エラーハンドリングをラップします。
func (p *GCSPublisher) parseDestination(uri string) (string, string, error) {
	if uri == "" {
		return "", "", fmt.Errorf("保存先URI (DestinationURI) が空です")
	}

	bucket, object, err := remoteio.ParseGCSURI(uri)
	if err != nil {
		// エラーメッセージの整形をここに閉じ込められます
		return "", "", fmt.Errorf("無効な保存先URIです (%s): %w", uri, err)
	}
	return bucket, object, nil
}
