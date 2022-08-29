[TOC]


# epoll系列函数

\#include <sys/epoll.h>  
int epoll_create(int size);

- size指epoll实例的大小
- 成功时返回epoll文件描述符，失败时返回-1

调用epoll_create函数时创建的fd保存空间称为“epoll例程”  
通过size参数决定例程的大小，但这只是向操作系统提的“建议”（Linux2.6.8后该参数将被完全忽略）

int epoll_ctl(int epfd, int op, int fd, struct epoll_event* event);

- epfd用于注册监视对象的epoll例程的fd
- op用于指定监视对象的添加、删除或更改等操作
    + EPOLL_CTL_ADD将fd注册到epoll例程
    + EPOLL_CTL_DEL从epoll例程中删除fd
    + EPOLL_CTL_MOD更改注册的fd的关注事件发生情况
- fd需要注册的监视对象的fd
- event监事对象的事件类型
- 例子：epoll_ctl(A, EPOLL_CTL_ADD, B, C);“epoll例程A中注册文件描述符B，主要目的时监视参数C中的事件” 
- 成功时返回0，失败时返回-1
```c
struct epoll_event
{
    __uint32_t      event;
    epoll_data_t    data;
}
    typedef union epoll_data
    {
        void*       ptr;
        int         fd;
        __uint32_t  u32;
        __uint64_t  u64;
    }epoll_data_t;

//调用
struct epoll_event event;
......
event.events = EPOLLIN;
event.data.fd=sockfd;
epoll_ctl(epfd, EPOLL_CTL_ADD, sockfd, &event);
```
其中，epoll_event.events的选项有：
- EPOLLIN需要读取数据的情况
- EPOLLOUT输出缓冲为空，可用立即发送数据的情况
- EPOLLPRI收到OOB数据
- EPOLLDHIP断开连接或半关闭的情况，这在边缘触发方式下非常有用
- EPOLLERR发生错误
- EPOLLET以边缘触发的方式得到事件通知
- EPOLLONESHOT发生一次事件后，相应的fd不再收到事件通知。因此需要向epoll_ctl函数的第二个参数传递EPOLL_CTL_MOD，再次设置事件
- 可以通过位运算同时传递多个上述参数：|

int epoll_wait(int epfd, struct epoll_event* events, int maxevents, int timeout)

- epfd事件发生监视范围的epoll例程的fd
- events保存发生事件的文件描述符集合的地址
- maxevents第二个参数中可以保存的最大事件数
- timeout以1/1000秒为单位的等待事件，传递-1时，一直等待直到发生事件
- 成功时返回发生事件的文件描述符数，失败时返回-1
- 第二个参数所指缓冲需要动态分配
```c
int event_cnt;
struct epoll_event* ep_events;
......
ep_events = malloc(sizeof(struct epoll_event)*EPOLL_SIZE);//EPOLL_SIZE是宏常量
......
event_cnt = epoll_wait(epfd, ep_events, EPOLL_SIZE, -1);
```

# epoll服务器端示例代码

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/epoll.h>

#define BUF_SIZE 100
#define EPOLL_SIZE 50
void error_handling(char* buf);

int main(int argc, char const *argv[])
{
    int serv_sock, clnt_sock;
    struct sockaddr_in serv_adr, clnt_adr;//socket.h
    socklen_t adr_sz;
    int str_len;
    char buf[BUF_SIZE];

    struct epoll_event* ep_events;//epoll.h
    struct epoll_event event;
    int epfd, event_cnt;

    if(argc!=2){
        printf("Usage : %s <port>\n", argv[0]);
        exit(1);
    }

    serv_sock=socket(PF_INET, SOCK_STREAM, 0);
    memset(&serv_adr, 0, sizeof(serv_adr));
    serv_adr.sin_family=AF_INET;
    serv_adr.sin_addr.s_addr=htonl(INADDR_ANY);
    serv_adr.sin_port=htons(atoi(argv[1]));

    if(bind(serv_sock, (struct sockaddr*)&serv_adr, sizeof(serv_adr))==-1)
        error_handling("bind() error");
    if(listen(serv_sock, 5)==-1)
        error_handling("listen() error");

    epfd=epoll_create(EPOLL_SIZE);
    ep_events=malloc(sizeof(struct epoll_event)*EPOLL_SIZE);

    event.events=EPOLLIN;
    event.data.fd=serv_sock;
    epoll_ctl(epfd, EPOLL_CTL_ADD, serv_sock, &event);

    while(1)
    {
        event_cnt=epoll_wait(epfd, ep_events, EPOLL_SIZE, -1);
        if(event_cnt==-1)
        {
            puts("epoll_wait() error");
            break;
        }

        for (int i = 0; i < event_cnt; ++i)
        {
            if(ep_events[i].data.fd==serv_sock)
            {
                adr_sz=sizeof(clnt_adr);
                clnt_sock=accept(serv_sock, (struct sockaddr*)&clnt_adr, &adr_sz);
                event.events=EPOLLIN;
                event.data.fd=clnt_sock;
                epoll_ctl(epfd, EPOLL_CTL_ADD, clnt_sock, &event);
                printf("connected client: %d \n", clnt_sock);
            }
            else
            {
                str_len=read(ep_events[i].data.fd, buf, BUF_SIZE);//?
                if(str_len==0)//close request
                {
                    epoll_ctl(epfd, EPOLL_CTL_DEL, ep_events[i].data.fd, NULL);
                    close(ep_events[i].data.fd);
                    printf("closed client: %d \n",ep_events[i].data.fd);
                }
                else
                {
                    write(ep_events[i].data.fd, buf, str_len);//echo
                }
            }
        }
    }
    close(serv_sock);
    close(epfd);
    return 0;
}

void error_handling(char* msg)
{
    fputs(msg, stderr);
    fputc('\n', stderr);
    exit(1);
}


```

# 条件触发与边缘触发

条件触发：每当收到客户端数据时，都会注册该事件，并多次调用epoll_wait()  

边缘触发：接收数据时仅注册一次该事件  
为了实现边缘触发，需要

- 通过errno变量验证错误原因
- 为了完成非阻塞IO，更改套接字特性

为了在发生错误时获得额外的信息，Linux申明了这个全局变量  
\#include <error.h>  
int errno  
每种函数发生错误时，保存到errno中的值都不同，没有必要记住所有可能的值  
此处需要知道的是：read()发现输入缓冲中没有数据可读时返回-1，同时在errno中保存EAGAIN常量

将套接字改为非阻塞方式的函数  
\#include <fcntl.h>  
int fcntl(int filedes, int cmd, ...);

- filedes需要修改的fd
- cmd函数调用的目的
- 成功时返回cmd参数相关值，失败时返回-1
- 如果向第二个参数传递F_GETFL，可用获得第一个参数所指的fd属性（int型）
- 如果传递F_SETFL，可以更改fd属性。若希望将fd改为非阻塞模式，需要：
    + int flag = fcntl(fd, F_GETFL, 0);//获得属性
    + fcntl(fd, F_SETFL, flag|O_NONBLOCK);//添加O_NONBLOCK属性




