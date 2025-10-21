# go-toc

`go-toc` 是一个用 Go 编写的 Markdown 文件 TOC（目录）生成工具。它可以自动扫描 Markdown 文件的标题，生成带锚点的目录，同时支持重复标题自动编号和控制输出方式。

---

## 功能特点

* 自动生成 Markdown 文件的 TOC（目录）
* 支持指定生成目录的最大标题级别（1~6）
* 支持重复标题自动序号化（如 `-2`, `-3`）
* 支持直接输出到终端或写入文件
* 保留原文档结构，不修改代码块、列表或引用内容
* 完全兼容 GitHub Markdown 锚点规则（空格转 `-`，大小写转换，小标点处理）

---

## 安装

### 方法 1：使用 Go 命令运行

```bash
go run main.go -i input.md -o output.md
```

### 方法 2：通过 `go install` 安装（推荐）

```bash
go install github.com/wsk20/go-toc@latest
```

然后在命令行中直接运行：

```bash
go-toc -i input.md -o output.md
```

---

## 使用说明

```text
Usage of go-toc:
  -i string
        输入 Markdown 文件 (default "input.md")
  -o string
        输出 Markdown 文件 (default "output.md")
  -levels int
        生成目录时包含的最大标题级别（1~6，默认6）
  -stdout
        是否直接输出 TOC 到控制台而不写文件
  -modify-title
        是否修改正文重复标题，自动序号化为 -2, -3
```

---

## 命令行示例

1. **生成 TOC 到文件**

```bash
go-toc -i input.md -o output.md
```

2. **直接输出 TOC 到控制台**

```bash
go-toc -i input.md -stdout
```

3. **只生成 1~6 级标题的 TOC**

```bash
go-toc -i input.md -o output.md -levels=3
```

4. **处理重复标题自动编号**

```bash
go-toc -i input.md -o output.md -modify-title
```

---

## 注意事项

* 默认输入文件为 `input.md`，如果不存在会提示帮助信息。
* 默认输出文件为 `output.md`，如果不指定 `-stdout` 参数，会生成新文件。
* 代码块内、列表、块引用内的 `#` 不会被当作标题。
* 重复标题会根据出现顺序自动添加编号，保持 TOC 与正文一致。
* 使用 `-levels` 可控制生成目录的最大标题深度（如 `-levels=3` 仅生成 `#`、`##`、`###`）。

---

## License

MIT License