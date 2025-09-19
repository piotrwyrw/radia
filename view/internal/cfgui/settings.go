package cfgui

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

type SettingsPanel struct {
	objects  map[reflect.Value]fyne.CanvasObject // Maps fields to their controls
	children []*SettingsPanel
}

func (s *SettingsPanel) SetDefaultValues() {
	for k, v := range s.objects {
		v, ok := v.(*widget.Entry)
		if !ok {
			continue
		}
		v.SetText(fmt.Sprintf("%v", k.Interface()))
	}

	for _, child := range s.children {
		child.SetDefaultValues()
	}
}

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
		return
	}

	logrus.Fatalf("Could not detect appropriate entry widget value type for: %v\n", field.Kind())
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

// CreateSettingsPanel -Create a settings panel via reflection
func CreateSettingsPanel(s interface{}) (fyne.CanvasObject, *SettingsPanel, error) {
	sValue := reflect.ValueOf(s)

	if sValue.Kind() == reflect.Ptr {
		sValue = sValue.Elem()
	}
	if sValue.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("expected struct got %T", s)
	}

	sType := sValue.Type()

	form := container.New(layout.NewFormLayout())

	panel := SettingsPanel{objects: make(map[reflect.Value]fyne.CanvasObject), children: make([]*SettingsPanel, 0)}

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
			obj, child, err := CreateSettingsPanel(fieldValue.Addr().Interface())
			if err != nil {
				return nil, nil, err
			}

			// Add accordion category, don't generate the field itself
			form.Add(layout.NewSpacer())
			form.Add(util.Do(func() fyne.CanvasObject {
				accordion := widget.NewAccordion(widget.NewAccordionItem(displayName, obj))
				accordion.OpenAll()
				return accordion
			}))

			panel.children = append(panel.children, child)

			continue
		}

		control, err := createAndBindControl(fieldValue)
		if err != nil {
			return nil, nil, err
		}
		form.Add(widget.NewLabel(displayName))
		form.Add(control)
		panel.objects[fieldValue] = control

	}
	return form, &panel, nil
}
