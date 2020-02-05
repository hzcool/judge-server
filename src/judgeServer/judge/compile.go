package judge
import (
	// "fmt";
	"os/exec"
	"bytes"
)

func compile(command string) (string,bool) {
	cmd := exec.Command("/bin/bash","-c",command)
    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Run();err!=nil {
		return out.String(),false
	}
	return out.String(),true
}
