package plugin

import (
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/snowflake"
)

func init() {
	if err := udf.RegisterGlobalUDSCreator("snowflake_id", udf.UDSCreatorFunc(snowflake.NewState)); err != nil {
		panic(err)
	}
	if err := udf.RegisterGlobalUDF("snowflake_id", udf.UnaryFunc(snowflake.Snowflake)); err != nil {
		panic(err)
	}
}
