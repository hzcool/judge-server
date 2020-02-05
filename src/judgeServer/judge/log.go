package judge


import (
	"os"
	"path/filepath"
	"log"
)

func Log(logfile string, message string) error {
	file,err := os.OpenFile(filepath.Join(LOG_DIR_PATH,logfile),os.O_APPEND|os.O_WRONLY|os.O_CREATE,0666)
	if err!=nil {
		return err
	}
	defer file.Close()
	logger := log.New(file, "", log.LstdFlags|log.Lshortfile) 
	logger.Println(message)
	return nil
}