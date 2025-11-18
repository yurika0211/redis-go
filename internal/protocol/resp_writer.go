package protocol

import (
	"fmt"
	"io"
)

/**
 * Write a simple string to the writer
 * @param s string
 * @param w io.Writer 写入器
 */
func WriteSimpleString(w io.Writer, s string) error {
	_, err := w.Write([]byte("+" + s + "\r\n"))
	return err
}

/**                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               
 * Write a bulk string to the writer
 * @param s string
 * @param w io.Writer
 */
func WriteBulkString(w io.Writer, s string) error {
	_, err := w.Write([]byte("$" + fmt.Sprint(len(s)) + "\r\n" + s + "\r\n"))
	return err
}
