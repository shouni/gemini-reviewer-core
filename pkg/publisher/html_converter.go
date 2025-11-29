package publisher

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/shouni/go-utils/timeutil"
)

// =================================================================
// ヘルパーロジック (共通)
// =================================================================

// convertMarkdownToHTML の実体
// どちらの Publisher からも呼び出せるように、レシーバなしの関数として定義しています。
func convertMarkdownToHTML(ctx context.Context, data ReviewData) (io.Reader, error) {
	// NewMarkdownToHtmlRunner は md_adapter.go で定義されている
	markdownRunner, err := NewMarkdownToHtmlRunner(ctx)
	if err != nil {
		return nil, err
	}

	nowJST := timeutil.NowJST()
	reviewTimeStr := nowJST.Format("2006/01/02 15:04:05 JST")

	summaryMarkdown := fmt.Sprintf(
		"レビュー対象リポジトリ: `%s`\n\nブランチ差分: `%s` ← `%s`\n\nレビュー実行日時: *%s*\n\n",
		data.RepoURL,
		data.BaseBranch,
		data.FeatureBranch,
		reviewTimeStr,
	)

	var buffer bytes.Buffer
	buffer.WriteString("# " + reviewTitle + "\n\n")
	buffer.WriteString(summaryMarkdown + "\n\n")
	buffer.WriteString(data.ReviewMarkdown) // 本文を追加

	// 実際の MarkdownToHtmlRunner の Run メソッドを実行
	return markdownRunner.Run(ctx, buffer.Bytes())
}
