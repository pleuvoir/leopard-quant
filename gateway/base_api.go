package gateway

type BaseApi struct {
	UnmarshalerOptions *UnmarshalerOptions
	ApiOptions         *ApiOptions
}

func NewBaseApi(opts ...apiOptions) *BaseApi {
	m := &BaseApi{}
	m.UnmarshalerOptions = &UnmarshalerOptions{}
	m.ApiOptions = NewApiOptions(opts...)
	return m
}

func NewBaseApiWithOptions(options *ApiOptions) *BaseApi {
	m := &BaseApi{}
	m.UnmarshalerOptions = &UnmarshalerOptions{}
	m.ApiOptions = options
	return m
}

func (m *BaseApi) WithUnmarshalerOption(opts ...UnmarshalerOption) *BaseApi {
	for _, opt := range opts {
		opt(m.UnmarshalerOptions)
	}
	return m
}

type Sub interface {
	Subscribe(symbol string, c ApiCallback) error
	CancelSubscribe(symbol string) error
	Name() string
}
