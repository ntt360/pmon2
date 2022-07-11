# pmon2 开发

## 环境构建

```shell
# 运行init_dev.sh 初始化开发配置文件
sh init_dev.sh
```
`init_dev.sh` 会负责生成项目的开发配置文件，以及相关目录、环境变量等。

## 测试开发

```shell
sudo PMON2_CONF=config/config-dev.yml ./bin/pmond
sudo PMON2_CONF=config/config-dev.yml ./bin/pmon2 exec bin/test 
```

因为 `pmon2` 启动进程使用的是fork/exec,所以需要sudo或root级别权限。
