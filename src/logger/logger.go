package logger


import (
	"fmt"
	"log"
	"os"
)

type logger struct {
	Info *log.Logger
	Debug *log.Logger
	Warning *log.Logger
	Error *log.Logger
}

var Logger logger = logger{}

func init() {
	// NOTE: not a good approach
	fmt.Println("Init Logger")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Ltime | log.Lmicroseconds)
	Logger.Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Debug = log.New(os.Stdout, "[Debug] ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Warning = log.New(os.Stdout, "[Warning] ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Error = log.New(os.Stdout, "[Error] ", log.Ldate|log.Ltime|log.Lshortfile)
}

