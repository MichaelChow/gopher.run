---
title: "6.6 cwgo 代码生成工具、git、docker"
date: 2025-05-24T02:24:00Z
draft: false
weight: 6006
---

# 6.6 cwgo 代码生成工具、git、docker

- 官方文档：[http://www.cloudwego.io/zh/docs/cwgo/overview/](http://www.cloudwego.io/zh/docs/cwgo/overview/)
- 安装：
    ```shell
    go install github.com/cloudwego/cwgo@latest
    GO111MODULE=on go install github.com/cloudwego/thriftgo@lates
    brew install protobuf
    ```
- 命令行的命令自动补全：[http://www.cloudwego.io/zh/docs/cwgo/tutorials/auto-completion/](http://www.cloudwego.io/zh/docs/cwgo/tutorials/auto-completion/)
    ![](/images/1fd24637-29b5-80e9-9879-dc8e20c852ed/image_1fe24637-29b5-8085-989d-e6f4d37d741a.jpg)


# cwgo脚手架: 自动生成项目结构layout

- thrift版：
    ```shell
     cwgo server --type RPC --module github.com/cloudwego/biz-demp/gomall/demo/demo_thrift --service demo_thrift --idl ../../idl/echo.thrift
     
     go mod tidy
     go work use .
     
     go run .
    ```
    ![](/images/1fd24637-29b5-80e9-9879-dc8e20c852ed/image_1fe24637-29b5-80ae-a72b-fc22745ce1a6.jpg)




- makefile：一键代码生成
    ```makefile
    .PHONY: gen-demo-proto
    gen-demo-proto:
    	@cd demo/demo_proto && cwgo server -I ../../idl --type RPC --module github.com/cloudwego/biz-demp/gomall/demo/demo_proto --service demo_proto --idl ../../idl/echo.proto
    .PHONY: gen-demo-thrift
    gen-demo-thrift:
    	@cd demo/demo_thrift && cwgo server -I ../../idl --type RPC --module github.com/cloudwego/biz-demp/gomall/demo/demo_thrift --service demo_thrift --idl ../../idl/echo.thrift
    ```
    - 运行makefile：
        ```shell
        make gen-demo-thrift
        ```
- IDL拓展
    - hertZ对IDL的拓展：[http://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/annotation/](http://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/annotation/)
        - 字段注解: 标记字段从query、header、cookie、form等获取
        - 方法注解：get、post、put 
        - client注解












