package main

//go:generate statik -src=./public

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/rakyll/statik/fs"

	_ "github.com/npflan/speedtest/statik"
)

func empty(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Add("Cache-Control", "post-check=0, pre-check=0")
		w.Header().Set("Pragma", "no-cache")
	}
	b := make([]byte, 1<<20)
	for {
		_, e := r.Body.Read(b)
		if e == io.EOF || e == io.ErrUnexpectedEOF {
			break
		} else if e != nil {
			panic(e)
		}
	}
	defer r.Body.Close()
	w.WriteHeader(200)
}

var garbageBuf []byte = make([]byte, 1<<20)

func garbage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Add("Cache-Control", "post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=random.dat")
	w.Header().Set("Content-Transfer-Encoding", "binary")

	chunks := 4
	if ckSize := r.FormValue("ckSize"); ckSize != "" {
		if iSize, err := strconv.ParseInt(ckSize, 10, 32); err == nil {
			if iSize > 1024 {
				chunks = 1024
			} else {
				chunks = int(iSize)
			}
		}
	}

	w.WriteHeader(200)
	for i := 0; i < chunks; i++ {
		w.Write(garbageBuf)
	}
}

func ip(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.RemoteAddr))
}

func do() error {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/empty", empty)
	http.HandleFunc("/garbage", garbage)
	http.HandleFunc("/getIP", ip)
	http.Handle("/", http.FileServer(statikFS))
	return http.ListenAndServe(":8080", nil)
}

func main() {
	if err := do(); err != nil {
		panic(err)
	}
}
