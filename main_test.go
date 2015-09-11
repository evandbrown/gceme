package main

import (
	"google.golang.org/cloud/compute/metadata"
	"testing"
)

func TestGCE(t *testing.T) {
	i := newInstance()
	if !metadata.OnGCE() && i.Error != "Not running on GCE" {
		t.Error("Test not running on GCE, but error does not indicate that fact.")
	}
}
