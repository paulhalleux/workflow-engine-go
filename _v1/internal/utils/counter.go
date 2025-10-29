package utils

type Counter struct {
	value int
}

func NewCounter() *Counter {
	return &Counter{value: 0}
}

func (c *Counter) Increment() {
	c.value++
}

func (c *Counter) Decrement() {
	c.value--
}

func (c *Counter) GetValue() int {
	return c.value
}
