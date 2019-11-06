package events

import "github.com/gingerxman/eel/event"

var ORDER_FINISHED *event.Event
var ORDER_CANCLLED *event.Event
var ORDER_SETTLED *event.Event
var ORDER_WAITING_CONFIRM *event.Event
var ORDER_REFUNDED *event.Event

func init()  {
	ORDER_FINISHED = event.NewEvent("order:finished", "order")
	ORDER_CANCLLED = event.NewEvent("order:canceled", "order")
	ORDER_SETTLED = event.NewEvent("order:settled", "order")
	ORDER_WAITING_CONFIRM = event.NewEvent("order:waiting_confirm", "order")
	ORDER_REFUNDED = event.NewEvent("order:refunded", "order")
}