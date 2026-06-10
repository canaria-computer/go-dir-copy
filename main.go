package main

import (
	"errors"
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
	concurrency := flag.Int("n", 16, "concurrency number")
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "Usage: copydir [-n concurrency] <source> <destination>")
		os.Exit(1)
	}
	start := time.Now()
	if err := copyDirSyncPall(filepath.Clean(flag.Arg(0)), filepath.Clean(flag.Arg(1)), *concurrency); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	fmt.Printf("Process time: %s\n", time.Since(start))
}

func copyDirSyncPall(srcDir, dstDir string, concurrency int) error {
	srcDir = filepath.ToSlash(srcDir)
	var tasks []string

	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {

			return nil
		}
		tasks = append(tasks, filepath.ToSlash(path))
		return nil
	})
	if err != nil {
		return err
	}

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
		sem  = make(chan struct{}, concurrency)

		bufPool = sync.Pool{
			New: func() any {
				return make([]byte, 1*1024*1024)
			},
		}
	)

	wg.Add(len(tasks))
	for _, task := range tasks {
		dstFile := filepath.Join(dstDir, strings.TrimPrefix(task, srcDir))
		srcAbsFile, err := filepath.Abs(task)
		if err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("abs path %s: %w", task, err))
			mu.Unlock()
			wg.Done()
			continue
		}

		sem <- struct{}{}
		fmt.Println("Copying...", srcAbsFile)
		go func(src, dst string) {
			defer func() {
				wg.Done()
				<-sem
			}()

			b := bufPool.Get().([]byte)
			defer bufPool.Put(b)

			if err := copyFileWithTimeStamp(src, dst, b); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(srcAbsFile, dstFile)
	}
	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("%d error(s) during copy:\n%w", len(errs), errors.Join(errs...))
	}
	return nil
}

func copyFileWithTimeStamp(srcFile, dstFile string, buf []byte) error {
	srcStat, err := os.Stat(srcFile)
	if err != nil {
		return err
	}

	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	if err := os.MkdirAll(filepath.Dir(dstFile), 0755); err != nil {
		return err
	}

	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.CopyBuffer(dst, src, buf); err != nil {
		return err
	}

	return os.Chtimes(dstFile, srcStat.ModTime(), srcStat.ModTime())
}
