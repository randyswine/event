package event

import (
	"fmt"
	"time"
)

// Event - интерфейс события.
type Event interface {
	Name() string    // Name() возвращает тип события.
	Data() any       // Data() возвращает данные события.
	Date() time.Time // Date() возвращает время события.
}

// EventListener - интерфейс объекта-подписчика на события.
// Применять в том случае, если метода обработки должен иметь некоторый котекст.
type EventListener interface {
	Handle(Event) error
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

// HandlerCallback - тип функции обратного вызова для обработки события.
type HandlerCallback func(Event) error
