package myip

import "testing"

func TestGetPublicIP(t *testing.T) {
	if ip, err := GetPublicIP(); err != nil {
		t.Fatal(err)
	} else {
		t.Log(ip)
	}
}
