package command

import (
	"net"
	"strings"

	"com.ityurika/go-redis-clone/internal/db"
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
			protocol.WriteError(conn, "ERR wrong number of arguments for 'echo' command")
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
	"HMSET": func(conn net.Conn, args []string) {
		// args: [HMSET, key, field, value, field, value, ...]
		if len(args) < 4 || (len(args)-2)%2 != 0 {
			protocol.WriteError(conn, "ERR wrong number of arguments for 'hmset' command")
			return
		}
		key := args[1]
		for i := 2; i < len(args); i += 2 {
			field := args[i]
			value := args[i+1]
			ok := db.GetDB().HSetField(key, field, value)
			if !ok {
				protocol.WriteError(conn, "ERR failed to set hash")
				return
			}
		}
		protocol.WriteSimpleString(conn, "OK")
	},
	"HGET": func(conn net.Conn, args []string) {
		// args: [HGET, key, field]
		if len(args) != 3 {
			protocol.WriteError(conn, "ERR wrong number of arguments for 'hget' command")
			return
		}
		val, ok := db.GetDB().HGetField(args[1], args[2])
		if !ok {
			// nil bulk reply -> empty bulk
			protocol.WriteBulkString(conn, "")
			return
		}
		protocol.WriteBulkString(conn, val)
	},
	"SADD": func(conn net.Conn, args []string) {
		if len(args) < 3 {
			protocol.WriteError(conn, "ERR wrong number of arguments for 'sadd' command")
			return
		}
		ok := db.GetDB().SADD(args[1], args[2])
		if !ok {
			protocol.WriteError(conn, "ERR failed to add member to set")
			return
		}
		protocol.WriteBulkString(conn, "OK")
	}, 
	"SMEMBERS": func(conn net.Conn, args []string) {
		if len(args) != 2 {
			protocol.WriteError(conn, "ERR wrong number of arguments for 'smembers' command")
			return
		}
		members, ok := db.GetDB().SMEMBERS(args[1])
		if !ok {
			protocol.WriteError(conn, "ERR failed to get members from set")
			return
		}
		protocol.WriteArray(conn, members)
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
