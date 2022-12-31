 # mmap

 当使用mmap来映射文件时，kernel使用懒加载把虚拟内存初始化，但实际物理内存并没有初始化。只有在实际使用时才会通过page fault实现物理内存的数据加载。
 可以通过以下实验来验证此场景

1. 创建一个磁盘大小为5G左右的文件
 > head -c 5000000000  /dev/zero > mmapfile.txt

2. 执行mmap.go文件，可以通过top查看进程内存使用情况

 > Ref: https://ostechnix.com/create-files-certain-size-linux/

3. 使用curl来访问mmap文件不同的offset和size，观察内存变化
 
> http://192.168.0.114:8001/mmap?offset=0&size=2000000000

在访问mmap接口后，通过free -h 可发现 page/buff的使用量在增加，可以得知是page cache的使用量增加：
```
root@xiaomi:~# free -h
              total        used        free      shared  buff/cache   available
Mem:            11G        1.2G        8.0G        912K        2.7G         10G
Swap:          4.0G        2.0M        4.0G

# 调用几次的前后cache变化
root@xiaomi:~# free -h
              total        used        free      shared  buff/cache   available
Mem:            11G        1.2G        5.1G        912K        5.6G         10G
Swap:          4.0G        2.0M        4.0G
```

通过容器启动并对内存加以限制（docker run --name mmap-server -v $PWD/mmapfile.txt:/root/mmapfile.txt  --memory=1024m -d -p 8001:8001  mmap-server:v0.1），
发现使用的内存为page cache，达不到限制内存的目的。



## 构建
> ./build.sh


## TODO
可针对一个很大的磁盘文件(如20G)，做mmap，通过web server 的形式每次去调用offset和size来获取offset，观察对应资源变化
- webserver
- go_client
- 使用sar观察page fault情况（ebpf可能也可以）

