# Go 言語でファイルをコピー

Go 言語で並行処理することでファイルを早くコピーすることを目的にして開発されたコピープログラム

## 実験結果

JavaScript/TypeScript フレームワークの Qwik の `node_module` フォルダをSSDからHDDにコピーしたときの実行結果

| Method     | ERT (ms) |
| ---------- | -------- |
| Go         | 32165    |
| XCOPY      | 68935    |
| PowerShell | 72919    |
| ROBOCOPY   | 96275    |

![実験結果をグラフにしたもの。視覚的に速くなったと主張している。](./ERT.gif)
