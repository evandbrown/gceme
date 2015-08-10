package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

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
	Error      string
}

const (
	html = `<!doctype html>
<html>
<head>
<!-- Compiled and minified CSS -->
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.97.0/css/materialize.min.css">

<!-- Compiled and minified JavaScript -->
<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.97.0/js/materialize.min.js"></script>
<title>Frontend Web Server</title>
</head>
<body>
<div class="container">
<div class="row">
<div class="col s3">&nbsp;</div>
<div class="col s6">
<table class="bordered striped hoverable">
  <thead>
	<tr>
		<th data-field="prop">Property</th>
		<th data-field="value">Value</th>
	</tr>
  </thead>

  <tbody>
	<tr>
	  <td>Name</td>
	  <td>{{.Name}}</td>
	</tr>
	<tr>
	  <td>ID</td>
	  <td>{{.Id}}</td>
	</tr>
	<tr>
	  <td>Hostname</td>
	  <td>{{.Hostname}}</td>
	</tr>
	<tr>
	  <td>Zone</td>
	  <td>{{.Zone}}</td>
	</tr>
	<tr>
	  <td>Project</td>
	  <td>{{.Project}}</td>
	</tr>
	<tr>
	  <td>InternalIP</td>
	  <td>{{.InternalIP}}</td>
	</tr>
	<tr>
	  <td>ExternalIP</td>
	  <td>{{.ExternalIP}}</td>
	</tr>
	  </tbody>
</table>
</div>
<div class="col s3">&nbsp;</div>
</div>
</div>
</html>`
)

func main() {

	frontend := flag.Bool("frontend", false, "run in frontend mode")
	port := flag.Int("port", 8080, "port to bind")
	backend := flag.String("backend-service", "", "hostname of backend server")
	flag.Parse()

	if *frontend {
		log.Println("Operating in frontend mode...")
		tpl, err := template.New("out").Parse(html)
		if err != nil {
			panic(err)
		}

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			resp, err := http.Get(*backend)
			if err != nil {
				fmt.Fprintf(w, "Error: %s\n", err.Error())
				return
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(w, "Error: %s\n", err.Error())
				return
			}

			i := &Instance{}
			json.Unmarshal([]byte(body), i)
			tpl.Execute(w, i)
		})

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))

	} else {
		log.Println("Operating in backend mode...")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			a := &assigner{}
			i := newInstance()
			i.Id = a.assign(metadata.InstanceID)
			i.Zone = a.assign(metadata.Zone)
			i.Name = a.assign(metadata.InstanceName)
			i.Hostname = a.assign(metadata.Hostname)
			i.Project = a.assign(metadata.ProjectID)
			i.InternalIP = a.assign(metadata.InternalIP)
			i.ExternalIP = a.assign(metadata.ExternalIP)

			if a.err != nil {
				i.Error = a.err.Error()
			}

			resp, err := json.Marshal(i)
			if err != nil {
				fmt.Fprint(w, err.Error())
			}

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
