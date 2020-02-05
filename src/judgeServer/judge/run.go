package judge
import (
	// "fmt"
	"os/exec"
	"bytes"
	// "io/ioutil"
	"encoding/json"
	
)

func run(command string) string {
	cmd := exec.Command("/bin/bash","-c",command)
    var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	return out.String()
}

//特判，退出码为0时通过检验，否则错误
func spj_run(command string) bool {
	cmd := exec.Command("/bin/bash","-c",command)
	if err:=cmd.Run();err!=nil {
		return false
	}
	return true
}

func run_one_task(task *case_config,resp *judgeResponse, env *env_config) {

	s := run(task.getRunCommand())
	var mp map[string]interface{}
	json.Unmarshal([]byte(s),&mp)
	
	cs := &singleCase {
		Cpu_time : int(mp["cpu_time"].(float64)),
		Memory : int(mp["memory"].(float64)),
		Result : int(mp["status"].(float64)),
		Info : mp["info"].(string),
		Id : task.test_case_id,
	}

	if cs.Result == 0 {
		if task.l.spj {
			if !spj_run(task.getSpjRunCommand()) {
				cs.Result = WRONG_ANSWER
			}
		} else if cs.Info != task.stripped_output_md5 {
			cs.Result = WRONG_ANSWER
		}
	}

	env.mutex.Lock()
	resp.Total++
	if cs.Result == ACCEPTED {
		resp.Pass++
	} else if resp.Result < cs.Result {
		resp.Result = cs.Result
	}
	if resp.Cpu_time < cs.Cpu_time {
		resp.Cpu_time =  cs.Cpu_time
	}
	if resp.Memory < cs.Memory {
		resp.Memory =  cs.Memory
	} 
	resp.Total_cpu_time += cs.Cpu_time
	resp.Cases = append(resp.Cases,cs)
	env.wg.Done()
	env.mutex.Unlock()
}



