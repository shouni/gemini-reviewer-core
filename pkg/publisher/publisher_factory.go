package publisher

import (
	"context"
	"fmt"
	"net/url"

	"github.com/shouni/go-remote-io/pkg/gcsfactory"
	"github.com/shouni/go-remote-io/pkg/s3factory"
)

// FactoryRegistry は、必要な外部依存関係のファクトリ群をまとめた構造体です。
// これらをNewPublisherに渡すことで、依存性の注入を実現します。
type FactoryRegistry struct {
	GCSFactory gcsfactory.Factory
	S3Factory  s3factory.Factory
}

// NewPublisher は、指定されたURIスキームに基づいて、適切な Publisher 実装を構築して返します。
// URIがどのスキームにも一致しない場合、エラーを返します。
func NewPublisher(ctx context.Context, uri string, registry FactoryRegistry) (Publisher, error) {
	// 1. URIをパースし、スキームを抽出
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("URIのパースに失敗しました: %w", err)
	}

	scheme := u.Scheme

	// 2. スキームに基づいてファクトリを選択し、Publisherを構築
	switch scheme {
	case "gs":
		if registry.GCSFactory == nil {
			return nil, fmt.Errorf("GCS URIが指定されましたが、GCS Factoryがnilです")
		}
		return NewGCSPublisher(registry.GCSFactory)

	case "s3":
		if registry.S3Factory == nil {
			return nil, fmt.Errorf("S3 URIが指定されましたが、S3 Factoryがnilです")
		}
		return NewS3Publisher(registry.S3Factory)

	default:
		return nil, fmt.Errorf("サポートされていないURIスキームです: %s (サポート: gs://, s3://)", scheme)
	}
}
