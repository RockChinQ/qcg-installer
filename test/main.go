package main

import (
	"fmt"
	"regexp"
	"strings"
)

func InputString(prompt string) string {
	print(prompt)
	str := ""

	fmt.Scanf("%s", &str)

	if len(str) < 1 || str == "\n" || str == "" || str == " " || str == "\r" || str == "\t" || str == "\b" {
		fmt.Scanf("%s", &str)
	}

	return str
}

func main() {
	fmt.Println("test0")

	re, err := regexp.Compile("sk[-][a-zA-Z0-9]{48}")
	if err != nil {
		panic(err)
	}
	for {
		input := InputString("请输入OpenAI账号的api_key: ")
		input = strings.Trim(input, " \n\r\t\b")
		fmt.Println("input:", input)
		// println("bytes:", []byte(input)[0])
		if re.MatchString(input) {
			println("api_key格式正确")
			break
		} else if input != "" && input != "\n" {
			println("api_key格式错误")
		}
	}

	qqn := 0
	print("请输入QQ号: ")
	fmt.Scanf("%d", &qqn)
	fmt.Scanf("%d", &qqn)
}
