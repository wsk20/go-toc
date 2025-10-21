// Package main
// Author: wsk20
// Created on: 2025-10-21 15:10:09
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// 命令行参数定义
	inputFile := flag.String("i", "input.md", "输入 Markdown 文件")                 // -i 指定输入文件，默认 input.md
	outputFile := flag.String("o", "output.md", "输出 Markdown 文件")               // -o 指定输出文件，默认 output.md
	levels := flag.String("levels", "1-6", "生成 TOC 的标题级别范围，例如 1-4")             // -levels 指定 TOC 标题级别范围
	stdout := flag.Bool("stdout", false, "是否直接输出 TOC 到控制台而不写文件")                // -stdout 如果为 true，直接输出 TOC 到终端
	modifyTitle := flag.Bool("modify-title", false, "是否修改正文重复标题，自动序号化为 -2, -3") // -modify-title 对重复标题进行自动序号化
	flag.Parse()                                                                // 解析命令行参数

	// 解析标题级别范围
	levelParts := strings.Split(*levels, "-")
	minLevel, maxLevel := 1, 6 // 默认标题级别范围
	if len(levelParts) == 2 {
		if l, err := strconv.Atoi(levelParts[0]); err == nil {
			if l > 6 {
				minLevel = 1
			} else {
				minLevel = l
			}
		}
		if l, err := strconv.Atoi(levelParts[1]); err == nil {
			if l > 6 || l == 1 {
				maxLevel = 6
			} else {
				maxLevel = l
			}
		}
		if minLevel >= maxLevel {
			minLevel = 1
			maxLevel = 6
		}
	}

	fName := *inputFile
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		fmt.Println("输入文件不存在，请指定正确文件或查看帮助信息：")
		flag.Usage()
		return
	}

	// 打开输入文件
	f, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer f.Close() // 确保文件在函数退出时关闭

	scanner := bufio.NewScanner(f) // 按行读取文件

	// 初始化变量
	lines := []string{}            // 存储修改后的正文内容
	toc := []string{}              // 存储 TOC 条目
	inCodeBlock := false           // 是否在代码块中
	titleCount := map[string]int{} // 统计标题出现次数，用于重复标题编号

	// 逐行处理
	for scanner.Scan() {
		line := scanner.Text()

		// 检测代码块
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock // 遇到 ``` 切换代码块状态
			lines = append(lines, line)
			continue
		}

		trimmed := strings.TrimSpace(line) // 去掉首尾空格

		// 列表、块引用或代码块内的 # 不算标题
		if inCodeBlock || strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "> ") {
			lines = append(lines, line)
			continue
		}

		// 正则匹配 Markdown 标题
		headerRegex := regexp.MustCompile(`^(#{1,6})\s+(.*)`)
		// 知识点扩展：
		// 1. headerRegex.MatchString(trimmed) 只返回是否匹配布尔值
		// 2. headerRegex.FindStringSubmatch(line) 返回匹配结果数组：
		//    match[0] = 整行匹配，match[1] = # 数量，match[2] = 标题内容
		match := headerRegex.FindStringSubmatch(line)
		if match != nil {
			level := len(match[1]) // # 的数量即标题级别

			// 忽略不在 TOC 范围的标题
			if level < minLevel || level > maxLevel {
				lines = append(lines, line)
				continue
			}

			title := match[2]

			// 处理重复标题
			count := titleCount[title]
			if count > 0 {
				anchorSuffix := fmt.Sprintf("-%d", count+1) // 重复标题加后缀
				if *modifyTitle {
					title = title + anchorSuffix
					line = fmt.Sprintf("%s %s", match[1], title) // 修改正文标题
				}
				// 生成锚点，空格替换为 -, 全部小写
				anchor := strings.ReplaceAll(strings.ToLower(title), " ", "-")
				toc = append(toc, fmt.Sprintf("%s- [%s](#%s)", strings.Repeat("  ", level-minLevel), title, anchor))
			} else {
				anchor := strings.ReplaceAll(strings.ToLower(title), " ", "-")
				toc = append(toc, fmt.Sprintf("%s- [%s](#%s)", strings.Repeat("  ", level-minLevel), title, anchor))
			}
			titleCount[match[2]]++ // 记录标题出现次数

			lines = append(lines, line)
		} else {
			// 普通行
			lines = append(lines, line)
		}
	}

	// 拼接 TOC 字符串
	tocStr := "# 目录\n\n" + strings.Join(toc, "\n") + "\n"

	// 如果使用 stdout 参数直接输出
	if *stdout {
		fmt.Print(tocStr)
		return
	}

	// 写入输出文件
	out, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("无法创建输出文件:", err)
		return
	}
	defer out.Close()

	out.WriteString(tocStr)
	for _, line := range lines {
		out.WriteString(line + "\n")
	}

	fmt.Println("TOC 已生成到", *outputFile)
	if *modifyTitle {
		fmt.Println("正文重复标题已自动序号化为 -2, -3，以匹配 TOC")
	}
}
