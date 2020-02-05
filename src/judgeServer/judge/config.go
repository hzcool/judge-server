package judge

import (
	"strconv"
	"strings"
	// "fmt"
)





//测评的语言相关配置
type lang_config struct {
	lang string
	spj bool
	max_cpu_time int
	max_real_time int
	max_memory int 
	src_path string
	exe_path string
	
	cmd map[string]string
	
	//当spj为true时,需要下面3个信息
	spj_lang string //需要从testcase info文件读入
	spj_src_path string 
	spj_exe_path string
}



func (l *lang_config) getCompileCommand() string  {
	compile_command := l.cmd["compile_command"]
	return strings.Replace(strings.Replace(compile_command,"{src_path}",l.src_path,-1),"{exe_path}",l.exe_path,-1)
}

//特判程序的编译指令
func (l *lang_config) getSpjCompileCommand() string  {
	spj_cmd := spj_language[l.spj_lang]
	spj_compile_command := spj_cmd["compile_command"]
	return strings.Replace(strings.Replace(spj_compile_command,"{spj_src_path}",l.spj_src_path,-1),"{spj_exe_path}",l.spj_exe_path,-1)
}



//一个测试用例的所需信息，当测试题目为spj， output_data_path和stripped_output_md5不需要
type case_config struct{
	input_data_path string
	output_data_path string	
	stripped_output_md5 string
	output_path string
	max_output_size int
	test_case_id int 
	l *lang_config
}

func (c *case_config) getRunCommand() string {
	return RUNNER_PATH + 
		   " \"" +  c.l.lang + "\" " + 
		   "\"" + strings.Replace(c.l.cmd["run_command"],"{exe_path}",c.l.exe_path,-1) + "\" " + 
		   strconv.Itoa(c.l.max_cpu_time) + " " + 
		   strconv.Itoa(c.l.max_real_time) + " " +
		   strconv.Itoa(c.l.max_memory) + " " + 
		   strconv.Itoa(c.max_output_size) + " " +
		   strings.Replace(c.l.cmd["exe_file"],"{exe_path}",c.l.exe_path,-1) + " " + 
		   c.input_data_path + " " + 
		   c.output_path;
}

func (c *case_config) getSpjRunCommand() string {
	args := c.l.spj_exe_path + " " + c.output_path + " " + c.input_data_path
	switch c.l.spj_lang {
	case "Python2":
		return "python2 " + args
	case "Python3":
		return "python3 " + args
	}
	return args
}





//编译错误，系统错误时的返回信息
type errResponse struct {
	Err ErrType `json:"err"`
	Info string `json:"info"`
}


//单样例结果
type singleCase struct {
	Cpu_time int `json:"cpu_time"`  //运行时间
	Memory int `json:"memory"` //运行内存
	Result int `json:"result"` // 运行结果
	Info string `json:"info"` //运行信息
	Id int `json:"id"`       //测试样例号
}

//测评返回的结果
type judgeResponse struct{
	Compile_info string `json:"compile_info"`  //编译信息
	Total int `json:"total"`   //测试样例数
	Pass int `json:"pass"`     //通过的样例数量
	Result int `json:"result"` //测试结果
	Cpu_time int `json:"cpu_time"` //最大的样例时间
	Memory int `json:"memory"` //最大内存消耗
	Total_cpu_time int `json:"Total_cpu_time"` //总的cpu时间消耗
	Cases []*singleCase `json:"cases"` 
	Err ErrType `json:"err"` //此时的Err为OK 
}






