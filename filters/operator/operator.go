package operator

import (
	"github.com/FooSoft/goldsmith"
)

type Operator interface {
	goldsmith.Filter
}

func And(filters ...goldsmith.Filter) Operator {
	return &operatorAnd{filters}
}

type operatorAnd struct {
	filters []goldsmith.Filter
}

func (*operatorAnd) Name() string {
	return "operator"
}

func (filter *operatorAnd) Accept(file *goldsmith.File) (bool, error) {
	for _, f := range filter.filters {
		accept, err := f.Accept(file)
		if err != nil {
			return false, err
		}
		if !accept {
			return false, nil
		}
	}

	return true, nil
}

func Not(filter goldsmith.Filter) Operator {
	return &operatorNot{filter}
}

type operatorNot struct {
	filter goldsmith.Filter
}

func (*operatorNot) Name() string {
	return "operator"
}

func (filter *operatorNot) Accept(file *goldsmith.File) (bool, error) {
	accept, err := filter.filter.Accept(file)
	return !accept, err
}

func Or(filters ...goldsmith.Filter) Operator {
	return &operatorOr{filters}
}

type operatorOr struct {
	filters []goldsmith.Filter
}

func (*operatorOr) Name() string {
	return "operator"
}

func (filter *operatorOr) Accept(file *goldsmith.File) (bool, error) {
	for _, f := range filter.filters {
		accept, err := f.Accept(file)
		if err != nil {
			return false, err
		}
		if accept {
			return true, nil
		}
	}

	return false, nil
}
