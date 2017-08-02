package main

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	addr     = flag.String("listen", ":8888", "port to listen")
	chatAddr = flag.String("chat_addr", "localhost:3333", "chat tcp addr to proxy pass")
)

func main() {
	flag.Parse()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("can not get os working directory: %v", err)
	}

	web := http.FileServer(http.Dir(wd + "/web"))

	http.Handle("/", web)
	http.Handle("/web/", http.StripPrefix("/web/", web))
	http.Handle("/ws", upstream("chat", "tcp", *chatAddr))

	log.Printf("proxy is listening on %q", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func upstream(name, network, addr string) http.Handler {
	if conn, err := net.Dial(network, addr); err != nil {
		log.Printf("warning: test upstream %q error: %v", name, err)
	} else {
		log.Printf("upstream %q ok", name)
		conn.Close()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		peer, err := net.Dial(network, addr)
		if err != nil {
			log.Printf("dial upstream error: %v", err)
			w.WriteHeader(502)
			return
		}
		if err := r.Write(peer); err != nil {
			log.Printf("write request to upstream error: %v", err)
			w.WriteHeader(502)
			return
		}
		hj, ok := w.(http.Hijacker)
		if !ok {
			w.WriteHeader(500)
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			w.WriteHeader(500)
			return
		}

		log.Printf(
			"serving %s < %s <~> %s > %s",
			peer.RemoteAddr(), peer.LocalAddr(), conn.RemoteAddr(), conn.LocalAddr(),
		)

		go func() {
			defer peer.Close()
			defer conn.Close()
			io.Copy(peer, conn)
		}()
		go func() {
			defer peer.Close()
			defer conn.Close()
			io.Copy(conn, peer)
		}()
	})
}

func indexHandler(wd string) (http.Handler, error) {
	index, err := os.Open(wd + "/web/index.html")
	if err != nil {
		return nil, err
	}
	stat, err := index.Stat()
	if err != nil {
		return nil, err
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "", stat.ModTime(), index)
	}), nil
}
