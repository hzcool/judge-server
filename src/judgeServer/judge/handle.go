package judge

import (
	"strconv"
	"fmt"
	// "io"
	"net/http"
	// "bytes"
	"encoding/json"
	"sync"
	"io/ioutil"
	"os"
	"path/filepath"
)


type env_config struct {
	tmp_dir string
	wg sync.WaitGroup
	mutex sync.Mutex
}

var (
	ch chan *env_config
)

//初始化
func Init() {
	ch = make(chan *env_config,MAX_TASKS)
	for i:=1;i<=MAX_TASKS;i++ {
		dir := filepath.Join(TEMP_FILE_DIR,strconv.Itoa(i))
		os.MkdirAll(dir,os.ModePerm)
		ch <- &env_config {
			tmp_dir:dir,
		}
	}
}

//删除一个文件夹下的所有文件
func RemoveContents(dir string) error {
    d,err := os.Open(dir)
    if err != nil {
        return err
    }
    defer d.Close()
    names,err := d.Readdirnames(-1)
    if err != nil {
        return err
    }
    for _,name := range names {
        err = os.RemoveAll(filepath.Join(dir,name))
        if err != nil {
            return err
        }
    }
    return nil
}
 

//检查头信息和请求类型
func check(w http.ResponseWriter, req *http.Request) string {	
	if req.Method != "POST" {
		return "Invalid Method"
	}
	contentTypes := req.Header["Content-Type"]
	if contentTypes==nil || contentTypes[0]!="application/json" {
		return "Invalid Content-Type"
	} 
	tokens := req.Header["Access-Token"]
	if(tokens==nil || tokens[0]!=TOKEN) {
		return "Wrong Access-Token"
	}
	return ""
}

//检查请求信息
func checkJsonBody(post *map[string]interface{}) string {
	if _,ok := (*post)["lang"].(string); !ok {
		return "no lang field or wrong lang type"
	} 
	if _,ok := (*post)["max_cpu_time"].(float64); !ok {
		return "no max_cpu_time field or wrong max_cpu_time type"
	}
	if _,ok := (*post)["max_memory"].(float64); !ok {
		return "no max_memory field or wrong max_memory type"
	}
	if _,ok := (*post)["test_case"].(string); !ok {
		return "no test_case field or wrong test_case type"
	}
	if _,ok := (*post)["src"].(string); !ok {
		return "no src field or wrong src type"
	}
	return ""
}

//检查测试数据配置的info文件
func checkJsonInfo(info *map[string]interface{}) string {
	// fmt.Println(reflect.TypeOf((*info)["spj"]))
	if _,ok := (*info)["spj"].(bool); !ok {
		return "no spj field or wrong spj type"
	}
	if _,ok := (*info)["test_cases"].(map[string]interface{}); !ok {
		return "no test_cases"
	}
	if (*info)["spj"].(bool) {
		//特判题
		if spj_lang,ok := (*info)["spj_lang"].(string); !ok {
			return "no spj_lang field or wrong spj_lang type"
		} else if spj_lang!="C" && spj_lang != "C++"{
			return "the spj_lang is not supported yet"
		}
	} 
	return ""
}


func Ping(w http.ResponseWriter, req *http.Request) {
	if errInfo:=check(w,req); errInfo!="" {
		responseJson(w, &errResponse{INVALID_HEADER,errInfo})
		return
	}
	data, _ := ioutil.ReadFile(filepath.Join(BASE_PATH,"/judgeServer/judge/ping.txt"))
	responseJson(w, &errResponse{OK,string(data)})
}


// "{src,lang,max_cpu_time,max_memory,test_case}"
func Judge(w http.ResponseWriter, req *http.Request) {
	_env := <-ch
	defer func() {
		_ = RemoveContents(_env.tmp_dir)
		ch <- _env
		if err:=recover(); err!=nil {
			responseJson(w, &errResponse{SYSTEM_EXCEPTION, fmt.Sprintf("%s", err)})
		}
	}()
	if errInfo:=check(w,req); errInfo!="" {
		responseJson(w, &errResponse{INVALID_HEADER,errInfo})
		return
	}
	
	post,err1 := decodeJson(req.Body)
	if err1 != nil {
		responseJson(w, &errResponse{INVALID_JSON_REQUEST,err1.Error()})
		return
	}

	if errInfo:=checkJsonBody(&post); errInfo!="" {
		responseJson(w, &errResponse{INVALID_HEADER,errInfo})
		return
	}

	
	var l lang_config
	var ok bool

	
	l.lang = post["lang"].(string)
	l.cmd, ok = language[l.lang]
	if (!ok) {
		responseJson(w, &errResponse{INVALID_JSON_REQUEST,"this language is not supported yet"});
		return
	}
	l.src_path = filepath.Join(_env.tmp_dir,l.cmd["src_name"])
	l.exe_path = filepath.Join(_env.tmp_dir,l.cmd["exe_name"])
	//保存用户代码
	if err2:=ioutil.WriteFile(l.src_path,[]byte(post["src"].(string)),0777); err2!=nil {
		responseJson(w, &errResponse{STORE_SRC_FAILED,err2.Error()})
		return
	}
	
	//编译
	compile_info,success := compile(l.getCompileCommand())
	if !success {
		responseJson(w, &errResponse{COMPILE_ERROR,compile_info})
		return
	}


	l.max_cpu_time = int(post["max_cpu_time"].(float64))
	if l.max_cpu_time<CPU_TIME_MIN || l.max_cpu_time > CPU_TIME_MAX {
		responseJson(w, &errResponse{INVALID_CONFIG,"invalid max_cpu_time"})
		return
	}
	l.max_real_time = 2*l.max_cpu_time
	l.max_memory = int(post["max_memory"].(float64))
	if l.max_memory<MEMORY_MIN || l.max_memory > MEMORY_MAX {
		responseJson(w, &errResponse{INVALID_CONFIG,"invalid max_memory"})
		return
	}
	

	buf ,err3 := ioutil.ReadFile(filepath.Join(TEST_CASE_DIR,post["test_case"].(string),"info"))
	if err3!=nil {
		responseJson(w, &errResponse{INVALID_JSON_REQUEST,"the test_case not found"})
		return
	}

	var info map[string]interface{}
	if err4 := json.Unmarshal(buf,&info); err4 != nil {
		responseJson(w, &errResponse{INVAILD_TESTCASE_INFO_FILE,"invalid info file"})
	}
	
	if errInfo:=checkJsonInfo(&info); errInfo!="" {
		responseJson(w, &errResponse{INVAILD_TESTCASE_INFO_FILE,errInfo})
		return
	}
	
	cases := info["test_cases"].(map[string]interface {})
	case_dir := filepath.Join(TEST_CASE_DIR,post["test_case"].(string))
	
	l.spj = info["spj"].(bool)
	if l.spj { 
		l.spj_lang = info["spj_lang"].(string)
		l.spj_exe_path = filepath.Join(case_dir,spj_language[l.spj_lang]["exe_name"])
		if _,err:=os.Stat(l.spj_exe_path); err!=nil { 
			l.spj_src_path = filepath.Join(case_dir,spj_language[l.spj_lang]["src_name"])
			//编译spj_src文件
			fmt.Println("haha")
			if spj_compile_info,success2 := compile(l.getSpjCompileCommand()); !success2 {
				responseJson(w, &errResponse{COMPILE_ERROR,"special judge file:\n"+spj_compile_info})
				return
			}
		}
	}
	
	resp := &judgeResponse{
		Compile_info : compile_info,
		Total : 0,
		Pass : 0,
		Cpu_time : 0,
		Memory : 0,
		Total_cpu_time : 0,
		Err : OK,
	}

	//跑样例
	for id,v := range cases {
		mp := v.(map[string]interface {})
		_id,_ := strconv.Atoi(id)
		task := &case_config {
			input_data_path : filepath.Join(case_dir,mp["input_name"].(string)),
			output_path : filepath.Join(_env.tmp_dir,"result" + id),
			test_case_id : _id,
			max_output_size : DEFAULT_OUTPUT_SIZE,
			l : &l,
		}
		if !l.spj {
			task.output_data_path = filepath.Join(case_dir, mp["output_name"].(string))
			task.stripped_output_md5 = mp["stripped_output_md5"].(string)
			if osize,ok := mp["output_size"].(float64); ok {
				task.max_output_size = 2*int(osize)
			}
		} 
		_env.wg.Add(1)
		go run_one_task(task,resp,_env)
	}
	_env.wg.Wait()
	responseJson(w,resp)
}

