#include <unistd.h>
#include <sys/types.h>
#include <sys/resource.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <cstdlib>
#include <cstdio>
#include <seccomp.h>
#include <thread>
#include <string>
#include <fcntl.h>
#include <ctime>
#include <fstream>
#include <sstream>
#include <cstring>

#include "config.h"
#include "md5.h"

//去除字符串尾部字符
void Trim(string &str) {
    for(int i=str.size()-1;i>=0;i--) {
        if(!isblank(str[i])) {
            str.erase(i+1);
            break;
        } 
    }
}

//读取文件内容，并去除每行末尾的空字符
bool get_content_trimed(string& file_path,string& content) {
    ifstream file(file_path.c_str());
    if(!file) return false;
    string line;
    while(getline(file,line)) {
        Trim(line);
        content += line;
    }
    file.close();
    return true;
}

//系统调用限制
bool rules_init(char* exec_file) {
    int syscalls_blacklist[] = {SCMP_SYS(clone),
                            SCMP_SYS(fork), SCMP_SYS(vfork),
                            SCMP_SYS(kill), 
                            #ifdef __NR_execveat
                            SCMP_SYS(execveat)
                            #endif
    };
    int syscalls_blacklist_length = sizeof(syscalls_blacklist) / sizeof(int);
    scmp_filter_ctx ctx = seccomp_init(SCMP_ACT_ALLOW);
    if(!ctx) {
        // puts("初始化限制系统调用出错");
        return false;
    }

    for (int i = 0; i < syscalls_blacklist_length; i++) {
        if (seccomp_rule_add(ctx, SCMP_ACT_KILL, syscalls_blacklist[i], 0) != 0) {
            // puts("限制命令出错");
            return false;
        }
    }

    seccomp_rule_add(ctx,SCMP_ACT_KILL,SCMP_SYS(execve),1,SCMP_A0(SCMP_CMP_NE,(scmp_datum_t)exec_file));
    
    // use SCMP_ACT_KILL for socket, python will be killed immediately
    if (seccomp_rule_add(ctx, SCMP_ACT_ERRNO(EACCES), SCMP_SYS(socket), 0) != 0) {
        return false;
    }

    if (seccomp_load(ctx) != 0) {
        return false;
    }
    return true;
}

//资源限制
bool resource_init(runner_config *runner) {

    rlim_t cpu_time = runner->max_cpu_time/1000 + 1;
    if(setrlimit(RLIMIT_CPU,new rlimit{cpu_time,cpu_time+1})) {
        // puts("限制cpu时间错误");
        return false;
    } 

    rlim_t output_size = runner->max_output_size;
    if(setrlimit(RLIMIT_FSIZE,new rlimit{output_size,output_size})) {
        // puts("限制输出文件大小错误");
        return false;
    }
    
   
    

    if(runner->lang != "Java") {
        rlim_t memory = runner->max_memory;
        if(setrlimit(RLIMIT_AS,new rlimit{memory,memory})) {
            // puts("限制空间错误");
            return false;
        }

        rlim_t stack_size = 128*1024*1024;
        if(setrlimit(RLIMIT_STACK,new rlimit{stack_size,stack_size})) {
            // puts("限制栈大小错误");
            return false;
        }
    }

    
    return true;
}


//子进程监视，超过real_time就kill
void watch(pid_t child,int sec) {
    sleep(sec);
    // puts("child_process killed");
    kill(child,SIGXCPU);
}


//执行用户代码
void child_process(runner_config *runner) {
    if(!freopen(runner->input_path.c_str(),"r",stdin)) exit(open_input_file_error);
    if(!freopen(runner->output_path.c_str(),"w",stdout)) exit(open_output_file_error);
    if(!resource_init(runner)) exit(resource_limit_error);
    char* exec_file = (char*)runner->exe_path.c_str();
    if(runner->lang != "Java") {
        if(!rules_init(exec_file)) exit(seccomp_init_error);
    }
    stringstream ss(runner->run_command);
    char* argv[20]; int k = 0;
    string tmp;
    while(ss>>tmp) {
        argv[k] = new char[tmp.size()+1];
        strcpy(argv[k],(char*)tmp.c_str());
        k++;
    }
    argv[k] = NULL;
    char* envp[] = {NULL};
    execve(exec_file,argv,envp);
    exit(system_call_error);
}

void run(runner_config *runner,result_config *result)
{
    pid_t pid = fork();
    if(pid==0) {
        child_process(runner);
    } else {
        thread monitor(watch,pid,runner->max_real_time/1000);
        monitor.detach();
        int status;
        struct rusage resource_usage;
        wait4(pid,&status,0,&resource_usage);
        // printf("%d %d %d\n",status,WEXITSTATUS(status),WIFEXITED(status));
        

        result->cpu_time = max((int)(resource_usage.ru_utime.tv_sec * 1000 +
                                resource_usage.ru_utime.tv_usec / 1000)-2,0);
        result->memory = resource_usage.ru_maxrss * 1024;
        // printf("time memory : %d %d\n",result->cpu_time,result->memory);

        int retStatus = WEXITSTATUS(status);
        if(retStatus >= 20) {  
            result->error(retStatus);
            return;
        }
        
        if(retStatus == RUNTIME_ERROR) {  
            result->status = RUNTIME_ERROR;
            return;
        }

        int signal =  WTERMSIG(status);
        if(signal == SIGXFSZ) {
            result->status = OUTPUT_LIMIT_EXCEEDED;
            return;
        }
        if(signal == SIGXCPU || result->cpu_time > runner->max_cpu_time) {
            result->status = TIME_LIMIT_EXCEEDED;
            if(result->cpu_time < runner->max_cpu_time) 
                result->cpu_time = runner->max_real_time;
            return;
        }
        if(result->memory > runner->max_memory) {
            result->status = MEMORY_LIMIT_EXCEEDED;
            return;
        }
        if(!WIFEXITED(status)) {  //非正常退出
            result->status = RUNTIME_ERROR;
            return;
        }
        string content;
        if(!get_content_trimed(runner->output_path,content)) {
            result->error(get_result_file_content_error);
            return;
        }
        result->info = md5(content);
    }
}


// g++ -std=c++11 -O3 -lm runner.cpp md5.cpp -o runner -lseccomp -lpthread  
int main(int argc,char* argv[]) {
    auto runner = new runner_config();
    auto result = new result_config();
    if(!runner->parse(argc,argv)) {
        result->error(unlegal_args_error);
    } else {
        // runner->show();
        run(runner,result);
    }
    puts(result->get_json().c_str());
    return 0;
}