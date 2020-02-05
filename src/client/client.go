package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
)

type lang_config struct{
	Lang string `json:"lang"`
	Max_cpu_time int `json:"max_cpu_time"`
	Max_memory int `json:"max_memory"`
	Test_case string `json:"test_case"`
	Src string `json:"src"`
}


var (
	SERVICE_URL = os.Getenv("SERVICE_URL")
	BASE =  filepath.Join(os.Getenv("GOPATH"),"src")
	cpp_path = filepath.Join(BASE,"client/main.cpp")
	c_path = filepath.Join(BASE,"client/main.c")
	py3_path = filepath.Join(BASE,"client/main.py")
	py2_path = filepath.Join(BASE,"client/py2.py")
	java_path = filepath.Join(BASE,"client/Main.java")
	spj_cpp_path = filepath.Join(BASE,"client/spj_main.cpp")
)


var  c_config = lang_config {
	Lang:"C",
	Max_cpu_time: 1000,
	Max_memory: 134217728,
	Test_case: "normal_problem",
}

var  cpp_config = lang_config {
	Lang:"C++",
	Max_cpu_time: 1000,
	Max_memory: 134217728,
	Test_case: "normal_problem",
}



var  py3_config = lang_config {
	Lang:"Python3",
	Max_cpu_time: 1000,
	Max_memory: 134217728,
	Test_case: "normal_problem",
}

var  py2_config = lang_config {
	Lang:"Python2",
	Max_cpu_time: 1000,
	Max_memory: 134217728,
	Test_case: "normal_problem",
}

var java_config = lang_config {
	Lang:"Java",
	Max_cpu_time: 1000,
	Max_memory: 134217728,
	Test_case: "normal_problem",
}

var  spj_cpp_config = lang_config {
	Lang:"C++",
	Max_cpu_time: 1000,
	Max_memory: 134217728,
	Test_case: "spj_problem",
}

func ping() {
	client := &http.Client{}
	url := SERVICE_URL+"/ping"
	req,_ := http.NewRequest("POST",url,nil)


	
	req.Header.Add("Content-Type","application/json")
	req.Header.Add("Access-Token",os.Getenv("ACCESS_TOKEN"))


	res,_ := client.Do(req)
	defer res.Body.Close()
	
	var mp  map[string]interface{}
	json.NewDecoder(res.Body).Decode(&mp)
	fmt.Println("ping result : ")
	
	fmt.Println(mp["info"].(string))
}

func main()  {

	//ping
	ping()

	//C
	// src,_ := ioutil.ReadFile(c_path)
	// c_config.Src = string(src)
	// post,_ := json.Marshal(c_config)

	
	// C++
	// src,_ := ioutil.ReadFile(cpp_path)
	// cpp_config.Src = string(src)
	// cpp_config.Lang = "C++14"
	// post,_ := json.Marshal(cpp_config)

	//python2
	src,_ := ioutil.ReadFile(py2_path)
	py2_config.Src = string(src)
	post,_ := json.Marshal(py2_config)

	//python3
	// src,_ := ioutil.ReadFile(py3_path)
	// py3_config.Src = string(src)
	// post,_ := json.Marshal(py3_config)

	//Java
	// src,_ := ioutil.ReadFile(java_path)
	// java_config.Src = string(src)
	// post,_ := json.Marshal(java_config)

	//spj C++
	// src,_ := ioutil.ReadFile(spj_cpp_path)
	// spj_cpp_config.Src = string(src)
	// post,_ := json.Marshal(spj_cpp_config)


	buffer:= bytes.NewBuffer(post)
	client := &http.Client{}
	url := SERVICE_URL+"/judge"


	req,_ := http.NewRequest("POST",url,buffer)
	req.Header.Add("Content-Type","application/json")
	req.Header.Add("Access-Token",os.Getenv("ACCESS_TOKEN"))


	res,_ := client.Do(req)
	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Println("***********************************************")
	fmt.Println(string(resBody))
}