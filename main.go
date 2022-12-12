package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	osname, arch := determineEnvironment()
	println("OS:" + osname + " Arch:" + arch)

	proxyString := flag.String("p", "", "proxy string")
	flag.Parse()

	mcl_file := downloadMCLInstaller(osname, arch, *proxyString)
	installMCL(osname, arch, mcl_file, *proxyString)

	python_achive_file := downloadPython(osname, arch, *proxyString)
	installPython(osname, arch, python_achive_file, *proxyString)

	cloneSource()
	makeConfig()

	writeLaunchScript(osname, arch)
	println("安装完成!")
	println("请先运行run-mirai.bat登录qq号成功之后，保持运行状态，运行run-bot.bat")
}

// 确定OS和架构
func determineEnvironment() (osname, arch string) {
	return runtime.GOOS, runtime.GOARCH
}

func downloadPython(osname, arch, proxy string) string {
	python_url := ""
	if osname == "windows" {
		if arch == "386" {
			python_url = "https://www.python.org/ftp/python/3.10.9/python-3.10.9-embed-win32.zip"
		} else if arch == "amd64" {
			python_url = "https://www.python.org/ftp/python/3.10.9/python-3.10.9-embed-amd64.zip"
		} else {
			panic("不支持的操作系统、架构组合")
		}
	} else if osname == "linux" {
		python_url = "https://www.python.org/ftp/python/3.10.9/Python-3.10.9.tar.xz"
	}

	println("下载Python环境:" + python_url)
	return DownloadFile(python_url, "./python", proxy)
}

func installPython(osname, arch, achive_file, proxy string) {
	println("安装Python环境")
	if osname == "windows" {
		//解压归档文件
		DeCompress(achive_file, "./python/")
		//下载pip
		println("下载pip")
		pip_url := "https://bootstrap.pypa.io/get-pip.py"
		pip_file := DownloadFile(pip_url, "./python/", proxy)
		//安装pip
		println("安装pip")
		RunCMDPipe("安装pip", ".", "./python/python.exe ", pip_file)
		ReplaceStringInFile("./python/python310._pth", "#import site", "import site")

		//安装依赖
		println("安装依赖")
		RunCMDPipe("安装依赖", ".", "./python/Scripts/pip.exe ", "install", "pymysql", "yiri-mirai", "openai", "colorlog", "func_timeout")
		RunCMDPipe("安装依赖", ".", "./python/Scripts/pip.exe ", "install", "websockets", "--upgrade")

	}
}

func downloadMCLInstaller(osname, arch, proxy string) string {
	mcl_url := ""
	if osname == "windows" {
		if arch == "386" {
			mcl_url = "https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-windows-x86.exe"
		} else if arch == "amd64" {
			mcl_url = "https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-windows-amd64.exe"
		} else {
			panic("不支持的操作系统、架构组合")
		}
	} else if osname == "linux" {
		if arch == "386" {
			mcl_url = "https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-linux-amd64-musl"
		} else if arch == "amd64" {
			mcl_url = "https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-linux-amd64-musl"
		} else if arch == "arm" {
			mcl_url = "https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-linux-arm-musl"
		} else {
			panic("不支持的操作系统、架构组合")
		}
	}

	println("下载MCL安装器:" + mcl_url)
	return DownloadFile(mcl_url, "./mirai", proxy)
}

func installMCL(osname, arch, installer_file, proxy string) {
	println("安装mirai, 建议全部选项直接回车")
	installer_file = strings.ReplaceAll(installer_file, "mirai/", "")
	println(installer_file)
	if osname == "windows" {
		RunCMDPipe("安装mirai", "./mirai", installer_file)
	} else if osname == "linux" {
		RunCMDPipe("安装mirai", "chmod", "+x", installer_file)
		RunCMDPipe("安装mirai", "./mirai", installer_file)
	}
	RunCMDTillStringOutput("安装mirai", "./mirai", "I/main: mirai-console started successfully.", "./java/bin/java", "-jar", "mcl.jar")
	RunCMDPipe("安装mirai", "./mirai", "./java/bin/java", "-jar", "mcl.jar", "--update-package", "net.mamoe:mirai-api-http", "--channel", "stable-v2", "--type", "plugin")
	RunCMDTillStringOutput("安装mirai", "./mirai", "I/main: mirai-console started successfully.", "./java/bin/java", "-jar", "mcl.jar")
}

func cloneSource() {
	println("下载源代码")
	RunCMDPipe("下载源代码", ".", "git", "clone", "https://github.com/RockChinQ/QChatGPT")
}

func makeConfig() {
	println("生成配置文件")
	RunCMDPipe("生成配置文件", "./QChatGPT", "../python/python", "main.py")
	// RunCMDPipe("./QChatGPT", "../python/python", "main.py", "init_db")
	mirai_api_http_config := `adapters:
  - ws
debug: true
enableVerify: true
verifyKey: yirimirai
singleMode: false
cacheSize: 4096
adapterSettings:
  ws:
    host: localhost
    port: 8080
    reservedSyncId: -1`
	ioutil.WriteFile("./mirai/config/net.mamoe.mirai-api-http/setting.yml", []byte(mirai_api_http_config), 0644)

	println("=============================================")

	api_key := ""
	print("请输入OpenAI账号的api_key: ")
	fmt.Scanf("%s", &api_key)
	ReplaceStringInFile("./QChatGPT/config.py", "openai_api_key", api_key)

	qqn := 0
	print("请输入QQ号: ")
	fmt.Scanf("%d", &qqn)
	fmt.Scanf("%d", &qqn)
	ReplaceStringInFile("./QChatGPT/config.py", "1234567890", strconv.Itoa(qqn))
}

func writeLaunchScript(osname, arch string) {
	println("生成启动脚本")
	if osname == "windows" {
		ioutil.WriteFile("./run-mirai.bat", []byte(`cd mirai/
java\bin\java -jar mcl.jar`), 0644)
		ioutil.WriteFile("./run-bot.bat", []byte(`cd QChatGPT
..\python\python.exe main.py`), 0644)
	} else if osname == "linux" {
		ioutil.WriteFile("./run-mirai.sh", []byte(`cd mirai/
java/bin/java -jar mcl.jar`), 0644)
		ioutil.WriteFile("./run-bot.sh", []byte(`cd QChatGPT
../python/python main.py`), 0644)
	}
}
