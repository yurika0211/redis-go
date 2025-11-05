package command

import (
	"net"
	"strings"

	"com.ityurika/go-redis-clone/internal/protocol"
)

var commands = map[string]func(net.Conn, []string){
	"PING": func(conn net.Conn, args []string) {
		if len(args) > 1 {
			protocol.WriteBulkString(conn, args[1])
		} else {
			protocol.WriteSimpleString(conn, "PONG")
		}
	},
	"ECHO": func(conn net.Conn, args []string) {
		if len(args) != 2 {
			protocol.WriteSimpleString(conn, "ERR wrong number of arguments for 'echo' command")
			return
		}
		protocol.WriteBulkString(conn, args[1])
	},
}

func HandleCommand(conn net.Conn, cmd string, args []string) {
	f, ok := commands[strings.ToUpper(cmd)]
	if !ok {
		protocol.WriteSimpleString(conn, "ERR unknown command '"+cmd+"'")
		return
	}
	f(conn, args)
}
