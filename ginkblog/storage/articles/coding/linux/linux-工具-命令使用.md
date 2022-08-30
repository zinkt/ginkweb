

## 阿里云服务器设置swap分区

1. 使用top命令查看当前系统是否有交换分区。

   运行top命令后，可以在KiB Swap打头的那一行，看到交换分区相关信息。如果显示KiB Swap:  0 total 就说明没有交换分区。从top命令中退出使用“q”键。

2. 首先创建用户交换分区的文件。

   [root@wbl~]# dd if=/dev/zero of=/mnt/swap bs=1M count=1024   //创建一个G得交换空间

3. init初始化分区文件

   [root@wbl~]# mkswap /mnt/swap

4. 启动交换分区

   [root@wbl~]# swapon /mnt/swap

5. 设置开机自动挂载

   [root@wbl~]# vim /etc/fstab

   添加 /mnt/swap swap swap defaults 0 0

6. 设置使用swap分区的阀值

   [root@wbl~]# vim /etc/sysctl.conf

   修改文件中的vm.swappiness = 50，阿里云centos默认是0。

   [root@www ~]# sysctl vm.swappiness=20

   注:这个设置为50，是物理内存空间小于20%时，开始使用交换空间

## 快速统计内存占用

`ps -eo pid,rss,pmem,pcpu,vsz,args --sort=rss`