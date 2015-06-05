package goapi

import (
	"net"
	"bufio"
	"strings"
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
        
        str = strings.TrimSpace(str);
        
        line := strings.SplitN(str, " ", 2);
        
        if(strings.Index(line[0], ":")>0){
            line = strings.SplitN(str, ":", 2);
            
            line[0] = strings.ToLower(line[0]);
            
            key := strings.Split(line[0],"");
            
            key[0] = strings.ToUpper(key[0]);
            
            i := 0;
            for(i < len(key)){
                if(key[i] == "-"){
                    key[i+1] = strings.ToUpper(key[i+1]);
                }
                i++;
            }
            line[0] = strings.Join(key, "");
            
        }

        if(len(line)>1){
            key := strings.TrimSpace(line[0]);
            val := strings.TrimSpace(line[1]);
            
            header[key] = val;
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
