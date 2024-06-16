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

func main() {
	flag.Parse()
	s := time.Now()
	copyDirSyncPall(filepath.Clean(flag.Arg(0)), filepath.Clean(flag.Arg(1)))
	fmt.Printf("process time: %s\n", time.Since(s))
}

func copyDirSyncPall(srcDir, dstDir string) error {
	srcDir = filepath.ToSlash(srcDir)
	task := []string{}
	err := filepath.Walk(srcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		task = append(task, path)
		return nil
	})

	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(task))
	for _, v := range task {
		v = filepath.ToSlash(v)
		dstFile := filepath.Join(dstDir, strings.Replace(v, srcDir, "", 1))
		srcAbsFile, _ := filepath.Abs(v)
		fmt.Println("Copying...", srcAbsFile)
		go func(s, d string) {
			err := copyFileWithTimeStamp(s, d)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(srcAbsFile, dstFile)
	}
	wg.Wait()
	return nil
}

func copyFileWithTimeStamp(srcFile, dstFile string) error {
	// 作成日時データ読み取り
	srcStat, err := os.Stat(srcFile)
	// ファイルを読み取る
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()
	// コピー先ディレクトリを作成
	err = os.MkdirAll(filepath.Dir(dstFile), 0755)
	if err != nil {

	}
	// コピー先ファイルを作成
	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	// ファイル中身を同期する
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	// 作成日時などを書き換える
	err = os.Chtimes(dstFile, srcStat.ModTime(), srcStat.ModTime())
	if err != nil {
		return err
	}
	return nil
}
