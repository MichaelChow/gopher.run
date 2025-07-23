package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SearchIssues函数发出一个HTTP请求，然后解码返回的JSON格式的结果
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	// 用户提供的查询条件可能包含类似?和&之类的特殊字符，为了避免对URL造成冲突，用url.QueryEscape来对查询中的特殊字符进行转义操作
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	//!-
	// For long-term stability, instead of http.Get, use the
	// variant below which adds an HTTP request header indicating
	// that only version 3 of the GitHub API is acceptable.
	//
	//   req, err := http.NewRequest("GET", IssuesURL+"?q="+q, nil)
	//   if err != nil {
	//       return nil, err
	//   }
	//   req.Header.Set(
	//       "Accept", "application/vnd.github.v3.text-match+json")
	//   resp, err := http.DefaultClient.Do(req)
	//!+

	// We must close resp.Body on all execution paths.
	// (Chapter 5 presents 'defer', which makes this simpler.)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	// movie.go中使用了json.Unmarshal函数来将JSON格式的字符串解码为字节slice
	// 这里使用了基于流式的解码器json.Decoder，它可以从一个输入流解码JSON数据，尽管这不是必须的
	// 相对应的，还有一个针对输出流的json.Encoder编码对象
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}
