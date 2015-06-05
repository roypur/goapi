package main
import (
    "api"
    "fmt"
)
func main(){
    api.Listen(cont, ":4356")
}

func cont(req api.Request){
    req.Write("");
    req.Close();
    fmt.Println(req.Body);
}


