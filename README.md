# 🤖 Gemini Reviewer Core

[![Language](https://img.shields.io/badge/Language-Go-blue)](https://golang.org/)
[![Go Version](https://img.shields.io/github/go-mod/go-version/shouni/gemini-reviewer-core)](https://golang.org/)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/shouni/gemini-reviewer-core)](https://github.com/shouni/gemini-reviewer-core/tags)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Gemini Reviewer Core** は、Google Gemini API を活用し、Gitリポジトリのブランチ間の差分を分析してAIコードレビューを自動生成するための**コアライブラリ**です。クリーンなアーキテクチャ設計により、Git操作、AI通信、プロンプト管理といった外部依存性を完全に分離しています。

## ✨ 技術スタック (Technology Stack)

| 要素 | 技術 / ライブラリ | 役割 |
| :--- | :--- | :--- |
| **言語** | **Go (Golang)** | ツールの開発言語。クロスプラットフォームでの高速な実行を実現します。 |
| **Git 操作** | **go-git** | クローン、フェッチ、**3-dot diff** (共通祖先からの差分) の取得まですべてを Go のコード内で完結させ、**`ssh-key-path`に基づくSSH認証とホストキー検証の設定**を統合しました。 |
| **I/O 連携** | **`github.com/shouni/go-remote-io`** | GCSとローカルファイルシステムへのI/O操作を抽象化し、**GCSへのレビュー結果保存**を実現します。 |
| **Markdown to HTML** | **`github.com/shouni/go-text-format`** | AIが出力したMarkdown形式のレビュー結果を、スタイル付きの完全なHTMLドキュメントに**変換・レンダリング**するために使用します。|
| **AI モデル** | **Google Gemini API** | 取得したコード差分を分析し、レビューコメントを生成するために使用します。**（温度設定による応答制御を適用済み）** |

-----

## ✨ 主要な機能と特徴

### 1\. 🔍 高度なGit差分分析

* **SSHネイティブ対応:** `go-git` を使用し、外部のSSHコマンドに依存せず、**`ssh-key-path`** 指定によるセキュアなSSH認証とホストキー検証をサポートします。
* **正確な差分取得:** ベースブランチとフィーチャーブランチ間のマージベース（共通祖先）を基準とした **3-dot diff (`A...B`)** を取得し、正確でクリーンな差分のみをAIに提供します。

### 2\. 🧱 責務の厳密な分離

* **アダプターパターンの採用:** Git操作、Gemini APIへの通信は、それぞれ独立した**アダプター**（`pkg/adapters`）として実装されており、コアロジックはこれらのインターフェース（ポート）に依存します。
* **プロンプト管理の集中:** プロンプトテンプレート (`prompt_release.md` など) の埋め込みとデータ注入は **`pkg/prompts`** パッケージに集約され、AIとの対話戦略を一元管理します。

### 3\. 💾 柔軟なI/Oと出力

* **リモートI/Oサポート:** `go-remote-io` を活用し、ローカルファイルへの出力に加え、**Google Cloud Storage (GCS) へのレビュー結果の保存**に対応します。
* **HTMLレンダリング:** AIのMarkdown出力を、組み込みテンプレートとCSSを使用した完全なHTMLドキュメントに変換し、視認性の高いレビュー結果を提供します。

-----

## 📐 ライブラリ構成

このライブラリは、クリーンアーキテクチャに基づき、**コアロジックを外部のインフラストラクチャから分離**しています。特に、**アダプター**層 (`pkg/adapters`) と**プロンプト**層 (`pkg/prompts`) が、外部依存の責務を負います。

```
gemini-reviewer-core
├── pkg
│   ├── adapters  # 外部システムへの接続層 (Port and Adapter パターン)
│   │   ├── gemini_adapter.go # Gemini API通信の実装
│   │   ├── git_service.go    # go-gitを使用したGit操作の実装
│   │   └── html_runner.go    # go-text-formatを使用したHTML変換の実装
│   └── prompts   # AIプロンプトのデータとロジック管理
│       ├── template_builder.go # テンプレートの選択とデータ注入ロジック
│       ├── template_data.go    # プロンプトの入力データ構造
│       └── prompt_*.md         # go:embed される生プロンプトファイル
```

-----

### 📜 ライセンス (License)

このプロジェクトは [MIT License](https://opensource.org/licenses/MIT) の下で公開されています。
