
## 解题思路

我使用的是邓公视频中所提示的解法，即讲过的归并，简单来说就是两步：  

1. 将坐标对按x排序：mergeSort()。此处使用的是归并排序，当然也可以使用其他排序
2. 求逆序对数：invsOf()。在对y的归并排序过程中，记录顺序对个数

## AC代码

```cpp
#include <cstdio>
#include <cctype>
#define MAXN 4000000
struct dot
{
    int x;
    int y;
} dots[MAXN];
long long cnt = 0;
dot* b = new dot[2000000];      //临时数组，减少每次new的时间消耗
template <typename T>
inline T readstdin()
{
    T input = 0;    short sign = 1;    char ch = 0;
    while (! isdigit(ch))
    {
        if (ch == '-')  sign = -1;
        ch = getchar();
    }
    while (isdigit(ch))
    {
        input = (input << 3) + (input << 1) + (ch - '0');
        ch = getchar();
    }
    return sign * input;
}
void mergeSort(int lo, int mi, int hi)
{
    if (hi - lo < 2)
        return;
    mergeSort(lo, (lo + mi) >> 1, mi);
    mergeSort(mi, (mi + hi) >> 1, hi);
    dot *a = dots + lo;
    int lb = mi - lo;
    for (int i = 0; i < lb; i++)
        b[i] = a[i];
    int lc = hi - mi;
    dot *c = dots + mi;
    for (int i = 0, j = 0, k = 0; j < lb;)
    {
        if ((j < lb) && (lc <= k || (b[j].x < c[k].x)))
            a[i++] = b[j++];
        if ((k < lc) && (c[k].x < b[j].x))
            a[i++] = c[k++];
    }
}
void invsOf(int lo, int hi)
{
    if (hi - lo < 2)
        return;
    invsOf(lo, (lo + hi) >> 1);
    invsOf((lo + hi) >> 1, hi);
    int mi = (lo+hi) >> 1;
        dot *a = dots + lo;
    int lb = mi - lo;
    for (int i = 0; i < lb; i++)
        b[i] = a[i];
    int lc = hi - mi;
    dot *c = dots + mi;
    for (int i = 0, j = 0, k = 0; j < lb ;)
    {
        if ((j < lb) && (lc <= k || (b[j].y < c[k].y)))
        {
            if (k < lc)
                cnt += lc - k;
            a[i++] = b[j++];
        }
        if ((k < lc) && (c[k].y < b[j].y))
        {
            a[i++] = c[k++];
        }
    }
}
int main(int argc, char const *argv[])
{
    long long n;
    scanf("%lld", &n);
    for (int i = 0; i < n; i++)
    {
        dots[i].x = readstdin<int>();
        dots[i].y = readstdin<int>();
    }
    mergeSort(0, n >> 1, n);
    invsOf(0, n);
    printf("%lld", cnt);
    return 0;
}
```

## 解题历程

一开始我第一反应的直观解法就是暴力对比，当然这个O(n^2)算法只能拿到45分，后续的检测点都会超时    
但看完邓公的提示后依然有点困惑，主要在于如何排除x的干扰，从而专注于求y的顺序对个数    
且在想到需要先按x排序后，我仍没想明白：用归并法求y的顺序对时，x的顺序已被打乱，对x的排序还有何意义呢  
之后彻底想通：尽管x的顺序会被打乱，但在求y的顺序对数的归并过程中，[lo,mi),[mi,hi)两个序列中各元素的x值仍是前者大于后者的    
这就给归并求y顺序对个数创造了条件   

并且在归并（排序、求顺序对数）的过程中，可以优化的地方是（邓公在视频中留下的问题）：  
由于c序列是原本就存在于a序列中的，因此当b序列全部按序置入a中（即lb<=j）后，a序列就已经完成merge了   
