package health

import (
	"errors"
	"github.com/siddontang/go-mysql/canal"
	"sync/atomic"
	"time"
)

const HeartbeatTable = "SystemEvents"

//Main struct of our daemon
type Health struct {
	start          time.Time //Daemon start time
	canal          *canal.Canal
	heartbeatTable string //Where is our heartbeats?

	Version        string         //Current daemon version
	Lifetime       uint64         //How many times are we alive?
	Errors         uint64         //Count of errors
	Warnings       uint64         //Count of warnings
	Heartbeat      int64          //How many heartbeats were?
	BinLogPosition uint32         //Current binlog position
	BinLogFile     string         //Current binlog file
	Additional     interface{}    //Additional info
	CacheState     map[string]int //How many elements we have in cache now
}

// Create and return new Health object
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

// Set canal to health object
func (health *Health) SetCanal(canal *canal.Canal) {
	health.canal = canal
}

// Return current daemon Health
func (health *Health) Health() error {
	health.Lifetime = uint64(time.Since(health.start).Seconds())
	if health.canal == nil {
		return errors.New("health.canal is empty")
	}

	health.BinLogPosition = health.canal.SyncedPosition().Pos
	health.BinLogFile = health.canal.SyncedPosition().Name

	return nil
}

//Check heartbeat table and update heartbeat counter
func (health *Health) CanalListener(e *canal.RowsEvent) error {
	if health.heartbeatTable == e.Table.Name {
		if len(e.Rows) > 1 {
			return health.updateHeartbeat(e.Rows[1][1].(string))
		}
	}
	return nil
}

// Convert mysql datetime ("yyyy-mm-dd hh:ii:ss") to unix_timestamp
func (health *Health) updateHeartbeat(datetime string) error {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", datetime, time.Local)
	if err == nil {
		health.Heartbeat = t.Unix()
	}
	return err
}

// Increase warning counter
func (health *Health) incWarning() {
	atomic.AddUint64(&health.Warnings, 1)
}

// Increase error counter
func (health *Health) incError() {
	atomic.AddUint64(&health.Errors, 1)
}
