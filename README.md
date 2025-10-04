
# 🪵 `loger` — Custom Logging & Error Wrapping for Go

This package provides a **lightweight logging and error tracking system** with:

* Unique session IDs for each log instance
* Optional info logs (configurable via env var)
* Automatic file and line number tagging for errors
* Structured output for easy tracing

---

## 📦 Installation

Simply include it in your Go project (as a local package or module import):

```go
import "github.com/saravanan611/loger"
```

---

## ⚙️ Initialization

Before logging, initialize a log session using:

```go
l := loger.Init()
```

Each logger instance (`l`) will have its own **unique UID**, shown in every log line.
You can enable or disable info logs using an environment variable:

```bash
export InfoFlog=Y   # Enable info logs
# or
export InfoFlog=N   # Disable info logs
```

---

## 🧠 How It Works

* **Info logs** are only printed when `InfoFlog=Y`.
* **Error logs** always print to `stderr`.
* Each log line includes:

  * A prefix (`INFO:` or `ERROR:`)
  * Date/time
  * Filename and line number
  * Session UID (unique per Init)
* You can wrap errors using `loger.Return(err)` to automatically add file and line metadata.

---

## 🧩 Example Usage

### ✅ Basic Info and Error Logs

```go
package main

import (
	"errors"
	"github.com/saravanan611/loger"
)

func main() {
	l := loger.Init()

	l.Info("Application started successfully")
	
	err := openFile()
	if err != nil {
		l.Err(err)
	}
}

func openFile() error {
	return loger.Return(errors.New("Failed to open file"))
}
```

### 🧾 Output (Example)

```
INFO: 2025/10/04 08:20:13 main.go:10: [b9f9a46ebbcf4b10b4edb3d1c9a65c94] Application started successfully
ERROR: 2025/10/04 08:20:13 openfile.go:15: [b9f9a46ebbcf4b10b4edb3d1c9a65c94] test/openfile.go:15 [b9f9a46ebbcf4b10b4edb3d1c9a65c94] Failed to open file
```

---

## ⚡ Wrapping Errors Everywhere

You can safely wrap any returned error in your project using:

```go
return loger.Return(err)
```

It automatically adds the file name and line number to help you **trace exactly where the error came from.**

Example:

```go
if dbConn == nil {
    return loger.Return(errors.New("database connection is nil"))
}
```

---

## 🧩 Advanced Example

```go
func processData(l *loger.LogStruct) error {
	data, err := os.ReadFile("missing.txt")
	if err != nil {
		return loger.Return(fmt.Errorf("read failed: %w", err))
	}
	l.Info("Processing data:", string(data))
	return nil
}

func main() {
	l := loger.Init()
	if err := processData(l); err != nil {
		l.Err(err)
	}
}
```

### Output:

```
ERROR: 2025/10/04 08:31:52 process.go:12: [9d51aa4d40c7444f8b2f0b197b04c229] test/process.go:12 [9d51aa4d40c7444f8b2f0b197b04c229] read failed: open missing.txt: no such file or directory
```

---

## 🧩 Summary

| Feature                 | Description                                                  |
| ----------------------- | ------------------------------------------------------------ |
| **Unique ID**           | Each `Init()` call creates a session ID shown in all logs    |
| **File Info in Errors** | Automatically shows where the error occurred                 |
| **Env Control**         | Toggle info logs via `InfoFlog=Y`                            |
| **Clean API**           | `l.Info()`, `l.Err()`, and `loger.Return()` are all you need |

---

## 🧰 Best Practices

✅ Always wrap your returned errors:

```go
return loger.Return(err)
```

✅ Always log using the same session instance:

```go
l := loger.Init()
```

✅ Disable info logs in production by setting:

```bash
export InfoFlog=N
```
---
