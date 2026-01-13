package persistence

import (
	"encoding/binary"
	"fmt"
	"os"

	"com.ityurika/go-redis-clone/internal/data"
	"com.ityurika/go-redis-clone/internal/db"
)

type RDB struct {
	filename string
}

/*
 * 创建RDB实例
 * @param filename RDB文件名
 */
func CreateRDBInstance(filename string) *RDB {
	fmt.Println("Creating RDB instance...")
	return &RDB{
		filename: filename,
	}
}

/**
 * 保存数据库到RDB文件，可以通过goroutine实现异步存储
 * @param kv 数据库实例
 */
func (r *RDB) Save(kv *db.DB) error {
	fmt.Println("Saving database to dump file...")
	f, err := os.Create(r.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	//Write RDB file header
	if _, err := f.Write([]byte("RDB0")); err != nil {
		return err
	}

	//Write the number of keys
	n := uint32(kv.GetStoreLength())
	if err := binary.Write(f, binary.BigEndian, n); err != nil {
		return err
	}

	//Write each key-value
	for k, v := range kv.Store {
		// --- 写入 Key ---
		keyBytes := []byte(k)
		if err := binary.Write(f, binary.BigEndian, uint32(len(keyBytes))); err != nil {
			return err
		}
		if _, err := f.Write(keyBytes); err != nil {
			return err
		}

		// --- 写入 Value ---
		// 1. 先把 Value 转换为字节数组（这里建议调用你定义的序列化方法，或者暂时用 string）
		valStr := fmt.Sprintf("%v", v)
		valBytes := []byte(valStr)

		// 2. 写入 Value 的长度 (这是 Load 函数读取时需要的第一个 4 字节)
		if err := binary.Write(f, binary.BigEndian, uint32(len(valBytes))); err != nil {
			return err
		}

		// 3. 写入 Value 的具体内容
		if _, err := f.Write(valBytes); err != nil {
			return err
		}
	}
	return nil
}

//TODO：现在只有一个save方法，完整的机制应该是fork子进程，然后父进程继续执行，在go当中是相当于在每次创建的时候就异步调用一次bsave方法，把数据存在数据库当中。

/**1
 * 加载RDB文件到数据库
 * @param kv 数据库实例
 */
func (r *RDB) Load(kv *db.DB) error {
	fmt.Println("Loading RDB file to database...")
	f, err := os.Open(r.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	header := make([]byte, 4)
	if _, err := f.Read(header); err != nil {
		return err
	}
	if string(header) != "RDB0" {
		return fmt.Errorf("invalid rdb file")
	}

	var n uint32
	if err := binary.Read(f, binary.BigEndian, &n); err != nil {
		return err
	}

	for i := uint32(0); i < n; i++ {
		var keyLen uint32
		//二进制的读写
		if err := binary.Read(f, binary.BigEndian, &keyLen); err != nil {
			return err
		}
		keyBuf := make([]byte, keyLen)
		if _, err := f.Read(keyBuf); err != nil {
			return err
		}

		var vallen uint32
		if err := binary.Read(f, binary.BigEndian, &vallen); err != nil {
			return err
		}
		ValBuf := make([]byte, vallen)
		if _, err := f.Read(ValBuf); err != nil {
			return err
		}
		kv.Store[string(keyBuf)] = &data.StringValue{Val: string(ValBuf)}
	}
	return nil
}
