package main
import (
    "github.com/roypur/goapi/src"
    "fmt"
)
func main(){
    goapi.Listen(cont, ":4356")
}

func cont(req api.Request){
    req.Write("");
    req.Close();
    fmt.Println(req.Body);
}


