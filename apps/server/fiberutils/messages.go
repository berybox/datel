package fiberutils

import "github.com/gofiber/fiber/v2"

// Messages messages to be passed to the top of the page
type Messages []struct {
	Text string
	Type string
}

// CreateMessages creates new message list
func CreateMessages() *Messages {
	return &Messages{}
}

// PullFromCtx retrieves messages from the Fiber context and adds them to the current messages
func (m *Messages) PullFromCtx(c *fiber.Ctx) *Messages {
	if msgs, ok := c.Locals("Msgs").(*Messages); ok {
		*m = append(*m, *msgs...)
		c.Locals("Msgs", nil)
		return m
	}
	return m
}

// PutToCtx puts messages to the Fiber context
func (m *Messages) PutToCtx(c *fiber.Ctx) {
	c.Locals("Msgs", m)
}

func (m *Messages) addMessage(text, typ string) *Messages {
	*m = append(*m, Messages{{Text: text, Type: typ}}[0])
	return m
}

// AddSuccess adds message of type "success"
func (m *Messages) AddSuccess(text string) *Messages {
	return m.addMessage(text, "success")
}

// AddInfo adds message of type "info"
func (m *Messages) AddInfo(text string) *Messages {
	return m.addMessage(text, "info")
}

// AddDanger adds message of type "danger"
func (m *Messages) AddDanger(text string) *Messages {
	return m.addMessage(text, "danger")
}

// AddWarning adds message of type "warning"
func (m *Messages) AddWarning(text string) *Messages {
	return m.addMessage(text, "warning")
}
