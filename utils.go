package main

import (
	"archive/zip"
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	pb "github.com/cheggaaa/pb/v3"
	gogit "github.com/go-git/go-git/v5"
)

func DownloadFileOrPrepared(file_nick, file_url, save_path, proxy string) string {
	//检查save_path是否存在，不存在则创建
	if _, err := os.Stat(save_path); os.IsNotExist(err) {
		os.Mkdir(save_path, os.ModePerm)
	}

	//检查目标文件是否已存在
	spt := strings.Split(file_url, "/")
	local_file_name := spt[len(spt)-1]

	if exist, _ := exists(save_path + "/" + local_file_name); exist {
		return save_path + "/" + local_file_name
	}
	if exist, _ := exists(local_file_name); exist {
		_, err := copy(local_file_name, save_path+"/"+local_file_name)
		if err != nil {
			panic(err)
		}
		return save_path + "/" + local_file_name
	}

	isLocal, fileName := DownloadFileWrapper(file_nick, file_url, save_path, proxy)
	if isLocal {
		_, err := copy(fileName, save_path+"/"+fileName)
		if err != nil {
			panic(err)
		}
		return save_path + "/" + fileName
	}
	return fileName
}

//返回是否是本地文件
func DownloadFileWrapper(file_nick, file_url, save_path, proxy string) (bool, string) {

	for {
		//TODO 提前检查是否已存在
		// input := InputString("是否自动下载" + file_nick + ",如您已提前下载并放置在本目录请输入n。(y/n):")
		input := "y"
		if input == "y" {
			fileName, err := DownloadFile(file_url, save_path, proxy)
			if err != nil {
				println("[ERR]" + file_nick + "下载失败,建议自行下载文件放置在本目录后输入n继续之后步骤。链接 " + file_url)
				for {
					input := InputString("是否尝试重新下载？(y/n)")
					if input == "n" {
						//TODO 文件名称检查
						spt := strings.Split(file_url, "/")
						return true, spt[len(spt)-1]
					} else if input == "y" {
						break
					} else {
						print("\n")
					}
				}
			} else {
				return false, fileName
			}
		} else if input == "n" {
			spt := strings.Split(file_url, "/")
			return true, spt[len(spt)-1]
		} else {
			print("\n")
		}
	}
}
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
func copy(src, dst string) (int64, error) {
	f, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func DownloadFile(file_url, save_path, proxy string) (string, error) {

	// 下载文件到path/并使用pb进度条
	// 解析/后的文件名字
	urlMap := strings.Split(file_url, "/")
	fileName := urlMap[len(urlMap)-1]

	// 解析带? = fileName 的文件名字
	if strings.Contains(fileName, "=") {
		splitName := strings.Split(fileName, "=")
		fileName = splitName[len(splitName)-1]
	}

	client := http.DefaultClient
	client.Timeout = time.Second * 60 * 10 //设置超时时间

	req, err := http.NewRequest(http.MethodGet, file_url, nil)
	if err != nil {
		return "", err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err == nil {
			tr.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	resp, err := (&http.Client{
		Transport: tr,
	}).Do(req)

	if err != nil {
		return "", err
	}

	// 判断get url的状态码, StatusOK = 200
	if resp.StatusCode == http.StatusOK {
		log.Printf("[INFO] 正在下载: [%s]", fileName)
		fmt.Print("\n")

		downFile, err := os.Create(save_path + "/" + fileName)
		if err != nil {
			return "", err
		}
		// 不要忘记关闭打开的文件.
		defer downFile.Close()

		// 获取下载文件的大小
		i, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
		sourceSiz := int64(i)
		source := resp.Body

		// 创建一个进度条
		bar := pb.New(int(sourceSiz)).SetRefreshRate(time.Millisecond*10).Set(pb.Bytes, true).SetWriter(os.Stdout)

		bar.SetWidth(80)

		bar.Start()

		barWriter := bar.NewProxyWriter(downFile)

		io.Copy(barWriter, source)
		bar.Finish()

		fmt.Print("\n")
		log.Printf("[INFO] [%s]下载成功.", fileName)
		return save_path + "/" + fileName, nil
	} else {
		fmt.Print("\n")
		log.Printf("[ERROR] [%s]下载失败,%s.", fileName, resp.Status)
		return "", errors.New("文件状态码不正确")
	}
}
func IsDir(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		log.Println(err)
		return false
	}
	return s.IsDir()
}
func CreateDir(dirName string) bool {
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
func EnsurePath(path string) {
	if !IsDir(path) {
		CreateDir(path)
	}
}

//解压
func DeCompress(zipFile, dest string) (string, error) {
	//检查目标目录是否存在，不存在则创建
	EnsurePath(dest)
	// 打开zip文件
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return "", err
	}

	defer func() {
		err := reader.Close()
		if err != nil {
			log.Fatalf("[unzip]: close reader %s", err.Error())
		}
	}()

	var (
		first string // 记录第一次的解压的名字
		order int    = 0
	)

	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return "", err
		}
		filename := filepath.Join(dest, file.Name)
		println("解压" + filename)
		//记录第一次的名字
		if order == 0 {
			first = filename
		}
		order += 1
		//fmt.Println(getDir(filename))
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(filename, 0755)
			if err != nil {
				return "", err
			}
		} else {
			w, err := os.Create(filename)
			if err != nil {
				return "", err
			}
			//defer w.Close()
			_, err = io.Copy(w, rc)
			if err != nil {
				return "", err
			}
			iErr := w.Close()
			if iErr != nil {
				log.Fatalf("[unzip]: close io %s", iErr.Error())
			}
			fErr := rc.Close()
			if fErr != nil {
				log.Fatalf("[unzip]: close io %s", fErr.Error())
			}
		}
	}
	return first, nil
}

func RunCMDPipe(task_label, dir, cmd_str string, args ...string) (string, error) {
	cmd := exec.Command(cmd_str, args...)
	cmd.Dir = dir
	//显示运行的命令
	fmt.Println("@"+dir, cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	result := ""
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadByte()
		if err2 != nil || io.EOF == err2 {
			break
		}

		ch := string(line)

		result += ch
		fmt.Print(ch)
		ts := int(time.Now().Unix()-1670949827) - start_time
		if ch == "\r" || ch == "\n" {
			fmt.Print("[" + strconv.Itoa(ts) + ":" + task_label + "] ")
		}
	}

	cmd.Wait()
	fmt.Print("\033[m")
	fmt.Println()
	return result, nil
}

func RunCMDTillStringOutput(task_label, dir, ending, cmd_str string, args ...string) (string, error) {
	cmd := exec.Command(cmd_str, args...)
	cmd.Dir = dir
	//显示运行的命令
	fmt.Println("@"+dir, cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)
	var result string
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadByte()
		if err2 != nil || io.EOF == err2 {
			break
		}

		ch := string(line)

		fmt.Print(ch)
		ts := int(time.Now().Unix()-1670949827) - start_time
		if ch == "\r" || ch == "\n" {
			fmt.Print("[" + strconv.Itoa(ts) + ":" + task_label + "] ")
		}

		result += ch
		if strings.Contains(result, ending) {
			cmd.Process.Kill()
			break
		}
	}

	cmd.Wait()
	fmt.Print("\033[m")
	fmt.Println()
	return result, nil
}

func ReplaceStringInFile(filename, oldStr, newStr string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	newContents := strings.Replace(string(bytes), oldStr, newStr, -1)
	err = ioutil.WriteFile(filename, []byte(newContents), 0664)
	return err
}

func GitClone(repo, dir string) error {
	_, err := gogit.PlainClone(dir, false, &gogit.CloneOptions{
		URL:      repo,
		Progress: os.Stdout,
	})
	return err
}

func InputString(prompt string) string {
	print(prompt)
	str := ""

	fmt.Scanf("%s", &str)

	if len(str) < 1 || str == "\n" || str == "" || str == " " || str == "\r" || str == "\t" || str == "\b" {
		fmt.Scanf("%s", &str)
	}

	return str
}
