# pmon2
golang进程管理工具(golang process manager)，专门用于 go 常驻进程管理 （daemon）

## cli 命令

### 1. 启动进程（start）

```bash
pmon start [./二进制文件名] [参数]
```
启动进程，可以传入若干参数，参数说明如下：

```shell
--name 进程名称，如果不传递，则以二进制文件名作为默认名称
--log  指定进程运行日志，不指定则使用pmon默认进程保存日志路径
-- arg1 arg2 arg2 额外参数列表，多个参数以空格分隔
```

### 2. 查看进程列表（list）

```bash
pmon list
```

### 3. 停止进程

```bash
pmon stop [进程名称或进程ID]
```

### 4. 重启进程

```bash
pmon reload [进程名称或进程ID]
```

### 5. 删除进程

```bash
pmon del [进程名称或ID]
```

### 6. 日志操作

查看所有进程日志流

```bash
pmon log
```

查看单个进程日志流

```bash
pmon log [进程名称 or 进程id]
```
