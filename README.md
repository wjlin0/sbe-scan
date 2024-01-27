<h4 align="center"> sbe-scan 是一个用Go语言编写的 SpringBoot ENV 利用工具，可以用来获取 SpringBoot 应用的配置信息，包括配置文件中的配置项、环境变量、JVM参数等。 </h4>
<p align="center">
<a href="https://github.com/wjlin0/sbe-scan/releases/"><img src="https://img.shields.io/github/release/wjlin0/sbe-scan" alt=""></a> 
<a href="https://github.com/wjlin0/sbe-scan" ><img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/wjlin0/sbe-scan"></a>
<a href="https://github.com/wjlin0/sbe-scan/releases"><img src="https://img.shields.io/github/downloads/wjlin0/sbe-scan/total" alt=""></a> 
<a href="https://github.com/wjlin0/sbe-scan"><img src="https://img.shields.io/github/last-commit/wjlin0/sbe-scan" alt=""></a> 
<a href="https://blog.wjlin0.com/"><img src="https://img.shields.io/badge/wjlin0-blog-green" alt=""></a>
</p>

# 安装

sbe-scan 需要`go 1.21`才能完成安装 执行以下命令

```shell
go install github.com/wjlin0/sbe-scan/cmd/sbe-scan@latest
```
或者
安装完成的二进制文件在[release](https://github.com/wjlin0/sbe-scan/releases)中下载
- [macOS-arm64](https://github.com/wjlin0/sbe-scan/releases/download/v0.0.3/sbe-scan_0.0.3_macOS_arm64.zip)

- [macOS-amd64](https://github.com/wjlin0/sbe-scan/releases/download/v0.0.3/sbe-scan_0.0.3_macOS_amd64.zip)

- [linux-amd64](https://github.com/wjlin0/sbe-scan/releases/download/v0.0.3/sbe-scan_0.0.3_linux_amd64.zip)

- [windows-amd64](https://github.com/wjlin0/sbe-scan/releases/download/v0.0.3/sbe-scan_0.0.3_windows_amd64.zip)

- [windows-386](https://github.com/wjlin0/sbe-scan/releases/download/v0.0.3/sbe-scan_0.0.3_windows_386.zip)


# 使用
```shell
sbe-scan -help
```
```text
sbe-scan is a tool to scan spring boot env.

Usage:
  sbe-scan [flags]

Flags:
INPUT:
   -url, -u string[]  URL to scan
   -list string[]     File containing list of URLs to scan

OUTPUT:
   -o, -output string  Directory to save results (default "output/")

CONFIG:
   -eu, -env-url string[]            URL to get env
   -ju, -jolokia-url string[]        URL to get jolokia
   -jlu, -jolokia-list-url string[]  URL to get jolokia list
   -en, -env-name string[]           env name to get env
   -m, -method string[]              method to get env (support methods one)
   -header string[]                  Headers to use for enumeration

LIMIT:
   -timeout int          time to wait in seconds before timeout (default 10)
   -t, -thread int       Number of concurrent threads (default 10) (default 10)
   -rl, -rate-limit int  Rate limit for enumeration speed (n req/sec)

DEBUG:
   -debug                           Enable debugging
   -p, -proxy string[]              list of http/socks5 proxy to use (comma separated or file input)
   -irt, -input-read-timeout value  timeout on input read (default 3m0s)
   -no-stdin                        disable stdin processing

UPDATE:
   -update  Update tool


Examples:
Run sbe-scan on a single targets
        $ sbe-scan -url https://example.com
Run sbe-scan on a list of targets
        $ sbe-scan -list list.txt
Run sbe-scan on a single targets with env-url
        $ sbe-scan -url https://example.com -eu /actuator/env
Run sbe-scan on a single targets with jolokia-list-url
        $ sbe-scan -url https://example.com -jlu /actuator/jolokia/list
Run sbe-scan on a single targets a proxy server
        $ sbe-scan -url https://example.com  -proxy http://127.0.0.1:7890
          
```

use pathScan to collect targets and pass them to sbe-scan via standard input

```shell
pathScan -ue quake -uq 'app:"springboot"' -uc -silent -ul 200 | sbe-scan
```
```text
➜  ~ pathScan -ue quake -uq 'app:"springboot"' -uc -silent -ul 200 | sbe-scan

            __
   _____   / /_   ___           _____  _____  ____ _   ____
  / ___/  / __ \ / _ \ ______  / ___/ / ___/ / __  /  / __ \
 (__  )  / /_/ //  __//_____/ (__  ) / /__  / /_/ /  / / / /
/____/  /_.___/ \___/        /____/  \___/  \__,_/  /_/ /_/

				wjlin0.com

慎用。你要为自己的行为负责
开发者不承担任何责任，也不对任何误用或损坏负责.
[INF] Current sbe-scan version v0.0.3 (latest)
[INF] Loaded 109 targets from input
[INF] Running all methods
[INF] find url https://example.com  write to /Users/wjl/output/5051ebf7c0761402e8d072c9a50cb0f1.application.json
```
