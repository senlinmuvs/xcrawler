package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

const (
	version = "1.0.1"

	Rule_RefPipe = "@"
)

var (
	help                    bool
	v                       bool
	printBody               bool
	printReqHeader          bool
	printRespHeader         bool
	debug                   bool
	header                  string
	headerFile              string
	url                     string
	params                  string
	nextUrlParamsRule       string
	firstPipeParams         string
	totalPage               int
	endParams               string
	stopCont                string
	stopOnRespSamePre       bool
	stopOnJsonSelectorEmpty bool
	selector                string
	jsonSelector            string
	regex                   string
	file                    string
	dir                     string
	sleepRule               string
	longSleepRule           string
	proxy                   string
	horStyle                bool
	compatStyle             bool
	quote                   bool
	quoteSub                bool
	sep                     string
	timeout                 int
	retry                   int

	input    string
	inputSep string
	//
	inputItems []string

	httpClient = resty.New()
	reqCount   = 0
	startTime  = int64(0)
)

func init() {
	flag.BoolVar(&help, "help", false, "print help info")
	flag.BoolVar(&v, "v", false, "show version")
	flag.BoolVar(&printBody, "pb", false, "print body")
	flag.BoolVar(&printReqHeader, "qh", false, "print request header")
	flag.BoolVar(&printRespHeader, "ph", false, "print response header")
	flag.BoolVar(&debug, "X", false, "debug")
	flag.StringVar(&header, "hd", "", "设置请求头k1=v1,k2=v2")
	flag.StringVar(&headerFile, "hf", "", "请求头在此文件")
	flag.StringVar(&url, "u", "", "url或起始url模版")
	flag.StringVar(&params, "p", "", "初始参数,即第一页url模版的参数,多个逗号分隔")
	flag.StringVar(&nextUrlParamsRule, "np", "", "下一页的url填充参数的变化规律,多个逗号分隔。\n+n:表示在上一页参数基础上加n")
	flag.StringVar(&firstPipeParams, "pp", "", "此参数只在有管道输入时有效，表示只有第一个管道数据使用这个初始参数，之后的都使用-p的初始参数")
	flag.IntVar(&totalPage, "tp", 0, "总共翻多少页")
	flag.StringVar(&endParams, "ep", "", "最后一个参数")
	flag.StringVar(&stopCont, "sc", "", "停止内容：当内容包含此字符串时停止翻页")
	flag.BoolVar(&stopOnRespSamePre, "ss", false, "当此页内容与上一页内容相同时停止翻页")
	flag.BoolVar(&stopOnJsonSelectorEmpty, "se", false, "当json选择器内容为空时停止翻页")
	flag.StringVar(&selector, "s", "", "html选择器,多个逗号分隔")
	flag.StringVar(&jsonSelector, "js", "", "json选择器,多个逗号分隔")
	flag.StringVar(&regex, "r", "", "正则表达式,多个逗号分隔")
	flag.StringVar(&file, "f", "", "从文件中取url,每行一个")
	flag.StringVar(&dir, "d", "", "保存目录,默认当前")
	flag.StringVar(&inputSep, "is", "\t", "pipe input separeator")
	flag.StringVar(&sleepRule, "sr", "", "10,1 表示每10次请求休息1秒\n1-10,1-5 表示每1-10次请休息随机1-5秒")
	flag.StringVar(&longSleepRule, "lsr", "", "100,10 表示每100秒后休息10秒\n50-100,10-20 表示每50-100秒后息随机休息10-20秒")
	flag.StringVar(&proxy, "px", "", "proxy")
	flag.IntVar(&timeout, "to", 10, "timeout seconds")
	flag.IntVar(&retry, "rt", 3, "retry count")
	flag.BoolVar(&horStyle, "hs", false, "是否横列输出")
	flag.BoolVar(&compatStyle, "cs", false, "是否开启紧凑输出，多个空格和换行替换为一个")
	flag.BoolVar(&quote, "q", false, "是否打引号")
	flag.BoolVar(&quoteSub, "qs", false, "选择器的选择项有多个时，每个结果是否打引号")
	flag.StringVar(&sep, "sp", "\n", "一个选择器对应多个结果时的分隔符")
	flag.Usage = usage
}

func main() {
	fileInfo, _ := os.Stdin.Stat()
	if (fileInfo.Mode() & os.ModeNamedPipe) == os.ModeNamedPipe {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			input += s.Text() + "\n"
		}
		if input != "" {
			inputItems = strings.Split(strings.Trim(input, " \n"), "\n")
		}
	}

	flag.Parse()
	if v {
		fmt.Printf("crawler %s\n", version)
	}
	if help {
		flag.Usage()
	}
	if longSleepRule != "" {
		startTime = CurMills()
	}
	if url != "" {
		if inputItems != nil {
			for i := 0; i < len(inputItems); i++ {
				downWebPage(i, inputItems[i])
			}
		} else {
			downWebPage(0, nil)
		}
	}
}

func fillUrlTmpl(tmpl string, params []string) string {
	for _, p := range params {
		tmpl = strings.Replace(tmpl, "{}", p, 1)
	}
	return tmpl
}

func downWebPage(i int, pipeParam interface{}) {
	page := 0
	var paramsArr []string
	var nextUrlParamsRuleArr []string
	curUrl := url
	var preRespBytes []byte
	for {
		if debug {
			fmt.Println(CurTimeStr(), fmt.Sprintf("%d/%d", i+1, len(inputItems)), page+1)
		}
		if paramsArr != nil {
			curUrl = fillUrlTmpl(url, paramsArr)
		} else {
			curParams := params
			if firstPipeParams != "" && i == 0 {
				curParams = firstPipeParams
			}
			if curParams != "" {
				paramsArr = Split(curParams, ",")
				if inputItems != nil {
					for i, p := range paramsArr {
						if p == Rule_RefPipe {
							paramsArr[i] = pipeParam.(string)
						}
					}
				}
				curUrl = fillUrlTmpl(url, paramsArr)
			} else {
				if pipeParam != nil {
					paramsArr = []string{pipeParam.(string)}
					curUrl = fillUrlTmpl(url, paramsArr)
				}
			}
		}
		checkSleep()
		if debug {
			fmt.Println(curUrl)
		}
		//
		req := httpClient.R()
		httpClient.SetTimeout(time.Duration(timeout) * time.Second)
		httpClient.SetRetryCount(retry)
		if proxy != "" {
			httpClient.SetProxy(proxy)
		}
		setHeader(req)
		if printReqHeader {
			printHeader("Req  Header", req.Header)
		}
		resp, e := req.Get(curUrl)
		reqCount += 1
		if e != nil {
			panic(fmt.Errorf("downWebPage err %s", e.Error()))
		}
		if printRespHeader {
			printHeader("Resp Header", resp.Header())
		}
		page += 1
		datas := resp.Body()
		body := string(datas)
		if debug {
			fmt.Println(len(body), "bytes", resp.StatusCode())
		}
		if printBody {
			fmt.Println(body)
		}
		if stopOnRespSamePre {
			if preRespBytes != nil && ByteArrEq(datas, preRespBytes) {
				break
			}
			preRespBytes = datas
		}
		if selector != "" {
			doSelector(body)
		} else if jsonSelector != "" {
			stop := doJsonSelector(datas)
			if stop {
				break
			}
		}
		if endParams != "" {
			if ArrEqIg(strings.Split(endParams, ","), paramsArr, "@") {
				break
			}
		}
		if nextUrlParamsRule == "" {
			break
		} else {
			if totalPage > 0 {
				if page == totalPage {
					break
				}
			}
			if stopCont != "" {
				if strings.Contains(body, stopCont) {
					break
				}
			}
			//
			nextUrlParamsRuleArr = Split(nextUrlParamsRule, ",")
			if len(nextUrlParamsRuleArr) != len(paramsArr) {
				panic("-p与-np参数个数不一致")
			}
			for i, pr := range nextUrlParamsRuleArr {
				if pr == Rule_RefPipe {
					if pipeParam != nil {
						paramsArr[i] = pipeParam.(string)
					}
				} else {
					n := getRulePlusN(pr)
					if n != 0 {
						paramsArr[i] = strconv.Itoa(ToInt(paramsArr[i]) + n)
					}
				}
			}
		}
	}
}

func setHeader(req *resty.Request) {
	if headerFile != "" {
		lines := ReadLines(headerFile)
		for _, line := range lines {
			arr := strings.Split(line, ":")
			if len(arr) > 1 {
				req.Header[arr[0]] = []string{strings.Trim(arr[1], " ")}
			}
		}
	}
	if header != "" {
		arr := strings.Split(header, ",")
		for _, ar := range arr {
			kv := strings.Split(ar, ":")
			if len(kv) > 1 {
				req.Header[kv[0]] = []string{strings.Trim(kv[1], " ")}
			}
		}
	}
}

func doSelector(body string) {
	doc, e := goquery.NewDocumentFromReader(strings.NewReader(body))
	if e != nil {
		panic(fmt.Errorf("body doc err %s", e.Error()))
	}
	selArr := strings.Split(selector, ",")
	resArr := [][]string{}
	for _, sel := range selArr {
		ind := strings.Index(sel, "|")
		xfunc := ""
		if ind >= 0 {
			xfunc = sel[ind+1:]
			sel = sel[:ind]
		}
		attr := ""
		attrI0 := strings.Index(sel, "[")
		attrI1 := 0
		if attrI0 >= 0 {
			attrI1 = strings.Index(sel, "]")
			if attrI1 >= 0 && attrI1 > attrI0 {
				attr = SubUnicode(sel, attrI0+1, attrI1)
				sel = SubUnicode(sel, 0, attrI0)
			}
		}
		arr := []string{}
		doc.Find(sel).Each(func(i int, s *goquery.Selection) {
			x := ""
			if attr == "" {
				x = strings.TrimSpace(s.Text())
			} else {
				x = strings.TrimSpace(s.AttrOr(attr, ""))
			}
			if compatStyle {
				d1, _ := hex.DecodeString("0a20")
				x = strings.ReplaceAll(x, string(d1), " ")
				d2, _ := hex.DecodeString("200a")
				x = strings.ReplaceAll(x, string(d2), "")
				x = strings.ReplaceAll(x, " ", " ")
				x = strings.ReplaceAll(x, " \n", "\n")
				x = strings.ReplaceAll(x, "\n ", "\n")
				x = Compact(x, " ", "\n")
				x = strings.ReplaceAll(x, string(d1), "\n")
			}
			if xfunc != "" {
				x = doFunc(x, xfunc)
			}
			arr = append(arr, x)
		})
		resArr = append(resArr, arr)
	}
	if horStyle {
		printHorStyle(resArr)
	} else {
		printVerStyle(resArr)
	}
}
func doFunc(x, xfunc string) string {
	xfunc = Compact(xfunc, " ")
	arr := strings.Split(xfunc, " ")
	if len(arr) > 2 {
		if arr[0] == "cut" {
			x = strings.Trim(x, arr[1])
			ar := strings.Split(x, arr[1])
			ind := ToInt(arr[2])
			if ind >= 0 && len(ar) > ind {
				return ar[ind]
			}
		}
	}
	return x
}
func doJsonSelector(data []byte) (stop bool) {
	resArr := [][]string{}
	arr := strings.Split(jsonSelector, ",")
	for _, js := range arr {
		ar := strings.Split(js, ".")
		lastKey := ""
		if len(ar) > 1 {
			lastKey = ar[len(ar)-1]
			ar = ar[:len(ar)-1]
		}
		arr := []string{}
		jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, e error) {
			if e != nil {
				panic(fmt.Errorf("json err %s", e.Error()))
			}
			if len(value) > 0 {
				var v string
				v, e = jsonparser.GetString(value, lastKey)
				if e != nil {
					panic(fmt.Errorf("json k %s err %s", lastKey, e.Error()))
				}
				arr = append(arr, v)
			}
		}, ar...)
		if len(arr) > 0 {
			resArr = append(resArr, arr)
		} else {
			stop = true
		}
	}
	if horStyle {
		printHorStyle(resArr)
	} else {
		printVerStyle(resArr)
	}
	return
}

func getRulePlusN(r string) int {
	if len(r) > 1 {
		n := ToInt(r[1:])
		return n
	}
	return 0
}

func checkSleep() {
	if sleepRule != "" {
		srArr := strings.Split(sleepRule, ",")
		scArr := strings.Split(srArr[0], "-")
		sleepWhenReqCount := 0
		if len(scArr) > 1 {
			s0 := ToInt(scArr[0])
			s1 := ToInt(scArr[1])
			if s0 >= 0 && s1 >= 0 && s1 > s0 {
				sleepWhenReqCount = Rand(s1-s0) + s0
			}
		} else {
			sleepWhenReqCount = ToInt(srArr[0])
		}
		//
		secArr := strings.Split(srArr[1], "-")
		sleepSecond := 0
		if len(secArr) > 1 {
			s0 := ToInt(secArr[0])
			s1 := ToInt(secArr[1])
			if s0 >= 0 && s1 >= 0 && s1 > s0 {
				sleepSecond = Rand(s1-s0) + s0
			}
		} else {
			sleepSecond = ToInt(srArr[1])
		}
		if sleepWhenReqCount > 0 && sleepSecond > 0 {
			if reqCount >= sleepWhenReqCount {
				if debug {
					fmt.Println("sleep", sleepSecond, "sec", "req", sleepWhenReqCount)
				}
				time.Sleep(time.Duration(sleepSecond) * time.Second)
				reqCount = 0
			}
		}
	}
	//
	if longSleepRule != "" {
		cur := CurMills()
		arr := strings.Split(longSleepRule, ",")
		if len(arr) < 2 {
			panic(fmt.Errorf("-lsr参数错误 %s", longSleepRule))
		}
		dur := int64(0)
		durArr := strings.Split(arr[0], "-")
		if len(durArr) > 1 {
			d0 := ToInt(durArr[0])
			d1 := ToInt(durArr[1])
			if d0 >= 0 && d1 >= 0 && d1 > d0 {
				dur = int64(Rand(d1-d0) + d0)
			}
		}
		if dur > 0 && cur-startTime >= dur*1000 {
			secArr := strings.Split(arr[1], "-")
			sec := 0
			if len(secArr) > 1 {
				s0 := ToInt(secArr[0])
				s1 := ToInt(secArr[1])
				if s0 >= 0 && s1 >= 0 && s1 > s0 {
					sec = Rand(s1-s0) + s0
				}
			}
			if debug {
				fmt.Println("sleep", sec, "sec", "dur", dur, "sec")
			}
			time.Sleep(time.Duration(sec) * time.Second)
			startTime = cur
		}
	}
}

func printHeader(tag string, header http.Header) {
	maxLen := 10
	for k := range header {
		le := len(k)
		if maxLen < le {
			maxLen = le
		}
	}
	fmt.Println("------------------", tag, "------------------")
	for k, v := range header {
		fmt.Printf("%"+strconv.Itoa(maxLen)+"s %s\n", k, v)
	}
	fmt.Println("-------------------------------------------------")
}
func printHorStyle(arr [][]string) {
	l := len(arr)
	for i, ar := range arr {
		le := len(ar)
		item := ""
		for j, a := range ar {
			if quoteSub && len(ar) > 1 {
				a = strconv.Quote(a)
			}
			if j == le-1 {
				item += a
			} else {
				item += a + sep
			}
		}
		if quote {
			item = strconv.Quote(item)
		}
		if item != "" {
			fmt.Print(item)
		}
		if i != l-1 {
			fmt.Print(",")
		}
	}
	fmt.Println()
}
func printVerStyle(arr [][]string) {
	for _, ar := range arr {
		le := len(ar)
		item := ""
		for j, a := range ar {
			if quoteSub && len(ar) > 1 {
				a = strconv.Quote(a)
			}
			if j == le-1 {
				item += a
			} else {
				item += a + sep
			}
		}
		if quote {
			item = strconv.Quote(item)
		}
		if item != "" {
			fmt.Print(item)
		}
		fmt.Println()
	}
}
func usage() {
	fmt.Printf("xcrawler version: %s\nOptions:\n", version)
	flag.PrintDefaults()
}
