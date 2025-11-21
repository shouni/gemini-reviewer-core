package publisher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/shouni/go-remote-io/pkg/factory"
	"github.com/shouni/go-remote-io/pkg/remoteio"
)

// 定数定義
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

// OutputWriter は go-remote-io の Writer が満たすべきインターフェースです。
type OutputWriter interface {
	WriteToGCS(ctx context.Context, bucketName, objectPath string, reader io.Reader, contentType string) error
}

// GCSPublisher
type GCSPublisher struct {
	writer OutputWriter
}

// NewGCSPublisher
func NewGCSPublisher(ioFactory factory.Factory) (*GCSPublisher, error) {
	writer, err := ioFactory.NewOutputWriter()
	if err != nil {
		return nil, fmt.Errorf("OutputWriterの生成に失敗しました: %w", err)
	}

	w, ok := writer.(OutputWriter)
	if !ok {
		return nil, fmt.Errorf("writer が OutputWriter インターフェースを実装していません")
	}

	return &GCSPublisher{
		writer: w,
	}, nil
}

// Publish メイン処理
func (p *GCSPublisher) Publish(ctx context.Context, uri string, data ReviewData) error {
	bucketName, objectPath, err := remoteio.ParseGCSURI(uri)
	if err != nil {
		return err
	}

	htmlReader, err := p.convertMarkdownToHTML(ctx, data)
	if err != nil {
		return fmt.Errorf("HTML変換に失敗しました: %w", err)
	}

	slog.Info("GCSへアップロード開始", "bucketName", bucketName, "objectPath", objectPath)
	if err := p.writer.WriteToGCS(ctx, bucketName, objectPath, htmlReader, contentTypeHTML); err != nil {
		return fmt.Errorf("GCSへの書き込みに失敗しました: %w", err)
	}

	return nil
}

// convertMarkdownToHTML (ヘルパー)
func (p *GCSPublisher) convertMarkdownToHTML(ctx context.Context, data ReviewData) (io.Reader, error) {
	markdownRunner, err := NewMarkdownToHtmlRunner(ctx)
	if err != nil {
		return nil, err
	}

	// data の情報を使ってレポートのヘッダーを作成
	summaryMarkdown := fmt.Sprintf(
		"レビュー対象リポジトリ: `%s`\n\nブランチ差分: `%s` ← `%s`\n\n",
		data.RepoURL,
		data.BaseBranch,
		data.FeatureBranch,
	)

	var buffer bytes.Buffer
	buffer.WriteString("# " + reviewTitle + "\n\n")
	buffer.WriteString(summaryMarkdown + "\n\n")
	buffer.WriteString(data.ReviewMarkdown) // 本文を追加

	return markdownRunner.Run(ctx, buffer.Bytes())
}
