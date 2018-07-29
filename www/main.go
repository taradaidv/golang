package main
import (
        "log"
        "net/http"
    "os"
    "fmt"
    "path/filepath"
    "io/ioutil"
        "encoding/xml"
    
)



type Recurlyservers struct {
	XMLName     xml.Name `xml:"nmaprun"`
	Scanner		string   `xml:"scanner,attr"`
	Version     string   `xml:"version,attr"`
    Startstr	string   `xml:"startstr,attr"`
    Svs			[]host	`xml:"host"`

}

type host struct {
    XMLName		xml.Name `xml:"host"`
    Starttime	string   `xml:"starttime"`
}


func handler(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.Write([]byte("This is an example server.\n"))
}

func fooHandler (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write([]byte(r.Header.Get("User-Agent")))

}

func fileHandler (w http.ResponseWriter, r *http.Request) {
        
        var files []string

    root := "/tmp"
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        files = append(files, path)
        return nil
    })
    if err != nil {
        panic(err)
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    
    for _, file := range files {
        w.Write([]byte(file+"<br>"))
    }
    
}

func readHandler (w http.ResponseWriter, r *http.Request) {
        
      data, err := ioutil.ReadFile("./nmap/result.xml")

  if err != nil {
	
		
	fmt.Fprintf(w, "err: %v", err)
	
	}
		
w.Write([]byte(data))


    
}


func xmlHandler (w http.ResponseWriter, r *http.Request) {
        
      file, err := os.Open("./nmap/result.xml") // For read access.     
    if err != nil {
        fmt.Printf("error: %v", err)
        return
    }
    defer file.Close()
    data, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Printf("error: %v", err)
        return
    }
    v := Recurlyservers{}
    err = xml.Unmarshal(data, &v)
    if err != nil {
        fmt.Printf("error: %v", err)
        return
    }

w.Header().Set("Content-Type", "text/html; charset=utf-8")
fmt.Fprintf(w, "%v", v)
	
    
}



func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Fprintf(w, htmlStr)
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Hello, POST method. ParseForm() err: %v", err)
			return
		}

		// Post form from website
		switch r.FormValue("post_from") {
		case "web":
			fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
			s := r.FormValue("key")
			fmt.Fprintf(w, "key = %s, len = %v\n", s, len(s))

		case "client":
			fmt.Fprintf(w, "Post from client! r.PostForm = %v\n", r.PostForm)

		default:
			fmt.Fprintf(w, "Unknown post source:-(\n")
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
        http.HandleFunc("/", handler)
        http.HandleFunc("/1", fooHandler)
        http.HandleFunc("/2", fileHandler)
        http.HandleFunc("/3", readHandler)
        http.HandleFunc("/4", xmlHandler)
        http.HandleFunc("/hello", helloHandler)

        log.Printf("About to listen on 8443.")
        err := http.ListenAndServeTLS("127.0.0.1:8443", "./tls/CERTIFICATE.key", "./tls/PRIVATE.key", nil)
        log.Fatal(err)
}



var htmlStr = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8" />
</head>
<body>
  <div>
      <form method="POST" action="/hello">
          <input name="post_from" type="text" value="web" >
          <input name="key" type="text" value="Hello, -">
	  <input type="submit" value="submit" /hello>
      </form>
  </div>
</body>
</html>
`