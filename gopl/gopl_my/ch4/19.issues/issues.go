// Issues prints a table of GitHub issues matching the search terms.
// See page 112.
package main

import (
	"fmt"
	"log"
	"os"

	github "gopher.run/go/src/ch4/18.github"
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s %s %s\n",
			item.Number, item.User.Login, item.Title, item.HTMLURL, item.CreatedAt)
	}
}

/*
//!+textoutput
$ go build gopl.io/ch4/issues
$ ./issues repo:golang/go is:open json decoder
92 issues:
#48298     dsnet encoding/json: add Decoder.DisallowDuplicateFields https://github.com/golang/go/issues/48298 2021-09-09 19:39:33 +0000 UTC
#69449   gazerro encoding/json: Decoder.Token does not return an error f https://github.com/golang/go/issues/69449 2024-09-13 15:50:14 +0000 UTC
#5901        rsc encoding/json: allow per-Encoder/per-Decoder registrati https://github.com/golang/go/issues/5901 2013-07-17 16:39:15 +0000 UTC
#56733 rolandsho encoding/json: add (*Decoder).SetLimit https://github.com/golang/go/issues/56733 2022-11-14 18:51:33 +0000 UTC
#6647    btracey x/pkgsite: display type kind of each named type https://github.com/golang/go/issues/6647 2013-10-23 17:19:48 +0000 UTC
#42571     dsnet encoding/json: clarify Decoder.InputOffset semantics https://github.com/golang/go/issues/42571 2020-11-13 00:09:09 +0000 UTC
#11046     kurin encoding/json: Decoder internally buffers full input https://github.com/golang/go/issues/11046 2015-06-03 19:25:08 +0000 UTC
#67525 mateusz83 encoding/json: don't silently ignore errors from (*Deco https://github.com/golang/go/pull/67525 2024-05-20 14:10:55 +0000 UTC
#58649 nabokihms encoding/json: show nested fields path if DisallowUnkno https://github.com/golang/go/issues/58649 2023-02-22 23:20:53 +0000 UTC
#43716 ggaaooppe encoding/json: increment byte counter when using decode https://github.com/golang/go/pull/43716 2021-01-15 08:58:39 +0000 UTC
#36225     dsnet encoding/json: the Decoder.Decode API lends itself to m https://github.com/golang/go/issues/36225 2019-12-19 22:26:12 +0000 UTC
#26946    deuill encoding/json: clarify what happens when unmarshaling i https://github.com/golang/go/issues/26946 2018-08-12 18:19:01 +0000 UTC
#29035    jaswdr proposal: encoding/json: add error var to compare  the  https://github.com/golang/go/issues/29035 2018-11-30 11:21:31 +0000 UTC
#61627    nabice x/tools/gopls: feature: CLI syntax for renaming by iden https://github.com/golang/go/issues/61627 2023-07-28 06:40:34 +0000 UTC
#34543  maxatome encoding/json: Unmarshal & json.(*Decoder).Token report https://github.com/golang/go/issues/34543 2019-09-25 22:13:24 +0000 UTC
#32779       rsc encoding/json: memoize strings during decode https://github.com/golang/go/issues/32779 2019-06-25 21:08:30 +0000 UTC
#40128  rogpeppe proposal: encoding/json: garbage-free reading of tokens https://github.com/golang/go/issues/40128 2020-07-09 07:58:19 +0000 UTC
#40982   Segflow encoding/json: use different error type for unknown fie https://github.com/golang/go/issues/40982 2020-08-22 21:07:03 +0000 UTC
#59053   joerdav proposal: encoding/json: add a generic Decode function https://github.com/golang/go/issues/59053 2023-03-15 16:20:31 +0000 UTC
#65691  Merovius encoding/xml: Decoder does not reject xml-ProcInst prec https://github.com/golang/go/issues/65691 2024-02-13 10:33:20 +0000 UTC
#14750 cyberphon encoding/json: parser ignores the case of member names https://github.com/golang/go/issues/14750 2016-03-10 13:04:44 +0000 UTC
#40127  rogpeppe encoding/json: add Encoder.EncodeToken method https://github.com/golang/go/issues/40127 2020-07-09 07:52:47 +0000 UTC
#16212 josharian encoding/json: do all reflect work before decoding https://github.com/golang/go/issues/16212 2016-06-29 16:07:36 +0000 UTC
#41144 alvaroale encoding/json: Unmarshaler breaks DisallowUnknownFields https://github.com/golang/go/issues/41144 2020-08-31 14:27:19 +0000 UTC
#64847 zephyrtro encoding/json: UnmarshalJSON methods of embedded fields https://github.com/golang/go/issues/64847 2023-12-22 17:08:52 +0000 UTC
#56332    gansvv encoding/json: clearer error message for boolean like p https://github.com/golang/go/issues/56332 2022-10-19 19:30:20 +0000 UTC
#43513 Alexander encoding/json: add line number to SyntaxError https://github.com/golang/go/issues/43513 2021-01-05 10:59:27 +0000 UTC
#22752  buyology proposal: encoding/json: add access to the underlying d https://github.com/golang/go/issues/22752 2017-11-15 23:46:13 +0000 UTC
#33835     Qhesz encoding/json: unmarshalling null into non-nullable gol https://github.com/golang/go/issues/33835 2019-08-26 10:27:12 +0000 UTC
#33854     Qhesz encoding/json: unmarshal option to treat omitted fields https://github.com/golang/go/issues/33854 2019-08-27 00:20:25 +0000 UTC
//!-textoutput
*/
