package operator

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Fc struct{
	Name string
	Reg *regexp.Regexp
	Value *Fc
	Content string
	Ext []string
}

var (
	FUNC_COUNT          = regexp.MustCompile(`^COUNT\((.*?)\)$`)
	FUNC_DISTINCT       = regexp.MustCompile(`^DISTINCT\((.*?)\)$`)
	FUNC_SUM            = regexp.MustCompile(`^SUM\((.*?)\)$`)
	FUNC_MAX            = regexp.MustCompile(`^MAX\((.*?)\)$`)
	FUNC_MIN            = regexp.MustCompile(`^MIN\((.*?)\)$`)
	FUNC_AVG            = regexp.MustCompile(`^AVG\((.*?)\)$`)
	FUNC_FROM_UNIXTIME  = regexp.MustCompile(`^FROM_UNIXTIME\((.*?)\)$`)
	FUNC_UNIX_TIMESTAMP = regexp.MustCompile(`^UNIX_TIMESTAMP\((.*?)\)$`)
)

var funcMap = map[string]*regexp.Regexp{
	"COUNT":          FUNC_COUNT,
	"DISTINCT":       FUNC_DISTINCT,
	"SUM":            FUNC_SUM,
	"MAX":            FUNC_MAX,
	"MIN":            FUNC_MIN,
	"AVG":            FUNC_AVG,
	"FROM_UNIXTIME":  FUNC_FROM_UNIXTIME,
	"UNIX_TIMESTAMP": FUNC_UNIX_TIMESTAMP,
}

var timeTempMap = map[string]string{
	"%Y": "2006",
	"%M": "01",
	"%D": "02",
	"%H": "15",
	"%I": "04",
	"%S": "05",
	"%y": "2006",
	"%m": "01",
	"%d": "02",
	"%h": "15",
	"%i": "04",
	"%s": "05",
	"%EM": "Jan",
	"%em": "Jan",
}

const (
	MAX_UINT64 = ^uint64(0)
	MIN_UINT64 = uint64(0)
	EXT_DECHAR = "::"
)

func (fc *Fc) parseFunc() {
	for fn,reg := range funcMap{
		ret := reg.FindAllStringSubmatch(fc.Content, -1)
		if len(ret) > 0 {
			// 处理函数扩展入参
			ext := strings.Split(ret[0][1], EXT_DECHAR)
			f := &Fc{
				Name:   fn,
				Reg:     reg,
				Value:   nil,
				Content: ext[0],
				Ext:ext,
			}
			fc.Value = f
			f.parseFunc()
		}
	}
	return
}

func (f *Fc) _unix_timestamp(data []string) []string {
	ret := []string{}
	for _, v := range data {
		timeSchema := "2006-01-02 15:04:05"
		if len(f.Ext) > 1 {
			format := strings.Replace(f.Ext[1], "\"", "", -1)
			format = strings.Replace(format, "'", "", -1)
			timeSchema = timeFormatMap(format)
		}
		loc, _ := time.LoadLocation("Asia/Shanghai")                //设置时区
		tt, _ := time.ParseInLocation(timeSchema, v, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
		ti := strconv.FormatInt(tt.Unix(), 10)
		ret = append(ret, ti)
	}
	return ret
}

func (f *Fc) _from_unixtime(data []string) []string {
	ret := []string{}
	for _, v := range data {
		timeSchema := "2006-01-02 15:04:05"
		if len(f.Ext) > 0 {
			format := strings.Replace(f.Ext[1], "\"", "", -1)
			format = strings.Replace(format, "'", "", -1)
			timeSchema = timeFormatMap(format)
		}
		ts, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			ts = 0
		}
		tm := time.Unix(ts, 0).Format(timeSchema)
		ret = append(ret, tm)
	}
	return ret
}

func timeFormatMap(template string) string {
	for t, v := range timeTempMap {
		template = strings.Replace(template, t, v, -1)
	}
	return template
}

func (f *Fc) _avg(data []string) float64 {
	ret := float64(0)
	for _, v := range data {
		//intNum, _ := strconv.Atoi(v)
		//int64Num := uint64(intNum)
		float64Num, _ := strconv.ParseFloat(v, 64)
		//float64Num, _ = strconv.ParseFloat(fmt.Sprintf("%.10f", float64Num), 64)
		ret += float64Num
	}
	return ret / float64(len(data))
}

func (f *Fc) _min(data []string) float64 {
	ret := float64(MAX_UINT64)
	for _, v := range data {
		float64Num, _ := strconv.ParseFloat(v, 64)
		//float64Num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64Num), 64)
		if float64Num < ret {
			ret = float64Num
		}
	}
	return ret
}

func (f *Fc) _max(data []string) float64 {
	ret := float64(MIN_UINT64)
	for _, v := range data {
		float64Num, _ := strconv.ParseFloat(v, 64)
		//float64Num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64Num), 64)
		if float64Num > ret {
			ret = float64Num
		}
	}
	return ret
}

func (f *Fc) _count(data []string) int64 {
	return int64(len(data))
}

func (f *Fc) _sum(data []string) float64 {
	ret := float64(0)
	for _, v := range data {
		//int64, err := strconv.ParseInt(v, 10, 64)
		float64Num, err := strconv.ParseFloat(v, 64)
		if err != nil {
			float64Num = 0
		}
		ret += float64Num
	}
	return ret
}

func (f *Fc) _distinct(data []string) []string {
	distinctMap := map[string]int{}
	ret := []string{}
	for _, v := range data {
		if _, ok := distinctMap[v]; ok != true{
			distinctMap[v] = 1
			ret = append(ret, v)
			//ret[i] = v
		//} else {
		//	ret = append(ret, "NULL")
		}
	}
	return ret
}

func (f *Fc) _get(key string, data []map[string]string) []string {
	ret := []string{}
	for _, v := range data {
		ret = append(ret, v[key])
	}
	return ret
}
