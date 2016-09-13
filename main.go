package main
//https://github.com/cedriczirtacic

import "fmt"
import "flag"
import "os"
import "os/exec"
import "log"
import "net/http"
import "sync"

// html parser
import "golang.org/x/net/html"

var wait sync.WaitGroup
var dir string

func help(e int) {
    // print usage and help
    fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
    flag.PrintDefaults()
    os.Exit(e);
}

// lets make it easy and use our friendly wget
func thread_download( f string, t_number int ) int {
    defer wait.Done()
    var cmd *exec.Cmd
    var uri string

    fmt.Printf("thread %d: downloading %s\n", t_number, f)
    // URIs in 4chan always start without "http:", so lets check
    if f[0:2] == "//" {
        uri = fmt.Sprintf("http:%s", f)
    }else{
        uri = f
    }

    // I assume that this is the path to your wget binary
    cmd = exec.Command("/usr/bin/wget", "-P", dir, uri);
    err := cmd.Run()

    if err != nil {
        log.Fatal(err)
    }
    return 0
}

func main() {
    var (
        url string
        filetype string
        helpflag bool
        threadnum int
        files []string
    )
    // get params
    flag.StringVar(&url, "url", "", "Url to get data and parse")
    flag.StringVar(&filetype, "filetype", "any", "Get this specific filetype")
    flag.StringVar(&dir, "dir", "./", "Download files in dir")
    flag.IntVar(&threadnum, "threadnum", 5, "Number of threads to run for download")
    flag.BoolVar(&helpflag, "h", false, "Display usage")

    flag.Parse()
    // there are arguments? help? url is setted?
    if  len(os.Args) == 1 || helpflag || url == "" {
        help(0)
    }
    // check if is a real url
    if len(url) < 8 || url[0:7] != "http://" {
        fmt.Fprintf(os.Stderr, "Error: provide a real url.\n")
        help(1)
    }

    // now we make the GET call for html
    r, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
        os.Exit(2)
    }
    // once we got all data, close the reader
    defer r.Body.Close()

    // parsing begins
    var t *html.Tokenizer = html.NewTokenizer(r.Body)
    for {
        var tt html.TokenType = t.Next()
        switch tt {
            case html.ErrorToken:
                // we have an error? jmp cont!
                goto cont
            case html.StartTagToken:
                var (
                    tag []byte 
                    attr bool
                )
                tag, attr = t.TagName()
                // we have an <a> tag?
                if string(tag) == "a" && attr == true {
                    var (
                        attrs,val []byte
                        mattr bool
                    )
                    // get attributes!
                    for attrs, val , mattr = t.TagAttr(); mattr == true; attrs, val , mattr = t.TagAttr() {
                        // classes "fileThumb" have our file URI
                        if string(attrs) == "class" && string(val) == "fileThumb" && mattr == true {
                            attrs, val , mattr = t.TagAttr()
                            // if found and has href attribute inside, store url in files array
                            if string(attrs) == "href" {
                                // we want an specific kind of file?
                                if filetype != "any" {
                                    name := string(val)
                                    // check for extension/filetype
                                    if name[len(name) - len(filetype):len(name)] == filetype {
                                        files = append(files, name)
                                    }
                                // no? then grab all files
                                }else{
                                    files =  append(files, string(val))
                                }
                            }
                        }
                    }
                }
            //case html.EndTagToken:
        }
    }

// jumped here to continue with the process
cont:
    if len(files) <= 0 {
        fmt.Fprintf(os.Stderr, "Error: couldn't get the filenames! Is the filetype OK?\n")
        os.Exit(3)
    }

    fmt.Printf("Number of files retrieved: %d\n", len(files))
    for i := len(files); i > 0; {
        if i < threadnum {
            threadnum = i
        }
        for j := threadnum; j > 0; j-- {
            wait.Add(1)
            go thread_download(files[(i-1)], j)
            i--
        }
        wait.Wait()
    }
}

