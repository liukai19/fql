package syntax

import (
	"errors"
	"fmt"
	"fql/util"
	"regexp"
	"strconv"
	"strings"
)

const (
	RELATION_AND = "AND"
	RELATION_OR  = "OR"

	ORDER_METHOD_DESC = "DESC"
	ORDER_METHOD_ASC  = "ASC"

	EXPS_EQ    = "="
	EXPS_NE    = "!="
	EXPS_GT    = ">"
	EXPS_LT    = "<"
	EXPS_GTE   = ">="
	EXPS_LTE   = "<="
	EXPS_IN    = "IN"
	EXPS_NOTIN = "NOTIN"
	EXPS_LIKE  = "LIKE"

	FIELDS_ALL = "*"

	FIELD_PRXFIX = "@"
)

var ExpsMap = []string{
	EXPS_NE,
	EXPS_GTE,
	EXPS_LTE,
	EXPS_NOTIN,
	EXPS_IN,
	EXPS_EQ,
	EXPS_GT,
	EXPS_LT,
	EXPS_LIKE,
}

var ExpsItem = map[string]int{
	EXPS_NE:   1,
	EXPS_GTE:  1,
	EXPS_LTE:  1,
	EXPS_EQ:   2,
	EXPS_GT:   2,
	EXPS_LT:   2,
	EXPS_LIKE: 2,
}

var ExpsSet = map[string]int{
	EXPS_NOTIN: 1,
	EXPS_IN:    1,
}

var OrderMethod = []string{
	ORDER_METHOD_DESC,
	ORDER_METHOD_ASC,
}

type Syntax struct {
	Fields *Fields
	Source *Source
	Where  *Where
	Group  *GroupBy
	Order  *OrderBy
	Limit  *Limit
}

type Fields struct {
	Value []*Field
}
type Field struct {
	Name   string
	Remark string
}

type Where struct {
	Conds     []*Conds
	Relation  string
	ValueHash map[string]string
}
type Conds struct {
	Field     *Field
	Exps      string
	ItemValue string
	SetValue  []string
}

type GroupBy struct {
	Value []*Field
}

type OrderBy struct {
	Field  *Field
	Method string
}

type Limit struct {
	Limit  int
	Offset int
}

type Source struct {
	FilePath string
}

func ParseFields(str string) (*Fields, error) {
	ele := strings.Split(str, ",")
	if len(ele) <= 0 {
		return nil, errors.New("Parse Fields Syntax Fail.Near:" + str)
	}
	eles := []*Field{}
	for _, e := range ele {
		f := &Field{}
		e = strings.ToUpper(e)
		if -1 != strings.Index(e, "AS") {
			e2 := strings.Split(e, "AS")
			f.Name = e2[0]
			f.Remark = e2[1]
		} else {
			f.Name = e
			f.Remark = e
		}
		if ok := eleVerify(f.Name); ok != true {
			return nil, errors.New("Parse Fields Syntax Fail.Near:" + str)
		}
		eles = append(eles, f)
	}
	return &Fields{Value: eles}, nil
}

func ParseSource(str string) (*Source, error) {
	if len(str) <= 0 {
		return nil, errors.New("Parse Source Syntax Fail.Near:" + str)
	}
	return &Source{FilePath: str}, nil
}

func ParseLimit(str string) (*Limit, error) {
	if len(str) <= 0 {
		return nil, errors.New("Parse Limit Syntax Fail.Near:" + str)
	}
	l := &Limit{}
	if -1 != strings.Index(str, ",") {
		r := strings.Split(str, ",")
		limit, err := strconv.Atoi(r[0])
		if err != nil {
			return nil, errors.New("Parse Limit Syntax Fail.Near:" + str)
		}
		if limit < 0 {
			return nil, errors.New("Parse Limit Syntax Fail.Near:" + str)
		}
		offset, err := strconv.Atoi(r[1])
		if err != nil {
			return nil, errors.New("Parse Limit Syntax Fail.Near:" + str)
		}
		if offset < 1 {
			return nil, errors.New("Parse Limit Syntax Fail.Near:" + str)
		}
		l.Limit = limit
		l.Offset = offset
	} else {
		offset, err := strconv.Atoi(str)
		if err != nil {
			return nil, errors.New("Parse Limit Syntax Fail.Near:" + str)
		}
		if offset < 1 {
			return nil, errors.New("Parse Limit Syntax Fail.Near:" + str)
		}
		l.Limit = 0
		l.Offset = offset
	}
	return l, nil
}

func ParseGroup(str string) (*GroupBy, error) {
	ns := strings.ToUpper(str)
	ns = strings.Replace(ns, "BY", "", -1)
	ele := strings.Split(ns, ",")
	if len(ele) <= 0 {
		return nil, errors.New("Parse GroupBy Syntax Fail.Near:" + str)
	}
	eles := []*Field{}
	for _, e := range ele {
		if ok := eleVerify(e); ok != true {
			return nil, errors.New("Parse GroupBy Syntax Fail.Near:" + str)
		}
		f := &Field{
			Name:   e,
			Remark: e,
		}
		eles = append(eles, f)
	}
	return &GroupBy{Value: eles}, nil
}

func ParseOrder(str string) (*OrderBy, error) {
	ns := strings.ToUpper(str)
	ns = strings.Replace(ns, "BY", "", -1)
	o := &OrderBy{}
	for _, m := range OrderMethod {
		if -1 != strings.Index(ns, m) {
			o.Method = m
			break
		}
	}
	if len(o.Method) <= 0 {
		return nil, errors.New("Parse OrderBy Syntax Fail.Near:" + str)
	}
	e := strings.Replace(ns, o.Method, "", -1)
	if ok := eleVerify(e); ok != true {
		return nil, errors.New("Parse OrderBy Syntax Fail.Near:" + str)
	}
	o.Field = &Field{
		Name:   e,
		Remark: e,
	}
	return o, nil
}

func replaceStringValue(str string) (string, map[string]string) {
	valueHashMap := map[string]string{}
	reg := regexp.MustCompile(`\'(.*?)\'`)
	ret := reg.FindAllString(str, -1)
	if len(ret) > 0 {
		for _, v := range ret {
			hash := util.Md5V(v)
			valueHashMap[hash] = strings.Replace(v, "'", "", -1)
			str = strings.Replace(str, v, hash, -1)
		}
	}
	reg = regexp.MustCompile(`\"(.*?)\"`)
	ret = reg.FindAllString(str, -1)
	if len(ret) > 0 {
		for _, v := range ret {
			hash := util.Md5V(v)
			valueHashMap[hash] = strings.Replace(v, "\"", "", -1)
			str = strings.Replace(str, v, hash, -1)
		}
	}
	return str, valueHashMap
}

func ParseWhere(str string) (*Where, error) {
	w := &Where{
		Conds:     nil,
		Relation:  RELATION_AND,
		ValueHash: map[string]string{},
	}
	ns := strings.ToUpper(str)

	// 处理 string value
	ns, w.ValueHash = replaceStringValue(ns)

	if len(strings.Split(ns, RELATION_AND)) > 1 {
		w.Relation = RELATION_AND
	}
	if len(strings.Split(ns, RELATION_OR)) > 1 {
		w.Relation = RELATION_OR
	}
	conds := strings.Split(ns, w.Relation)
	for _, v := range conds {
		c := &Conds{}
		for _, e := range ExpsMap {
			if len(strings.Split(v, e)) > 1 {
				c.Exps = e
				break
			}
		}
		if c.Exps == "" {
			return nil, errors.New("Parse Where Syntax Fail.Near:" + str)
		}

		if _, ok := ExpsItem[c.Exps]; ok != false {
			condsVal := strings.Split(v, c.Exps)
			c.Field = &Field{
				Name:   condsVal[0],
				Remark: condsVal[0],
			}
			c.ItemValue = condsVal[1]
		}

		if _, ok := ExpsSet[c.Exps]; ok != false {
			condsVal := strings.Split(v, c.Exps)
			c.Field = &Field{
				Name:   condsVal[0],
				Remark: condsVal[0],
			}
			condsVal[1] = strings.Replace(condsVal[1], "(", "", -1)
			condsVal[1] = strings.Replace(condsVal[1], ")", "", -1)
			c.SetValue = strings.Split(condsVal[1], ",")
		}
		if len(c.Field.Remark) <= 0 {
			return nil, errors.New("Parse Where Syntax Fail.Near:" + str)
		}
		if ok := eleVerify(c.Field.Name); ok != true {
			return nil, errors.New("Parse Where Syntax Fail.Near:" + str)
		}
		if len(c.ItemValue) <= 0 && len(c.SetValue) <= 0 {
			return nil, errors.New("Parse Where Syntax Fail.Near:" + str)
		}
		w.Conds = append(w.Conds, c)
	}
	return w, nil
}

// 校验字段是否有效
func eleVerify(e string) bool {
	// 字段会有函数，目前先不校验
	return true
	if len(e) <= 0 {
		return false
	}
	if isOk, _ := regexp.MatchString(`^\@\d+$`, e); isOk {
		return true
	}
	return false
}

func (f *Fields) AddIndexField(num int) {
	for i, v := range f.Value {
		if v.Name == "*" {
			f.Value = append(f.Value[:i], f.Value[i+1:]...)
			for ni := 1; ni <= num; ni++ {
				f.Value = append(f.Value, &Field{
					Name:   fmt.Sprintf("%s%d", FIELD_PRXFIX, ni),
					Remark: fmt.Sprintf("%s%d", FIELD_PRXFIX, ni),
				})
			}
			break
		}
	}
}
