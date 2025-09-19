package ui

import (
	"fmt"
	"reflect"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/piotrwyrw/otherproj/internal/util"
)

func controlForType(t reflect.Type) (fyne.CanvasObject, error) {
	if t.Kind() == reflect.String {
		return util.Do(func() fyne.CanvasObject {
			entry := widget.NewEntry()
			entry.MultiLine = false
			return entry
		}), nil
	}

	if util.IsInteger(t) {
		return util.Do(func() fyne.CanvasObject {
			entry := widget.NewEntry()
			entry.MultiLine = false
			entry.Validator = func(s string) error {
				_, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid integer %s", s)
				}
				return nil
			}
			return entry
		}), nil
	}

	if util.IsFloat(t) {
		return util.Do(func() fyne.CanvasObject {
			entry := widget.NewEntry()
			entry.MultiLine = false
			entry.Validator = func(s string) error {
				_, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return fmt.Errorf("invalid float %s", s)
				}
				return nil
			}
			return entry
		}), nil
	}

	return nil, fmt.Errorf("could not match control for type %v", t)
}

// Create a settings panel via reflection
func createSettingsPanel(s interface{}) (fyne.CanvasObject, error) {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct got %T", s)
	}

	form := container.New(layout.NewFormLayout())

	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		typeField := t.Field(i)
		name := typeField.Tag.Get("ui")

		// Ignore marked fields
		if name == "-" {
			continue
		}

		if name == "" {
			name = typeField.Name
		}

		if typeField.Type.Kind() == reflect.Struct {
			obj, err := createSettingsPanel(field.Interface())
			if err != nil {
				return nil, err
			}

			// Add accordion category, don't generate the field itself
			//form.Add(widget.NewLabel(name))
			form.Add(layout.NewSpacer())
			form.Add(util.Do(func() fyne.CanvasObject {
				accordion := widget.NewAccordion(widget.NewAccordionItem(name, obj))
				accordion.OpenAll()
				return accordion
			}))
			continue
		}

		form.Add(widget.NewLabel(name))
		ctl, err := controlForType(field.Type())
		if err != nil {
			return nil, err
		}
		form.Add(ctl)
	}
	return form, nil
}
