package daemon

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/inkah-trace/common"
	"github.com/satori/go.uuid"
)

func Unmarshall(message []byte, e *inkah.Event) error {
	sm := string(message)
	sp := strings.Split(sm, ":")

	tsRaw := sp[0]

	i, err := strconv.ParseInt(tsRaw, 10, 64)
	if err != nil {
		em := fmt.Sprintf("Unable to parse event timestamp %s as Int: %s", tsRaw, err)
		return errors.New(em)
	}

	hn, err := os.Hostname()
	if err != nil {
		em := fmt.Sprintf("Unable to obtain hostname: %s", err)
		return errors.New(em)
	}

	ts := time.Unix(i, 0)

	e.Id = uuid.NewV4().String()
	e.Timestamp = ts
	e.Hostname = hn
	e.ServiceName = sp[1]
	e.TraceId = sp[2]
	e.SpanId = sp[3]
	e.ParentSpanId = sp[4]
	e.EventType = inkah.EventType(sp[5])
	e.Data = make([]*inkah.EventData, 0)

	if len(sp) >= 8 {
		for i := 6; i < len(sp); i = i + 2 {
			ed := inkah.EventData{}

			key := sp[i]
			var value string
			if len(sp) >= i+2 {
				value = sp[i+1]
			}

			ed.Key = key
			ed.Value = value

			e.Data = append(e.Data, &ed)
		}
	}

	return nil
}
