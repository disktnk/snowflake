package plugin

import (
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/snowflake"
)

func init() {
	if err := bql.RegisterGlobalUDSCreator("snowflake_id", bql.UDSCreatorFunc(snowflake.NewState)); err != nil {
		panic(err)
	}
	if err := udf.RegisterGlobalUDF("snowflake_id", udf.UnaryFunc(snowflake.Snowflake)); err != nil {
		panic(err)
	}
}
