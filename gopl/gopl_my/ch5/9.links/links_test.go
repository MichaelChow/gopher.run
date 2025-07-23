package links_test

import (
	"fmt"
	"testing"

	links "gopher.run/go/src/ch5/9.links"
)

func TestLinks(t *testing.T) {
	links, err := links.Extract("https://www.taobao.com")
	if err != nil {
		t.Errorf("Extract: %v", err)
	}
	for _, link := range links {
		fmt.Println(link)
	}

}
