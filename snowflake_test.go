package snowflake

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"testing"
	"time"
)

func TestSnowflake(t *testing.T) {
	ctx := core.NewContext(&core.ContextConfig{})

	{
		s, err := NewState(ctx, data.Map{
			"machine_id": data.Int(601),
		})
		if err != nil {
			t.Fatal(err)
		}
		if err := ctx.SharedStates.Add("test_snowflake", "snowflake_id", s); err != nil {
			t.Fatal(err)
		}
	}

	Convey("Given a snowflake state", t, func() {
		Convey("when calling the snowflake function within the same millisecond", func() {
			var (
				now time.Time
				v   data.Value
				err error
			)
			for {
				now = time.Now()
				v, err = Snowflake(ctx, data.String("test_snowflake"))
				if err != nil {
					So(err, ShouldBeNil)
				}
				if time.Now().Sub(now)/time.Millisecond == 0 {
					break
				}
			}
			So(err, ShouldBeNil)

			id, err := data.ToInt(v)
			So(err, ShouldBeNil)

			Convey("the value sholud contain the current millisecond", func() {
				elapse := now.UnixNano()/int64(time.Millisecond) -
					pluginPublishTime.UnixNano()/int64(time.Millisecond)
				So(id&(((1<<41)-1)<<22), ShouldEqual, elapse<<22)
			})

			Convey("the value should contain the machine id", func() {
				So((id>>12)&((1<<10)-1), ShouldEqual, 601)
			})
		})

		Convey("when calling the snowflake function multiple times within the same millisecond", func() {
			var ids []int64
			for {
				now := time.Now()
				var a []data.Value
				for i := 0; i < 3; i++ {
					v, err := Snowflake(ctx, data.String("test_snowflake"))
					if err != nil {
						So(err, ShouldBeNil)
					}
					a = append(a, v)
				}
				if time.Now().Sub(now)/time.Millisecond == 0 {
					for _, v := range a {
						id, err := data.ToInt(v)
						So(err, ShouldBeNil)
						ids = append(ids, id)
					}
					break
				}
			}

			Convey("the ids should have the same prefix", func() {
				for i := 1; i < len(ids); i++ {
					So(ids[i] & ^((1<<12)-1), ShouldEqual, ids[0] & ^((1<<12)-1))
				}
			})

			Convey("ids should differ 1", func() {
				for i := 1; i < len(ids); i++ {
					So(ids[i]-ids[i-1], ShouldEqual, 1)
				}
			})
		})
	})
}
