# i18n

This is a simple i18n library which uses JSON [resources](http://www.github.com/chrstphlbr/resource) as input. 

An example of such a JSON resource is:
```json
{
	"key": {
		"language1": "value",
		"language2": "value",
		"language3": "value"
	}
}
```

The following lines of code show how to use the default manager with the default directory (./files) for looking up a value for a key and a specific language. 
```go
manager := i18n.Manager()

value, err := manager.Get("key", "language")
if err != nil {
	// do error handling
}

// use value

```

Other repositories (other directories, multiple repositories, ...) with the default manager implementation are available with the consrtuctor function NewDefaultI18nManager. If you want to access other resources (e.g. network resources, ...) feel free to provide the NewDefaultI18nManager function with own implementations of resource.Repository and resource.Adapter.