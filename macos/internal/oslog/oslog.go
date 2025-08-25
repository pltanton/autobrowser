package oslog

/*
#cgo LDFLAGS:
#include <os/log.h>
#include <stdlib.h>

static inline os_log_t create_log(const char* subsystem, const char* category) {
    return os_log_create(subsystem, category);
}

static inline void write_log(os_log_t log, uint8_t typ, const char* msg) {
    os_log_with_type(log, (os_log_type_t)typ, "%{public}s", msg);
}
*/
import "C"
import (
	"log/slog"
	"unsafe"
)

type Logger struct {
	log C.os_log_t
}

func NewLogger(subsystem, category string) *Logger {
	csub := C.CString(subsystem)
	ccat := C.CString(category)
	defer C.free(unsafe.Pointer(csub))
	defer C.free(unsafe.Pointer(ccat))

	return &Logger{log: C.create_log(csub, ccat)}
}

func (l *Logger) Log(level slog.Level, msg string) {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))
	C.write_log(l.log, C.uchar(osLevel(level)), cmsg)
}

func osLevel(l slog.Level) uint8 {
	switch {
	case l >= slog.LevelError:
		return 16 // ERROR
	case l >= slog.LevelInfo:
		return 1 // INFO
	case l >= slog.LevelDebug:
		return 2 // DEBUG
	default:
		return 0 // DEFAULT
	}
}
