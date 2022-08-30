## 解题思路

跟着课程走，很容易想到使用链表，实现起来比较容易。
其中需要注意的是插入之后多个球相邻同色消除的问题，以及消除一种颜色之后，又有需要消除的情况。我的处理方法是递归调用这个“检测——消除函数”。



## AC代码（链表方法）

原本我选择每次操作都直接输出操作结果，但只能拿到95分，最后一个检测点TLE  
将所有结果先存入一个buf，最后一次性输出，减少IO时间后，成功AC

```cpp
#include <iostream>
#include <string>
using namespace std;
template <typename T>
struct ListNode
{
    T data;
    ListNode* pred;     //predecessor
    ListNode* succ;     //successor
    ListNode(){}
    ListNode(T e, ListNode* prev = NULL, ListNode* next = NULL):data(e),pred(prev),succ(next){}
};
template <typename T>
struct List
{
    List()
    {
        head=new ListNode<T>;
        tail=new ListNode<T>;
        head->succ = tail;
        head->pred = tail;
        tail->pred = head;
        tail->succ = head;
        _size=0;
    }
    ~List(){}
    //只读
    int size(){return _size;}
    bool empty(){return _size<=0;}
    ListNode<T>* begin(){return head->succ;}
    ListNode<T>* end(){return tail;}
    //可写
    ListNode<T>* push_back(T const& e)
    {
        ListNode<T>* p_tmp = new ListNode<T>(e,tail->pred,tail);
        p_tmp->pred->succ=p_tmp;
        tail->pred=p_tmp;
        _size++;
        return p_tmp;
    }
    ListNode<T>* insert(ListNode<T>* p, T const& e)
    {
        ListNode<T>* p_tmp = new ListNode<T>(e,p->pred,p);
        p_tmp->pred->succ=p_tmp;
        p->pred=p_tmp;
        _size++;
        return p_tmp;
    }
    ListNode<T>* erase(ListNode<T>* pos)
    {
        ListNode<T>* tmp=pos;
        pos->pred->succ=pos->succ;
        pos->succ->pred=pos->pred;
        pos=pos->succ;
        delete tmp;
        _size--;
        return pos;
    }
    ListNode<T>* erase(ListNode<T>* first, ListNode<T>* last)
    {
        while (first!=last)
            first=erase(first);
        return last;
    }
private:
    int _size;
    ListNode<T>* head; 
    ListNode<T>* tail;
};
List<char> zuma;
char output[200000000];
int cur = 0;
void traverse_and_print()
{
    if(zuma.size()==0)
    {    
        output[cur++]='-';
        output[cur++]='\n';
    }
    else
    {
        auto p = zuma.begin();
        while (p!=zuma.end())
        {
            output[cur++]=p->data;
            p=p->succ;
        }
            output[cur++]='\n';
    }
}
void check_del(ListNode<char>* pmiddle)
{
    int count=0;
    auto pfront=pmiddle;
    auto itr_before_begin = zuma.begin()->pred;
    while(pfront!=itr_before_begin)
    {
        if((pfront->pred)->data==pmiddle->data)
        {
            count++;
            pfront=pfront->pred;
        }
        else break;
    }
    auto pback = pmiddle;
    while(pback!=zuma.end())
    {
        if((pback->succ)->data==pmiddle->data)
        {
            count++;
            pback=pback->succ;
        }
        else break;
    }
    if(count>=2)
    {
        pback=pback->succ;
        zuma.erase(pfront,pback);
        check_del(pback);
    }
    else
        return;
}
int main(int argc, char const *argv[])
{
    ios::sync_with_stdio(false);
    cin.tie(0);
    string input;
    getline(cin,input);
    for (int i = 0; i < input.size(); i++)
        zuma.push_back(input[i]);
    int n,tmp_pos;
    char tmpc;
    cin >> n;
    for (int i = 0; i < n; i++)
    {
        cin >> tmp_pos >> tmpc;
        auto p = zuma.begin();
        for (int i = 0; i < tmp_pos; i++)
            p=p->succ;
        check_del(zuma.insert(p,tmpc));
        traverse_and_print();
        if(i==n-1)
        {
            output[cur]='\0';
            cout << output;
            cur=0;
        }
    }
    return 0;
}
```

## 踩坑记录

因为这门课的OJ平台不能使用stl（毕竟就是在讲数据结构嘛），自己实现list时函数总是出问题  
1. ListNode<char>* 这个指针++不是指向下一个节点！只是指向下一个连续储存的ListNode大小内存位置，其内容是完全无法确定的。这么憨的错误属实应该牢记。正确做法是 p=p->succ
