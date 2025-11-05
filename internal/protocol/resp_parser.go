package protocol

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)


/**
 * Parse an array from a reader.
 * @param r *bufio.Reader 读取器
 * @return []string 解析出的数组
 * @return error 错误信息
 */
func ParseArray(r *bufio.Reader) ([]string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(line, "*") {
		return nil, fmt.Errorf("expected *, got %q", line)
	}

	countStr := strings.TrimSuffix(line[1:], "\r\n")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, count)
	for i := 0; i < count; i++ {
		lengthLine, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		
		if !strings.HasPrefix(lengthLine, "$") {
			return nil, fmt.Errorf("expected $, got %q", lengthLine)
		}

		dataLine, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		// Trim trailing CRLF from data line
		result = append(result, strings.TrimSuffix(dataLine, "\r\n"))
	}
	return result, nil
}
