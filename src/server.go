package goapi

import (
	"net"
	"bufio"
	"strings"
	"crypto/tls"
)

type Request struct{
    Header map[string]string
    Body string
    Method string
    Path string
    Version string
    Resp string // cross origin allow all for current request
    Conn net.Conn
}

func Listen(handler func(req Request), listen string, certFile string, keyFile string){

    tlsPair,_ := tls.LoadX509KeyPair(certFile, keyFile)

    tlsConfig := &tls.Config{Certificates : []tls.Certificate{tlsPair}}

    ln,_ := tls.Listen("tcp", listen, tlsConfig)
    
    
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
            line = strings.SplitN(str, ":", 2); // removes colon after header key
            
            line[0] = strings.ToLower(line[0]);
            
            key := strings.Split(line[0],"");
            
            key[0] = strings.ToUpper(key[0]);
            
            i := 0;
            
            //sets first letter after - to upper-case
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
    
    var resp string = "";
            
    if(len(header["Access-Control-Request-Method"])>0){
        resp += "Access-Control-Allow-Methods: " + header["Access-Control-Request-Method"] + "\r\n";
    }
    if(len(header["Access-Control-Request-Headers"])>0){
        resp += "Access-Control-Allow-Headers: " + header["Access-Control-Request-Headers"] + "\r\n";
    }
    if(len(header["Origin"])>4){
        resp += "Access-Control-Allow-Origin: " + header["Origin"] + "\r\n";
        resp += "Access-Control-Allow-Credentials: true";
    }else{
        resp += "Access-Control-Allow-Origin: *\r\n";
    }
    
    req.Resp = resp   
       
    return *req;
}
