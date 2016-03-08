package plugin

import (
	"github.com/sensorbee/snowflake"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
)

func init() {
	if err := udf.RegisterGlobalUDSCreator("snowflake_id", udf.UDSCreatorFunc(snowflake.NewState)); err != nil {
		panic(err)
	}
	if err := udf.RegisterGlobalUDF("snowflake_id", udf.UnaryFunc(snowflake.Snowflake)); err != nil {
		panic(err)
	}
}
