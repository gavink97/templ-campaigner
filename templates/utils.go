package templates

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/a-h/templ"
	"github.com/gavink97/templ-campaigner/internal/contacts"
)

type ContactDetails []contacts.Contact

func TemplateConstructor(contacts *[]contacts.Contact, name string) templ.Component {
	list := ContactDetails(*contacts)
	component, err := list.CallMethodByName(name)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	return component
}

func (c *ContactDetails) CallMethodByName(name string) (templ.Component, error) {
	method := reflect.ValueOf(c).MethodByName(name)
	if !method.IsValid() {
		slog.Error(fmt.Sprintf("%s is not a valid method", name))
		return nil, nil
	}

	results := method.Call(nil)

	if component, ok := results[0].Interface().(templ.Component); ok {
		return component, nil
	} else {
		return nil, nil
	}
}

func (c *ContactDetails) Preview() contacts.Contact {
	for _, user := range *c {
		return user
	}
	return contacts.Contact{}
}
