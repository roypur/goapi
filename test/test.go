package main
import (
    "github.com/roypur/goapi/src"
    "fmt"
)
func main(){
    goapi.Listen(cont, ":4356", 
"/home/roypur/ssl/startssl/purser_it/crt.pem", 
"/home/roypur/ssl/startssl/purser_it/key.pem")
}

func cont(req goapi.Request){
    req.Write("\r\n");
    req.Close();
    fmt.Println(req.Body);
    fmt.Println(req.Header);
    fmt.Println(req.Method);
    fmt.Println(req.Path);
    fmt.Println(req.Version);
}


