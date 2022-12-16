 # mmap

 当使用mmap来映射文件时，kernel使用懒加载把虚拟内存初始化，但实际物理内存并没有初始化。只有在实际使用时才会通过page fault实现物理内存的数据加载。
 可以通过以下实验来验证此场景

 1. 创建一个磁盘大小为5G左右的文件
 > head -c 5000000000  /dev/zero > mmapfile.txt

 2.执行mmap.go文件，可以通过top查看进程内存使用情况

 > Ref: https://ostechnix.com/create-files-certain-size-linux/


