# Snowflake UDF for SensorBee

This UDF generates IDs using Snowflake algorithm proposed by Twitter.
The ID is a 63 bit integer consisting of:

* 41 bit of timestamp in milliseconds (works for 69 years)
* 10 bit of machine id which must manually be assigned to each machine
* 12 bit of counter used for IDs having the same timestamp (generates up to 4096 IDs within a millisecond)

# Usage

## Registering plugin

Just import plugin package from an application:

```
import (
    _ "pfi/sensorbee/snowflake/plugin"
)
```

Or, register the user defined state and UDF manually to bql package.

## Using UDF from BQL

```
-- Create a user defined state for snowflake UDF.
CREATE STATE event_id_seq TYPE snowflake_id WITH machine_id=1

-- Assign IDs to an event sequence. IDs will be generated based on
-- the state 'event_id_seq'
CREATE STREAM events_with_id AS SELECT snowflake_id('event_id_seq'), * FROM events;
```

# TODO

* Provide timestamp_offset which allows users to generate ID available for 69 years from "now"
