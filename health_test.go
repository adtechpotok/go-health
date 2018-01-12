package health

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"fmt"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/schema"
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
	//assert := assert.New(t)

	newRowEvent := func(table string, action string, rows [][]interface{}) *canal.RowsEvent {
		e := new(canal.RowsEvent)
		t := new(schema.Table)
		t.Name = table

		e.Table = t
		e.Action = action
		e.Rows = rows
		fmt.Printf("%v\n", rows)

		return e
	}

	vs := make([]interface{}, 2)

	testCases := []struct {
		table string
		event *canal.RowsEvent
	}{
		{
			HeartbeatTable,
			newRowEvent(HeartbeatTable, canal.UpdateAction, [][]interface{}{vs})},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("tableName:%s, event:%v", tc.table, tc.event), func(t *testing.T) {
			health := New("1.0.0")
			health.CanalListener(tc.event)
		})
	}
}
