package health

import (
	"errors"
	"time"
	"github.com/siddontang/go-mysql/canal"
	"sync/atomic"
)

const HeartbeatTable = "SystemEvents"

type Health struct {
	start          time.Time
	canal          *canal.Canal
	heartbeatTable string

	Version        string
	Lifetime       uint64
	Errors         uint64
	Warnings       uint64
	Heartbeat      int64
	BinLogPosition uint32
	BinLogFile     string
	Additional     interface{}
}

// создает и возвращает новый объект здоровья
func New(version string) *Health {
	return &Health{
		start:          time.Now(),
		heartbeatTable: HeartbeatTable,
		//
		Errors:   0,
		Warnings: 0,
		Version:  version,
	}
}

func (health *Health) SetCanal(canal *canal.Canal) {
	health.canal = canal
}

// актуальное состояние "здоровья" демона
func (health *Health) Health() error {
	health.Lifetime = uint64(time.Since(health.start).Seconds())
	if health.canal == nil {
		return errors.New("health.canal is empty")
	}

	health.BinLogPosition = health.canal.SyncedPosition().Pos
	health.BinLogFile = health.canal.SyncedPosition().Name

	return nil
}

func (health *Health) CanalListener(e *canal.RowsEvent) error {
	if health.heartbeatTable == e.Table.Name {
		if len(e.Rows) > 1 {
			return health.updateHeartbeat(e.Rows[1][1].(string))
		}
	}
	return nil
}

// принимает datetime из mysql в формате "yyyy-mm-dd hh:ii:ss" и преобразовыват в unix_timestamp
func (health *Health) updateHeartbeat(datetime string) error {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", datetime, time.Local)
	if err == nil {
		health.Heartbeat = t.Unix()
	}
	return err
}

// потоко-безопасно увеличивает счетчик предупреждений
func (health *Health) incWarning() {
	atomic.AddUint64(&health.Warnings, 1)
}

// потоко-безопасно увеличивает счетчик ошибок
func (health *Health) incError() {
	atomic.AddUint64(&health.Errors, 1)
}
