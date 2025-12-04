package publisher

import (
	"context"
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

// NewPublisherAndSigner は、URIに基づいてPublisherとURLSignerを初期化します。
func NewPublisherAndSigner(ctx context.Context, targetURI string) (Publisher, remoteio.URLSigner, error) {
	registry := FactoryRegistry{}
	var urlSigner remoteio.URLSigner
	var err error

	// Publisherの動的生成
	publisher, err := NewPublisher(targetURI, registry)
	if err != nil {
		return nil, nil, fmt.Errorf("パブリッシャーの初期化に失敗しました: %w", err)
	}

	// GCSまたはS3のどちらか必要なファクトリのみを初期化し、RegistryとSignerを設定
	if remoteio.IsGCSURI(targetURI) {
		gcsFactory, err := gcsfactory.NewGCSClientFactory(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("GCSクライアントファクトリの初期化に失敗しました: %w", err)
		}
		registry.GCSFactory = gcsFactory

		signer, err := gcsFactory.NewGCSURLSigner()
		if err != nil {
			return nil, nil, fmt.Errorf("GCS URL Signerの取得に失敗しました: %w", err)
		}
		urlSigner = signer

	} else if remoteio.IsS3URI(targetURI) {
		s3Factory, err := s3factory.NewS3ClientFactory(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("S3クライアントファクトリの初期化に失敗しました (URI: %s): %w", targetURI, err)
		}
		registry.S3Factory = s3Factory

		signer, err := s3Factory.NewS3URLSigner()
		if err != nil {
			return nil, nil, fmt.Errorf("S3 URL Signerの取得に失敗しました: %w", err)
		}
		urlSigner = signer

	} else {
		return nil, nil, fmt.Errorf("未対応のストレージURIです: %s", targetURI)
	}

	return publisher, urlSigner, nil
}
