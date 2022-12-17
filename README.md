# Installer for QChatGPT

为[QChatGPT项目](https://github.com/RockChinQ/QChatGPT)使用Go语言编写的一键部署脚本，自动化部署所需依赖  

- 注意：下载的Python和mirai均为免安装版，不影响系统其他环境

## 使用方法

- **部署过程中遇到任何问题，请先在[QChatGPT](https://github.com/RockChinQ/QChatGPT/issues)或[qcg-installer](https://github.com/RockChinQ/qcg-installer/issues)的issue里进行搜索，若找不到请前往：交流、答疑群: `204785790`**
    - **进群提问前请您`确保`已经找遍文档和issue均无法解决**
    - **进群提问前请您`确保`已经找遍文档和issue均无法解决**
    - **进群提问前请您`确保`已经找遍文档和issue均无法解决**

### 1. 注册OpenAI账号

参考以下文章

> [只需 1 元搞定 ChatGPT 注册](https://zhuanlan.zhihu.com/p/589470082)  
> [手把手教你如何注册ChatGPT，超级详细](https://guxiaobei.com/51461)

注册成功后请前往[个人中心](https://beta.openai.com/account/api-keys)查看`api_key`  

### 2. 安装器

- 从[Release页面](https://github.com/RockChinQ/qcg-installer/releases/latest)下载可执行文件，若无法访问请到[Gitee](https://gitee.com/RockChin/qcg-installer/releases/latest)   
- 保存到电脑上某个空目录，直接运行，等待配置环境
- 完毕后根据提示输入`api-key`和`QQ号`  
- 到此安装完成

**常见问题**

<details>
<summary>📵网络状况不好，下载失败？</summary>

解决方法:

- 若您有网络代理可用于提速，可在启动安装器时提供参数`-p <代理地址>`,如：
```
qcg-installer-0.1-windows-x64.exe -p http://localhost:7890
```

- 也可以提前下载所需文件，安装器运行中将不再进行下载，此功能适用于安装器版本`0.7`以上
    - Windows系统，下载以下文件并放置在安装器同目录，**请勿**重命名
        - [python-3.10.9-embed-amd64.zip](https://www.python.org/ftp/python/3.10.9/python-3.10.9-embed-amd64.zip)
        - [get-pip.py](https://bootstrap.pypa.io/get-pip.py)
        - [mcl-installer-a02f711-windows-amd64.exe](https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-windows-amd64.exe)
    - Linux系统，下载以下文件并放置在安装器同目录，**请勿**重命名
        - [Python-3.10.9.tgz](https://www.python.org/ftp/python/3.10.9/Python-3.10.9.tgz)
        - [get-pip.py](https://bootstrap.pypa.io/get-pip.py)
        - [mcl-installer-a02f711-linux-amd64-musl](https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-linux-amd64-musl)
</details>
    
### 3. 运行程序

之后每次重启之后均需要按照以下步骤启动程序

#### i. 启动mirai
- 运行`run-mirai.bat`(Windows) 或`./run-mirai.sh`(Linux) 启动mirai
- 并输入`login <QQ号> <QQ密码>`根据提示登录账号([登录教程](https://yiri-mirai.wybxc.cc/tutorials/01/configuration#4-%E7%99%BB%E5%BD%95-qq))

#### ii. 运行主程序

- 登录完成后运行`run-bot.bat`(Windows) 或 `./run-bot.sh`(Linux) 启动主程序  

**常见问题**

- mirai登录提示`QQ版本过低`，见[此issue](https://github.com/RockChinQ/QChatGPT/issues/38)
- 运行`run-bot.bat`闪退请见[此解决方案](https://github.com/RockChinQ/qcg-installer/issues/2)
- 若启动后提示安装`uvicorn`或`hypercorn`，请**不要**安装，会导致不明原因bug

## 目前支持的平台和架构

- Windows x64
- CentOS x64
    - 以及其他使用`yum`作为包管理器的操作系统
- Ubuntu x64
    - 以及其他使用`apt`作为包管理器的操作系统
- Raspbian arm64
