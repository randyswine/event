package event

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Event - интерфейс события.
type Event interface {
	Name() string    // Name() возвращает тип события.
	Data() any       // Data() возвращает данные события.
	Date() time.Time // Date() возвращает время события.
}

// HandlerCallback - тип функции обратного вызова для обработки события.
type HandlerCallback func(Event) error

// EventListener - интерфейс объекта-подписчика на события.
// Применять в том случае, если метода обработки должен иметь некоторый котекст.
type EventListener interface {
	Handle(Event) error
}

type handler struct {
	eventName  string
	callbackId string
	callback   HandlerCallback
}

// Unsub() удаляет из таблицы подписок callback.
func (h *handler) Unsub() {
	if h == nil {
		return
	}
	delete(subscribeRegister[h.eventName], h.callbackId)
	if len(subscribeRegister[h.eventName]) == 0 {
		delete(subscribeRegister, h.eventName)
	}
	return
}

var (
	rwmx              sync.RWMutex                  // RMutex таблицы подписок.
	subscribeRegister map[string]map[string]handler // Таблица подписок. Ключом является имя события.
)

func init() {
	// Инициализация таблицы подписок.
	subscribeRegister = make(map[string]map[string]handler)
}

// On() позволяет подписаться на события определенного типа через функцию обратного вызова.
func On(eventName string, callback HandlerCallback) *handler {
	if callback == nil {
		return nil
	}
	defer rwmx.Unlock()
	rwmx.Lock()
	guid, err := uuid.NewUUID()
	if err != nil {
		panic(fmt.Errorf("Unsubscribe error: %v", err))
	}
	sguid := guid.String()
	newHandler := handler{
		eventName:  eventName,
		callbackId: sguid,
		callback:   callback,
	}
	if len(subscribeRegister[eventName]) == 0 {
		subscribeRegister[eventName] = make(map[string]handler)
	}
	subscribeRegister[eventName][sguid] = newHandler
	return &newHandler
}

// Subscribe() позволяет объекту поджписаться на события определенного типа.
func Subscribe(eventName string, listener EventListener) *handler {
	return On(eventName, listener.Handle)
}

// FireEvent() инициирует вызов в отдельной рутине всех функций обратного вызова для заданного типа события.
func FireEvent(e Event) {
	defer rwmx.RUnlock()
	rwmx.RLock()
	var wg sync.WaitGroup
	for _, handler := range subscribeRegister[e.Name()] {
		wg.Add(1)
		go func(e Event, callback HandlerCallback) {
			defer func() {
				wg.Done()
			}()
			callback(e)
		}(e, handler.callback)
	}
	wg.Wait()
}

// basicEvent - базовая имплементация Event.
type basicEvent struct {
	name string    // Тип события.
	data any       // Данные события.
	date time.Time // Время события.
}

// NewEvent() возвращает basicEvent.
func NewEvent(name string, data any) Event {
	return basicEvent{
		name: name,
		data: data,
		date: time.Now(),
	}
}

// Name() возвращает тип события.
func (e basicEvent) Name() string {
	return e.name
}

// Data() возвращает данные события.
func (e basicEvent) Data() any {
	return e.data
}

// Date() возвращает время события.
func (e basicEvent) Date() time.Time {
	return e.date
}

// String() возвращает строчное представление события.
func (e basicEvent) String() string {
	return fmt.Sprintf("%s(%s): %v", e.name, e.date, e.data)
}
