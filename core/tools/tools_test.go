package tools

import (
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEnsureDir(t *testing.T) {
	Convey("EnsureDir test", t, func() {
		dir1, _ := EnsureDir("~/.smg")
		So(dir1, ShouldEqual, "/Users/user/.smg")
		dir2, _ := EnsureDir("/a/b/c")
		So(dir2, ShouldEqual, "/a/b/c")
	})
}

func TestJoinSlash(t *testing.T) {
	Convey("JoinSlash test", t, func() {
		So(JoinSlash("a"), ShouldEqual, "a")
		So(JoinSlash(""), ShouldEqual, "")
		So(JoinSlash("", "a"), ShouldEqual, "a")
		So(JoinSlash("a", "b"), ShouldEqual, "a/b")
		So(JoinSlash("a", "b", "c/d"), ShouldEqual, "a/b/c/d")
	})
}

func TestKeySplit(t *testing.T) {
	Convey("KeySplit test", t, func() {
		So(SplitSlash("a"), ShouldResemble, []string{"a"})
		So(SplitSlash("a/b"), ShouldResemble, []string{"a", "b"})
	})
}

func TestMesgeH(t *testing.T) {
	Convey("MesgeH test", t, func() {
		a := H{"a": 1, "b": 1, "c": 1}
		b := H{"a": 2}
		c := H{"c": 3}
		So(MergeH(a, b, c), ShouldResemble, H{"a": 2, "b": 1, "c": 3})
	})
}

func TestAbsDir(t *testing.T) {
	wd, _ := os.Getwd()
	Convey("AbsDir test", t, func() {
		So(AbsDir("~/.smg"), ShouldEqual, "/Users/user/.smg")
		So(AbsDir("/Users/user/.smg"), ShouldEqual, "/Users/user/.smg")
		So(AbsDir("./smg"), ShouldEqual, path.Join(wd, "smg"))
		So(AbsDir("~/.smg/http.yaml"), ShouldEqual, "/Users/user/.smg/http.yaml")
	})
}
