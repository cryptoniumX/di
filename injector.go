package di

import (
	"fmt"
	"reflect"
)

func (container *Container) Inject(servicePtr interface{}) error {
	ptrValue := reflect.ValueOf(servicePtr)

	// Ensure that servicePtr is a pointer to a struct
	if ptrValue.Kind() != reflect.Ptr || ptrValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Invalid pointer type")
	}

	structValue := ptrValue.Elem()
	// Iterate through the fields of the struct
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Type().Field(i)
		fieldValue := structValue.Field(i)
		depdencyName, ok := field.Tag.Lookup("di")
		if ok {
			dependency, err := invokeByName(container, field.Type.String(), depdencyName)
			if err != nil {
				return err
			}

			if !fieldValue.CanSet() {
				return fmt.Errorf("Field is not settable %s", field.Name)
			}

			if dependency == nil {
				return fmt.Errorf(
					"Dependency not found. Field: %s, Dependency: %s",
					field.Name,
					depdencyName,
				)
			}

			fieldValue.Set(reflect.ValueOf(dependency))
		}
	}

	return nil

}
