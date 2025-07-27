---
title: "二、Go类型系统"
date: 2024-11-25T01:47:00Z
draft: false
weight: 2000
---

# 二、Go类型系统

- 计算机底层全是bit，而**实际操作则是基于大小固定的单元中的数值，称为字（word），****如整数、浮点数、比特数组、内存地址****等**；进而构成更大的聚合类型；Go的数据类型宽泛，向下匹配硬件特性，向上满足程序员所需；
- Go语言将数据类型分为四类：
    | **基础类型（basic type）** | number、boolean、string，是Go语言世界数据的原子。 | 
    | --- | --- | 
    | **组合类型（aggregate type）** | array**、**struct，由基础类型组合，值由内存中的一组变量构成，是Go语言世界数据的分子。 | 
    | **引用类型（reference type）** | pointer、slice、map、function、**channel**，都**是间接指向****程序变量或状态，**操作所引用数据的全部效果会遍及该数据的全部引用； | 
    | **接口类型（interface type）** | interface | 


## 📚 目录

- [2.1 number、boolean](2-1-number-boolean/)
- [2.2 string](2-2-string/)
- [2.4 array、slice、append](2-4-array-slice-append/)
- [2.5 pointer](2-5-pointer/)
- [2.6 map、make、new](2-6-map-make-new/)
- [2.7 泛型 generics](2-7-泛型-generics/)
- [2.8 struct、组合嵌套、json](2-8-struct-组合嵌套-json/)
- [2.9 func declaration、递归、多值返回、可变参数](2-9-func-declaration-递归-多值返回-可变参数/)
- [2.10 函数变量、匿名函数、闭包](2-10-函数变量-匿名函数-闭包/)
- [2.11 error、panic、recover、defer ](2-11-error-panic-recover-defer-/)
- [2.12 method declaration、指针接收者、struct组合嵌套、封装](2-12-method-declaration-指针接收者-struct组合嵌套-封装/)
- [2.13 interface type、实现、interface变量的值、interface的类型断言](2-13-interface-type-实现-interface变量的值-interface的类型断言/)
- [2.14 常用interface](2-14-常用interface/)

























