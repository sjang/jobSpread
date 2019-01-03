package jobs

import (
	"fmt"
	"jobSpread/work"
	"testing"
)

func TestTrackMeta(t *testing.T) {
	worker := work.MyWorker()
	defer work.ReturnWorker(worker)

	fmt.Println(*t)
//	TrackMeta(worker, []int64{1,2,3})
}
