package publisher

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/shouni/go-remote-io/pkg/remoteio"
	"github.com/shouni/go-remote-io/pkg/s3factory"
)

// =================================================================
// S3 Publisher の定義
// =================================================================

// S3OutputWriter は go-remote-io の Writer が満たすべき S3 書き込みインターフェースです。
type S3OutputWriter interface {
	WriteToS3(ctx context.Context, bucketName, objectPath string, reader io.Reader, contentType string) error
}

// S3Publisher は AIレビュー結果をS3バケットに公開する責務を担います。
type S3Publisher struct {
	writer S3OutputWriter
}

// NewS3Publisher は新しい S3Publisher インスタンスを作成します。
func NewS3Publisher(factory s3factory.Factory) (*S3Publisher, error) {
	writer, err := factory.NewOutputWriter()
	if err != nil {
		return nil, fmt.Errorf("OutputWriterの生成に失敗しました: %w", err)
	}

	w, ok := writer.(S3OutputWriter)
	if !ok {
		return nil, fmt.Errorf("writer が S3OutputWriter インターフェースを実装していません")
	}

	return &S3Publisher{
		writer: w,
	}, nil
}

// Publish メイン処理
// S3Publisher は Publisher インターフェースを満たします。
func (p *S3Publisher) Publish(ctx context.Context, uri string, data ReviewData) error {
	// remoteio.ParseS3URI を使用してバケット名とオブジェクトパスをパース
	bucketName, objectPath, err := remoteio.ParseS3URI(uri)
	if err != nil {
		return err
	}

	htmlReader, err := convertMarkdownToHTML(ctx, data)
	if err != nil {
		return fmt.Errorf("HTML変換に失敗しました: %w", err)
	}

	slog.Info("S3へアップロード開始", "bucketName", bucketName, "objectPath", objectPath)
	if err := p.writer.WriteToS3(ctx, bucketName, objectPath, htmlReader, contentTypeHTML); err != nil {
		return fmt.Errorf("S3への書き込みに失敗しました: %w", err)
	}

	return nil
}
