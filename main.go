package main

import (
	"encoding/json"
	"fmt"
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
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		i := newInstance()
		i.Id = assign(metadata.InstanceID)
		i.Zone = assign(metadata.Zone)
		i.Name = assign(metadata.InstanceName)
		i.Hostname = assign(metadata.Hostname)
		i.Project = assign(metadata.ProjectID)
		i.InternalIP = assign(metadata.InternalIP)
		i.ExternalIP = assign(metadata.ExternalIP)

		resp, err := json.Marshal(i)
		if err != nil {
			fmt.Fprint(w, err.Error())
		}

		fmt.Fprintf(w, "%s", resp)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func assign(getVal func() (string, error)) string {
	s, err := getVal()
	if err != nil {
		return "Error: unable to retrieve metadata for this value"
	}
	return s
}

func newInstance() *Instance {
	var i = new(Instance)
	return i
}
