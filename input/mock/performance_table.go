/*
Copyright [2018] [jc3wish]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mock

import (
	"context"
	pluginDriver "github.com/brokercap/Bifrost/plugin/driver"
	"github.com/brokercap/Bifrost/sdk/pluginTestData"
	"math/rand"
	"time"
)

type PerformanceTable struct {
	SchemaName          string
	TableName           string
	TableDataCount      int
	InterVal            time.Duration
	TableRowsEventCount int
	DeleteEventRatio    int
	BatchSize           int
	ch                  chan *pluginDriver.PluginDataType
	event               *pluginTestData.Event
}

func (t *PerformanceTable) GetSchemaName() string {
	return t.SchemaName
}

func (t *PerformanceTable) GetTableName() string {
	return t.TableName
}

func (t *PerformanceTable) Start(ctx context.Context, ch chan *pluginDriver.PluginDataType) {
	t.ch = ch
	t.event = pluginTestData.NewEvent()
	t.event.SetSchema(t.SchemaName)
	t.event.SetTable(t.TableName)

	var count int
	var halfDataCount int
	if t.TableDataCount >= t.TableRowsEventCount {
		t.TableRowsEventCount = t.TableDataCount
		halfDataCount = t.TableDataCount
	} else {
		halfDataCount = t.TableDataCount / 2
	}

	timer := time.NewTimer(t.InterVal)
	t.Batch(&count, halfDataCount)
	for {
		select {
		case <-timer.C:
			t.Batch(&count, halfDataCount)
		case <-ctx.Done():
			return
		}
	}

}

func (t *PerformanceTable) Batch(count *int, halfDataCount int) {
	for i := 0; i < t.BatchSize; i++ {
		if *count < halfDataCount {
			t.Callback(t.event.GetTestInsertData())
			*count++
			continue
		}
		n := rand.Intn(100) + 1
		if n <= t.DeleteEventRatio {
			t.Callback(t.event.GetTestDeleteData())
			*count--
		}
		if *count >= t.TableDataCount {
			t.Callback(t.event.GetTestUpdateData(true))
			continue
		}
		if n <= 50 {
			t.Callback(t.event.GetTestUpdateData(true))
			continue
		}
		t.Callback(t.event.GetTestInsertData())
		*count++
	}
}

func (t *PerformanceTable) Callback(data *pluginDriver.PluginDataType) {
	if t.ch == nil {
		return
	}
	t.ch <- data
}
