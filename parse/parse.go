package parse

import (
	"fmt"
	"fql/operator"
	"fql/syntax"
	"strings"
)

type Parser struct {
	sql       string // 待解析的查询
	delimiter string // 文本分割符
	i         int    // 当前所在查询字符串中的位置
	token     []string
	Query     *Query
	Syntax    *syntax.Syntax
	Data      []map[string]string
}

const (
	DEFAULT_DELIMITER = " "
)

type Query struct {
	Fields string
	Source string
	Where  string
	Group  string
	Order  string
	Limit  string
}

var q = &Query{}

var keepTokens = map[string]*string{
	"FROM":   &q.Source,
	"WHERE":  &q.Where,
	"GROUP":  &q.Group,
	"ORDER":  &q.Order,
	"SELECT": &q.Fields,
	"LIMIT":  &q.Limit,
}

func GetInstance(sql string, delimiter string) *Parser {
	return &Parser{
		sql:       sql,
		delimiter: delimiter,
		i:         0,
		Query:     q,
		Syntax:    &syntax.Syntax{},
		Data:      []map[string]string{},
	}
}

// step1 记号切分
// step2 Query分区
// step3 语法解析
// step4 算子计算
// step5 输出结果
func (p *Parser) Parse() error {
	// 获取记号
	p.getTokens()

	// Query 分区
	p.queryPartition()

	// 语法解析
	p.syntax()

	// 算子计算
	p.operator()

	// 结果输出
	p.output()

	return nil
}

func (p *Parser) output() {
	title := []string{}
	for _,f := range p.Syntax.Fields.Value{
		title = append(title, f.Remark)
	}
	fmt.Println(strings.Join(title, " | "))

	if len(p.Data) <= 0 {
		fmt.Println("not data.")
	} else {
		for _,v:=range p.Data{
			item := []string{}
			for _,f:=range title{
				item = append(item, v[f])
			}
			fmt.Println(strings.Join(item, " | "))
		}
	}
}

func (p *Parser) operator() {
	op := &operator.Oper{
		Syntax:    p.Syntax,
		Delimiter: p.delimiter,
		Data:      nil,
	}
	op.Process()
	p.Data = op.Data
}

func (p *Parser) syntax() {
	if len(p.Query.Fields) > 0 {
		if r, err := syntax.ParseFields(p.Query.Fields); err != nil {
			panic(err.Error())
		} else {
			p.Syntax.Fields = r
		}
	}
	if len(p.Query.Source) > 0 {
		if r, err := syntax.ParseSource(p.Query.Source); err != nil {
			panic(err.Error())
		} else {
			p.Syntax.Source = r
		}
	}
	if len(p.Query.Where) > 0 {
		if r, err := syntax.ParseWhere(p.Query.Where); err != nil {
			panic(err.Error())
		} else {
			p.Syntax.Where = r
		}
	}
	if len(p.Query.Group) > 0 {
		if r, err := syntax.ParseGroup(p.Query.Group); err != nil {
			panic(err.Error())
		} else {
			p.Syntax.Group = r
		}
	}
	if len(p.Query.Order) > 0 {
		if r, err := syntax.ParseOrder(p.Query.Order); err != nil {
			panic(err.Error())
		} else {
			p.Syntax.Order = r
		}
	}
	if len(p.Query.Limit) > 0 {
		if r, err := syntax.ParseLimit(p.Query.Limit); err != nil {
			panic(err.Error())
		} else {
			p.Syntax.Limit = r
		}
	}
}

func (p *Parser) queryPartition() {
	tokenIndex := 0
	for p.i < len(p.token) {
		t := p.pop()
		if _, ok := keepTokens[strings.ToUpper(t)]; ok != false {
			tokenIndex = p.i - 1
		} else {
			*keepTokens[strings.ToUpper(p.token[tokenIndex])] += t
		}
	}
	//fmt.Println(keepTokens)
}

// 获取记号
func (p *Parser) getTokens() {
	p.token = strings.Split(p.sql, " ")
}

func (p *Parser) pop() string {
	token := p.token[p.i]
	p.i += 1
	return token
}
