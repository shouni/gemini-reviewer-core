# 🤖 Gemini Reviewer Core

[![Language](https://img.shields.io/badge/Language-Go-blue)](https://golang.org/)
[![Go Version](https://img.shields.io/github/go-mod/go-version/shouni/gemini-reviewer-core)](https://golang.org/)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/shouni/gemini-reviewer-core)](https://github.com/shouni/gemini-reviewer-core/tags)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Gemini Reviewer Core** は、Google Gemini API を活用し、Gitリポジトリのブランチ間の差分を分析してAIコードレビューを自動生成するための**コアライブラリ**です。

本ライブラリは、**CLIツール**や**Webアプリケーション**の共通基盤として設計されています。Git操作、AI通信、プロンプト生成といった「ビジネスロジック」に加え、HTML変換や結果保存を行う公開層（Publisher）を提供します。これにより、利用側はインフラの実装詳細を意識することなく、レビュー結果の生成から公開までを一貫して行えます。

---

## 🎯 本ライブラリの価値 (Value Proposition)

本ライブラリは、AIレビュー機能を単なるスクリプトとしてではなく、**大規模なアプリケーションにも組み込める持続可能な基盤**として提供します。

* **⚡ 開発効率の最大化:** Git操作、SSH認証、AIとの対話、結果の整形・公開といった複雑なプロセスをすべて抽象化し、利用者はコアなビジネスロジックの記述に集中できます。
* **📐 アーキテクチャの堅牢性:** 厳格なクリーンアーキテクチャと**依存性逆転の原則 (DIP)** に基づいて構築されているため、テストが容易で、クラウドストレージの変更やAIモデルの切り替えといった**拡張に対する耐性**を持っています。
* **🌐 マルチクラウド対応:** **GCS/S3**への公開を抽象化されたインターフェース (`Publisher`) で一元管理できるため、利用環境を選びません。

---

## ✨ 技術スタック (Technology Stack)

| 要素 | 技術 / ライブラリ | 役割 |
| :--- | :--- | :--- |
| **言語** | **Go (Golang)** | ライブラリの開発言語。 |
| **Git 操作** | **go-git** (`adapters`層) | クローン、フェッチ、**3-dot diff** (共通祖先からの差分) の取得まですべてを Go のコード内で完結させ、**SSH認証とホストキー検証の設定**を統合しました。 |
| **AI モデル** | **Google Gemini API** (`adapters`層) | 取得したコード差分を分析し、レビューコメントを生成するために使用します。 |
| **Markdown to HTML** | **`go-text-format`** (`publisher`層) | AIが出力したMarkdown形式のレビュー結果を、スタイル付きの完全なHTMLドキュメントに**変換**するために使用します。|
| **ストレージ操作** | **`go-remote-io`** (`publisher`層) | **GCS/S3** 等のクラウドストレージへのアップロード処理を抽象化し、CLIとWebアプリで接続処理（Factory）を共通化するために使用します。 |
| **プロンプト管理** | **`text/template` + `embed`** | レビューモード（Release/Detail）に応じたプロンプトテンプレートをバイナリに埋め込み、動的に生成します。 |

---

## ✨ 主要な機能と特徴

### 1. 🔍 高度なGit差分分析

* **SSHネイティブ対応:** `go-git` を使用し、外部のSSHコマンドに依存せず、認証情報を注入可能な設計になっています。
* **正確な差分取得:** ベースブランチとフィーチャーブランチ間のマージベース（共通祖先）を基準とした **3-dot diff (`A...B`)** を取得し、正確でクリーンな差分のみをAIに提供します。

### 2. 🧱 責務の厳密な分離と抽象化

* **アダプターパターンの採用:** Git操作とGemini APIへの通信は、それぞれ独立した**アダプター**（`pkg/adapters`）として実装されており、外部システムへの依存をカプセル化しています。
* **依存性逆転の原則 (DIP):** すべての外部連携（Git, AI, Storage）はインターフェースに依存しており、**ストラテジーパターン**を利用することで高い拡張性とテスト容易性を実現しています。
* **プロンプト管理の集中:** プロンプトテンプレート (`prompt_release.md` など) の埋め込みとデータ注入は **`pkg/prompts`** パッケージに集約され、AIとの対話戦略を一元管理します。

### 3. 🎨 HTML変換機能の提供

* **高品質なHTMLレンダリング:** AIが生成したMarkdownテキストを受け取り、組み込みテンプレートとCSSを使用して、即座にブラウザで閲覧可能な**スタイリング済みの完全なHTMLドキュメント文字列**を返却する機能を提供します。

### 4. 📤 結果の公開・保存 (Publisher)

* **クラウドベンダーの統合:** GCSと**Amazon S3**への公開をサポート。抽象的な **`Publisher` インターフェース**により、ロジックを共通化しています。
* **変換と保存の統合:** MarkdownからHTMLへの変換ロジックと、その結果をクラウドストレージへ保存する処理を **Publisher** 層として提供しています。
* **依存性の注入 (DI):** ストレージ接続のファクトリー (`go-remote-io`) を外部から注入できる設計により、CLI（都度接続）と Web Worker（コネクション再利用）の両方で効率的に動作します。

---

## 📐 ライブラリ構成

このライブラリは、クリーンアーキテクチャに基づき、**コアロジックを再利用可能なコンポーネントとして提供**します。

```text
gemini-reviewer-core
├── pkg
│   ├── adapters  \# 外部システムへの接続層 (Port and Adapter パターン)
│   │   ├── gemini\_adapter.go \# Gemini API通信の実装
│   │   └── git\_service.go    \# go-gitを使用したGit操作の実装
│   │
│   ├── publisher \# 結果の変換・出力・公開層
│   │   ├── publisher.go      \# Publisherインターフェースとデータ構造 (コア抽象化)
│   │   ├── publisher\_factory.go \# URIスキームに基づくPublisherのファクトリ
│   │   ├── gcs\_publisher.go  \# GCSへの公開実装
│   │   ├── s3\_publisher.go   \# S3への公開実装
│   │   ├── html\_converter.go \# GCS/S3で共通利用するHTML変換ロジック
│   │   └── md\_adapter.go     \# 外部Markdownライブラリのアダプター
│   │
│   └── prompts   \# AIプロンプトのデータとロジック管理
│       ├── template\_builder.go \# テンプレートの選択とデータ注入ロジック
│       ├── template\_data.go    \# プロンプトの入力データ構造
│       └── template\_manager.go \# go:embed される生プロンプトファイルの管理
````

-----

### 📜 ライセンス (License)

このプロジェクトは [MIT License](https://opensource.org/licenses/MIT) の下で公開されています。
