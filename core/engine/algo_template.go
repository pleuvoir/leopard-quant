package engine

import . "leopard-quant/core/engine/model"

type AlgoTemplate interface {
	Init()
	onBar(bar Bar)
	updateTick(tick Tick)
	updateOrder(order Order)
	updateTrade(trade Trade)
	updateTimer()
	onTimer()
	onStart()
	onStop()
	subscribe(symbol string)
	unsubscribe(symbol string)
	sendOrder(symbol string)
	buy()
	sell()
	cancelOrder()
	cancelAllOrder()
}
