package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const concurrency = 3

func main() {
	// コマンドライン引数をパース
	flag.Parse()
	// プロセス開始時間を記録
	start := time.Now()
	// ディレクトリのコピーを開始
	if err := copyDirSyncPall(filepath.Clean(flag.Arg(0)), filepath.Clean(flag.Arg(1))); err != nil {
		fmt.Println("Error:", err)
	}
	// 経過時間を表示
	fmt.Printf("Process time: %s\n", time.Since(start))
}

// ディレクトリを同期的かつ並行にコピーする
func copyDirSyncPall(srcDir, dstDir string) error {
	// パスを標準化
	srcDir = filepath.ToSlash(srcDir)
	var tasks []string

	// コピー元ディレクトリを再帰的に探索し、全ファイルパスを取得
	err := filepath.Walk(srcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// ディレクトリはスキップ
		if info.IsDir() {
			return nil
		}
		tasks = append(tasks, path)
		return nil
	})

	if err != nil {
		return err
	}

	// 同期処理のためのWaitGroupを設定
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	// 同時実行制御用チャネル
	sem := make(chan struct{}, concurrency)

	// 各ファイルを並行でコピー
	for _, task := range tasks {
		sem <- struct{}{}
		task = filepath.ToSlash(task)
		// コピー先のファイルパスを設定
		dstFile := filepath.Join(dstDir, strings.Replace(task, srcDir, "", 1))
		srcAbsFile, _ := filepath.Abs(task)
		fmt.Println("Copying...", srcAbsFile)
		go func(src, dst string) {
			if err := copyFileWithTimeStamp(src, dst); err != nil {
				fmt.Println(err)
			}
			wg.Done()
			<-sem
		}(srcAbsFile, dstFile)
	}
	// すべてのコピーが完了するまで待機
	wg.Wait()
	return nil
}

// ファイルをコピーし、タイムスタンプを維持する
func copyFileWithTimeStamp(srcFile, dstFile string) error {
	// 元ファイルの情報を取得
	srcStat, err := os.Stat(srcFile)
	if err != nil {
		return err
	}

	// 元ファイルを開く
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	// コピー先ディレクトリを作成
	if err := os.MkdirAll(filepath.Dir(dstFile), 0755); err != nil {
		return err
	}

	// コピー先ファイルを作成
	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	// ファイル内容をコピー
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	// コピー先ファイルのタイムスタンプを元ファイルと同じに設定
	if err := os.Chtimes(dstFile, srcStat.ModTime(), srcStat.ModTime()); err != nil {
		return err
	}
	return nil
}
