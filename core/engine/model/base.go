package model

type Tick struct {
	Symbol string
}

func (t Tick) UpdateTick(tick Tick) {

}

type Order struct {
	Id string
}

func (o Order) IsActive() bool {
	return false
}

type Trade struct {
	Id      string
	OrderId string
}

type Bar struct {
}

// RunMode 运行模式
type RunMode struct {
	Code string
	Name string
}

var (
	Live        = RunMode{"live", "实盘"}
	DryRun      = RunMode{"dryRun", "试运行"}
	BackTesting = RunMode{"backTesting", "回测"}
)
