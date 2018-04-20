package health

import (
	"fmt"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/schema"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSyncInc(t *testing.T) {
	assert := assert.New(t)
	health := New("1.0.0")

	const len uint64 = 10
	var i uint64
	for ; i < len; i++ {
		health.incWarning()
		health.incError()
	}

	assert.Equal(len, health.Warnings)
	assert.Equal(len, health.Errors)
}

func TestUpdateHeartbeat(t *testing.T) {
	assert := assert.New(t)
	ts := time.Now()
	testCases := []struct {
		value    string
		expected int64
		isError  bool
	}{
		{ts.Format("2006-01-02 15:04:05"), ts.Unix(), false},
		{ts.Format("2006-01-02"), 0, true},
		{ts.Format("abcde"), 0, true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("value:%s, expected:%d, error:%v", tc.value, tc.expected, tc.isError), func(t *testing.T) {
			health := New("1.0.0")

			err := health.updateHeartbeat(tc.value)
			assert.Equal(tc.expected, health.Heartbeat)
			assert.Equal(tc.isError, err != nil)
		})
	}
}

func TestCanalListener(t *testing.T) {
	assert := assert.New(t)
	ts := time.Now()

	newRowEvent := func(table string, action string, rows [][]interface{}) *canal.RowsEvent {
		e := new(canal.RowsEvent)
		t := new(schema.Table)
		t.Name = table

		e.Table = t
		e.Action = action
		e.Rows = rows
		return e
	}

	testCases := []struct {
		event    *canal.RowsEvent
		expected int64
	}{
		{
			newRowEvent(
				HeartbeatTable,
				canal.UpdateAction,
				[][]interface{}{{"Heartbeat", "2018-01-12 22:08:33"}, {"Heartbeat", ts.Format("2006-01-02 15:04:05")}}),
			ts.Unix(),
		},
		{
			newRowEvent(
				HeartbeatTable,
				canal.UpdateAction,
				[][]interface{}{{"Heartbeat", "2018-01-12 22:08:33"}}),
			0,
		},
		{
			newRowEvent(
				"someTable",
				canal.UpdateAction,
				[][]interface{}{{"Heartbeat", "2018-01-12 22:08:33"}, {"Heartbeat", ts.Format("2006-01-02 15:04:05")}}),
			0,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("event:%v, expected:%d", tc.event, tc.expected), func(t *testing.T) {
			health := New("1.0.0")
			health.CanalListener(tc.event)
			assert.Equal(tc.expected, health.Heartbeat)
		})
	}
}

func TestHealthEmptyCanal(t *testing.T) {
	assert := assert.New(t)

	health := New("1.0.0")
	err := health.Health()

	assert.EqualError(err, "health.canal is empty")
	assert.Equal(uint64(time.Since(health.start).Seconds()), health.Lifetime)
}

func TestHealthSetCanal(t *testing.T) {
	assert := assert.New(t)
	health := New("1.0.0")
	canal := new(canal.Canal)
	canal.AddDumpDatabases("12")
	health.SetCanal(canal)
	assert.Equal(canal, health.canal)
}

