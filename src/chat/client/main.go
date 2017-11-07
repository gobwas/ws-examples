package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/jroimartin/gocui"
)

var (
//	endpoint = flag.String("")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s <endpoint> [flags]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	u, err := url.ParseRequestURI(args[0])
	if err != nil {
		log.Fatal(err)
	}

	d := wsutil.DebugDialer{
		Dialer: &ws.Dialer{
			Header: func(w io.Writer) {
				(http.Header{
					"X-Upgrade": []string{"hello"},
				}).Write(w)
			},
			OnStatusError: func(status int, reason []byte, r io.Reader) {
				log.Printf("status error: %d %s", status, reason)
				bts, err := ioutil.ReadAll(r)
				if err != nil {
					log.Printf("can no read %d status response: %v", err)
					return
				}
				log.Printf("resp:\n%s", bts)
			},
		},

		OnRequest: func(req []byte) {
			log.Printf("conn req:\n%s====", req)
		},
		OnResponse: func(res []byte) {
			log.Printf("conn res:\n%s====", res)
		},
	}
	conn, br, hs, err := d.Dial(context.Background(), u.String())
	if err != nil {
		log.Fatalf("can not connect: %v", err)
	}

	log.Printf("handshake: %+v %v", hs, br)

	for {
		f, err := ws.ReadFrame(br)
		if err != nil {
			panic(err)
		}
		log.Printf("FRAME: %v %q", f, f.Payload)
	}

	return

	var (
		in     = make(chan []byte, 1)
		out    = make(chan []byte, 1)
		errors = make(chan error, 1)
	)
	go func() {
		var i int
		w := wsutil.NewWriter(conn, ws.StateClientSide, ws.OpText)
		for msg := range out {
			publish(w, i, msg)
			if err := w.Flush(); err != nil {
				errors <- err
				return
			}
			i++
		}
	}()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	go func() {
		for msg := range in {
			g.Execute(func(g *gocui.Gui) error {
				v, err := g.View("msg")
				if err == nil {
					_, err = v.Write(msg)
				}
				return err
			})
		}
	}()
	g.SetManagerFunc(layout)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		bts, err := ioutil.ReadAll(v)
		if err == nil {
			out <- bts
			v.Clear()
			v.SetCursor(0, 0)
		}
		return err
	})

	if br != nil {
		//log.Println("Has buffered data")
		for br.Buffered() > 0 {
			f, err := ws.ReadFrame(br)
			log.Println(string(f.Payload), err)
			if err != nil {
				errors <- err
				return
			}
			in <- f.Payload
		}
		//log.Println("released")
		ws.PutReader(br)
	}
	br = bufio.NewReader(conn)
	for {
		f, err := ws.ReadFrame(br)

		//msg, err := wsutil.ReadServerText(br)
		log.Println(string(f.Payload), err)
		if err != nil {
			errors <- err
			return
		}
		in <- f.Payload
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		errors <- err
	}
	log.Println("DONE", <-errors)
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("input", 0, maxY-4, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
		v.Title = "Your message"
		v.Editable = true
		v.Wrap = true
	}
	if v, err := g.SetView("msg", 0, 0, maxX-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
		v.Title = "Chat"
		v.Wrap = true
		v.Autoscroll = true
	}
	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	g.Close()
	return gocui.ErrQuit
}

func publish(w io.Writer, id int, text []byte) {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	w.Write([]byte(`{"id":`))
	w.Write([]byte(strconv.Itoa(id)))
	w.Write([]byte(`,"method":"publish","params":{"text":"`))
	w.Write(bytes.TrimSpace(text))
	w.Write([]byte(`","time":"`))
	w.Write([]byte(now))
	w.Write([]byte(`"}}`))
}

type rwConn struct {
	r io.Reader
	w io.Writer
	net.Conn
}

func (c rwConn) Read(p []byte) (int, error)  { return c.r.Read(p) }
func (c rwConn) Write(p []byte) (int, error) { return c.w.Write(p) }
