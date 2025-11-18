package command

import (
	"net"
	"strings"

	"com.ityurika/go-redis-clone/internal/protocol"
	"com.ityurika/go-redis-clone/internal/db"
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
	"SET": func(conn net.Conn, args []string) {
		if len(args) != 3 {
			protocol.WriteSimpleString(conn, "ERR wrong number of arguments for 'set' command")
			return
		}
		db.GetDB().SetString(args[1], args[2])
		protocol.WriteSimpleString(conn, "OK")
	},
	"GET": func(conn net.Conn, args []string) {
		if len(args) != 2 {
			protocol.WriteBulkString(conn, "ERR wrong number of arguments for 'get' command")
			return
		}
		val, ok := db.GetDB().GetString(args[1])
		if !ok {
			protocol.WriteBulkString(conn, "")
			return
		}
		protocol.WriteBulkString(conn, val)
	},
}

/**
 * HandleCommand handles the command from the client.
 * @param conn net.Conn 连接
 * @param cmd string 命令
 * @param args []string 参数
 */
func HandleCommand(conn net.Conn, cmd string, args []string) {
	f, ok := commands[strings.ToUpper(cmd)]
	if !ok {
		protocol.WriteSimpleString(conn, "ERR unknown command '"+cmd+"'")
		return
	}
	f(conn, args)
}
