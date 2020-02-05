# judge-server  

一个用于搭建online-judger的测评核心服务模块，实现语言是go。  
支持的语言：  
+ C11  (gcc7.4.0)
+ C++11,C++14,C++17 (g++7.4.0)
+ Python2 (Python2.7.17)
+ Python3 (Python3.6.9)
+ java (javac13.0.2)  

# 快速搭建  
+ 1. 依赖工具  
 docker和docker-compose， 版本不要太旧即可，可自行尝试安装
+ 2. git clone 该项目到你的某个文件夹内  
    `git clone https://github.com/hzcool/judge-server.git /usr/local/dest`
+ 3. 进入到该目录内，输入启动命令即可  
    `docker-compose up -d`  
+ 4. 测试，输入  
`curl -H "Content-Type:application/json" -H "Access-Token:123" -X POST http://0.0.0.0:8000/ping`   
出现返回的语言配置信息，即启动成功

# 使用方法  
## 1.头部信息，必须指定两个头部:  
`"Content-Type:application/json"`  
`"Access-Token:123"`    
`Access-Token`是连接的验证码，这个可以在`docker-compose.yml`进行修改，调用的web方法全部使用POST 

## 2.API  
> ## `ping`
> ### URL
> + `http://0.0.0.0:8000/ping`  
> ### 参数
> + 不需要参数  
> 
> 
> ### 返回值  
> + 返回值是个json数据  
> ``` json
> {
>      "err" : 0, 
>       "info" : "
>           1.编译选项 : [C , C++ , C++14 , C++17 , Python2 , Python3 , Java] 
>
>           2. 编译器版本 
>           C : gcc7.4.0
>           C++ : g++7.4.0
>           Python2 : Python2.7.14
>           Python3 : Python3.6.9
>           Java : Java13.0.2    
>       "
> }
> ```
  
  
> ## `judge`
> ### URL  
>+ `http://0.0.0.0:8000/judge`
>
>### 参数
>+ 需要上传一个json数据如下
>```golang
>{
>       "lang" : "C++",   //支持的语言是ping返回值中编译选项的其中之一即可
>       "max_cpu_time" : 1000 , //运行的最大时间限制
>       "max_memory" : 128*1024*1024, //最大内存限制
>       "test_case" : "normal_problem", //测试的题目号
>       "src" : "code",  //src是测试的源代码 
>}
>```
> 
> ### 返回值
> + 1.出现编译错误、上传json错误、系统错误等，返回值格式如下  
> ```golang
>   {
>       "err" : 4 , //4对应的是编译错误
>       "info" : "reason" //错误原因
>   }
>``` 
>  
>+ 2.编译成功时，err域为0
>```golang
>{
>       "err" : 0,
>       "compile_info" : "waring", //返回编译的警告信息
>       "result" : 0 , //结果，0表示全部正确
>       "total" : 10 , //测试样例总数
>       "pass" : 10, // 通过的样例数
>       "cpu_time" : 100, // 运行时间，取决于所有测试样例中耗时最长的，单位ms
>       "memory" : 20*1024*1024, // 运行内存，取决于所有测试样例中消耗内存最大的，单位字节
>       "total_cpu_time" : 1000, //所有样例的运行时间之和
>       "cases" : [   //所有的样例结果，组成一个json对象数组
>               {  
>                   "id" : 1, // 测试样例编号
>                   "cpu_time" : 100, //该样例的运行时间
>                   "memory" : 20*1024*1024, //该样例的内存消耗
>                   "result" : 0, //该样例的运行结果
>                   "info" : "39b494a1de04ddb493c47bcdc6f67e76", //运行信息，如果result为0，即正确的情况下是用户输出结果的md5，否则是错误信息
>                 }
>               ]
>}
>```

# 错误信息
+ ### 返回值中的**err**域，可能的结果如下
```golang
type ErrType int 
const (
	OK ErrType = 0   //没有错误
	INVALID_HEADER ErrType = 1 // 请求头错误
 	INVALID_JSON_REQUEST ErrType = 2 // 请求的json数据不对
	STORE_SRC_FAILED ErrType = 3 // 系统保存上传的测试代码出错
	COMPILE_ERROR ErrType = 4  // 编译错误
	INVALID_CONFIG ErrType = 5 // 时间、内存的设置不在合理范围内
	INVAILD_TESTCASE_INFO_FILE ErrType = 6 //测试数据的info描述文件不对
	SYSTEM_EXCEPTION ErrType = 7 //系统异常
)
```

+ ### 测试样例的**result**域可能的结果如下
```golang
const (
	ACCEPTED int = 0         
	WRONG_ANSWER int = 1    
	TIME_LIMIT_EXCEEDED int = 2
    MEMORY_LIMIT_EXCEEDED int = 3
    OUTPUT_LIMIT_EXCEEDED int = 4
    RUNTIME_ERROR int = 5
    SYSTEM_ERROR int = 6
)
```




  




    

