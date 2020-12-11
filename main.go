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
}

var (
	Dir    *string
	TmpDir string
	Port   *int
	Css    string
)

func loadPaste(pasteid string) (*Page, error) {
	content, err := ioutil.ReadFile(filepath.Join(TmpDir, pasteid))
	if err != nil {
		return nil, err
	}
	datas := strings.SplitN(string(content), "\n", 2)
	times := strings.SplitN(datas[0], "|", 2)
	return &Page{
		Content:   datas[1],
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

func addPaste(body, pasteid string, eol int) {
	fmt.Printf("addPaste called: %s - %s\n", pasteid, body)
	tmpfn := filepath.Join(TmpDir, pasteid)
	t := time.Now()
	t_body := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	tf := t.Add(time.Duration(eol) * time.Minute)
	tf_body := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d",
		tf.Year(), tf.Month(), tf.Day(),
		tf.Hour(), tf.Minute(), tf.Second())
	body = fmt.Sprintf("%s|%s\n",
		t_body, tf_body) + body
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
	body := r.FormValue("content")
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
		addPaste(body, pasteid.String(), eol)
	}()
	http.Redirect(w, r, "/view/"+pasteid.String(), http.StatusFound)
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	cssReader := strings.NewReader(Css)
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	io.Copy(w, cssReader)
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

func main() {
	Dir = flag.String("dir", "/tmp/", "Directory where temporary dir will be created and received paste file")
	Port = flag.Int("port", 8000, "Listening port")

	tmpdir, err := ioutil.TempDir(*Dir, "*.paste")
	if err != nil {
		log.Panic(err)
	}
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
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/raw/", rawHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/css/", cssHandler)
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", *Port), nil))
}
