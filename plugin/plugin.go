package plugin

import (
	"pfi/nobu/snowflake"
	"pfi/sensorbee/sensorbee/bql"
)

func init() {
	bql.RegisterGlobalUDSCreator("snowflake", bql.UDSCreatorFunc(snowflake.NewState))
	// TODO: register UDF
}
