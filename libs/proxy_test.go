package lib

import (
	"strings"
	"testing"
)

func TestGetProxy(t *testing.T) {
	proxies, err := getProxy()
	if err != nil {
		t.Errorf(err.Error())
	}

	first := proxies[0]
	if strings.Count(first, ".") != 3 {
		t.Errorf("Error in ip %v", first)
	}

	if !strings.Contains(first, ":") {
		t.Errorf("Error in port %v", first)
	}
}
