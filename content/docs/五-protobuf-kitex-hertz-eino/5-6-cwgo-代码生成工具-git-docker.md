---
title: "5.6 cwgo 代码生成工具、git、docker"
date: 2025-05-24T02:24:00Z
draft: false
weight: 5006
---

# 5.6 cwgo 代码生成工具、git、docker

- 官方文档：[http://www.cloudwego.io/zh/docs/cwgo/overview/](http://www.cloudwego.io/zh/docs/cwgo/overview/)
- 安装：
    ```shell
    go install github.com/cloudwego/cwgo@latest
    GO111MODULE=on go install github.com/cloudwego/thriftgo@lates
    brew install protobuf
    ```
- 命令行的命令自动补全：[http://www.cloudwego.io/zh/docs/cwgo/tutorials/auto-completion/](http://www.cloudwego.io/zh/docs/cwgo/tutorials/auto-completion/)
    ![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/7df6e8b3-5a82-4d0d-bd83-df7758c91c97/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4667P6TO2F3%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005411Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIQCj5CFkG8F4DX8AnAmgOW5Sipqdwn2VdrWrlIFxXMvWRgIgW1fUIA%2F0TnRFk7VzdWz7QymRIDSDTyrUA%2FFPkkjrZ88qiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDF%2FhM4RZjI3gc7dsNCrcA7IQpm7AKj%2F0nHo64fqma81SjkK4tJVQdszv0NJXRbDQu%2BdBCo%2FncfzinS%2ByB6isZY8QpUe32MVhRoNazgcxY66dj18AWJXclRub%2Ba8WGfMOUg540OiLoRd36hjDzMs2%2BqgI9XPSQiluA73MDNaKjcpUs5gUSv%2F1WkMKcICLyV2pamDw%2B4hjzDTUVUQYf4S5syBxSd5YeZDd%2BBkMEKYdavmE%2BoERMuLUpinPQP348RI5MYWrhgeGvCdNnSGNqSa%2FvTrz0sY%2FgtKCBz6sNnnzeDPCB%2BHrMl4rDbtZxMhlKAc%2FYp8LsZqC2HsP64N%2FGOv4bdj35m87dEhXbKWMQ90zRMa81SsMK0rr7tVNKuqBj%2Bu90bAAu%2B2G0jH%2F%2F9dP%2By%2BAHWdkdxDgLLrhNOuhTJ%2BAbo0mgXRvy947VO30f%2Bg2mxnbg3n5gzkJopacMSWjmcV5kv2E%2BuXFvx95pQW01%2FLAEIGITQoLfYk0VYyg76OjU9HQd82Fg9Xz%2F%2BPBBKIIg%2FB%2FUE1jIiLj5%2FrilrkFD9Q0xkpTWG%2FgyZigbGrqo%2BBn34rNryEJYcoYeVBTsVH%2FScT2degvfzsgLdSExZZ2nyHuCdYINgcEa8U7U4tcUfhJ9iF0xrLjuh%2BTzfWnsNTMMJK668MGOqUBTykgSSWun%2Bq%2BKLLOz06QeF5efw%2BXM5xNA%2BbL8XV9Ber8ErdGsySocIeB%2FbPs3IDxgZMz9T4DqmlVNV8XDIWBoF%2BxFm%2B%2BIAC6zTOsVpd9JSMQm1dg2qN6YRb1fpy1NbhQUaRe3Kmz8V2TOTFGAKb1h7atjIAa9j6EWrfReynnnifhKWjST4j4g2wmxOaEdBESlPzajZmXVtVu468hMg%2FhrWb%2BwLiQ&X-Amz-Signature=56f4616f37762b7f7df79d20f48c59209a1916424c3cd2825b049d67e902dda2&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)
# cwgo脚手架: 自动生成项目结构layout

- thrift版：
    ```shell
     cwgo server --type RPC --module github.com/cloudwego/biz-demp/gomall/demo/demo_thrift --service demo_thrift --idl ../../idl/echo.thrift
     
     go mod tidy
     go work use .
     
     go run .
    ```
    ![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/972e5086-42ab-4809-b8b8-e31ff12bf5fe/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466UZDTAK2U%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005412Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIFf%2FWXB%2FM68C14tPGSgbU4vboOvdyoqc5SDQaY87Qs9PAiEAhMIO%2F2%2FPO%2B5nZ5dgKxGkWK5fNrZO6Dw0URIAP4dF51sqiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDDZAxfs202KHCMVnJCrcA4LxgonR3mqYjIESn6oNuxgy2CCJlnXKLDCpdio%2FCm4xqVc13ISp0hTyAT4PCA3lJ%2FSF%2BaZRLXOKyqtYWrof7723sc7vZuHExhuczR1Xc6Alg0iRVvPRVmF9qXwwBsU5sd527UgzNK3IFTNux8VDqXjdbNzeXKv18vqLbLFNYUQ8stv%2F3Wmyo2la2zaPGRThABGB8Wsdrvpob029UtA%2BM3Qx1yIS8SLb%2BMJpntpkZATtMd%2FDej5BvNax%2FpQoP3doTxKfVSny7Pqt7RJogUfLNigs9KXagcTX4hOV6p1WbUQRW3ZYndSRCQ67C3%2BOZV5mUP0GpiBUuiyYuPVV4AV6TSQUPvR4P0W%2BEMwU2KikAnmTHLd3bdc9hiqQIdM4l56HXyUA69KsxbuHJy7eip2vODJ9kjmnwHHMNsMHOr3zv%2F8QtH3usHX8qGc7Nw2b4fgXBQuGxJ4vWpBKL7nhehBvBL4749Dl4KgX8VOP9%2FfL7NhES6w7ZLY79Wt5o99PCtbMC5P8y9QChu%2F3A88RMJSCqzx9LiAE0bTVDP2mFO%2BPZNR4t73KGxvgjyAp7Jy0NKd9rkzJwtFLjhAx3qjRPT27O12jkgNrc19kYEvWxuKqE9fmT9nv%2FVuAVbaTg0x7MMC668MGOqUBAwGZFGw7h%2Feq7wPQF8yks4arQmBkYGLuCvjFU72%2B5F8C5FQ9WceCnMt8eCbwcEH3bsB2t163ByOpma4wPOf%2FPPvlLot6NE5X0JiPl1MByDVQHe7YNJOeM%2FexReasugJ%2FubV1VlOtCH52QqANkZypCtHVGzNajvywOqDe2PvnbI5%2FigYDXOngglZBFv5VRVP3xPGgGtEMSpymONfLVTvQkqyFt9MO&X-Amz-Signature=89e4dc93563622bcfdc730f07f92cf62ef9e92b93b491e05a0b4c699657f0acd&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)




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












