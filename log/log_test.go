package log

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	Convey("日志模块", t, func() {
		So(1, ShouldEqual, 1)
		So(2, ShouldEqual, 2)
		So(3, ShouldEqual, 3)
	})
}
