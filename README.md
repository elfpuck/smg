# SMG
> 什么鬼, 可配置化命令行工具, 支持 http、mysql、consul、exec等操作。

# 配置文件
```

```
## smgEchoServer
> 返回请求的内容,用于排查与调试
```
smg system es 
```
## Variables
> 在操作过程中, 可以使用 全局变量、flags、args等数据
### 使用args
* 使用args[n] 参数
```
{{ index .args n }}
```
* args 合并 
```
{{join .args "sep" }}
```
## smg.yaml
```yaml
# 全局变量
variables:              
name:                      
version:                  
desc:
# 命令 map[string]command
command: 
  # 命令名称
  xxcommand:
    # 别名 []string
    aliases:
    usage:
    desc:
    # args参数描述
    argsUsage:
    # args 最小个数 int
    argsMin:
    category: 
    # 子命令 map[string]command
    subCommand:
    # flag参数 map[string]flag 
    flag:
      # flag 名称
      xxflag:
        aliases:
        usage:
        required:
        # 默认值
        value:
        # 使用环境变量 []string
        envVars:
    # 运行的命令
    action:
      # 命令类型 http、exec、mysql、consul
      type:
      exec:
        # 执行shell []string
        script: 
      mysql:
        # mysql 地址
        dsn:
        # 执行语句
        query:
      # consul 
      consul:
        # 输出文件,默认不填, std
        output:
        # 输出后json格式优化与取值,使用gjson实现
        resultPath:
        address:
        token:
        prefix:
        # consul 执行方法, GET、SET、LIST、DELETE
        method:
        key:
        value: 
      http:
        # 输出文件,默认不填, std
        output:
        # 输出后json格式优化与取值,使用gjson实现
        resultPath:
        method:
        url: 
        # map[string][]string 
        header:
          key: []
        data-raw:
        # []string 如: key=value
        data-urlencode:
          - key=value 
        # []string 如: file[0]=@"./smg.yaml"
        form-data:
          - file[0]=@"./smg.yaml"
```
