package main

import(
    "fmt"
    "judgeServer/judge"
    "net/http"
)


func main(){
    
    judge.Init()
    http.HandleFunc("/ping",judge.Ping)
    http.HandleFunc("/judge",judge.Judge)
    fmt.Println("已开启测评姬")
    http.ListenAndServe(judge.SERVICE_PORT, nil)
}