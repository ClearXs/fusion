package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slog"
)

type Predicate = string

const (
	Or  Predicate = "$or"
	And Predicate = "$and"
	Nor Predicate = "$nor"
)

type Logical interface {
	// GetPredicate 获取逻辑谓词
	GetPredicate() Predicate
	// Append 追加用于查询的数据
	Append(ele bson.E) Logical
	// AppendArray 追加用于查询的数据 bson array
	AppendArray(ele bson.D) Logical
	// AppendLogical 追加逻辑谓词
	AppendLogical(other Logical) Logical
	// ToBson 获取加上逻辑谓词的bson数据
	ToBson() bson.E
	// ToJsonString to json
	ToJsonString() string
}

type LogicalImpl struct {
	predicate  Predicate
	additional []bson.E
}

// NewLogical 创建默认的逻辑词，默认是And
func NewLogical() Logical {
	return NewLogicalAnd()
}

// NewLogicalDefault 创建默认的逻辑词，默认是And
func NewLogicalDefault(ele bson.E) Logical {
	return NewLogicalAndDefault(ele)
}

// NewLogicalDefaultArray 创建默认的逻辑词，默认是And
func NewLogicalDefaultArray(ele bson.D) Logical {
	return NewLogicalAndDefaultArray(ele)
}

// NewLogicalDefaultLogical 创建默认的逻辑词，默认是And
func NewLogicalDefaultLogical(other Logical) Logical {
	return NewLogicalAndDefaultLogical(other)
}

// NewLogicalOr 创建逻辑or谓词
func NewLogicalOr() Logical {
	return &LogicalImpl{predicate: Or, additional: make([]bson.E, 0)}
}

// NewLogicalOrDefault 创建逻辑or谓词
func NewLogicalOrDefault(ele bson.E) Logical {
	logical := NewLogicalOr()
	logical.Append(ele)
	return logical
}

// NewLogicalOrDefaultArray 创建逻辑or谓词
func NewLogicalOrDefaultArray(ele bson.D) Logical {
	logical := NewLogicalOr()
	logical.AppendArray(ele)
	return logical
}

// NewLogicalOrDefaultLogical 创建逻辑or谓词
func NewLogicalOrDefaultLogical(other Logical) Logical {
	logical := NewLogicalOr()
	logical.AppendLogical(other)
	return logical
}

// NewLogicalAnd 创建逻辑and谓词
func NewLogicalAnd() Logical {
	return &LogicalImpl{predicate: And, additional: make([]bson.E, 0)}
}

// NewLogicalAndDefault 创建逻辑and谓词
func NewLogicalAndDefault(ele bson.E) Logical {
	logical := NewLogicalAnd()
	logical.Append(ele)
	return logical
}

// NewLogicalAndDefaultArray 创建逻辑and谓词
func NewLogicalAndDefaultArray(ele bson.D) Logical {
	logical := NewLogicalAnd()
	logical.AppendArray(ele)
	return logical
}

// NewLogicalAndDefaultLogical 创建逻辑and谓词
func NewLogicalAndDefaultLogical(other Logical) Logical {
	logical := NewLogicalAnd()
	logical.AppendLogical(other)
	return logical
}

// NewLogicalNor 创建逻辑nor谓词
func NewLogicalNor() Logical {
	return &LogicalImpl{predicate: Nor, additional: make([]bson.E, 0)}
}

// NewLogicalNorDefault 创建逻辑nor谓词
func NewLogicalNorDefault(ele bson.E) Logical {
	logical := NewLogicalNor()
	logical.Append(ele)
	return logical
}

// NewLogicalNorDefaultArray 创建逻辑nor谓词
func NewLogicalNorDefaultArray(ele bson.D) Logical {
	logical := NewLogicalNor()
	logical.AppendArray(ele)
	return logical
}

// NewLogicalNorDefaultLogical 创建逻辑nor谓词
func NewLogicalNorDefaultLogical(other Logical) Logical {
	logical := NewLogicalNor()
	logical.AppendLogical(other)
	return logical
}

func (logic *LogicalImpl) GetPredicate() Predicate {
	return logic.predicate
}

func (logic *LogicalImpl) Append(ele bson.E) Logical {
	logic.additional = append(logic.additional, ele)
	return logic
}

func (logic *LogicalImpl) AppendArray(ele bson.D) Logical {
	logic.additional = append(logic.additional, ele...)
	return logic
}

func (logic *LogicalImpl) AppendLogical(other Logical) Logical {
	logic.additional = append(logic.additional, other.ToBson())
	return logic
}

func (logic *LogicalImpl) ToBson() bson.E {
	if len(logic.additional) == 0 {
		return bson.E{}
	}
	return bson.E{Key: logic.predicate, Value: logic.additional}
}

func (logic *LogicalImpl) ToJsonString() string {
	b := logic.ToBson()
	jsonBytes, err := bson.MarshalExtJSON(b, true, true)
	if err != nil {
		slog.Error("Failed to marshal bson to json", "err", err, "bson", b)
		return ""
	}
	return string(jsonBytes)
}
