package github_test

import (
	"fmt"
	"testing"

	github "gopher.run/go/src/ch4/18.github"
)

func TestSearch(t *testing.T) {
	result, _ := github.SearchIssues([]string{"repo:golang/go", "is:open", "json", "decoder"})
	fmt.Printf("%#v", result)
	// &github.IssuesSearchResult{TotalCount:92, Items:[]*github.Issue{(*github.Issue)(0x1400009ed90), (*github.Issue)(0x1400009ee00), (*github.Issue)(0x1400009ee70), (*github.Issue)(0x1400009eee0), (*github.Issue)(0x1400009ef50), (*github.Issue)(0x1400009efc0), (*github.Issue)(0x1400009f030), (*github.Issue)(0x1400009f0a0), (*github.Issue)(0x1400009f110), (*github.Issue)(0x1400009f180), (*github.Issue)(0x1400009f1f0), (*github.Issue)(0x1400009f260), (*github.Issue)(0x1400009f2d0), (*github.Issue)(0x1400009f340), (*github.Issue)(0x1400009f3b0), (*github.Issue)(0x1400009f420), (*github.Issue)(0x1400009f490), (*github.Issue)(0x1400009f500), (*github.Issue)(0x1400009f570), (*github.Issue)(0x1400009f5e0), (*github.Issue)(0x1400009f650), (*github.Issue)(0x1400009f6c0), (*github.Issue)(0x1400009f730), (*github.Issue)(0x1400009f7a0), (*github.Issue)(0x1400009f810), (*github.Issue)(0x1400009f880), (*github.Issue)(0x1400009f8f0), (*github.Issue)(0x1400009f960), (*github.Issue)(0x1400009f9d0), (*github.Issue)(0x1400009fa40)}}
	fmt.Printf("%d issues: \n", result.TotalCount)
}
