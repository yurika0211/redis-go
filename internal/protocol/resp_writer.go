package protocol

import (
	"fmt"
	"io"
)

/**
 * Write a simple string to the writer
 * @param s string
 * @param w io.Writer 写入器
 * @return error
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

/**
 * WriteError writes a RESP Error to the writer: -ERR message\r\n
 */
func WriteError(w io.Writer, s string) error {
	_, err := w.Write([]byte("-" + s + "\r\n"))
	return err
}

/**
 * WriteArray writes a RESP Array to the writer: *len(arr)\r\n +arr[i]\r\n
 */
func WriteArray(w io.Writer, arr []string) error {
	_, err := w.Write([]byte("*" + fmt.Sprint(len(arr)) + "\r\n"))
	if err != nil {
		return err
	}
	for _, s := range arr {
		err := WriteBulkString(w, s)
		if err != nil {
			return err
		}
	}
	return nil
}

