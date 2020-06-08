## PL0 compiler

### 概述

该项目是编译原理课程的课程作业，为简化版pl0语言写一个编译器，生成中间代码并解释执行

- 题目要求和实验报告见 `docs` 文件夹
- 两个pl0程序示例见 `bin` 文件夹

### 使用方法

使用 go 1.13+ 版本都可以正常编译

#### 编译

```
go build -ldflags "-w -s" -trimpath -o pl0compiler.exe .
```

#### 运行

```
pl0compiler.exe pl0程序文件名

例如
pl0compiler.exe bin/maxif.txt
```
