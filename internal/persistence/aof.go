package persistence

import (
	"bufio"
	"fmt"
	"os"
)

type AOF struct {
	filename *os.File
	writer *bufio.Writer
}

type CommandExecuter func(cmd string)

var executor CommandExecuter

func RegisterExecutor(f CommandExecuter) {
	executor = f
}

/**
 * Open AOF file
 * @param filename(aof.txt)
 * @return
 */
func OpenAOF(filename string) (*AOF, error) {
	fmt.Println("Opening AOF file")
	f, err := os.OpenFile (
		filename,
		os.O_CREATE | os.O_APPEND | os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}
	return &AOF{
		filename: f,
		writer: bufio.NewWriter(f),
	},nil

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
	if executor != nil {
		executor(cmd[0])
	}
	return aof.writer.Flush()
}


