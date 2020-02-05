package judge

//测评语言
var (
	language  =  map[string]map[string]string {
		"C" : {
			"src_name" : "main.c",
			"exe_name" : "main",
			"exe_file" : "{exe_path}",
			"compile_command" : "gcc -std=c11 -O2 -lm {src_path} -o {exe_path} 2>&1",
			"run_command" : "{exe_path}", 		
		},

		"C++" : {
			"src_name" : "main.cpp",
			"exe_name" : "main",
			"exe_file" : "{exe_path}",
			"compile_command" : "g++ -std=c++11 -O2 -lm {src_path} -o {exe_path} 2>&1",
			"run_command" : "{exe_path}", 		
		},
		
		"C++14" : {
			"src_name" : "main.cpp",
			"exe_name" : "main",
			"exe_file" : "{exe_path}",
			"compile_command" : "g++ -std=c++14 -O2 -lm {src_path} -o {exe_path} 2>&1",
			"run_command" : "{exe_path}", 		
		},
		
		"C++17" :{
			"src_name" : "main.cpp",
			"exe_name" : "main",
			"exe_file" : "{exe_path}",
			"compile_command" : "g++ -std=c++17 -O2 -lm {src_path} -o {exe_path} 2>&1",
			"run_command" : "{exe_path}", 		
		},

		"Python2" : {
			"src_name" : "solution.py",
			"exe_name" : "solution.pyc",
			"exe_file" : "/usr/bin/python2",
			"compile_command" : "python2 -m py_compile {src_path}",
			"run_command" : "python2 {exe_path}",
		},

		"Python3" : {
			"src_name" : "solution.py",
			"exe_name" : "__pycache__/solution.cpython-36.pyc",
			"exe_file" : "/usr/bin/python3",
			"compile_command" : "python3 -m py_compile {src_path}",
			"run_command" : "python3 {exe_path}",
		},

		"Java" : {
			"src_name" : "Main.java",
			"exe_name" : "",
			"exe_file" : "/usr/bin/java",
			"compile_command" : "javac -encoding UTF8 {src_path} -d {exe_path} 2>&1",
			"run_command" : "java -cp {exe_path} -Djava.security.manager -Dfile.encoding=UTF-8 -Djava.awt.headless=true -XX:+UseSerialGC -Xss64m Main",
		},
	} 

	//spj-judge
	spj_language = map[string]map[string]string {
		"C" : {
			"src_name" : "spj.c",
			"exe_name" : "spj",
			"compile_command" : "gcc -std=c11 -O2 -lm {spj_src_path} -o {spj_exe_path}",
			"run_command" : "{spj_exe_path} {output} {data_input}", //执行文件，用户输出, 测试数据输入
		},
		"C++" : {  //默认17
			"src_name" : "spj.cpp",
			"exe_name" : "spj",
			"compile_command" : "g++ -std=c++17 -O2 -lm {spj_src_path} -o {spj_exe_path}",
			"run_command" : "{spj_exe_path} {output} {data_input}", //执行文件，用户输出, 测试数据输入
		},
		"Python2" : {
			"src_name" : "spj.py",
			"exe_name" : "spj.pyc",
			"compile_command" : "python2 -m py_compile {spj_src_path}",
			"run_command" : "python2 {spj_exe_path} {output} {data_input}", //执行文件，用户输出, 测试数据输入
		},
		"Python3" : {
			"src_name" : "spj.py",
			"exe_name" : "__pycache__/spj.cpython-36.pyc",
			"compile_command" : "python3 -m py_compile {spj_src_path}",
			"run_command" : "python3 {spj_exe_path} {output} {data_input}", //执行文件，用户输出, 测试数据输入
		},
	}
)