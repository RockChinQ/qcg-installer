package main

import (
	"archive/zip"
	"bufio"
	"crypto/tls"
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
)

func DownloadFile(file_url, save_path, proxy string) string {
	//检查save_path是否存在，不存在则创建
	if _, err := os.Stat(save_path); os.IsNotExist(err) {
		os.Mkdir(save_path, os.ModePerm)
	}

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
		log.Fatal(err)
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
		log.Fatal(err)
	}

	// 判断get url的状态码, StatusOK = 200
	if resp.StatusCode == http.StatusOK {
		log.Printf("[INFO] 正在下载: [%s]", fileName)
		fmt.Print("\n")

		downFile, err := os.Create(save_path + "/" + fileName)
		if err != nil {
			log.Fatal(err)
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
		return save_path + "/" + fileName
	} else {
		fmt.Print("\n")
		log.Printf("[ERROR] [%s]下载失败,%s.", fileName, resp.Status)
		return ""
	}
}

//解压
func DeCompress(zipFile, dest string) (string, error) {
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
	fmt.Println(cmd.Args)

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
		if ch == "\r" || ch == "\n" {
			fmt.Print("[" + task_label + "] ")
		}
	}

	cmd.Wait()
	fmt.Print("\033[m")
	return result, nil
}

func RunCMDTillStringOutput(task_label, dir, ending, cmd_str string, args ...string) (string, error) {
	cmd := exec.Command(cmd_str, args...)
	cmd.Dir = dir
	//显示运行的命令
	fmt.Println(cmd.Args)

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
		if ch == "\r" || ch == "\n" {
			fmt.Print("[" + task_label + "] ")
		}

		result += ch
		if strings.Contains(result, ending) {
			cmd.Process.Kill()
			break
		}
	}

	cmd.Wait()
	fmt.Print("\033[m")
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
