package main

import (
	"fmt"
	quic "github.com/zhangjiaxinghust/mp-quic"
	"github.com/zhangjiaxinghust/mp-quic/h2quic"
	"io"
	"net/http"
	"strings"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

type requestHandler struct {
	endpoint	string
}

func (f *requestHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	upath := request.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		request.URL.Path = upath
	}

	multipath := true

	quicConfig := &quic.Config{
		CreatePaths: multipath,
	}

	hclient := &http.Client{
		Transport: &h2quic.RoundTripper{QuicConfig: quicConfig},
	}

	addr := f.endpoint + request.URL.Path
	rsp, err := hclient.Get(addr)
	if err != nil {
		panic(err)
	}

	//body := &bytes.Buffer{}
	//_, err = io.Copy(body, rsp.Body)
	/*if err != nil {
		panic(err)
	}*/
	if(addr[len(addr)-1] == '/'){
		fmt.Println(addr + " is a dir")
		io.WriteString(writer, "<html><head><style>img{width:40px;height:40px;}</style></head><body>")


		io.Copy(writer, rsp.Body)
		io.WriteString(writer, "</body></html>")
	}else{
		fmt.Println(addr + " is a file")
		io.Copy(writer, rsp.Body)
	}


}

func main() {

	ip :="https://192.168.122.15:6121"

	http.Handle(
		"/",
		&requestHandler{ip},
	)

	http.ListenAndServe(":8888", nil)
}