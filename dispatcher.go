package event

import "sync"

// dispatcher реализует таблицу подписок и рассылку событий подписчикам.
type dispatcher struct {
	mx                sync.Mutex                   // Доступ к таблице подписок блокируется.
	subscribeRegister map[string][]HandlerCallback // Таблица подписок. Ключом является имя события.
}

// New() инициализирует таблицу подписок и возвращает указатель на диспетчер событий.
func New() *dispatcher {
	return &dispatcher{
		subscribeRegister: make(map[string][]HandlerCallback, 0),
	}
}

// On() позволяет подписаться на события определенного типа через функцию обратного вызова.
func (d *dispatcher) On(eventName string, callback HandlerCallback) {
	defer d.mx.Unlock()
	d.mx.Lock()
	d.subscribeRegister[eventName] = append(d.subscribeRegister[eventName], callback)
}

// Subscribe() позволяет объекту поджписаться на события определенного типа.
func (d *dispatcher) Subscribe(eventName string, listener EventListener) {
	d.On(eventName, listener.Handle)
}

// FireEvent() инициирует вызов в отдельной рутине всех функций обратного вызова для заданного типа события.
func (d *dispatcher) FireEvent(e Event) {
	defer d.mx.Unlock()
	d.mx.Lock()
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
