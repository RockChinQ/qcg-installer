package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var start_time = int(time.Now().Unix() - 1670949827)

var version = "0.12"

func main() {
	showVersion := flag.Bool("v", false, "show version")

	proxyString := flag.String("p", "", "proxy string")
	flag.Parse()

	if *showVersion {
		fmt.Println("QChatGPT installer\nVersion: " + version)
		return
	}
	println(strconv.Itoa(start_time))
	osname, arch := determineEnvironment()
	println("OS:" + osname + " Arch:" + arch)

	go func() {
		resp, err := http.Get("http://rockchin.top:18989/report?osname=" + osname + "&arch=" + arch + "&timestamp=" + strconv.FormatInt(time.Now().Unix(), 10) + "&version=" + version + "&mac=0&message=start")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
	}()

	python_achive_file := downloadPython(osname, arch, *proxyString)
	installPython(osname, arch, python_achive_file, *proxyString)

	go func() {
		resp, err := http.Get("http://rockchin.top:18989/report?osname=" + osname + "&arch=" + arch + "&timestamp=" + strconv.FormatInt(time.Now().Unix(), 10) + "&version=" + version + "&mac=0&message=done_python")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
	}()

	mcl_file := downloadMCLInstaller(osname, arch, *proxyString)
	installMCL(osname, arch, mcl_file, *proxyString)

	go func() {
		resp, err := http.Get("http://rockchin.top:18989/report?osname=" + osname + "&arch=" + arch + "&timestamp=" + strconv.FormatInt(time.Now().Unix(), 10) + "&version=" + version + "&mac=0&message=done_mcl")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
	}()

	cloneSource()
	makeConfig(osname)

	writeLaunchScript(osname, arch)
	go func() {
		resp, err := http.Get("http://rockchin.top:18989/report?osname=" + osname + "&arch=" + arch + "&timestamp=" + strconv.FormatInt(time.Now().Unix(), 10) + "&version=" + version + "&mac=0&message=done_all")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
	}()
	println("===============????????????===============")
	if osname == "linux" {
		println("????????????run-mirai.sh??????qq?????????????????????????????????????????????run-bot.sh")
	} else if osname == "windows" {
		println("????????????run-mirai.bat??????qq?????????????????????????????????????????????run-bot.bat")
	}
	fmt.Printf("?????????????????????...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
	if osname == "windows" {
		os.Stdin.Read(b)
	}
}

// ??????OS?????????
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
			panic("???????????????????????????????????????")
		}
	} else if osname == "linux" {
		python_url = "https://www.python.org/ftp/python/3.10.9/Python-3.10.9.tgz"
	}

	return DownloadFileOrPrepared("Python????????????", python_url, "./python", proxy)
}

func installPython(osname, arch, achive_file, proxy string) {
	println("??????Python??????")
	if osname == "windows" {
		println(achive_file)
		//??????????????????
		_, err := DeCompress(achive_file, "./python/")
		if err != nil {
			panic(err)
		}
		//??????pip
		println("??????pip")
		pip_url := "https://bootstrap.pypa.io/get-pip.py"
		pip_file := DownloadFileOrPrepared("pip????????????", pip_url, "./python/", proxy)
		//??????pip
		println("??????pip")
		RunCMDPipe("??????pip", ".", "./python/python.exe ", pip_file)
		ReplaceStringInFile("./python/python310._pth", "#import site", "import site")

		//????????????
		println("????????????")
		RunCMDPipe("????????????", ".", "./python/Scripts/pip.exe ", "install", "yiri-mirai", "openai", "colorlog", "func_timeout", "-i", "http://pypi.douban.com/simple", "--trusted-host", "pypi.douban.com") //-i http://pypi.douban.com/simple --trusted-host pypi.douban.com
		RunCMDPipe("????????????", ".", "./python/Scripts/pip.exe ", "install", "websockets", "--upgrade", "-i", "http://pypi.douban.com/simple", "--trusted-host", "pypi.douban.com")
		RunCMDPipe("????????????", ".", "./python/Scripts/pip.exe ", "install", "dulwich", "-i", "http://pypi.douban.com/simple", "--trusted-host", "pypi.douban.com")

	} else if osname == "linux" {
		// DeCompress(achive_file,"./python/")
		RunCMDPipe("??????Python??????", ".", "tar", "zxvf", achive_file, "-C", "./python")
		linux_installerCompiler()
		pwd, _ := RunCMDPipe("??????pwd", "./python/", "pwd")
		pwd = strings.Trim(pwd, "\n")
		RunCMDPipe("??????????????????", "./python/Python-3.10.9", "./configure", "--prefix="+pwd)
		RunCMDPipe("??????Python", "./python/Python-3.10.9", "make")
		RunCMDPipe("??????Python", "./python/Python-3.10.9", "make", "install")

		println("????????????")
		RunCMDPipe("????????????", ".", "python/bin/pip3", "install", "yiri-mirai", "openai", "colorlog", "func_timeout", "-i", "http://pypi.douban.com/simple", "--trusted-host", "pypi.douban.com")
		RunCMDPipe("????????????", ".", "python/bin/pip3", "install", "websockets", "--upgrade", "-i", "http://pypi.douban.com/simple", "--trusted-host", "pypi.douban.com")
		RunCMDPipe("????????????", ".", "python/bin/pip3", "install", "dulwich", "-i", "http://pypi.douban.com/simple", "--trusted-host", "pypi.douban.com")
	}
}

func linux_installerCompiler() {

	result, _ := RunCMDPipe("??????????????????", ".", "apt")
	print(result)
	if result == "" {
		result, _ := RunCMDPipe("??????????????????", ".", "yum")
		if result == "" {
			fmt.Println("????????????Linux????????????")
			os.Exit(-1)
		} else {
			RunCMDPipe("??????????????????", ".", "yum", "install", "zlib-devel", "bzip2-devel", "openssl", "openssl-devel", "ncurses-devel", "sqlite-devel",
				"readline-devel", "tk-devel", "gcc", "make", "readline", "libffi-devel", "-y") //zlib-devel bzip2-devel openssl openssl-devel ncurses-devel sqlite-devel readline-devel tk-devel gcc make readline libffi-devel -y
		}
	} else {
		RunCMDPipe("??????????????????", ".", "apt", "update")
		RunCMDPipe("??????????????????", ".", "apt", "install", "build-essential", "zlib1g-dev", "libncurses5-dev", "libgdbm-dev", "libnss3-dev", "libssl-dev", "libreadline-dev", "libffi-dev", "libsqlite3-dev", "wget", "libbz2-dev")
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
			panic("???????????????????????????????????????")
		}
	} else if osname == "linux" {
		if arch == "386" {
			mcl_url = "https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-linux-amd64-musl"
		} else if arch == "amd64" {
			mcl_url = "https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-linux-amd64-musl"
		} else if arch == "arm" {
			mcl_url = "https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-linux-arm-musl"
		} else {
			panic("???????????????????????????????????????")
		}
	}

	println("??????MCL?????????:" + mcl_url)
	return DownloadFileOrPrepared("MCL?????????", mcl_url, "./mirai", proxy)
}

func installMCL(osname, arch, installer_file, proxy string) {
	println("??????mirai")
	installer_file = strings.ReplaceAll(installer_file, "mirai/", "")
	println(installer_file)
	if osname == "windows" {
		RunCMDPipe("??????mirai", "./mirai", installer_file)
	} else if osname == "linux" {
		RunCMDPipe("??????mirai", "./mirai", "chmod", "+x", installer_file)
		RunCMDPipe("??????mirai", "./mirai", installer_file)
	}

	RunCMDTillStringOutput("??????mirai", "./mirai", "I/main: mirai-console started successfully.", "./java/bin/java", "-jar", "mcl.jar")
	RunCMDPipe("??????mirai", "./mirai", "./java/bin/java", "-jar", "mcl.jar", "--update-package", "net.mamoe:mirai-api-http", "--channel", "stable-v2", "--type", "plugin")
	RunCMDTillStringOutput("??????mirai", "./mirai", "I/main: mirai-console started successfully.", "./java/bin/java", "-jar", "mcl.jar")

	//????????????
	ReplaceStringInFile("./mirai/config/Console/AutoLogin.yml", "protocol: ANDROID_PHONE", "protocol: ANDROID_PAD")
}

func cloneSource() {
	println("???????????????")
	GitClone("https://gitee.com/RockChin/QChatGPT", "./QChatGPT")
	// RunCMDPipe("???????????????", ".", "git", "clone", "https://gitee.com/RockChin/QChatGPT")
}

func makeConfig(osname string) {
	println("??????????????????")
	if osname == "linux" {
		RunCMDTillStringOutput("??????????????????", "./QChatGPT", "??????????????????", "../python/bin/python3", "main.py")
	} else {
		RunCMDTillStringOutput("??????????????????", "./QChatGPT", "??????????????????", "../python/python.exe", "main.py")
	}
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

	println("===================????????????==================")

	re := regexp.MustCompile("^sk-[a-zA-Z0-9]{48}$")
	for {
		input := InputString("?????????OpenAI?????????api_key: ")

		if re.MatchString(input) {
			ReplaceStringInFile("./QChatGPT/config.py", "openai_api_key", input)
			break
		} else if input != "" && input != "\n" {
			println("api_key????????????")
		}
	}

	qqn := 0
	print("?????????QQ???: ")
	if osname == "windows" {
		fmt.Scanf("%d", &qqn)
	}
	fmt.Scanf("%d", &qqn)
	ReplaceStringInFile("./QChatGPT/config.py", "1234567890", strconv.Itoa(qqn))
}

func writeLaunchScript(osname, arch string) {
	println("??????????????????")
	if osname == "windows" {
		ioutil.WriteFile("./run-mirai.bat", []byte(`cd mirai/
java\bin\java -jar mcl.jar
pause`), 0644)
		ioutil.WriteFile("./run-bot.bat", []byte(`cd QChatGPT
..\python\python.exe main.py
pause`), 0644)
	} else if osname == "linux" {
		ioutil.WriteFile("./run-mirai.sh", []byte(`cd mirai/
java/bin/java -jar mcl.jar`), 0644)
		RunCMDPipe("??????????????????", ".", "chmod", "+x", "run-mirai.sh")
		ioutil.WriteFile("./run-bot.sh", []byte(`cd QChatGPT
../python/bin/python3 main.py`), 0644)
		RunCMDPipe("??????????????????", ".", "chmod", "+x", "run-bot.sh")
	}
}
