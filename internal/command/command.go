package command

import (
	"net"
	"strings"
	"os"
	"com.ityurika/go-redis-clone/internal/db"
	"com.ityurika/go-redis-clone/internal/protocol"
	"com.ityurika/go-redis-clone/internal/persistence"
	"bufio"
	"io"
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
	"LPUSH": func(conn net.Conn, args [] string) {
		if len(args) !=3 {
			protocol.WriteError(conn, "ERR wrong number of arguments for 'lpush' command")
			return
		}
		val, ok := db.GetDB().LPUSH(args[1], args[2]) 
		if !ok {
			protocol.WriteError(conn, "ERR failed to push value to list")
			return
		}
		protocol.WriteArray(conn, val)
		protocol.WriteSimpleString(conn, "OK")
	},
	"LGET": func(conn net.Conn, args []string) {
		if len(args) != 2 {
			protocol.WriteError(conn, "ERR wrong number of arguments for 'lget' command")
			return
		}
		val, ok :=db.GetDB().LGET(args[1])
		if !ok {
			protocol.WriteError(conn, "ERR failed to get list by index")
			return
		}
		protocol.WriteArray(conn, val)
	},
	"ZADD": func(conn net.Conn, args [] string) {
		if len(args) < 4 || (len(args)-2)%2 !=0 {
			protocol.WriteError(conn, "ERR wrong number of arguments for 'zadd' command")
			return
		}
		key := args[1]
		for i := 2; i < len(args); i += 2 {
			score := args[i]
			member := args[i+1]
			ok := db.GetDB().ZADD(key, score, member)
			if !ok {
				protocol.WriteError(conn, "ERR failed to add member to sorted set")
				return
			}
		}
		protocol.WriteSimpleString(conn, string("OK"))
	},
	"ZRANGE": func(conn net.Conn, args[] string) {
		if len(args) != 4 {
			protocol.WriteError(conn, "ERR wrong number of arguments for 'zrange' command")
			return
		}
		data, ok := db.GetDB().ZRANGE(args[1], args[2], args[3])
		if !ok {
			protocol.WriteError(conn, "ERR failed to get range from sorted set")
			return
		}
		protocol.WriteArray(conn, data)
		protocol.WriteSimpleString(conn, "OK")
	},
}

/**
 * HandleCommand handles the command from the client.
 * @param conn net.Conn 连接
 * @param cmd string 命令
 * @param args []string 参数
 */
func HandleCommand(conn net.Conn, cmd string, args []string) {
	ExecuteAOF(db.GetDB(), args)
	f, ok := commands[strings.ToUpper(cmd)]
	if !ok {
		protocol.WriteSimpleString(conn, "ERR unknown command '"+cmd+"'")
		return
	}
	f(conn, args)
}

/**
 * ExecuteAOF executes the command from AOF.
 * @param db *db.DB 数据库
 * @param cmd []string 命令
 */
func ExecuteAOF(db *db.DB, args [] string) {
	aof, err := persistence.OpenAOF("/media/shiokou/DevRepo/DevHub/Projects/2025-myapp/redis-golang/go-redis-server/log/aof.txt")
	if err != nil {
		panic(err)
	}
	switch strings.ToUpper(args[0]) {
		case "SET":
			db.SetString(args[1], args[2])
			aof.Append(args)
			protocol.WriteSimpleString(os.Stdout, "OK")
		// case "DEL":
		// 	db.Del(cmd[1])
		// 	persistence.Append(cmd)
	}
}

/**
 * Load AOF file
 */
func LoadAOF(filename string, conn net.Conn) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		args, err := protocol.ParseArray(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		HandleCommand(conn, args[0], args)
	}
	return nil
}