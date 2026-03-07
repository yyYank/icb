# icb - Isolated Clipboard for CLI 設計ドキュメント

## 概要

SSH接続先サーバー内で完結するクリップボード履歴ツール。
パイプでテキストを蓄積し、fzf風のTUIで履歴から選択して標準出力する。

## ユースケース

```bash
# 蓄積（パイプで渡すだけ）
echo "some text" | cb
cat main.go | cb

# 選択して出力
cb
cb | bash
cb > out.txt
```

## アーキテクチャ

```
icb/
├── main.go
├── cmd/
│   └── root.go       # cobra root（パイプ判定→蓄積 or TUI起動）
├── store/
│   └── store.go      # 履歴の読み書き
└── tui/
    └── tui.go        # bubbletea TUI
```

## 技術スタック

| 役割 | ライブラリ |
|---|---|
| CLIフレームワーク | `cobra` |
| TUIフレームワーク | `bubbletea` |
| TUIコンポーネント | `bubbles` (textinput, list) |
| ストレージ | JSON Lines |

## ストレージ仕様

- パス: `~/.cb_history`
- フォーマット: JSON Lines（1エントリ1行）

```json
{"id":"uuid","content":"echo hello","created_at":"2026-03-07T12:00:00Z"}
{"id":"uuid","content":"cat main.go | cb","created_at":"2026-03-07T12:01:00Z"}
```

- 複数行コンテンツは`\n`エスケープして1行に収める
- 上限: 1000件（超えたら古いものから削除）

## TUI仕様

fzfライクなインクリメンタルサーチUI。

```
> search query...
──────────────────────────────
  echo hello world
  cat main.go | cb
▶ ssh -i ~/.ssh/key user@host
  SELECT * FROM users LIMIT 10
──────────────────────────────
4/100  ↑↓で移動  Enterで選択  Ctrl+Cでキャンセル
```

### キーバインド

| キー | 動作 |
|---|---|
| `↑` / `↓` | 項目移動 |
| 文字入力 | インクリメンタルサーチ |
| `Enter` | 選択して標準出力 |
| `Ctrl+C` / `Esc` | キャンセル |

## コマンド仕様

### `icb` — パイプあり → 蓄積
標準入力を受け取って履歴に追記して終了。

```bash
echo "text" | cb
cat file.txt | cb
```

### `icb` — パイプなし → TUI起動
選択したエントリを標準出力して終了。

```bash
cb
cb | bash
cb > out.txt
```

### 判定ロジック
```go
stat, _ := os.Stdin.Stat()
if (stat.Mode() & os.ModeCharDevice) == 0 {
    // パイプ → 蓄積
} else {
    // TTY → TUI起動
}
```

## 将来的な拡張候補（v1では対象外）

- `cb list` — 一覧をプレーンテキストで出力
- `cb clear` — 履歴削除
- タグ・ラベル機能
- プレビューペイン（複数行の中身を右ペインに表示）
- 設定ファイル（上限件数、ストレージパスなど）

