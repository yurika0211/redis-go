package persistence

import (
	    "com.ityurika/go-redis-clone/internal/db"
		"os"
		"encoding/binary"
		"fmt"
)

type RDB struct {
	filename string
}


/*
 * 创建RDB实例
 * @param filename RDB文件名
*/
func CreateRDBInstance (filename string) *RDB {
	fmt.Println("Creating RDB instance...")
	return & RDB {
		filename: filename,
	}
}

/**
 * 保存数据库到RDB文件
 * @param kv 数据库实例
*/
func (r *RDB) Save(kv *db.DB) error {
	fmt.Println("Saving database to RDB file...")
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
		if err := binary.Write(f, binary.BigEndian, uint32(len(k))); err != nil {
			return err
		}
		if _, err := f.Write([]byte(k)); err != nil {
			return err
		}
		// Here we should serialize the value v based on its type
		// This is a simplified example, actual implementation may vary
		valStr := fmt.Sprintf("%v", v)
		if err := binary.Write(f, binary.BigEndian, v); err != nil {
			return err
		}
		if _ , err := f.Write([]byte(valStr)); err != nil {
			return err
		}
	}
	return nil
}

/**1
 * 加载RDB文件到数据库
 * @param kv 数据库实例
*/
func (r *RDB) Load(kv *db.DB) error {
	fmt.Println("Loading RDB file to database...")
	f, err := os.Open(r.filename)
	if err != nil { return err}
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
	}
	return nil
}