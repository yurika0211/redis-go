package persistence

import (
	"os"
	"testing"

	"com.ityurika/go-redis-clone/internal/data"
	"com.ityurika/go-redis-clone/internal/db"
)

// 测试 RDB 的 Save 和 Load 功能
func TestRDB_SaveLoad(t *testing.T) {
	// 创建临时文件
	tmpfile := "dump.rdb"
	defer os.Remove(tmpfile)

	// 初始化数据库
	d := db.GetDB()
	d.Store = make(map[string]data.Value) // 清空
	d.Store["foo"] = &data.StringValue{Val: "bar"}
	d.Store["num"] = &data.StringValue{Val: "123"} // 简单示例

	// 创建 RDB 实例
	rdb := CreateRDBInstance(tmpfile)

	// 保存数据库到 RDB 文件
	if err := rdb.Save(d); err != nil {
		t.Fatalf("RDB Save failed: %v", err)
	}

	// 清空数据库，模拟重启
	d.Store = make(map[string]data.Value)

	// 加载 RDB 文件
	if err := rdb.Load(d); err != nil {
		t.Fatalf("RDB Load failed: %v", err)
	}

	// 验证数据
	if strVal, ok := d.GetString("foo"); !ok || strVal != "bar" {
		t.Errorf("expected foo=bar, got %v", strVal)
	}
	if strVal, ok := d.GetString("num"); !ok || strVal != "123" {
		t.Errorf("expected num=123, got %v", strVal)
	}

	// 测试加载不存在的文件
	rdbBad := CreateRDBInstance("nonexistent.rdb")
	if err := rdbBad.Load(d); err == nil {
		t.Errorf("expected error when loading nonexistent file")
	}

	// 测试加载格式错误的文件
	badFile := "bad_rdb.rdb"
	os.WriteFile(badFile, []byte("BAD!"), 0644)
	defer os.Remove(badFile)
	rdbBad2 := CreateRDBInstance(badFile)
	if err := rdbBad2.Load(d); err == nil {
		t.Errorf("expected error when loading invalid rdb file")
	}
}
