package di

import (
	"fmt"
	"reflect"
	"sync"
)

type Container struct {
	services map[string]interface{}
	mutex    sync.RWMutex
}

func NewContainer() *Container {
	return &Container{
		services: make(map[string]interface{}),
	}
}

func (c *Container) Register(name string, service interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.services[name] = service
}

func (c *Container) Get(name string) (interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	service, exists := c.services[name]
	if !exists {
		return nil, fmt.Errorf("service '%s' not found", name)
	}

	return service, nil
}

func (c *Container) GetTyped(name string, target interface{}) error {
	service, err := c.Get(name)
	if err != nil {
		return err
	}

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Interface {
		return fmt.Errorf("target must be a pointer to an interface")
	}

	targetValue.Elem().Set(reflect.ValueOf(service))
	return nil
}

func (c *Container) Inject(target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	targetType := targetValue.Elem().Type()
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		injectTag := field.Tag.Get("inject")

		if injectTag != "" {
			service, err := c.Get(injectTag)
			if err != nil {
				return fmt.Errorf("failed to inject %s: %w", injectTag, err)
			}

			fieldValue := targetValue.Elem().Field(i)
			if fieldValue.CanSet() {
				fieldValue.Set(reflect.ValueOf(service))
			}
		}
	}

	return nil
}