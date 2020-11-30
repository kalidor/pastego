package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Page struct {
	Content   string
	Pasteid   string
	Url       string
	TimeStart string
	TimeStop  string
	Iv        string
	Key       string
}

var (
	Dir    *string
	TmpDir string
	Port   *int
	Css    string
	Js     [2]string
)

func loadPaste(pasteid string) (*Page, error) {
	content, err := ioutil.ReadFile(filepath.Join(TmpDir, pasteid))
	if err != nil {
		return nil, err
	}
	datas := strings.SplitN(string(content), "\n", 3)
	times := strings.SplitN(datas[0], "|", 2)
	return &Page{
		Iv:        datas[1],
		Content:   datas[2],
		Pasteid:   pasteid,
		TimeStart: times[0],
		TimeStop:  times[1],
	}, nil
}

func rawHandler(w http.ResponseWriter, r *http.Request) {
	pasteid := r.URL.Path[len("/raw/"):]
	if len(pasteid) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	content, err := ioutil.ReadFile(filepath.Join(TmpDir, pasteid))
	if err != nil {
		fmt.Fprintf(w, "Not found :/")
	} else {
		datas := strings.SplitN(string(content), "\n", 2)
		s := strings.NewReader(datas[1])
		io.Copy(w, s)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	pasteid := r.URL.Path[len("/view/"):]
	if len(pasteid) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Printf("Searching: %s", pasteid)
	paste, err := loadPaste(pasteid)
	if err != nil {
		fmt.Fprintf(w, "Not found :/")
	} else {
		t, _ := template.New("view").Parse(VIEW)
		paste.Url = r.Host
		fmt.Println(paste.Url)
		t.Execute(w, paste)
	}
}

func removePaste(paste string) {
	fmt.Printf("Removing: %s\n", paste)
	if err := os.Remove(paste); err != nil {
		log.Println("Cannot remove file.", err)
	}
}

func addPaste(body, pasteid, iv string, eol int) {
	fmt.Printf("addPaste called: %s - %s\n", pasteid, body)
	tmpfn := filepath.Join(TmpDir, pasteid)
	t := time.Now()
	tf := t.Add(time.Duration(eol) * time.Minute)
	body = fmt.Sprintf("%s|%s\n%s\n%s",
		t.Format("2006-01-02 15:03:00"),
		tf.Format("2006-01-02 15:03:00"),
		iv,
		strings.ReplaceAll(body, " ", "+"),
	)
	err := ioutil.WriteFile(tmpfn, []byte(body), 0600)
	if err != nil {
		fmt.Println("Cannot write file.", err)
	} else {
		fmt.Println(time.Now())
		select {
		case <-time.After(time.Duration(eol) * time.Minute):
			fmt.Println(time.Now())
			removePaste(tmpfn)
		}
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("ciphertext")
	iv := r.FormValue("iv")
	if len(body) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	eol, err := strconv.Atoi(r.FormValue("eol"))
	if err != nil {
		handler(w, r)
	}
	pasteid, _ := uuid.NewRandom()
	go func() {
		addPaste(body, pasteid.String(), iv, eol)
	}()
	http.Redirect(w, r, "/view/"+pasteid.String(), http.StatusFound)
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	cssReader := strings.NewReader(Css)
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	io.Copy(w, cssReader)
}
func jsHandler(w http.ResponseWriter, r *http.Request) {
	jsFile := r.URL.Path[len("/js/"):]
	log.Println("Got JS", jsFile)
	var JsReader *strings.Reader
	if jsFile == "aes-ctr-encrypt.js" {
		JsReader = strings.NewReader(Js[0])
	} else {
		JsReader = strings.NewReader(Js[1])
	}
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	io.Copy(w, JsReader)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, INDEX)
}

func LoadCss() {
	css, err := ioutil.ReadFile("css/bootstrap.min.css")
	if err != nil {
		log.Println("Cannot read CSS file", err)
	}
	Css = string(css)
}
func LoadJs() {
	js, err := ioutil.ReadFile("js/aes-ctr-encrypt.js")
	if err != nil {
		log.Println("Cannot read JSS file", err)
	}
	Js[0] = string(js)
	js, err = ioutil.ReadFile("js/aes-ctr-decrypt.js")
	if err != nil {
		log.Println("Cannot read JSS file", err)
	}
	Js[1] = string(js)
}

func main() {
	Dir = flag.String("dir", "/tmp/", "Directory where temporary dir will be created and received paste file")
	Port = flag.Int("port", 8000, "Listening port")

	tmpdir, err := ioutil.TempDir(*Dir, "*.paste")
	if err != nil {
		log.Panic(err)
	}
	log.Println(tmpdir)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Cleaning tmp directory")
		os.RemoveAll(tmpdir)
		os.Exit(0)
	}()
	TmpDir = tmpdir
	LoadCss()
	LoadJs()
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/raw/", rawHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/css/", cssHandler)
	http.HandleFunc("/js/", jsHandler)
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", *Port), nil))
}
