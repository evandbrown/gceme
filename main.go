package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	//"net/http/httputil"
	"runtime"
	//"strings"

	"github.com/GoogleCloudPlatform/gcloud-golang/compute/metadata"
)

type Instance struct {
	Id         string
	Name       string
	Hostname   string
	Zone       string
	Project    string
	InternalIP string
	ExternalIP string
	LBRequest  string
	ClientIP   string
	Error      string
}

var Version string = "version"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	showversion := flag.Bool("version", false, "display version")
	frontend := flag.Bool("frontend", false, "run in frontend mode")
	port := flag.Int("port", 8080, "port to bind")
	backend := flag.String("backend-service", "", "hostname of backend server")
	flag.Parse()

	if *showversion {
		fmt.Printf("Version %s\n", Version)
		return
	}

	if *frontend {
		log.Println("Operating in frontend mode...")
		tpl := template.Must(template.New("out").Parse(html))

		transport := http.Transport{DisableKeepAlives: false}
		client := &http.Client{Transport: &transport}
		req, _ := http.NewRequest(
			"GET",
			*backend,
			nil,
		)
		req.Close = false

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			i := &Instance{}
			resp, err := client.Do(req)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintf(w, "Error: %s\n", err.Error())
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error: %s\n", err.Error())
				return
			}
			err = json.Unmarshal([]byte(body), i)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error: %s\n", err.Error())
				return
			}
			tpl.Execute(w, i)
		})

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))

	} else {

		a := &assigner{}
		i := newInstance()
		i.Id = a.assign(metadata.InstanceID)
		i.Zone = a.assign(metadata.Zone)
		i.Name = a.assign(metadata.InstanceName)
		i.Hostname = a.assign(metadata.Hostname)
		i.Project = a.assign(metadata.ProjectID)
		i.InternalIP = a.assign(metadata.InternalIP)
		i.ExternalIP = a.assign(metadata.ExternalIP)
		resp, _ := json.Marshal(i)

		log.Println("Operating in backend mode...")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if a.err != nil {
				i.Error = a.err.Error()
			}

			//raw, _ := httputil.DumpRequest(r, true)
			//i.LBRequest = string(raw)

			fmt.Fprintf(w, "%s", resp)
		})

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
	}

}

type assigner struct {
	err error
}

func (a *assigner) assign(getVal func() (string, error)) string {
	if a.err != nil {
		return ""
	}
	s, err := getVal()
	if err != nil {
		a.err = err
	}
	return s
}

func newInstance() *Instance {
	var i = new(Instance)
	return i
}
