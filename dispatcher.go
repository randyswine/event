package event

import "sync"

var dispatcherInstance *dispatcher // Синглтон диспетчера.
var rwmx sync.RWMutex              // RMutex синглтона диспетчера.

// dispatcher реализует таблицу подписок и рассылку событий подписчикам.
type dispatcher struct {
	rwmx              sync.RWMutex                 // Доступ к таблице подписок блокируется.
	subscribeRegister map[string][]HandlerCallback // Таблица подписок. Ключом является имя события.
}

// New() инициализирует таблицу подписок и возвращает указатель на диспетчер событий.
func Dispatcher() *dispatcher {
	rwmx.RLock()
	instance := dispatcherInstance
	rwmx.RUnlock()
	if instance == nil {
		rwmx.Lock()
		dispatcherInstance = &dispatcher{
			subscribeRegister: make(map[string][]HandlerCallback, 0),
		}
		instance = dispatcherInstance
		rwmx.Unlock()
	}
	return instance
}

// On() позволяет подписаться на события определенного типа через функцию обратного вызова.
func (d *dispatcher) On(eventName string, callback HandlerCallback) {
	defer d.rwmx.Unlock()
	d.rwmx.Lock()
	d.subscribeRegister[eventName] = append(d.subscribeRegister[eventName], callback)
}

// Subscribe() позволяет объекту поджписаться на события определенного типа.
func (d *dispatcher) Subscribe(eventName string, listener EventListener) {
	d.On(eventName, listener.Handle)
}

// FireEvent() инициирует вызов в отдельной рутине всех функций обратного вызова для заданного типа события.
func (d *dispatcher) FireEvent(e Event) {
	defer d.rwmx.RUnlock()
	d.rwmx.RLock()
	var wg sync.WaitGroup
	for _, callback := range d.subscribeRegister[e.Name()] {
		wg.Add(1)
		go func(e Event, callback HandlerCallback) {
			callback(e)
			wg.Done()
		}(e, callback)
	}
	wg.Wait()
}
