# Gopher技术栈
Go技术栈学习平台

## 关于

gopher.run是一个专注于Go技术栈的学习平台，提供高质量的教程、最佳实践和实用技巧。

## 特性

- 清晰的文档结构
- 代码示例
- 最佳实践指南
- 持续更新

## 本地开发

```bash
# 克隆仓库
git clone https://github.com/MichaelChow/gopher.run.git

# 进入项目目录
cd gopher.run

# 安装依赖
go mod download

hugo mod get -u

hugo mod tidy

hugo mod graph

hugo --gc

hugo --minify

# 启动开发服务器
hugo server -D
```

## 贡献

欢迎提交 Pull Request 或创建 Issue 来帮助改进这个项目。

## 许可证

MIT License
