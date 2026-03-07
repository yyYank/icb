
# icb — Internal/Isolated Clipboard

[English](#english) / [日本語](#日本語)

---

<a name="english"></a>

## English

A standalone clipboard history tool for terminal environments.
Works entirely within your shell — no OS clipboard, no GUI, no dependencies.
Designed for use over SSH.

### Concept

When you SSH into a remote server, the system clipboard doesn't follow you.
`icb` gives you a persistent clipboard history that lives inside the terminal session itself.

```bash
echo "some text" | icb     # store
icb                        # pick from history → stdout
```

### Install

```bash
go install github.com/yyYank/icb@latest
```

Or download a binary from the [releases page](https://github.com/yyYank/icb/releases).

### Shell Integration

Add one line to your `~/.zshrc` or `~/.bashrc` to enable the `Ctrl+X I` keybinding.

**zsh**

```bash
echo 'eval "$(icb init)"' >> ~/.zshrc
```

**bash**

```bash
echo 'eval "$(icb init)"' >> ~/.bashrc
```

Then reload your shell:

```bash
source ~/.zshrc   # or ~/.bashrc
```

Now you can press `Ctrl+X I` at any point while typing a command to open the TUI and insert the selected entry at the cursor:

```
$ make OUT=<Ctrl+X I>  →  TUI opens  →  select entry  →  make OUT=some text
```

### Usage

**Store** — pipe anything into `icb` to save it to history.

```bash
echo "hello world" | icb
cat main.go | icb
curl https://example.com/script.sh | icb
```

**Pick & Paste** — run `icb` without arguments to open the TUI.

```bash
icb               # browse and select
icb | bash        # select and execute
icb > out.txt     # select and save to file
```

**TUI**

```
> search query...
──────────────────────────────────────
  echo hello world
  cat main.go | icb
▶ ssh -i ~/.ssh/key user@host
  SELECT * FROM users WHERE id = 1
──────────────────────────────────────
4/100  ↑↓ to move  Enter to select  Ctrl+C to cancel
```

| Key | Action |
|---|---|
| `↑` / `↓` | Move cursor |
| Type anything | Incremental search |
| `Enter` | Select → stdout |
| `Ctrl+C` / `Esc` | Cancel |
| `Ctrl+X I` | Insert at cursor (requires shell integration) |

### How it works

`icb` detects whether stdin is a pipe or a TTY at runtime.

- **Pipe** → reads stdin and appends to history
- **TTY** → opens the TUI to browse and select

History is stored in `~/.icb_history` as JSON Lines. Up to 1000 entries are kept; older entries are pruned automatically.

### Built with

- [cobra](https://github.com/spf13/cobra) — CLI framework
- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework

### License

MIT

---

<a name="日本語"></a>

## 日本語

ターミナル環境で完結するクリップボード履歴ツール。
OSのクリップボード、GUI、外部依存なし。SSH接続先でもそのまま使える。

### コンセプト

SSHでリモートサーバーに入ると、手元のクリップボードは使えない。
`icb` はシェル環境の中だけで動く、独立したクリップボード履歴を提供する。

```bash
echo "some text" | icb     # 蓄積
icb                        # 履歴から選択 → 標準出力
```

### インストール

```bash
go install github.com/yyYank/icb@latest
```

または[リリースページ](https://github.com/yyYank/icb/releases)からバイナリをダウンロード。

### シェルインテグレーション

`~/.zshrc` または `~/.bashrc` に1行追加するだけで `Ctrl+X I` キーバインドが使えるようになる。

**zsh**

```bash
echo 'eval "$(icb init)"' >> ~/.zshrc
```

**bash**

```bash
echo 'eval "$(icb init)"' >> ~/.bashrc
```

シェルを再読み込みする：

```bash
source ~/.zshrc   # または ~/.bashrc
```

コマンド入力中に `Ctrl+X I` を押すとTUIが開き、選んだ内容がカーソル位置に挿入される：

```
$ make OUT=<Ctrl+X I>  →  TUI起動  →  選択  →  make OUT=some text
```

### 使い方

**蓄積する** — パイプで渡すだけで履歴に保存される。

```bash
echo "hello world" | icb
cat main.go | icb
curl https://example.com/script.sh | icb
```

**選択して使う** — 引数なしで起動するとTUIが開く。

```bash
icb               # 履歴から選択
icb | bash        # 選択してそのまま実行
icb > out.txt     # 選択してファイルに保存
```

**TUI**

```
> 検索ワード...
──────────────────────────────────────
  echo hello world
  cat main.go | icb
▶ ssh -i ~/.ssh/key user@host
  SELECT * FROM users WHERE id = 1
──────────────────────────────────────
4/100  ↑↓で移動  Enterで選択  Ctrl+Cでキャンセル
```

| キー | 動作 |
|---|---|
| `↑` / `↓` | カーソル移動 |
| 文字入力 | インクリメンタルサーチ |
| `Enter` | 選択 → 標準出力 |
| `Ctrl+C` / `Esc` | キャンセル |
| `Ctrl+X I` | カーソル位置に挿入（シェルインテグレーション必要） |

### 仕組み

起動時に標準入力がパイプかTTYかを自動判定する。

- **パイプあり** → 標準入力を読み取って履歴に追記
- **パイプなし** → TUIを起動して履歴を閲覧・選択

履歴は `~/.icb_history` にJSON Lines形式で保存される。上限は1000件で、超えた分は古いものから自動削除。

### 使用ライブラリ

- [cobra](https://github.com/spf13/cobra) — CLIフレームワーク
- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUIフレームワーク

### ライセンス

MIT
