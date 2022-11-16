package model

type Tick struct {
}

type Order struct {
}

type Trade struct {
	Id string
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
