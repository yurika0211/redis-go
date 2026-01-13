package persistence

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type AOF struct {
	file     *os.File      // 实际的文件句柄
	writer   *bufio.Writer // 缓冲区
	filename string        // 文件路径字符串
	strategy FsyncStrategy // 刷盘策略
}

type CommandExecuter func(cmd string)

var executor CommandExecuter

type FsyncStrategy interface {
	OnAppend(file *os.File)
}

// 1. Always: 每次都强制刷盘
type AlwaysStrategy struct{}

func (s *AlwaysStrategy) OnAppend(file *os.File) {
	file.Sync() // 调用操作系统的 fsync
}

// 2. No: 啥都不干，等操作系统自己处理
type NoStrategy struct{}

func (s *NoStrategy) OnAppend(file *os.File) {}

// 3. Everysec: 后台协程每秒刷一次
type EverysecStrategy struct {
	lastSync time.Time
}

func (s *EverysecStrategy) OnAppend(file *os.File) {
	if time.Since(s.lastSync) > time.Second {
		file.Sync()
		s.lastSync = time.Now()
	}
}
func RegisterExecutor(f CommandExecuter) {
	executor = f
}

/**
 * Create AOF file
 * @param filepath(aof.txt)
 * @return
 */
func NewAOF(filename string) (*AOF, error) {
	//在目标地址创建aof.txt文件
	f, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}
	return &AOF{
		file:     f,                  // 句柄给 file
		filename: filename,           // 路径给 filename (string)
		writer:   bufio.NewWriter(f), // 缓冲区基于文件句柄
		strategy: &NoStrategy{},
	}, nil
}

/**
 * Open AOF file
 * @param filename(aof.txt)
 * @return
 */
func OpenAOF(filename string) (*AOF, error) {
	fmt.Println("Opening AOF file")
	f, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}
	return &AOF{
		file: f,
		filename: filename,
		writer:   bufio.NewWriter(f),
		strategy: &NoStrategy{},
	}, nil

}

/**
 * Append command to AOF file
 * @param cmd
 * @return
 */
func (aof *AOF) Append(cmd []string) error {
	fmt.Println("Appending to AOF file")
	fmt.Fprintf(aof.writer, "*%d\r\n", len(cmd))
	for _, arg := range cmd {
		fmt.Fprintf(aof.writer, "$%d\r\n", len(arg))
		fmt.Fprintf(aof.writer, "%s\r\n", arg)
	}
	aof.writer.Flush()
	aof.strategy.OnAppend(aof.file)
	return nil
}
