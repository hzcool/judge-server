#ifndef CONFIG_H
#define CONFIG_H
#include <cstdio>
#include <string>
using namespace std;

//./a.out "C++" "./test" 1000 2000 134217728 20971520 ./test 1.in 1.out
struct runner_config{
    string lang;    
    string run_command;
    int max_cpu_time;
    int max_real_time;
    int max_memory;
    int max_output_size;
    string exe_path;  //可执行文件路径
    string input_path; // 测试数据输入文件路径
    string output_path; //结果输出的文件路径

    bool parse(int argc,char* argv[]) {
        if(argc < 10) return false;
        lang = argv[1];
        run_command = argv[2];
        max_cpu_time = stoi(argv[3]);
        max_real_time = stoi(argv[4]);
        max_memory = stoi(argv[5]);
        max_output_size = stoi(argv[6]);
        exe_path = argv[7];
        input_path = argv[8];
        output_path = argv[9];
        return true;
    }
    void show() {
        printf("language : %s\n",lang.c_str());
        printf("run_command : %s\n",run_command.c_str());
        printf("max_cpu_time : %d\n",max_cpu_time);
        printf("max_real_time : %d\n",max_real_time);
        printf("max_memory : %d\n",max_memory);
        printf("max_output_size : %d\n",max_output_size);
        printf("exe_path : %s\n",exe_path.c_str());
        printf("input_path : %s\n",input_path.c_str());
        printf("output_path : %s\n",output_path.c_str());
    }
};

//用户代码运行结果的状态
enum {
    OK = 0,
    TIME_LIMIT_EXCEEDED = 2,
    MEMORY_LIMIT_EXCEEDED = 3,
    OUTPUT_LIMIT_EXCEEDED = 4,
    RUNTIME_ERROR = 5,
    SYSTEM_ERROR = 6,
};
//系统错误时的具体错误
enum SYSTEM_ERROR_ITEM {
    unlegal_args_error = 20,
    resource_limit_error = 21,
    open_input_file_error = 22,
    open_output_file_error = 23,
    get_result_file_content_error = 24,
    seccomp_init_error = 25,
    system_call_error = 26,
    system_exceptions = 27,
};

struct result_config {
    int status;
    int cpu_time;
    int memory;
    //当status为OK时，info为测评输出结果的md5编码
    //当status为SYSTEM_ERROR时，info是具体的错误信息
    //
    string info;
    string get_json() {
        return "{\"status\":" + to_string(status) + " , " + 
        "\"info\":" + "\"" + info + "\" , " +
        "\"cpu_time\":" + to_string(cpu_time) + " , " +
        "\"memory\":" + to_string(memory) + "}" ; 
    }
    void init(int _status,string _info) {
        status = _status;
        info = _info;
    }
    void error(int sig) {
        status = SYSTEM_ERROR;
        switch (sig) {
            case unlegal_args_error:
                info = "传入参数错误";
                break;
            
            case resource_limit_error:
                info = "资源限制错误";
                break;
            case open_input_file_error:
                info = "打开输入文件错误";
                break;
            case open_output_file_error:
                info = "打开输出文件错误";
                break;
            case get_result_file_content_error:
                info = "获取结果文件内容出错";
                break;
            case seccomp_init_error:
                info = "系统调用限制错误";
                break; 
            case system_call_error:
                info = "系统调用错误";
                break;
            default:
                info = "系统异常";
                break;
        }
    }
};



#endif