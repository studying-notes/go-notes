---
date: 2022-09-11T15:24:18+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 语言调度器源码分析之线程本地储存"  # 文章标题
url:  "posts/go/docs/internal/goroutine/scheduler/tls"  # 设置网页永久链接
tags: [ "Go", "tls" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 简介

**线程本地存储又叫线程局部存储，其英文为 Thread Local Storage，简称 TLS**，简而言之就是**线程私有的全局变量**。

有过多线程编程的读者一定知道，普通的全局变量在多线程中是共享的，一个线程对其进行了修改，所有线程都可以看到这个修改，而线程私有的全局变量与普通全局变量不同，线程私有全局变量是线程的私有财产，每个线程都有自己的一份副本，某个线程对其所做的修改只会修改到自己的副本，并不会修改到其它线程的副本。

下面用例子来说明一下多线程共享全局变量以及线程私有全局变量之间的差异，并对 gcc 的线程本地存储做一个简单的分析。

## 普通的全局变量

```c++
#include <stdio.h>
#include <pthread.h>

int g = 0;  // 1，定义全局变量g并赋初值0

void *start(void *arg)
{
  printf("start, g[%p] : %d\n", &g, g); // 4，子线程中打印全局变量g的地址和值

  g++; // 5，修改全局变量

  return NULL;
}

int main(int argc, char *argv[])
{
  pthread_t tid;

  g = 100;  // 2，主线程给全局变量g赋值为100

  pthread_create(&tid, NULL, start, NULL); // 3， 创建子线程执行start()函数
  pthread_join(tid, NULL); // 6，等待子线程运行结束

  printf("main, g[%p] : %d\n", &g, g); // 7，打印全局变量g的地址和值

  return 0;
}
```

简单解释一下，这个程序在注释 1 的地方定义了一个全局变量 g 并设置其初值为 0，程序运行后主线程首先把 g 修改成了 100（注释 2），然后创建了一个子线程执行 start() 函数（注释 3），start() 函数先打印出 g 的值（注释 4）确定在子线程中可以看到主线程对 g 的修改，然后修改 g 的值（注释 5）后线程结束运行，主线程在注释 6 处等待子线程结束后，在注释 7 处打印 g 的值确定子线程对 g 的修改同样可以影响到主线程对 g 的读取。

编译并运行程序：

```bash
gcc thread.c -o thread -lpthread
```

```bash
./thread
```

```
start, g[0x601064] : 100
main, g[0x601064] : 101
```

从输出结果可以看出，全局变量 g 在两个线程中的地址都是一样的，任何一个线程都可以读取到另一个线程对全局变量 g 的修改，这实现了全局变量 g 的多个线程中的共享。

## 线程私有全局变量

了解了普通的全局变量之后我们再来看通过线程本地存储 (TLS) 实现的线程私有全局变量。这个程序与上面的程序几乎完全一样，唯一的差别就是在定义全局变量 g 时增加了 __thread 关键字，这样 g 就变成了线程私有全局变量了。

```c++
#include <stdio.h>
#include <pthread.h>

__thread int g = 0;  // 1，这里增加了 __thread 关键字，把 g 定义成私有的全局变量，每个线程都有一个 g 变量

void *start(void *arg)
{
  printf("start, g[%p] : %d\n", &g, g); // 4，打印本线程私有全局变量 g 的地址和值

  g++; // 5，修改本线程私有全局变量 g 的值

  return NULL;
}

int main(int argc, char *argv[])
{
  pthread_t tid;

  g = 100;  // 2，主线程给私有全局变量赋值为 100

  pthread_create(&tid, NULL, start, NULL); // 3，创建子线程执行 start() 函数
  pthread_join(tid, NULL);  // 6，等待子线程运行结束

  printf("main, g[%p] : %d\n", &g, g); // 7，打印主线程的私有全局变量 g 的地址和值

  return 0;
}
```

运行程序看一下效果：

```bash
gcc -g thread.c -o thread -lpthread
```

```bash
./thread
```

```
start, g[0x7f0181b046fc] : 0
main, g[0x7f01823076fc] : 100
```

从输出结果可以看出：首先，全局变量 g 在两个线程中的地址是不一样的；其次 main 函数对全局变量 g 赋的值并未影响到子线程中 g 的值，而子线程对 g 都做了修改，同样也没有影响到主线程中 g 的值，这个结果正是我们所期望的，这说明，每个线程都有一个自己私有的全局变量 g。

这看起来很神奇，明明2个线程都是用的同一个全局变量名来访问变量但却像在访问不同的变量一样。

下面我们就来分析一下 gcc 到底使用了什么黑魔法实现了这个特性。对于像这种由编译器实现的特性，我们怎么开始研究呢？最快最直接的方法就是使用调试工具来调试程序的运行，这里我们使用 gdb 来调试。

```bash
gdb ./thread
```

首先在源代码的第 20 行（对应到源代码中的 g = 100）处下一个断点，然后运行程序，程序停在了断点处，反汇编一下 main 函数：

```bash
(gdb) b thread.c:20
```

```
Breakpoint 1 at 0x400793: file thread.c, line 20.
```

```bash
(gdb) r
```

```
Starting program: /home/bobo/study/c/thread

Breakpoint 1, at thread.c:20
20g = 100;
```

```bash
(gdb) disass
```

```
Dump of assembler code for function main:
  0x0000000000400775 <+0>:push   %rbp
  0x0000000000400776 <+1>:mov   %rsp,%rbp
  0x0000000000400779 <+4>:sub   $0x20,%rsp
  0x000000000040077d <+8>:mov   %edi,-0x14(%rbp)
  0x0000000000400780 <+11>:mov   %rsi,-0x20(%rbp)
  0x0000000000400784 <+15>:mov   %fs:0x28,%rax
  0x000000000040078d <+24>:mov   %rax,-0x8(%rbp)
  0x0000000000400791 <+28>:xor   %eax,%eax
=> 0x0000000000400793 <+30>:movl   $0x64,%fs:0xfffffffffffffffc
  0x000000000040079f <+42>:lea   -0x10(%rbp),%rax
  0x00000000004007a3 <+46>:mov   $0x0,%ecx
  0x00000000004007a8 <+51>:mov   $0x400736,%edx
  0x00000000004007ad <+56>:mov   $0x0,%esi
  0x00000000004007b2 <+61>:mov   %rax,%rdi
  0x00000000004007b5 <+64>:callq 0x4005e0 <pthread_create@plt>
  0x00000000004007ba <+69>:mov   -0x10(%rbp),%rax
  0x00000000004007be <+73>:mov   $0x0,%esi
  0x00000000004007c3 <+78>:mov   %rax,%rdi
  0x00000000004007c6 <+81>:callq 0x400620 <pthread_join@plt>
  0x00000000004007cb <+86>:mov   %fs:0xfffffffffffffffc,%eax
  0x00000000004007d3 <+94>:mov   %eax,%esi
  0x00000000004007d5 <+96>:mov   $0x4008df,%edi
  0x00000000004007da <+101>:mov   $0x0,%eax
  0x00000000004007df <+106>:callq 0x400600 <printf@plt>
  ......
```

程序停在了 g = 100 这一行，看一下汇编指令，

```
=> 0x0000000000400793 <+30>:movl   $0x64,%fs:0xfffffffffffffffc
```

这句汇编指令的意思是把常量 100(0x64) 复制到地址为 %fs:0xfffffffffffffffc 的内存中，可以看出全局变量 g 的地址为%fs:0xfffffffffffffffc，fs 是段寄存器，0xfffffffffffffffc 是有符号数 -4，所以全局变量 g 的地址为：

```
fs段基址 - 4
```

前面我们在讲段寄存器时说过段基址就是段的起始地址，为了验证 g 的地址确实是 fs 段基址 - 4，我们需要知道 fs 段基址是多少，虽然我们可以用 gdb 命令查看 fs 寄存器的值，但 fs 寄存器里面存放的是段选择子（segment selector）而不是该段的起始地址，为了拿到这个基地址，我们需要加一点代码来获取它，修改后的代码如下：

```c++
#include <asm/prctl.h>
#include <pthread.h>
#include <stdio.h>
#include <sys/prctl.h>
#include <unistd.h>

__thread int g = 0;

void print_fs_base() {
  unsigned long addr;
  int ret = arch_prctl(ARCH_GET_FS, &addr);//获取fs段基地址
  if (ret < 0) {
    perror("error");
    return;
  }

  printf("fs base addr: %p\n", (void *) addr);//打印fs段基址

  return;
}

void *start(void *arg) {
  print_fs_base();//子线程打印fs段基地址
  printf("start, g[%p] : %d\n", &g, g);

  g++;

  return NULL;
}

int main(int argc, char *argv[]) {
  pthread_t tid;

  g = 100;

  pthread_create(&tid, NULL, start, NULL);
  pthread_join(tid, NULL);

  print_fs_base();//main线程打印fs段基址
  printf("main, g[%p] : %d\n", &g, g);

  return 0;
}
```

代码中主线程和子线程都分别调用了print_fs_base() 函数用于打印 fs 段基地址，运行程序看一下：

```
fs base addr: 0x7f36757c8700
start, g[0x7f36757c86fc] : 0
fs base addr: 0x7f3675fcb700
main, g[0x7f3675fcb6fc] : 100
```

可以看到：

- 子线程fs段基地址为0x7f36757c8700，g的地址为0x7f36757c86fc，它正好是基地址-4
- 主线程fs段基地址为0x7f3675fcb700，g的地址为0x7f3675fcb6fc，它也是基地址-4

由此可以得出，gcc 编译器（其实还有线程库以及内核的支持）使用了 CPU 的 **fs 段寄存器来实现线程本地存储**，不同的线程中 fs 段基地址是不一样的，这样看似同一个全局变量但在不同线程中却拥有不同的内存地址，实现了线程私有的全局变量。
