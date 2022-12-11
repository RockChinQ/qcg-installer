package main

import (
	"flag"
	"runtime"
)

func main() {
	osname, arch := determineEnvironment()
	println("OS:" + osname + " Arch:" + arch)

	proxyString := flag.String("p", "", "proxy string")
	flag.Parse()

	achive_file := downloadPython(osname, arch, *proxyString)
	installPython(osname, arch, achive_file, *proxyString)
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
		RunCMDPipe("./python/python.exe ", pip_file)
		ReplaceStringInFile("./python/python310._pth", "#import site", "import site")

		//安装依赖
		println("安装依赖")
		RunCMDPipe("./python/Scripts/pip.exe ", "install", "pymysql", "yiri-mirai", "openai", "colorlog", "func_timeout")

		if arch == "386" {

		} else if arch == "amd64" {
			// TODO
		} else {
			panic("不支持的操作系统、架构组合")
		}
	}
}
