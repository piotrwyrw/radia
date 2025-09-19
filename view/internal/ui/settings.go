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
	"github.com/sirupsen/logrus"
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

func bindEntryWidget(entry *widget.Entry, field reflect.Value) {
	resetValue := func() {
		entry.SetText(fmt.Sprintf("%v", field.Interface()))
	}

	resetValue()

	if !field.CanSet() {
		entry.Disable()
		return
	}

	if field.Kind() == reflect.String {
		entry.OnChanged = func(str string) {
			field.SetString(str)
		}
		return
	}

	if util.IsFloat(field.Type()) {
		entry.OnChanged = func(str string) {
			value, err := strconv.ParseFloat(str, 64)
			if err != nil {
				resetValue()
				return
			}
			field.SetFloat(value)
		}
		return
	}

	if util.IsInteger(field.Type()) {
		entry.OnChanged = func(str string) {
			value, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				resetValue()
				return
			}
			field.SetInt(value)
			return
		}
	}

	logrus.Fatal("Could not detect appropriate entry widget value type for: %v\n", field.Kind())
}

func createAndBindControl(field reflect.Value) (fyne.CanvasObject, error) {
	ctl, err := controlForType(field.Type())
	if err != nil {
		return nil, err
	}

	f := field

	switch control := ctl.(type) {
	case *widget.Entry:
		bindEntryWidget(control, f)
		break
	default:
		logrus.Fatal("Unsupported control for type %s\n", f.Type)
		break
	}

	return ctl, nil
}

// Create a settings panel via reflection
func createSettingsPanel(s interface{}) (fyne.CanvasObject, error) {
	sValue := reflect.ValueOf(s)

	if sValue.Kind() == reflect.Ptr {
		sValue = sValue.Elem()
	}
	if sValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct got %T", s)
	}

	sType := sValue.Type()

	form := container.New(layout.NewFormLayout())

	for i := 0; i < sType.NumField(); i++ {
		fieldValue := sValue.Field(i)
		fieldType := sType.Field(i)

		displayName := fieldType.Tag.Get("ui")

		// Ignore marked fields
		if displayName == "-" {
			continue
		}

		if displayName == "" {
			displayName = fieldType.Name
		}

		if fieldType.Type.Kind() == reflect.Struct {
			obj, err := createSettingsPanel(fieldValue.Addr().Interface())
			if err != nil {
				return nil, err
			}

			// Add accordion category, don't generate the field itself
			form.Add(layout.NewSpacer())
			form.Add(util.Do(func() fyne.CanvasObject {
				accordion := widget.NewAccordion(widget.NewAccordionItem(displayName, obj))
				accordion.OpenAll()
				return accordion
			}))
			continue
		}

		control, err := createAndBindControl(fieldValue)
		if err != nil {
			return nil, err
		}
		form.Add(widget.NewLabel(displayName))
		form.Add(control)

	}
	return form, nil
}
