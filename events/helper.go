package events

// Instance is global instance
var Instance IEventbus

// init func
func init() {
	Instance = New()
}

// Subscribe to a topic
func Subscribe(topic string, handler func(v ...any)) {
	Instance.Subscribe(topic, handler)
}

// SubscribersCount returns subscribers count
func SubscribersCount(topic string) int {
	return Instance.SubscribersCount(topic)
}

// Unsubscribe from a topic
func Unsubscribe(topic string) {
	Instance.Unsubscribe(topic)
}

// Dispatch a topic
func Dispatch(topic string, v ...any) {
	Instance.Dispatch(topic, v...)
}

// DispatchAsync dispatches a topic asynchronously
func DispatchAsync(topic string, v ...any) {
	Instance.DispatchAsync(topic, v...)
}

// DispatchFunc dispatches a topic with a function
func DispatchFunc(topic string, fn func(v ...any), v ...any) {
	Instance.DispatchFunc(topic, fn, v...)
}

// DispatchFuncAsync dispatches a topic asynchronously with a function
func DispatchFuncAsync(topic string, fn func(v ...any), v ...any) {
	Instance.DispatchFuncAsync(topic, fn, v...)
}
