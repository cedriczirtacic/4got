# 4got

Simple way to download all thread files from a 4chan board.

## Help
```bash
Usage of ./4got:
  -dir string
    	Download files in dir (default "./")
  -filetype string
    	Get this specific filetype (default "any")
  -h	Display usage
  -threadnum int
    	Number of threads to run for download (default 5)
  -url string
    	Url to get data and parse
```
## Compile
```bash
go get golang.org/x/net/html
go build
```
## Todo
* Crawler? Maybe?
