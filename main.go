package main

import (
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
	"syscall"
	"time"
)

type Page struct {
	Content string
	Pasteid string
}

var TmpDir string

func loadPaste(pasteid string)(*Page, error) {
    content, err := ioutil.ReadFile(filepath.Join(TmpDir, pasteid))
	if err != nil {
		return nil, err
	}
    return &Page{Content: string(content), Pasteid: pasteid}, nil
}

func rawHandler(w http.ResponseWriter, r *http.Request) {
	pasteid := r.URL.Path[len("/raw/"):]
	if len(pasteid) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	f, err := os.Open(filepath.Join(TmpDir, pasteid))
	if err != nil {
		fmt.Fprintf(w, "Not found :/")
	} else {
		io.Copy(w, f)
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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, INDEX)
}

func main() {
	tmpdir, err := ioutil.TempDir("/tmp/", "*.paste")
	if err != nil {
		log.Fatal(err)
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
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/raw/", rawHandler)
	http.HandleFunc("/create", createHandler)
    //TODO handle CSS file...
	log.Fatal(http.ListenAndServe(":8080", nil))
}
