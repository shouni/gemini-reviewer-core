package publisher

import (
	"fmt"

	"github.com/shouni/go-remote-io/pkg/gcsfactory"
	"github.com/shouni/go-remote-io/pkg/remoteio"
	"github.com/shouni/go-remote-io/pkg/s3factory"
)

// FactoryRegistry は、必要な外部依存関係のファクトリ群をまとめた構造体です。
// これらをNewPublisherに渡すことで、依存性の注入を実現します。
type FactoryRegistry struct {
	GCSFactory gcsfactory.Factory
	S3Factory  s3factory.Factory
}

// NewPublisher は、指定されたURIスキームに基づいて、適切な Publisher 実装を構築して返します。
func NewPublisher(uri string, registry FactoryRegistry) (Publisher, error) {

	if remoteio.IsGCSURI(uri) {
		if registry.GCSFactory == nil {
			return nil, fmt.Errorf("GCS URIが指定されましたが、必要なGCS Factoryがnilです")
		}
		return NewGCSPublisher(registry.GCSFactory)
	}

	if remoteio.IsS3URI(uri) {
		if registry.S3Factory == nil {
			return nil, fmt.Errorf("S3 URIが指定されましたが、必要なS3 Factoryがnilです")
		}
		return NewS3Publisher(registry.S3Factory)
	}

	return nil, fmt.Errorf("サポートされていないURIスキームです: %s (サポート: gs://, s3://)", uri)
}
