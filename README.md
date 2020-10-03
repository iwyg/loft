# LO/FT
## Cascading logging for go

Loft exposes a standardized logging interface and a default implementation of a logging handler logging to `log.Logger` 


```bash
> go get gitlab.tmt.de/golang/loft
```

### Log cascading
Logging levels are defined in a cascading order

```sh
	Debug   - 1
	Info      0
	Notice    1
	Warn      2
	Error     3
	Fatal     4
	Emergency 5
)
```

The logger can be configured to forward specific log levels to different handlers. To do so, you need to set the logging handlers in ascending order (`Debug` -> `Emergency`).

```go

infoWriter := &Output{}

logger := New("testing", []Handler{
    NewStdHandler(Debug, os.Stdout, log.LstdFlags),
    NewStdHandler(Info, infoWriter, log.LstdFlags),
    NewStdHandler(Error, os.Stderr, og.LstdFlags),
})

type Output struct {
	buf []byte
	mu  sync.Mutex
}

func (w *Output) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return string(w.buf)
}

func (w *Output) Write(n []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buf = n
	return len(n), nil
}
```
