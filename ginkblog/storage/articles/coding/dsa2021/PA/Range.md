## 解题思路

本题有两个解法  
正常（课程期望）解法：将数轴坐标排序后，二分查找到两个点，下标相减即可  
前缀和：将坐标直接作为数组下标，值为[0,i]范围内存在点的个数。查询时直接取出边界点对应的值相减即可   

## 正常解法

1. 这里使用快排将坐标排序
2. 二分查找到对应的点，并将下标相减

```cpp
#include <cstdio>
#include <cctype>
#define MAXN 500000
int axis[MAXN];

template <typename T>
inline T readstdin()                //加快IO的读函数
{
    T input = 0;    short sign = 1;    char ch = 0;
    while (!isdigit(ch))
    {
        if(ch == '-')
            sign = -1;
        ch = getchar();
    }
    while (isdigit(ch))
    {
        input = (input << 3) + (input << 1) + (ch - '0');
        ch = getchar();
    }
    return input * sign;
}
int binSearch( int n, int e)        //二分查找到[0,n]范围内的，不大于e的最大值
{
    int lo = 0, hi = n;
    while (lo<hi)
    {
        int mi = (lo+hi)>>1;
        (e<axis[mi]) ? hi=mi : lo = mi + 1;
    }
    return --lo;
}
int getpartition(int lo,int hi){    //用于快排
    int mi=axis[lo];
    while(lo<hi)
    {
        while(lo<hi&&mi<=axis[hi])
            hi--;
        axis[lo]=axis[hi];
        while(lo<hi&&axis[lo]<=mi)
            lo++;
        axis[hi]=axis[lo];
    }
    axis[lo]=mi;
    return lo;
}

void quicksort(int lo,int hi){      //快排
    if(lo<hi)
    {
        int mi=getpartition(lo,hi);
        quicksort(lo,mi-1);
        quicksort(mi+1,hi);
    }
}

int main(int argc, char const *argv[])
{
    int n = readstdin<int>();
    int m = readstdin<int>();

    for (int i = 0; i < n; i++)
    {
        axis[i] = readstdin<int>();
    }
    quicksort(0,n-1);
    int a, b;
    for (int i = 0; i < m; i++)
    {
        a = readstdin<int>()-1;
        b = readstdin<int>();
        int t1 = binSearch(n,a);
        int t2 = binSearch(n,b);
        printf("%d\n",t2-t1);
    }
    return 0;
}

```


## 前缀和

1. 初始化：将每个输入的数轴坐标直接作为数组的下标，若该处有点，则值置为1，否则为0  
2. 生成前缀和：axis[i] = axis[i-1] 即每个点的值为[0,i]中有效点的个数
3. 计算：接收输入后直接取出对应下标的值相减即可

```cpp
#include <cstdio>
#include <cctype>
#define MAXN 10000001
int axis[MAXN] = {0};

template <typename T>
inline T readstdin()
{
    T input = 0;    short sign = 1;    char ch = 0;
    while (!isdigit(ch))
    {
        if(ch == '-')
            sign = -1;
        ch = getchar();
    }
    while (isdigit(ch))
    {
        input = (input << 3) + (input << 1) + (ch - '0');
        ch = getchar();
    }
    return input * sign;
}

int main(int argc, char const *argv[])
{
    int n = readstdin<int>();
    int m = readstdin<int>();
    for (int i = 0; i < n; i++)             //初始化
        axis[readstdin<int>()] = 1;
    for (int i = 1; i < MAXN; i++)          //生成前缀和，即统计[0,i]的点总数
        axis[i] += axis[i-1];
    for (int i = 0; i < m; i++)
    {
        int a = readstdin<int>() - 1;
        int b = readstdin<int>();
        printf("%d\n",axis[b]-axis[a]);
    }
    return 0;
}
```

## 踩坑

1. 最开始题目描述里给的示例让我一度以为坐标是按序输入的
2. 即使二分查找到了对应的下标，相减时也要仔细考虑各种情况并做出修补（+1-1之类的）
3. 前缀和长知识了