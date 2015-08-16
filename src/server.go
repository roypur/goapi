package goapi

import (
	"net"
	"bufio"
	"strings"
)

type Request struct{
    Header map[string]string
    Body string
    Method string
    Path string
    Version string
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

var Resp string = "";

const Ok string = "HTTP/1.1 200 OK\r\nContent-Type: text/html; charset=utf-8\r\n";
const Redirect string = "HTTP/1.1 301 Moved Permanently\r\n";

func parse(conn net.Conn)(Request){

    message := bufio.NewReader(conn);
    
    var str string = "init";
    
    var method string;
    var version string;
    var path string;
    

    header := make(map[string]string);

    for(len(str) > 3){
    
        str,_ = message.ReadString('\n');
        
        str = strings.TrimSpace(str);
        
        var isMethod bool = false;        
                
        line := strings.SplitN(str, " ", 3);
        
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
            
        }else{
            if(len(line)>2){
                method = strings.ToUpper(strings.TrimSpace(line[0]));
                path = strings.TrimSpace(line[1]);
                version = strings.TrimSpace(line[2]);
                isMethod = true;
            }
        }

        if((len(line)>1) && isMethod==false){
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
    req.Method = method;
    req.Version = version;
    req.Path = path;

    req.Conn = conn;
    
            
    if(len(header["Access-Control-Request-Method"])>0){
        Resp += "Access-Control-Allow-Methods: " + header["Access-Control-Request-Method"] + "\r\n";
    }
    if(len(header["Access-Control-Request-Headers"])>0){
        Resp += "Access-Control-Allow-Headers: " + header["Access-Control-Request-Headers"] + "\r\n";
    }
    if(len(header["Origin"])>4){
        Resp += "Access-Control-Allow-Origin: " + header["Origin"] + "\r\n";
        Resp += "Access-Control-Allow-Credentials: true";
    }else{
        Resp += "Access-Control-Allow-Origin: *\r\n";
    }
    
    return *req;
}
