# Monkey Patching

Monkey Pathching is changing program at run time, typically by replacing a function or variable.

## Example

Let's we want to create unit test if `os.Writefile` failed on this function:

```go
type Config struct {
	Host string
	Port string
}

func SaveConfig(filename string, cfg *Config) error {
    data, err := json.Marshall(cfg)
    if err != nil {
        return err
    }

    err = os.Writefile(filename, data, 0666)
    if err != nil {
        return err
    }

    return nil
}
```

We can do that by using monkey patching to decouple `os.Writefile` and create mock.

```go
func SaveConfig(filename string, cfg *Config) error {
    data, err := json.Marshall(cfg)
    if err != nil {
        return err
    }

    err = os.Writefile(filename, data, 0666)
    if err != nil {
        return err
    }

    return nil
}

type fileWriter func(filename string, data []byte, perm fs.FileMode) error

var writeFile fileWriter = os.WriteFile
```

in Unit test, we can mock `writeFile` like this:

```go
func TestSaveConfig(t *testing.T) {
    // restore mock
    defer func() {
		writeFile = os.WriteFile
	}()

	writeFile = func(filename string, data []byte, perm fs.FileMode) error {
		return errors.New("failed to write data to file")
	}

	config := &Config{
		Host: "127.0.0.1",
		Port: "8080",
	}
	err := SaveConfig("config.json", config)
	assert.Error(t, err, nil)
}
```

## Advantage

### Easy to implement

Monkey patching is easy to implement. In the example above, we just change `os.Writefile(..)` with `var writeFile fileWriter = os.WriteFile`, and it gives us the ability to replace `writeFile()` with mock, which will allow us to test happy paths and error scenarios.

### Allows testing of global and singleton

## Disadvantage

### Data race

We use monkey pathing to replace a global variable with a copy that performs the way we need it to perform a particular test. This causes a data race issue because global variables are shared and used all over the place.

### Verbose test

When we mock code using monkey patching, we need to restore the original function or data. This is why the code became longer.

## What are the ideal use cases for monkey patching?

- With code that relies on a singleton
- With code that currently has no test, no dependency injection, and where you want to add test with a minimum of changes
- To decouple two packages without having to change the dependent package

## Reference:

https://www.packtpub.com/product/hands-on-dependency-injection-in-go/9781789132762
