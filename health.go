package health

import (
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

// todo - фигово что это не в инициализации
func (health *Health) SetCanal(canal *canal.Canal) {
	health.canal = canal
}

// возвращает актуальное состояние "здоровья" демона
func (health *Health) Health() *Health {
	health.Lifetime = uint64(time.Since(health.start).Seconds())
	health.BinLogPosition = health.canal.SyncedPosition().Pos
	health.BinLogFile = health.canal.SyncedPosition().Name
	return health
}

func (health *Health) CanalListener(e *canal.RowsEvent) {
	if health.heartbeatTable == e.Table.Name {
		if len(e.Rows) > 1 {
			health.updateHeartbeat(e.Rows[1][1].(string))
		}
	}
}

// принимает datetime из mysql в формате "yyyy-mm-dd hh:ii:ss" и преобразовыват в unix_timestamp
func (health *Health) updateHeartbeat(datetime string) {
	t, err := time.Parse("2006-01-02 15:04:05", datetime)
	if err == nil {
		health.Heartbeat = t.Unix()
	}
}

func (health *Health) incWarning() {
	atomic.AddUint64(&health.Warnings, 1)
}

func (health *Health) incError() {
	atomic.AddUint64(&health.Errors, 1)
}
