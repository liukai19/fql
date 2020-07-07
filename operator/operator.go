package operator

import (
	"bufio"
	"errors"
	"fql/syntax"
	"fql/util"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Oper struct {
	Syntax    *syntax.Syntax
	Delimiter string
	Data      []map[string]string
	GroupData map[string][]map[string]string
}

var processField = map[string]int{}

var outputField = map[string]bool{}

func (o *Oper) Process() {
	o.Source().GroupBy().Field().AfterWhere().OrderBy().Limit().Output()
}

func (o *Oper) Output() {
	// 屏蔽 IsShow === false的列
	if len(o.Data) > 0{
		for i,v := range o.Data{
			for f,_ := range v{
				if false == outputField[f] {
					delete(o.Data[i], f)
				}
			}
		}
	}
}

func (o *Oper) Limit() *Oper {
	if o.Syntax.Limit != nil {
		if o.Syntax.Limit.Limit >= len(o.Data) {
			o.Data = []map[string]string{}
			return o
		}
		min := o.Syntax.Limit.Limit
		max := o.Syntax.Limit.Limit + o.Syntax.Limit.Offset
		if max >= len(o.Data) {
			max = len(o.Data)
		}
		o.Data = o.Data[min:max]
	}
	return o
}

func (o *Oper) Field() *Oper {
	if len(o.Data) > 0 {
		o.Data = o.setSxtractByField(o.Data)
	}
	// 如果有 group 将 group 处理后转成 o.Data 格式
	if len(o.GroupData) > 0 {
		o.Data = []map[string]string{}
		for _, list := range o.GroupData {
			r := o.setSxtractByField(list)
			o.Data = append(o.Data, r...)
		}
	}
	return o
}

func (o *Oper) setSxtractByField(data []map[string]string) []map[string]string {
	fRet := map[string][]string{}
	for _, i := range o.Syntax.Fields.Value {
		fRet[i.Remark] = fieldProcess(i.Name, data)
	}
	for f, _ := range processField {
		if _,ok := fRet[f];ok!=true{
			fRet[f] = fieldProcess(f, data)
		}
	}
	r := []map[string]string{}
	for i := 0; i < len(data); i++ {
		item := map[string]string{}
		// 如果有字段都为 null 那么不输出该行
		flag := 0
		for _, f := range o.Syntax.Fields.Value {
			outputField[f.Remark] = f.IsShow
			if i >= len(fRet[f.Remark]) {
				item[f.Remark] = "NULL"
				flag++
			} else {
				item[f.Remark] = fRet[f.Remark][i]
			}
		}
		if flag == 0 {
			r = append(r, item)
		}
	}
	return r
}

func fieldProcess(name string, source []map[string]string) []string {
	if len(source) <= 0 {
		return nil
	}
	fc := &Fc{
		Name:    "Common",
		Content: name,
	}
	// 字段函数解析
	fc.parseFunc()
	fcList := []*Fc{}
	f := fc.Value
	for f != nil {
		fc.Content = f.Content
		fcList = append(fcList, f)
		f = f.Value
	}
	fcList = append(fcList, fc)
	// 反转
	for i, j := 0, len(fcList)-1; i < j; i, j = i+1, j-1 {
		fcList[i], fcList[j] = fcList[j], fcList[i]
	}
	data := []string{}
	for _, v := range fcList {
		switch v.Name {
		case "COUNT":
			r := v._count(data)
			data = []string{}
			data = append(data, strconv.Itoa(int(r)))
			//data = nil
			break
		case "DISTINCT":
			data = v._distinct(data)
			break
		case "SUM":
			r := v._sum(data)
			data = []string{}
			data = append(data, strconv.Itoa(int(r)))
			//data = nil
			break
		case "MAX":
			r := v._max(data)
			data = []string{}
			data = append(data, strconv.Itoa(int(r)))
			//data = nil
			break
		case "MIN":
			r := v._min(data)
			data = []string{}
			data = append(data, strconv.Itoa(int(r)))
			//data = nil
			break
		case "FROM_UNIXTIME":
			data = v._from_unixtime(data)
			break
		case "UNIX_TIMESTAMP":
			data = v._unix_timestamp(data)
			break
		default:
			data = v._get(v.Content, source)
			break
		}
	}
	return data
}

func (o *Oper) OrderBy() *Oper {
	if o.Syntax.Order != nil {
		o.Data = util.ArraySort(o.Data, o.Syntax.Order.Field.Remark, o.Syntax.Order.Method)
	}
	return o
}

func (o *Oper) GroupBy() *Oper {
	if o.Syntax.Group == nil {
		return o
	}

	if len(o.Data) <= 0 {
		return o
	}
	o.GroupData = map[string][]map[string]string{}
	for _, v := range o.Data {
		key := []string{}
		for _, gf := range o.Syntax.Group.Value {
			key = append(key, v[gf.Name])
		}
		keyStr := util.Md5V(strings.Join(key, "-"))
		o.GroupData[keyStr] = append(o.GroupData[keyStr], v)
	}
	return o
}

func (o *Oper) AfterWhere() *Oper {
	if o.Syntax.Where == nil {
		return o
	}
	if len(o.Data) > 0 {
		o.Data = o.filterByMap(o.Data)
	}
	return o
}

func (o *Oper) filterByMap(data []map[string]string) []map[string]string {
	if len(o.Syntax.Where.Conds) > 0 {
		dataRet := []map[string]string{}
		for _, line := range data {
			condRet := 0
			for _, v := range o.Syntax.Where.Conds {
				d := line[v.Field.Remark]
				if len(v.SetValue) > 0 {
					values := []string{}
					for _, sv := range v.SetValue {
						if s, ok := o.Syntax.Where.ValueHash[sv]; ok != false {
							values = append(values, s)
						} else {
							values = append(values, sv)
						}
					}
					if r := setFiliter(d, values, v.Exps); r != false {
						condRet += 1
					}
				} else {
					value := v.ItemValue
					if sv, ok := o.Syntax.Where.ValueHash[v.ItemValue]; ok != false {
						value = sv
					}
					if r := itemFiliter(d, value, v.Exps); r != false {
						condRet += 1
					}
				}
			}
			if syntax.RELATION_AND == o.Syntax.Where.Relation {
				if condRet == len(o.Syntax.Where.Conds) {
					//data = append(data[:i], data[i+1:]...)
					dataRet = append(dataRet, line)
				}
			} else {
				if condRet > 0 {
					//data = append(data[:i], data[i+1:]...)
					dataRet = append(dataRet, line)
				}
			}
		}
		return dataRet
	} else {
		return data
	}
}

func (o *Oper) BeforeWhere(data []string) bool {
	if o.Syntax.Where == nil {
		return true
	}
	if len(o.Syntax.Where.Conds) > 0 {
		condRet := 0
		for _, v := range o.Syntax.Where.Conds {
			d, err := getEleFromLine(v.Field.Name, data)
			// 非@N 格式字段容错
			if err != nil {
				condRet += 1
				continue
			}
			if len(v.SetValue) > 0 {
				values := []string{}
				for _, sv := range v.SetValue {
					if s, ok := o.Syntax.Where.ValueHash[sv]; ok != false {
						values = append(values, s)
					} else {
						values = append(values, sv)
					}
				}
				if r := setFiliter(d, values, v.Exps); r != false {
					condRet += 1
				}
			} else {
				value := v.ItemValue
				if sv, ok := o.Syntax.Where.ValueHash[v.ItemValue]; ok != false {
					value = sv
				}
				if r := itemFiliter(d, value, v.Exps); r != false {
					condRet += 1
				}
			}
		}

		if syntax.RELATION_AND == o.Syntax.Where.Relation {
			if condRet == len(o.Syntax.Where.Conds) {
				return true
			} else {
				return false
			}
		} else {
			if condRet > 0 {
				return true
			} else {
				return false
			}
		}
	}
	return true
}

func (o *Oper) Source() *Oper {
	//parse.GetInstance("aaa", " ")
	//manager.Do("", "")
	f, err := os.Open(o.Syntax.Source.FilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	o.Data = []map[string]string{}
	// 只保留 Syntax 中出现的字段
	processField = o.extractField()
	// 处理"*" case 标识 只处理一次
	allFieldFlag := false
	for {
		line, _, err := rd.ReadLine()
		if err != nil || io.EOF == err {
			break
		}
		sline := strings.Split(string(line), o.Delimiter)
		// 基于 Where 过滤 前置过滤(存@N 的字段), 减少内存消耗和 CPU计算
		if r := o.BeforeWhere(sline); r != true {
			continue
		}
		// "*"
		if false == allFieldFlag {
			o.Syntax.Fields.AddIndexField(len(sline))
			allFieldFlag = true
			//  处理完"*"后 重新 计算Syntax 中出现的字段
			processField = o.extractField()
			if len(processField) <= 0 {
				panic("Not Found Read Fields")
			}
		}
		// 只保留 Syntax 中出现的字段
		mline := map[string]string{}
		for f, _ := range processField {
			mline[f], _ = getEleFromLine(f, sline)
		}
		o.Data = append(o.Data, mline)
	}
	return o
}

func (o *Oper) extractField() map[string]int {
	uniqueField := map[string]int{}
	if o.Syntax.Fields != nil {
		for _, v := range o.Syntax.Fields.Value {
			if fs := regFieldFiliter(v.Name); len(fs) > 0 {
				for _, f := range fs {
					uniqueField[f] = 1
				}
			}
		}
	}

	if o.Syntax.Group != nil {
		for _, v := range o.Syntax.Group.Value {
			if fs := regFieldFiliter(v.Name); len(fs) > 0 {
				for _, f := range fs {
					uniqueField[f] = 1
				}
			}
		}
	}

	if o.Syntax.Order != nil {
		if fs := regFieldFiliter(o.Syntax.Order.Field.Name); len(fs) > 0 {
			for _, f := range fs {
				uniqueField[f] = 1
			}
		}
	}

	if o.Syntax.Where != nil {
		for _, v := range o.Syntax.Where.Conds {
			if fs := regFieldFiliter(v.Field.Name); len(fs) > 0 {
				for _, f := range fs {
					uniqueField[f] = 1
				}
			}
		}
	}
	return uniqueField
}

func setFiliter(d1 string, d2 []string, exps string) bool {
	d1 = strings.ToUpper(d1)
	switch exps {
	case syntax.EXPS_NOTIN:
		for _, v := range d2 {
			if d1 == v {
				return false
			}
		}
		break
	case syntax.EXPS_IN:
		f := false
		for _, v := range d2 {
			if d1 == v {
				f = true
				break
			}
		}
		if f == false {
			return false
		}
		break
	}
	return true
}

func itemFiliter(d1 string, d2 string, exps string) bool {
	d1 = strings.ToUpper(d1)
	switch exps {
	case syntax.EXPS_EQ:
		if d1 == d2 {
			return true
		}
		break
	case syntax.EXPS_GT:
		d1r, _ := strconv.Atoi(d1)
		d2r, _ := strconv.Atoi(d2)
		if d1r > d2r {
			return true
		}
		break
	case syntax.EXPS_LT:
		d1r, _ := strconv.Atoi(d1)
		d2r, _ := strconv.Atoi(d2)
		if d1r < d2r {
			return true
		}
		break
	case syntax.EXPS_GTE:
		d1r, _ := strconv.Atoi(d1)
		d2r, _ := strconv.Atoi(d2)
		if d1r >= d2r {
			return true
		}
		break
	case syntax.EXPS_LTE:
		d1r, _ := strconv.Atoi(d1)
		d2r, _ := strconv.Atoi(d2)
		if d1r <= d2r {
			return true
		}
		break
	case syntax.EXPS_NE:
		if d1 != d2 {
			return true
		}
		break
	case syntax.EXPS_LIKE:
		if -1 != strings.Index(d1, d2) {
			return true
		}
		break
	}
	return false
}

func getEleFromLine(si string, data []string) (string, error) {
	si = strings.Replace(si, syntax.FIELD_PRXFIX, "", -1)
	i, err := strconv.Atoi(si)
	if err != nil {
		return "", errors.New("Parse Where Syntax Fail [" + si + "]")
	}
	if i > len(data) {
		return "NULL", nil
	} else {
		i = i - 1
		return data[i], nil
	}
}

func regFieldFiliter(f string) []string {
	reg := regexp.MustCompile(`\@\d+`)
	return reg.FindAllString(f, -1)
}

// 校验字段是否有效
func regEleVerify(e string) bool {
	if len(e) <= 0 {
		return false
	}
	if isOk, _ := regexp.MatchString(`^\@\d+$`, e); isOk {
		return true
	}
	return false
}
