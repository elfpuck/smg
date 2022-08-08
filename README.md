# SMG
> 什么鬼, 可配置化命令行工具, 支持 http、mysql、consul、exec等操作。

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
