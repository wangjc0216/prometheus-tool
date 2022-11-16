# wal tool

可通过命令来：
1)查看wal目录下segment文件;
2)查看segment文件具体内容，包括sample、series、exemplar、tombstone

```
# 包含对应dir目录下的segment文件信息
go run main.go --wal_dir <prom-wal-dir>

# 查看segment具体内容
go run main.go --segfile mock_data/wal/00000041
```


