package goapi

import (
	"net"
	"bufio"
	"strings"
	
	"fmt"
)

type Request struct{
    Header map[string]string
    Body string
    Conn net.Conn
}

func Listen(handler func(req Request), listen string){
    ln,_ := net.Listen("tcp", listen)
    
    for {
        conn,_ := ln.Accept()
        go handler(parse(conn));
    }
}

func (r Request) Write(text string){
    r.Conn.Write([]byte(text));
}

func (r Request) Close(){
    r.Conn.Close();
}

func parse(conn net.Conn)(Request){

    var resp string = "HTTP/1.1 200 OK\r\nAccess-Control-Allow-Origin: *\r\nAccess-Control-Allow-Methods: POST\r\nContent-Type: text/html; charset=utf-8\r\n";

    message := bufio.NewReader(conn);
    
    var str string = "init";

    header := make(map[string]string);

    for(len(str) > 3){
    
        str,_ = message.ReadString('\n');
        
        fmt.Println(str);
        
        val := strings.SplitN(strings.TrimSpace(str), " ", 2);
        
        fmt.Println(val);
        if(len(val)>1){
            header[val[0]] = val[1]
        }
    }
    
    body := make([]byte, message.Buffered());
    
    message.Read(body);
    
    req := new(Request);
    
    req.Body = string(body);
    req.Header = header;
    req.Conn = conn;
    
    conn.Write([]byte(resp));
    
    return *req;
}
