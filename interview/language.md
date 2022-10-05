---
date: 2022-09-27T11:26:20+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Golang 面试语言问题"  # 文章标题
url:  "posts/go/interview/language"  # 设置网页永久链接
tags: [ "Go", "language" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## select 是随机的还是顺序的

随机选择。

## 局部变量分配在栈还是堆

编译器会自动决定把一个变量放在栈还是放在堆，编译器会做逃逸分析，当发现变量的作用域没有跑出函数范围，就可以在栈上，反之则必须分配在堆。

## 协程调度原理

M(machine): 代表着真正的执行计算资源，可以认为它就是系统线程。

P(processor): 表示逻辑 processor，是线程 M 的执行的上下文。

G(goroutine): 调度系统的最基本单位 goroutine，存储了 goroutine 的执行 stack 信息、goroutine 状态以及 goroutine 的任务函数等。

## runtime 机制

runtime 负责管理任务调度，垃圾收集及运行环境。

## make 和 new 的区别

值类型：int，float，bool，string，struct 和 array。变量直接存储值，分配栈区的内存空间，这些变量所占据的空间在函数被调用完后会自动释放。

引用类型：slice，map，chan 和值类型对应的指针。变量存储的是一个地址（或者理解为指针），指针指向内存中真正存储数据的首地址。内存通常在堆上分配，通过 GC 回收。

对于引用类型的变量，我们不仅要声明变量，更重要的是，我们得手动为它分配空间。

new(T) 是为一个 T 类型的新值分配空间，并将此空间初始化为 T 的零值，并返回这块内存空间的地址，也就是 T 类型的指针 *T，该指针指向 T 类型值占用的那块内存。

make(T) 返回的是初始化之后的 T，且只能用于 slice，map，channel 三种类型。make(T, args) 返回初始化之后 T 类型的值，且此新值并不是 T 类型的零值，也不是 T 类型的指针 *T，而是 T 类型值经过初始化之后的引用。

## 安全读写共享变量的方式

- 互斥锁
- 管道 channel
- 原子操作

## 协程和线程和进程的区别

### 进程

进程是程序的一次执行过程，是程序在执行过程中的**分配和管理资源的基本单位**，每个进程都有自己的地址空间，**进程是系统进行资源分配和调度的一个独立单位**。

每个进程都有自己的独立内存空间，不同进程通过 IPC（Inter-Process Communication）进程间通信来通信。由于进程比较重量，占据独立的内存，所以上下文进程间的切换开销（栈、寄存器、虚拟内存、文件句柄等）比较大，但相对比较稳定安全。

### 线程

线程是进程的一个实体，线程是内核态，而且是 CPU 调度和分派的基本单位，它是比进程更小的能独立运行的基本单位。线程自己基本上不拥有系统资源，只拥有一点在运行中必不可少的资源(如程序计数器，一组寄存器和栈)，但是它可与同属一个进程的其他的线程共享进程所拥有的全部资源。

线程间通信主要通过**共享内存**，上下文切换很快，资源开销较少，但相比进程不够稳定容易丢失数据。

### 协程

协程是一种用户态的轻量级线程，协程的调度完全由用户控制。协程拥有自己的寄存器上下文和栈。

协程调度切换时，将寄存器上下文和栈保存到其他地方，在切回来的时候，恢复先前保存的寄存器上下文和栈，直接操作栈则基本没有内核切换的开销，可以不加锁的访问全局变量，所以上下文的切换非常快。

## CAS 无锁算法

CAS 算法（Compare And Swap）, 是原子操作的一种, CAS 算法是一种有名的无锁算法。

无锁编程，即不使用锁的情况下实现多线程之间的变量同步，也就是在没有线程被阻塞的情况下实现变量的同步，所以也叫非阻塞同步（Non-blocking Synchronization）。可用于在多线程编程中实现不被打断的数据交换操作，从而避免多线程同时改写某一数据时由于执行顺序不确定性以及中断的不可预知性产生的数据不一致问题。

该操作通过将内存中的值与指定数据进行比较，当数值一样时将内存中的数据替换为新的值。

Go 中的 CAS 操作是借用了 CPU 提供的原子性指令来实现。CAS 操作修改共享变量时候不需要对共享变量加锁，而是通过类似乐观锁的方式进行检查，本质还是不断的占用 CPU 资源换取加锁带来的开销（比如上下文切换开销）。