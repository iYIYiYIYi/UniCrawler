package crawler

import (
	"UniCrawler/cmd/util"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var visitFunc func(time_s, school, stype string, should_visit *bool) bool

func Init(f func(time_s, school, stype string, should_visit *bool) bool) {
    visitFunc = f
}

func DefaultInit() {
    visitFunc = func(time_s, school, stype string, should_visit *bool) bool {
        if time_s != "" && school != "" && stype != "" {
            st,err := util.StrToTimeExpand(time_s)
            lst,_ := util.StrToTimeExpand(util.GetLastTime())
            if err != nil {
                return false
            }
            result := st.After(lst)
            *should_visit = result
            return result
        }
        return false
    }
}

func judgeShouldVisit(time_s, school, stype string, should_visit *bool) bool {
    if visitFunc != nil {
        return visitFunc(time_s, school, stype, should_visit)
    }
    return true
}

func Start() {
    var visited = map[string]bool{}
    util.GetAllVisited(&visited)
    // Instantiate default collector
    c := colly.NewCollector(
        colly.AllowedDomains("muchong.com"),
    )
    extensions.RandomUserAgent(c)
    school_m := make(map[string]util.School, 0)
    schools := []*util.School{}

    // 我们认为匹配该模式的是该网站的详情页
    detailRegex, _ := regexp.Compile(`/t-\d*-1$`)
   // 匹配下面模式的是该网站的列表页
    listRegex, _ := regexp.Compile(`/bbs/kaoyan\.php(\?(page=\d*|)(&page=\d*))*$`)
    should_visit := true
    // 所有 a 标签，上设置回调函数
    c.OnHTML("html", func(e *colly.HTMLElement) {
        link := e.Request.URL.Path
        if e.Request.URL.RawQuery != "" {
            link += "?" + e.Request.URL.RawQuery
        }

        // 已访问过的详情页或列表页，跳过
        if visited[link] && (detailRegex.Match([]byte(link)) || listRegex.Match([]byte(link))) {
            return
        }

        // 既不是列表页，也不是详情页
        // 那么不是我们关心的内容，要跳过
        if !detailRegex.Match([]byte(link)) && !listRegex.Match([]byte(link)) {
            println("not match", link)
            return
        }

        // 因为大多数网站有反爬虫策略
        // 所以爬虫逻辑中应该有 sleep 逻辑以避免被封杀
        time.Sleep(time.Second)
        dom := e.DOM
        println("match", link)

        visited[link] = true
        if detailRegex.Match([]byte(link)) {
            key := "http://"+e.Request.URL.Host+link
            s, ok := school_m[key]
            if !ok {
                return
            }
            // 详细信息页
            tables := dom.Find("div.forum_Mix+div.t_fsz").Find("td")
            s.Detail = strings.Replace(tables.Text(), "\n", "",1)

            schools = append(schools, &s)
            if len(schools) >= 25 {
                util.AddSchools(schools)
                schools = schools[:0]
            }
            delete(school_m, key)
        }

        if listRegex.Match([]byte(link)) {
            // 列表页
            tbody := dom.Find("tbody.forum_body_manage tr")
            tbody.Each(func(i int, s *goquery.Selection) {
                a_href := s.Find("a.xmc_ft12")
                alink, _ := a_href.Attr("href")
                fmt.Println(alink)
                infos := strings.Split(strings.TrimSpace(s.Text()), "\n")
                title := infos[0]
                name := infos[2]
                stype := infos[3]
                tim := infos[5]
                if judgeShouldVisit(tim, name, stype, &should_visit) && !visited[alink] {
                    fmt.Println("Visiting "+ name + " on " + alink)
                    fmt.Println(infos)
                    school_m[alink] = util.School{
                        Name: name,
                        Link: alink,
                        PublishTime: tim,
                        Title: title,
                    }
                    c.Visit(alink)
                }
            })
            
            if should_visit {
                dom.Find("tr.smalltxt").Find("td>a").Each(func(i int, s *goquery.Selection) {
                    next_link, exist := s.Attr("href")
                    if exist {
                        page := s.Text()
                        page_i,err := strconv.Atoi(page)
                        if err == nil && page_i <= 350{
                            fmt.Println("Starting access next page:"+next_link)
                            c.Visit("http://"+e.Request.URL.Host+"/bbs/"+next_link)
                        }
                    }
                })
            } else {
                if len(schools) > 0 {
                    util.AddSchools(schools)
                }
                fmt.Println("Stop Visiting...")
            }
        }

        time.Sleep(time.Millisecond * 2)
    })

    c.SetRequestTimeout(25 * time.Second)
    err := c.Visit("http://muchong.com/bbs/kaoyan.php")
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    c.Wait()
    util.Notification()
}