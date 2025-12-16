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

/**
 * Open AOF file
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
func (aof *AOF) Append(cmd []string) error {
	fmt.Println("Appending to AOF file")
	fmt.Fprintf(aof.writer, "*%d\r\n", len(cmd))
	for _, arg := range cmd {
		fmt.Fprintf(aof.writer, "$%d\r\n", len(arg))
		fmt.Fprintf(aof.writer, "%s\r\n", arg)
	}
	return aof.writer.Flush()
}

