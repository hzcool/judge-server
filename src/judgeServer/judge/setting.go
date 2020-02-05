package judge

import (
	"os"
	"path/filepath"
)

var (
	TOKEN = os.Getenv("ACCESS_TOKEN")  //安全验证码
	SERVICE_PORT = os.Getenv("SERVICE_PORT")  //服务网址

	BASE_PATH =  filepath.Join(os.Getenv("GOPATH"),"src")
	TEST_CASE_DIR = filepath.Join(BASE_PATH,"test_case") //测试样例所在文件夹
	TEMP_FILE_DIR = filepath.Join(BASE_PATH,"tmp")  //代码测试时生成的临时文件存放的文件夹
	RUNNER_PATH = filepath.Join(BASE_PATH,"runner/runner") //跑用户代码的运行器路径
	LOG_DIR_PATH = filepath.Join(BASE_PATH,"log") //日志路径
)

const (
	MAX_TASKS = 6 //同一时刻的测试用户上限
	CPU_TIME_MIN = 100 //单样例最小时间设置0.1s
	CPU_TIME_MAX = 20000 //单样例最大时间设置20s
	MEMORY_MIN = 20*1024*1024 //单样例最小内存 20MB
	MEMORY_MAX = 1024*1024*1024 //单样例最大内存 1GB
	DEFAULT_OUTPUT_SIZE = 20*1024*1024  //默认文件输出大小上限,非spj题是标准答案文件大小的两倍
)



//错误信息
type ErrType int 
const (
	OK ErrType = 0
	INVALID_HEADER ErrType = 1
	INVALID_JSON_REQUEST ErrType = 2
	STORE_SRC_FAILED ErrType = 3
	COMPILE_ERROR ErrType = 4
	INVALID_CONFIG ErrType = 5
	INVAILD_TESTCASE_INFO_FILE ErrType = 6
	SYSTEM_EXCEPTION ErrType = 7
)

//结果信息
const (
	ACCEPTED int = 0
	WRONG_ANSWER int = 1
	TIME_LIMIT_EXCEEDED int = 2
    MEMORY_LIMIT_EXCEEDED int = 3
    OUTPUT_LIMIT_EXCEEDED int = 4
    RUNTIME_ERROR int = 5
    SYSTEM_ERROR int = 6
)
