# 実装前テストケース・ゲート（test-case-gate）

このルールは **実装に着手する前**に必ず通過する品質ゲートである。
目的は「コードを書く前に、振る舞いの期待値を網羅的に言語化し、仕様の穴を潰しておく」こと。
テストケースの列挙と仕様ギャップの解消が終わるまで、**実装フェーズ（`/kiro-impl` 等）に進んではならない**。

このゲートには 2 つのモードがある。対象作業がどちらかを最初に判定すること。

- **新規モード**：新しい機能・新しいコンポーネントをゼロから作る。
- **追加開発モード（Brownfield）**：既存コードベースへの変更・拡張・修正。

---

## Human Decision hold

Human Decision が残っている場合、AI は SDD を次の範囲までに制限する。

- discovery
- requirements
- research
- option 比較
- risk 整理

次のフェーズには進まない。

- design
- tasks
- implementation

AI は推奨案を提示してよい。ただし、人間が採用 option を明示するまでは、それを採用案として扱わない。

再開条件:

- 人間が採用 option を明示する。
- 未決事項が `Human Decision Required` として requirements / research / final report に残っている。
- design / tasks は、採用済み option を input として再実行する。

Human Decision hold 中の final report には、少なくとも次を記録する。

```text
Human Decision Required: yes
Allowed phase: discovery / requirements / research
Stopped before: design / tasks / implementation
Options considered: ...
Recommended option by AI: ...
Adopted option by human: 未決
Resume condition: 人間が採用 option を明示すること
```

未決事項を AI が仮置きした場合も、その判断は未決のままとし、実装準備完了とは扱わない。

---

## Implementation Readiness Check

design / tasks から implementation へ進む前に、ローカル検証環境が最低限回るか確認する。

readiness は実装後 validation の代替ではない。実装課題と環境課題を分けるための事前確認である。

repo から検証コマンドを発見する。見る順序は README / AGENTS / package manifest / Makefile / CI / `.kiro/steering` とする。

確認対象は repo に存在するものだけにする。

- Ruby repo: `ruby -v`, `bundle -v`, `bundle check`, test command（例: `bundle exec rspec`）
- Node repo: `node -v`, `npm -v` / `pnpm -v` / `yarn -v`, package scripts（例: `test`, `lint`, `typecheck`）
- その他の repo: README / Makefile / CI にある検証コマンド

safe な no-op / dry-run / targeted command があれば最小実行する。
重すぎる command、破壊的 command、安全な no-op が無い command は無理に実行せず、`Not Run` と理由を書く。

final report / SDD artifacts には、少なくとも次を記録する。

```text
Implementation readiness:
- Runtime: Passed / Failed / Not Applicable / Blocked - <reason>
- Package manager: Passed / Failed / Not Applicable / Blocked - <reason>
- Dependency check: Passed / Failed / Not Run / Not Applicable / Blocked - <reason>
- Test command: Passed / Failed / Not Run / Blocked - <command and reason>
- Lint command: Passed / Failed / Not Run / Not Applicable / Blocked - <command and reason>
- Typecheck command: Passed / Failed / Not Run / Not Applicable / Blocked - <command and reason>
- Environment readiness blocker: none / <blocker>
```

実装前から失敗している検証は `environment readiness blocker` として扱い、実装チケットの失敗と混ぜない。
`Blocked` が残る場合は、implementation へ進む前に環境整備チケットまたは Human Decision を作る。
docs-only / 小規模変更では、非該当項目を `Not Applicable` としてよい。

---

## 共通：ゲート通過の判定基準（Definition of Ready）

以下をすべて満たすまで実装に進まない。

1. テストケースが下記カテゴリで列挙されている（該当しないカテゴリは「N/A＋理由」を明記）。
2. 各テストケースに「入力 / 事前条件 / 期待結果」が書かれている（実行可能なレベルの具体性）。
3. テストケース列挙の過程で見つかった**仕様ギャップ（曖昧・未定義・矛盾）がすべて解消**されている。
4. 解消の結果を、**どの文書に何を追記したか**まで反映済み（下記「追記マップ」）。
5. Human Decision Required が残っていない。残っている場合は design / tasks / implementation へ進まない。
6. 実装前 readiness check の結果が final report / SDD artifacts に記録できる状態になっている。`Blocked` が残る場合は implementation へ進まない。

---

## 新規モード

### Step 1. テストケースを列挙する

最低限、以下の 4 区分で網羅する。各区分が空なら「無し＋理由」を書く（空欄禁止）。

- **正常系（Happy path）**：代表的な正しい入力で期待どおり動くこと。
- **異常系（Error path）**：不正入力・権限不足・依存先エラー・タイムアウト等で、定義どおりに失敗/回復すること。
- **境界値（Boundary）**：0 / 1 / 最大 / 最小 / 空 / null / 上限+1 / 下限-1、長さ・件数・桁・日付境界など。
- **エッジ（Edge）**：同時実行・重複・順序入れ替え・冪等性・部分失敗・リトライ・文字コード/ロケール・大量データ等。

各ケースは次の形式で書く：

```
- [区分] <ケース名>
  入力 / 事前条件: ...
  期待結果: ...
```

### Step 2. 仕様ギャップを検出する

列挙の途中で「期待結果が書けない」「複数の解釈がありうる」ケースは**仕様ギャップ**である。
すべて洗い出し、次のいずれかで解消する：

- 仕様（`requirements.md` / `design.md`）に決定を追記して曖昧さを消す。
- 決められない場合は質問として明示し、回答を得てから確定する（推測で埋めない）。

### Step 3. 仕様へ反映（追記マップ）

解消した内容を文書へ反映する。目安：

| 見つかった内容 | 追記先 |
|---|---|
| 機能の振る舞い・受け入れ条件 | `.kiro/specs/<feature>/requirements.md` |
| 設計上の決定・制約・I/F | `.kiro/specs/<feature>/design.md` |
| テストケース一覧そのもの | `.kiro/specs/<feature>/requirements.md`（受け入れ条件として）または tasks のテスト項目 |
| プロジェクト横断の規約・前提 | `.kiro/steering/conventions.md` |

→ Step 1〜3 がすべて埋まったら**ゲート通過**。実装へ進んでよい。

---

## 追加開発モード（Brownfield）

既存への変更は「既存の振る舞いを壊さないこと」が最重要。新規モードに**以下を追加**で行う。

### Step 0. キャッチアップ（実装前）

- `.kiro/steering/` を読んで、対象領域の現状（Product / Tech / Structure・社内規約）を把握する。
- 既存仕様・関連コード・既存テストを読み、**現状の振る舞い**を確認してから変更範囲を確定する。
- steering / 仕様と実コードの**乖離（validate-gap）**を洗い出す。乖離があれば、それ自体を仕様ギャップとして扱う。

### Step 1+. 追加で列挙すべきテストケース

新規モードの 4 区分に加えて：

- **回帰系（Regression）**：今回触る箇所の**既存の正しい振る舞いが維持される**ことを確認するケース。
  変更前から通っていた経路を明示し、変更後も期待結果が変わらないことを担保する。
- **連携境界系（Integration boundary）**：変更対象が他モジュール・外部サービス・DB・API と接する境界。
  後方互換（既存呼び出し元への影響）、契約（スキーマ/型/エラー）、データ移行・互換の有無を網羅する。

### Step 2+. 仕様ギャップ＝乖離の解消

Step 0 で見つけた「仕様と実装の乖離」を、どちらが正かを判断して解消する：

- 実装が正 → 仕様（steering / specs）を**実態に合わせて更新**。
- 仕様が正 → 変更タスクとして扱い、テストケースに反映。
- 判断できない → 質問として明示し、回答を得る。

### Step 3+. 追記マップ（Brownfield 版）

| 見つかった内容 | 追記先 |
|---|---|
| 既存振る舞いの実態（ドキュメント化されていなかったもの） | `.kiro/steering/`（Tech / Structure） |
| 後方互換・移行方針 | `.kiro/specs/<feature>/design.md` |
| 回帰で守るべき経路・連携契約 | `.kiro/specs/<feature>/requirements.md`（受け入れ条件） |
| 横断規約の補強 | `.kiro/steering/conventions.md` |

→ Step 0〜3+ がすべて埋まったら**ゲート通過**。

---

## やりすぎ防止

自明・低リスクな小変更（タイポ修正・コメント・設定値の軽微変更など）に本ゲートをフル適用しない。
その場合は「仕様なし直接実装」に逃がしてよい。ただし**その判断自体を一行残す**こと
（例：「低リスク変更のため test-case-gate を簡略化、回帰のみ確認」）。

---

## Commit ガード（自動 commit の禁止）

cc-sdd の `kiro-impl` は既定で、サブタスク完了ごとに selective git commit を行う
（`dispatch → review → commit → next`、手順 `e) Commit`）。
**この repo ではその自動 commit を行わない。**

実装フェーズ（`/kiro-impl` 等）では次を守る。

- サブタスク完了時に行うのは、コード変更・`tasks.md` のチェック更新・`git diff` の提示までとする。
- `git add` / `git commit` / `git push` は実行しない。ステージングもしない。
- commit は、ユーザーが「commit して」等と明示指示した時だけ実行する。指示があるまで作業ツリーは未 commit のまま残す。
- `kiro-impl` の手順 `e) Commit`（parent-only selective staging）は、本ルールにより無効とする。

最終報告には、少なくとも次を残す。

```text
Auto-commit: disabled (vt overlay rule)
Uncommitted changes: <files or none>
Commit performed: no（ユーザー明示指示時のみ yes）
```
