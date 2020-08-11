# pmon2
golang进程管理工具(golang process manager)，专门用于 go 常驻进程管理 （daemon）

<img src="http://p0.qhimg.com/t017d6cbb68aed4b693.png" style="width:750px" />

## 如何使用？

#### 运行进程 [run/exec]

```bash
sudo pmon2 run [./二进制文件名] [参数]
```
启动进程，可以传入若干参数，参数说明如下：

```shell
--name 进程名称，如果不传递，则以二进制文件名作为默认名称
--log  指定进程运行日志，不指定则使用pmon默认进程保存日志路径
-- arg1 arg2 arg2 额外参数列表，多个参数以空格分隔
--user 指定进程启动的用户
```

示例：

```bash
sudo pmon2 run ./bin/gin -- -prjHome=`pwd` --user ntt360
```

#### 查看列表  [ list/ls ]

```bash
sudo pmon2 ls
```

#### 启动进程  [ start ]

```bash
sudo pmon2 start [id or name]
```

#### 停止进程  [ stop ]

```bash
sudo pmon2 stop [id or name]
```

#### 重启进程 [ reload ]

仅仅重启配置文件，该命令需要进程文件配合，reload命令仅仅发送 SIGUSR2 信号

```bash
pmon reload [id or name]
```

#### 删除进程  [ del/delete ]

```bash
pmon del [id or name]
```

#### 查看详情  [ show/desc ]

```bash
sudo pmon2 show [id or name]
```

## 如何安装？

使用 `yum` 直接安装即可， `pmon2` 目前托管于导航自有 `yum`源：http://yum.nav.qihoo.net:8360，请自行安装该 `yum` 源。

```bash
sudo yum install pmon2
```

### 如何运行？

`centos7` 目前使用 `systemd` 来管理 `pmon2` 进程，`centos6` 使用 `initctl` ，首次安装启用请执行：

```bash
# CentOS7
sudo systemctl start pmond

# CentOS6
sudo initctl start pmond
```

## 未来规划

1. log 命令
2. alarm对接oding报警

## QA?

Q: 日志切割？

A: pmon2 自带一个 `logrotate` 日志切割配置文件，会默认切割 `/var/log/pmon2` 目录下的日志文件，如果你是自定义日志路径，请自行实现日志切割。



Q: 进程启动参数必须传绝对路径？

A: 启动进程是，如果你传递的参数中存在路径，请使用绝对路径，`pmon2` 启动进程会新起一个新的沙河环境以避免环境变量污染问题。

