package publisher

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/shouni/go-remote-io/pkg/gcsfactory"
	"github.com/shouni/go-remote-io/pkg/remoteio"
)

// =================================================================
// GCS Publisher の定義
// =================================================================

// GCSOutputWriter は go-remote-io の Writer が満たすべき GCS 書き込みインターフェースです。
type GCSOutputWriter interface {
	WriteToGCS(ctx context.Context, bucketName, objectPath string, reader io.Reader, contentType string) error
}

// GCSPublisher は AIレビュー結果をGCSに公開する責務を担います。
type GCSPublisher struct {
	writer GCSOutputWriter
}

// NewGCSPublisher は新しい GCSPublisher インスタンスを作成します。
func NewGCSPublisher(factory gcsfactory.Factory) (*GCSPublisher, error) {
	writer, err := factory.NewOutputWriter()
	if err != nil {
		return nil, fmt.Errorf("OutputWriterの生成に失敗しました: %w", err)
	}

	w, ok := writer.(GCSOutputWriter)
	if !ok {
		return nil, fmt.Errorf("writer が GCSOutputWriter インターフェースを実装していません")
	}

	return &GCSPublisher{
		writer: w,
	}, nil
}

// Publish メイン処理
// GCSPublisher は Publisher インターフェースを満たします。
func (p *GCSPublisher) Publish(ctx context.Context, uri string, data ReviewData) error {
	bucketName, objectPath, err := remoteio.ParseGCSURI(uri)
	if err != nil {
		return err
	}

	htmlReader, err := convertMarkdownToHTML(ctx, data)
	if err != nil {
		return fmt.Errorf("HTML変換に失敗しました: %w", err)
	}

	slog.Info("GCSへアップロード開始", "bucketName", bucketName, "objectPath", objectPath)
	if err := p.writer.WriteToGCS(ctx, bucketName, objectPath, htmlReader, contentTypeHTML); err != nil {
		return fmt.Errorf("GCSへの書き込みに失敗しました: %w", err)
	}

	return nil
}
