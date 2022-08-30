[TOC]
# 相关知识

sh：Bourne Shell
csh：C Shell
bash：Bourne Again SHell

系统中合法的shells：/etc/shells  
账户登录时的默认shell：/etc/passwd 每行的最后一个字段  
前一次登录执行过的命令：~/.bash_history  
别名：alias  
查看命令类型：type [-apt] name  

# 变量

## 变量的使用与设置

```bash
显示
echo $name
echo ${name}

赋值
name=zinkt
doublename="$zinkt+$zinkt" 双引号保有特殊字符特性
singlename='$zinkt+$zinkt' 单引号为纯文本
varname=$zinkt+\$zinkt 反斜杠转义
name=${name}zinkt 扩增内容

操作
$(uname -r) 获取内核版本号（在有带空格参数时，使用括号）
`` 反单引号内的命令会被先执行，因此也可以是
`uname -r`

export name 使变量变成环境变量
unset name 取消变量

可在~/.bashrc 中指定变量
```

## 环境变量的功能

### env下典型环境变量

env可列出shell下所有环境变量
- HOME 用户根目录
- SHELL 目前shell使用的程序
- HISTSIZE 历史命令记录的条数
- MAIL
- PATH 执行文件查找的路径，以:分隔
- LANG 语系
- RANDOM 0~32767随机数，/dev/random这个文件是随机数生成器

### set观察所有变量

- PS1 命令提示符的格式
- $ shell的PID
- ? 上一个执行命令的返回值

export 将自定义变量转换为环境变量，以供子进程使用  
因为子进程会继承父进程的环境变量，却不会继承自定义变量  

## 其它

### locale 语系变量
- /usr/lib/locale 语系文件
- /etc/locale.conf 系统默认语系定义

### read

read var 等待输入

### declare

declare [-aixr] var 声明变量类型
-a array类型
-i integer
-x 环境变量
-r readonly，不可被更改内容，也不能unset

### array

var[index]=content

### ulimit

限制用户的某些系统资源  


# 命令别名与历史命令

## 命令别名的设置 alias，unalias

alias ll='ls -alh'

## 历史命令

```bash
history [n] 显示最近n条
history [-c] 将目前shell中的history内容全部清除
history [-raw] histfiles
!! 执行上一条命令
!n 执行第n条
!vi 执行最近命令开头是vi的命令
```

# Bash Shell的操作环境

***路径与命令查找顺序***  
1. 以相对/绝对路径来执行命令，如/bin/ls或./ls
2. 有alias找到命令来执行
3. bash内置的（buildin）命令来执行
4. 通过$PATH顺序找到的第一个命令来执行

***bash的登录与欢迎信息***  

/etc/issue
/etc/issue.net 用于远程登录的欢迎信息

## bash的环境配置文件

***读取过程***  
- login shell
    + /etc/profile 整体设置
        * PATH
        * MAIL
        * USER
        * HOSTNAME
        * HISTSIZE
        * umask 包括root默认为022而一般用户为002等
        * 此后会依次调用以下文件
            - /etc/profile.d/*.sh
            - /etc/locale.conf
            - /usr/share/bash-completion/completions/*
    + ~/.bash_profile （login shell才读）
    + 或~/.bash_login
    + 或~/.profile （这三个只读一个）

### 其他关键指令/文件

***source***读入环境配置文件  
source 配置文件名  如：
source ~/.bashrc

***.bashrc***（non-login shell会读）  
即配置文件  

***/etc/man_db.conf***  
规范了使用man的使用，man page的路径到哪里取寻找

***~/.bash_history***  

***~/.bash_logout***

## 关于终端

***stty (setting tty)***  
stty [-a]  

***set***  
set [-uvCHhmBx]  

### 默认快捷键

- Ctrl + C 终止目前的指令
- Ctrl + D 输入结束
- Ctrl + M 代表回车
- Ctrl + S 暂停屏幕的输出
- Ctrl + Q 恢复屏幕的输出
- Ctrl + Z 暂停当前的命令

### 通配符与特殊符号

*通配符*  
- * 代表【0到无穷多个】字符
- ? 代表【一定有一个】字符
- [] 代表【一定有一个在括号内】的字符
- [-] 若有减号在中括号，代表【在编码顺序内的所有字符】。如[0-9]代表0和9间的所有数字
- [^] 表示【反向选择】，如[^abc]代表一定有一个字符，只要是非a、b、c就接受

*特殊字符*  
- #         注释，其后不执行
- |         pipe，分隔两个管道命令
- ;         连续执行分隔符
- &         任务管理，将命令变成后台任务
- >、>>      数据流重定向
- <、<<
- ''        其中代表纯文本
- ""        具有变量替换功能
- ``        其中为先执行的指令


# 数据流重定向

< 标准输入  
<< 标准输入  
\> 标准输出，覆盖  
\>> 标准输出，追加  
2> 标准错误输出  
2>> 标准错误输出  

*结果存到不同文件中*  
find /home -name .bashrc > list_right 2> list_error

*/dev/null 垃圾桶黑洞*  
丢弃错误信息，只显示正确信息  
find /home -name .bashrc 2> /dev/null

*写入同一个文件*  
find /home -name .bashrc > list 2>&1  

*示例*  
```bash
zinkt@zinkt-asus:~$ cat > catfile
testing
???!!!
zinkt@zinkt-asus:~$ cat catfile 
testing
???!!!

cat > catfile < ~/.bashrc # < 即用文本代替键盘输入
zinkt@zinkt-asus:~$ cat catfile 
# ~/.bashrc: executed by bash(1) for non-login shells.
# see /usr/share/doc/bash/examples/startup-files (in the package bash-doc)
# for examples
。。。

无论是正确还是错误信息都写入一个文件：
find /home -name .bashrc > list 2> list 错误，此时两股数据可能会交叉写入导致错误 
find /home -name .bashrc > list 2>&1  正确，将 2> 转到 1>
```

## 命令执行判断的依据

*$?*命令返回值  

cmd1 && cmd2  
若cmd1执行完毕且正确执行（$?=0），则执行cmd2；若cmd1执行完毕且为错误，则cmd2不执行  
cmd1 || cmd2  
若cmd1执行完毕且正确执行，则cmd2不执行；若$?不为0，则执行cmd2

ls /tmp/abc || mkdir /tmp/abc && touch /tmp/abc/hehe

总是会建立/tmp/abc/hehe

# 管道命令pipe

仅能够处理标准输出  
管道命令必须要能接收来自前一个命令的数据成为标准输入继续处理才行  
例：ll /etc | less  

## 选取命令 cut、grep

*cut*将同一行的数据进行分解
`echo $PATH | cut -d ':' -f 3,5` 以:分隔，选出第3，第5个  
`export | cut -c 12-` 切出12个字符后的字段

*grep*分析一行信息，若有我们想要的，就将此行拿出来  
`grep [-cinv] [--color=auto] '查找字符' filename`  
- c 计算找到的次数，count
- i 忽略大小写 ignore
- n 输出行号
- v 反向选择，即显示没有'查找字符'的行

## 排序命令 sort、wc、uniq

*sort*按某列（属性）排序  
`cat /etc/passwd | sort`默认以第一条信息排序  
`cat /etc/passwd | sort -t ':' -k 3`passwd用':'分隔，-k用第三栏排序  

*uniq*将重复数据仅列出一个显示（一般还搭配计数）  
例`last | cut -d ' ' -f1 | sort | uniq -c`

*wc*计算信息的整体数据

## 双向重定向tee

将数据流分别送到文件和屏幕  
tee [-a以append方式写入数据] file  
`last | tee last.list | cut -d ' ' -f1`

## 字符转化命令

*tr*删除一段信息中的文字，或进行文字信息替换  
*col*大多数用做将[tab]替换为空格  
*join*处理两个相关的数据文件  
*paste*将两个文件直接贴在一起，中间用[tab]分隔  
*expand*将[tab]转换为空格键

## 划分命令split

将一个大文件依据文件大小或行数来划分  
split [-bl] file PREFIX

- b后接欲划分成的文件大小，可接单位b,k,m等
- l以行数来划分


## 参数代换xargs

`find /usr/sbin/ -perm /7000 | xargs ls -l`找出目录下有特殊权限的文件名，并用ls -l列出其属性


## 减号【-】的用途

某些需要用到文件名（如tar）来处理时，stdin和stdout可用【-】来替代：  
`tar -cvf - /home | tar -xvf - -C /tmp/homeback`













