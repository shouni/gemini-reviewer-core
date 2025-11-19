# 🤖 Gemini Reviewer Core

[![Language](https://img.shields.io/badge/Language-Go-blue)](https://golang.org/)
[![Go Version](https://img.shields.io/github/go-mod/go-version/shouni/gemini-reviewer-core)](https://golang.org/)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/shouni/gemini-reviewer-core)](https://github.com/shouni/gemini-reviewer-core/tags)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Gemini Reviewer Core** は、Google Gemini API を活用し、Gitリポジトリのブランチ間の差分を分析してAIコードレビューを自動生成するための**コアライブラリ**です。

本ライブラリは、**CLIツール**や**Webアプリケーション**の共通基盤として設計されており、Git操作、AI通信、プロンプト生成、HTML変換といった「ビジネスロジック」のみを提供します。保存先（ローカル/クラウド）や実行環境への依存を含みません。

## ✨ 技術スタック (Technology Stack)

| 要素 | 技術 / ライブラリ | 役割 |
| :--- | :--- | :--- |
| **言語** | **Go (Golang)** | ライブラリの開発言語。 |
| **Git 操作** | **go-git** (`adapters`層) | クローン、フェッチ、**3-dot diff** (共通祖先からの差分) の取得まですべてを Go のコード内で完結させ、**SSH認証とホストキー検証の設定**を統合しました。 |
| **AI モデル** | **Google Gemini API** (`adapters`層) | 取得したコード差分を分析し、レビューコメントを生成するために使用します。 |
| **Markdown to HTML** | **`go-text-format`** (`adapters`層) | AIが出力したMarkdown形式のレビュー結果を、スタイル付きの完全なHTMLドキュメントに**変換**するために使用します。|
| **プロンプト管理** | **`text/template` + `embed`** | レビューモード（Release/Detail）に応じたプロンプトテンプレートをバイナリに埋め込み、動的に生成します。 |

-----

## ✨ 主要な機能と特徴

### 1\. 🔍 高度なGit差分分析

  * **SSHネイティブ対応:** `go-git` を使用し、外部のSSHコマンドに依存せず、認証情報を注入可能な設計になっています。
  * **正確な差分取得:** ベースブランチとフィーチャーブランチ間のマージベース（共通祖先）を基準とした **3-dot diff (`A...B`)** を取得し、正確でクリーンな差分のみをAIに提供します。

### 2\. 🧱 責務の厳密な分離

  * **アダプターパターンの採用:** Git操作、Gemini APIへの通信、HTML変換は、それぞれ独立した**アダプター**（`pkg/adapters`）として実装されており、利用側はインターフェースを通じてこれらを操作します。
  * **プロンプト管理の集中:** プロンプトテンプレート (`prompt_release.md` など) の埋め込みとデータ注入は **`pkg/prompts`** パッケージに集約され、AIとの対話戦略を一元管理します。

### 3\. 🎨 HTML変換機能の提供

  * **HTMLレンダリング:** AIが生成したMarkdownテキストを受け取り、組み込みテンプレートとCSSを使用して、即座にブラウザで閲覧可能な**完全なHTMLドキュメント文字列**を返却する機能を提供します。

-----

## 📐 ライブラリ構成

このライブラリは、クリーンアーキテクチャに基づき、**コアロジックを再利用可能なコンポーネントとして提供**します。

```
gemini-reviewer-core
├── pkg
│   ├── adapters  # 外部システムへの接続層 (Port and Adapter パターン)
│   │   ├── gemini_adapter.go # Gemini API通信の実装
│   │   ├── git_service.go    # go-gitを使用したGit操作の実装
│   │   └── html_runner.go    # go-text-formatを使用したHTML変換の実装
│   │
│   └── prompts   # AIプロンプトのデータとロジック管理
│       ├── template_builder.go # テンプレートの選択とデータ注入ロジック
│       ├── template_data.go    # プロンプトの入力データ構造
│       └── template_manager.go # go:embed される生プロンプトファイルの管理
```

-----

### 📜 ライセンス (License)

このプロジェクトは [MIT License](https://opensource.org/licenses/MIT) の下で公開されています。
