##Simple go api-server

see test/test.go for example usage.


###type Request

    type Request struct{
        Header map[string]string
        Body string
        Conn net.Conn
    }

###func Listen(handler func(req Request), listen string)
Listen for http connections


###func (r Request) Write(text string)
Write response to client

###func (r Request) Close()
Close connection
