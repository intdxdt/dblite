package dblite

type QueryOption struct {
	on                      *On
	where                   *Where
	funcColumnsPlaceholders FuncColumnsPlaceholders
}

type QueryOpt func(*QueryOption)

func NewQueryOption(opts ...QueryOpt) *QueryOption {
	var callback = func(cols string, holders string) (string, string) {
		return cols, holders
	}
	var opt = &QueryOption{nil, nil, callback}
	for _, fn := range opts {
		fn(opt)
	}
	return opt
}

func WithOn(on *On) QueryOpt {
	return func(o *QueryOption) {
		o.on = on
	}
}

func WithWhere(where *Where) QueryOpt {
	return func(o *QueryOption) {
		o.where = where
	}
}

func WithFuncColumnsPlaceholders(fn FuncColumnsPlaceholders) QueryOpt {
	return func(o *QueryOption) {
		o.funcColumnsPlaceholders = fn
	}
}

func (opt *QueryOption) hasWhereClause() bool {
	return opt.where != nil
}

func (opt *QueryOption) whereString() string {
	return opt.where.clause
}

func (opt *QueryOption) whereArguments() []any {
	return opt.where.arguments
}

func (opt *QueryOption) hasOn() bool {
	return opt.on != nil && opt.on.hasOn()
}

func (opt *QueryOption) onString() string {
	return opt.on.clause
}

func (opt *QueryOption) hasOnUpsertColumns() bool {
	return opt.hasOn() && opt.on.hasUpsertColumns()
}

func (opt *QueryOption) hasOnArguments() bool {
	return opt.hasOn() && opt.on.hasArguments()
}

func (opt *QueryOption) OnClause(db *Database, getColumnValues FuncGetColumnValues) (string, []any) {
	return opt.on.OnClause(db, getColumnValues)
}
