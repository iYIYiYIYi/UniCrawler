package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
)

type School struct {
	Link        string
	Name        string
	PublishTime string
	Title		string
	Detail      string
}

var last_time string
var counter int = 0

func InitDatabase() {
	last_t, err := os.OpenFile("last_time",os.O_RDWR|os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	var buf [128]byte
    var content []byte
    for {
        n, err := last_t.Read(buf[:])
        if err == io.EOF {
            // 读取结束
            break
        }
        if err != nil {
            fmt.Println("read file err ", err)
            return
        }
        content = append(content, buf[:n]...)
    }
	last_time = string(content)
	if last_time == "" {
		last_time = "2023-02-22 01:10"
	}
	fmt.Println("Last fresh time:"+last_time)
	last_t.Close()
}

var last_t_tmp string
func AddSchools(schools []*School) {
	csv_f, err := os.OpenFile("schools.csv",os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	
	if err != nil {
		panic("读取文件出错")
	}
	writer := csv.NewWriter(csv_f)
	for index, school:=range schools {
		line := []string{strconv.Itoa(index + 1) ,school.Name, school.Link, school.PublishTime, school.Title, school.Detail}
        err := writer.Write(line)
        if err != nil {
            panic(err)
        }
		if last_t_tmp == "" {
			last_t_tmp = school.PublishTime
		} else {
			ltt, _ := StrToTimeExpand(last_t_tmp)
			ptt, _ := StrToTimeExpand(school.PublishTime)
			if ltt.Before(ptt) {
				last_t_tmp = school.PublishTime
			}
		}
	}
	writer.Flush()
	counter += len(schools)
	csv_f.Close()
}

func Notification() {
	if counter <= 0 {
		return
	}
	beeep.Notify("新的消息",strconv.Itoa(counter)+"条调剂信息被添加...","")
	if last_t_tmp == "" {
		return
	}
	last_time = last_t_tmp
	counter = 0
	last_t, err := os.OpenFile("last_time",os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	last_t.Write([]byte(last_time))
	last_t.Close()
}

func GetLastTime() string {
	if last_time == "" {
		last_time = "2023-02-20 01:10"
	}
	return last_time
}

func GetAllVisited(visited *map[string]bool) {
	csv_f, err := os.OpenFile("schools.csv",os.O_RDWR|os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
        panic(err)
    }
	reader := csv.NewReader(csv_f)
    reader.FieldsPerRecord = -1
	for {
		item, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		school := new(School)
		school.Name = item[1]
		school.Link = item[2]
		school.PublishTime = item[3]
		school.Title = item[4]
		school.Detail = item[5]

		(*visited)[school.Link] = true
    }
	csv_f.Close()
}

func ReadAllSchools() *[]School{
	csv_f, err := os.OpenFile("schools.csv",os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
        panic(err)
    }
	visited := map[string]bool{}
	reader := csv.NewReader(csv_f)
    reader.FieldsPerRecord = -1
    record, err := reader.ReadAll()
    if err != nil {
        panic(err)
    }

	schools := make([]School, 0)
	for _, item := range record {
		if visited[item[2]] {
			continue
		}
		school := new(School)
		school.Name = item[1]
		school.Link = item[2]
		school.PublishTime = item[3]
		school.Title = item[4]
		school.Detail = item[5]
		schools = append(schools, *school)
		visited[school.Link] = true
    }
	csv_f.Close()
	return &schools
}

func ReadSchoolsWithKeywords(keywords []string) *[]School{
	csv_f, err := os.OpenFile("schools.csv",os.O_RDWR|os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
        panic(err)
    }
	reader := csv.NewReader(csv_f)
    reader.FieldsPerRecord = -1
	schools := make([]School, 0)
	visited := map[string]bool{}
	for {
		item, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if visited[item[2]] {
			continue
		}
		school := new(School)
		school.Name = item[1]
		school.Link = item[2]
		school.PublishTime = item[3]
		school.Title = item[4]
		school.Detail = item[5]

		for _, keyword := range keywords {
			if strings.Contains(school.Name, keyword) || strings.Contains(school.Title, keyword) || strings.Contains(school.Detail, keyword) {
				schools = append(schools, *school)
				visited[school.Link] = true
				break
			}
		}
    }
	csv_f.Close()
	return &schools
}

func ReadAfterTime(tim *time.Time) *[]School {
	csv_f, err := os.OpenFile("schools.csv",os.O_RDWR|os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
        panic(err)
    }
	reader := csv.NewReader(csv_f)
    reader.FieldsPerRecord = -1
	visited := map[string]bool{}
	schools := make([]School, 0)
	for {
		item, err := reader.Read()
		if item == nil || err == io.EOF {
			break
		}
		u,err := StrToTimeExpand(item[3])
		if err != nil || tim.After(u) {
			continue
		}

		if visited[item[2]] {
			continue
		}
		school := new(School)
		school.Name = item[1]
		school.Link = item[2]
		school.PublishTime = item[3]
		school.Title = item[4]
		school.Detail = item[5]
		schools = append(schools, *school)
		visited[school.Link] = true
	}
	csv_f.Close()
	return &schools
}

func StrToTime(time_s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", time_s)
}

func StrToTimeExpand(time_s string)  (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", (time_s + ":00"))
}