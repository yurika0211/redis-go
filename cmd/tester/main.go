package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func writeArray(conn net.Conn, arr []string) error {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("*%d\r\n", len(arr)))
	for _, s := range arr {
		b.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
	}
	_, err := conn.Write([]byte(b.String()))
	return err
}

// readReply reads a single RESP value and returns a human-friendly string
func readReply(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimRight(line, "\r\n")
	if len(line) == 0 {
		return "", nil
	}
	switch line[0] {
	case '+':
		return line[1:], nil
	case '-':
		return "ERR: " + line[1:], nil
	case '$':
		n, _ := strconv.Atoi(strings.TrimSpace(line[1:]))
		if n == -1 {
			return "", nil
		}
		buf := make([]byte, n+2)
		_, err := r.Read(buf)
		if err != nil {
			return "", err
		}
		return string(buf[:n]), nil
	case '*':
		cnt, _ := strconv.Atoi(strings.TrimSpace(line[1:]))
		parts := make([]string, 0, cnt)
		for i := 0; i < cnt; i++ {
			v, err := readReply(r)
			if err != nil {
				return "", err
			}
			parts = append(parts, v)
		}
		return strings.Join(parts, ","), nil
	default:
		return line, nil
	}
}

func expectContains(got, want string) bool {
	return strings.Contains(got, want)
}

func runTests(addr string) int {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		fmt.Printf("connect failed: %v\n", err)
		return 1
	}
	defer conn.Close()
	r := bufio.NewReader(conn)

	tests := 0
	fails := 0

	// PING
	tests++
	writeArray(conn, []string{"PING"})
	if v, err := readReply(r); err != nil || !(expectContains(v, "PONG")) {
		fmt.Printf("PING fail: got=%q err=%v\n", v, err)
		fails++
	} else {
		fmt.Println("PING ok")
	}

	// SET/GET
	tests++
	writeArray(conn, []string{"SET", "tkey", "tval"})
	if v, _ := readReply(r); !expectContains(v, "OK") {
		fmt.Printf("SET fail: %q\n", v)
		fails++
	} else {
		fmt.Println("SET ok")
	}
	tests++
	writeArray(conn, []string{"GET", "tkey"})
	if v, _ := readReply(r); v != "tval" {
		fmt.Printf("GET fail: %q\n", v)
		fails++
	} else {
		fmt.Println("GET ok")
	}

	// HMSET/HGET
	tests++
	writeArray(conn, []string{"HMSET", "hkey", "f1", "v1", "f2", "v2"})
	if v, _ := readReply(r); !expectContains(v, "OK") {
		fmt.Printf("HMSET fail: %q\n", v)
		fails++
	} else {
		fmt.Println("HMSET ok")
	}
	tests++
	writeArray(conn, []string{"HGET", "hkey", "f1"})
	if v, _ := readReply(r); v != "v1" {
		fmt.Printf("HGET fail: %q\n", v)
		fails++
	} else {
		fmt.Println("HGET ok")
	}

	// LPUSH / LGET
	tests++
	writeArray(conn, []string{"LPUSH", "lkey", "a"})
	// LPUSH may return an array and OK; read first reply
	if v, err := readReply(r); err != nil || v == "" {
		fmt.Printf("LPUSH fail (reply1): %q err=%v\n", v, err)
		fails++
	} else {
		fmt.Println("LPUSH reply ok")
	}
	// there may be extra simple reply OK; try to read non-blocking briefly
	conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	if v, err := readReply(r); err == nil && v != "" {
		// consume optional OK
	}
	conn.SetReadDeadline(time.Time{})
	tests++
	writeArray(conn, []string{"LGET", "lkey"})
	if v, _ := readReply(r); !expectContains(v, "a") {
		fmt.Printf("LGET fail: %q\n", v)
		fails++
	} else {
		fmt.Println("LGET ok")
	}

	// SADD / SMEMBERS
	tests++
	writeArray(conn, []string{"SADD", "skey", "m1"})
	if v, _ := readReply(r); !expectContains(v, "OK") {
		fmt.Printf("SADD fail: %q\n", v)
		fails++
	} else {
		fmt.Println("SADD ok")
	}
	tests++
	writeArray(conn, []string{"SMEMBERS", "skey"})
	if v, _ := readReply(r); !expectContains(v, "m1") {
		fmt.Printf("SMEMBERS fail: %q\n", v)
		fails++
	} else {
		fmt.Println("SMEMBERS ok")
	}

	fmt.Printf("\nSummary: %d tests, %d failures\n", tests, fails)
	if fails > 0 {
		return 2
	}
	return 0
}

func main() {
	addr := flag.String("addr", "127.0.0.1:6379", "address of redis server")
	flag.Parse()
	code := runTests(*addr)
	os.Exit(code)
}
