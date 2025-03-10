// Copyright 2023-2024 daviszhen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plan

import (
	"testing"
)

const (
	NOT_CAST      = 0
	EXPLICIT_CAST = 1
	IMPLICIT_CAST = 2
)

var typeMatrix = [120][120]int{}

func fillTypeMatrix() {
	//blob
	typeMatrix[LTID_BLOB][LTID_BIT] = EXPLICIT_CAST
	typeMatrix[LTID_BLOB][LTID_VARCHAR] = EXPLICIT_CAST
	//interval
	typeMatrix[LTID_INTERVAL][LTID_VARCHAR] = IMPLICIT_CAST
	//date
	typeMatrix[LTID_DATE][LTID_TIMESTAMP] = IMPLICIT_CAST
	typeMatrix[LTID_DATE][LTID_VARCHAR] = IMPLICIT_CAST
	//time
	typeMatrix[LTID_TIME][LTID_VARCHAR] = IMPLICIT_CAST
	//timestamp
	typeMatrix[LTID_TIMESTAMP][LTID_DATE] = EXPLICIT_CAST
	typeMatrix[LTID_TIMESTAMP][LTID_TIME] = EXPLICIT_CAST
	typeMatrix[LTID_TIMESTAMP][LTID_VARCHAR] = IMPLICIT_CAST
	//boolean
	typeMatrix[LTID_BOOLEAN][LTID_BIT] = EXPLICIT_CAST
	typeMatrix[LTID_BOOLEAN][LTID_DOUBLE] = EXPLICIT_CAST
	typeMatrix[LTID_BOOLEAN][LTID_FLOAT] = EXPLICIT_CAST
	typeMatrix[LTID_BOOLEAN][LTID_HUGEINT] = EXPLICIT_CAST
	typeMatrix[LTID_BOOLEAN][LTID_BIGINT] = EXPLICIT_CAST
	typeMatrix[LTID_BOOLEAN][LTID_INTEGER] = EXPLICIT_CAST
	typeMatrix[LTID_BOOLEAN][LTID_SMALLINT] = EXPLICIT_CAST
	typeMatrix[LTID_BOOLEAN][LTID_TINYINT] = EXPLICIT_CAST
	typeMatrix[LTID_BOOLEAN][LTID_UBIGINT] = EXPLICIT_CAST
	typeMatrix[LTID_BOOLEAN][LTID_UINTEGER] = EXPLICIT_CAST
	typeMatrix[LTID_BOOLEAN][LTID_USMALLINT] = EXPLICIT_CAST
	typeMatrix[LTID_BOOLEAN][LTID_UTINYINT] = EXPLICIT_CAST
	typeMatrix[LTID_BOOLEAN][LTID_DECIMAL] = EXPLICIT_CAST
	typeMatrix[LTID_BOOLEAN][LTID_VARCHAR] = IMPLICIT_CAST
	//bit
	typeMatrix[LTID_BIT][LTID_BLOB] = EXPLICIT_CAST
	typeMatrix[LTID_BIT][LTID_DOUBLE] = EXPLICIT_CAST
	typeMatrix[LTID_BIT][LTID_FLOAT] = EXPLICIT_CAST
	typeMatrix[LTID_BIT][LTID_HUGEINT] = EXPLICIT_CAST
	typeMatrix[LTID_BIT][LTID_BIGINT] = EXPLICIT_CAST
	typeMatrix[LTID_BIT][LTID_INTEGER] = EXPLICIT_CAST
	typeMatrix[LTID_BIT][LTID_UBIGINT] = EXPLICIT_CAST
	typeMatrix[LTID_BIT][LTID_UINTEGER] = EXPLICIT_CAST
	typeMatrix[LTID_BIT][LTID_VARCHAR] = IMPLICIT_CAST
	//double
	typeMatrix[LTID_DOUBLE][LTID_BOOLEAN] = EXPLICIT_CAST
	typeMatrix[LTID_DOUBLE][LTID_BIT] = EXPLICIT_CAST
	typeMatrix[LTID_DOUBLE][LTID_VARCHAR] = IMPLICIT_CAST
	//float
	typeMatrix[LTID_FLOAT][LTID_BOOLEAN] = EXPLICIT_CAST
	typeMatrix[LTID_FLOAT][LTID_BIT] = EXPLICIT_CAST
	typeMatrix[LTID_FLOAT][LTID_DOUBLE] = IMPLICIT_CAST
	typeMatrix[LTID_FLOAT][LTID_VARCHAR] = IMPLICIT_CAST
	//hugeint
	typeMatrix[LTID_HUGEINT][LTID_BOOLEAN] = EXPLICIT_CAST
	typeMatrix[LTID_HUGEINT][LTID_BIT] = EXPLICIT_CAST
	typeMatrix[LTID_HUGEINT][LTID_DOUBLE] = IMPLICIT_CAST
	typeMatrix[LTID_HUGEINT][LTID_FLOAT] = IMPLICIT_CAST
	typeMatrix[LTID_HUGEINT][LTID_DECIMAL] = IMPLICIT_CAST
	typeMatrix[LTID_HUGEINT][LTID_VARCHAR] = IMPLICIT_CAST
	//bigint
	typeMatrix[LTID_BIGINT][LTID_BOOLEAN] = EXPLICIT_CAST
	typeMatrix[LTID_BIGINT][LTID_BIT] = EXPLICIT_CAST
	typeMatrix[LTID_BIGINT][LTID_DOUBLE] = IMPLICIT_CAST
	typeMatrix[LTID_BIGINT][LTID_FLOAT] = IMPLICIT_CAST
	typeMatrix[LTID_BIGINT][LTID_HUGEINT] = IMPLICIT_CAST
	typeMatrix[LTID_BIGINT][LTID_DECIMAL] = IMPLICIT_CAST
	typeMatrix[LTID_BIGINT][LTID_VARCHAR] = IMPLICIT_CAST
	//integer
	typeMatrix[LTID_INTEGER][LTID_BOOLEAN] = EXPLICIT_CAST
	typeMatrix[LTID_INTEGER][LTID_BIT] = EXPLICIT_CAST
	typeMatrix[LTID_INTEGER][LTID_DOUBLE] = IMPLICIT_CAST
	typeMatrix[LTID_INTEGER][LTID_FLOAT] = IMPLICIT_CAST
	typeMatrix[LTID_INTEGER][LTID_HUGEINT] = IMPLICIT_CAST
	typeMatrix[LTID_INTEGER][LTID_BIGINT] = IMPLICIT_CAST
	typeMatrix[LTID_INTEGER][LTID_DECIMAL] = IMPLICIT_CAST
	typeMatrix[LTID_INTEGER][LTID_VARCHAR] = IMPLICIT_CAST
	//smallint
	typeMatrix[LTID_SMALLINT][LTID_BOOLEAN] = EXPLICIT_CAST
	typeMatrix[LTID_SMALLINT][LTID_BIT] = EXPLICIT_CAST
	typeMatrix[LTID_SMALLINT][LTID_DOUBLE] = IMPLICIT_CAST
	typeMatrix[LTID_SMALLINT][LTID_FLOAT] = IMPLICIT_CAST
	typeMatrix[LTID_SMALLINT][LTID_HUGEINT] = IMPLICIT_CAST
	typeMatrix[LTID_SMALLINT][LTID_BIGINT] = IMPLICIT_CAST
	typeMatrix[LTID_SMALLINT][LTID_INTEGER] = IMPLICIT_CAST
	typeMatrix[LTID_SMALLINT][LTID_DECIMAL] = IMPLICIT_CAST
	typeMatrix[LTID_SMALLINT][LTID_VARCHAR] = IMPLICIT_CAST
	//tinyint
	typeMatrix[LTID_TINYINT][LTID_BOOLEAN] = EXPLICIT_CAST
	typeMatrix[LTID_TINYINT][LTID_BIT] = EXPLICIT_CAST
	typeMatrix[LTID_TINYINT][LTID_DOUBLE] = IMPLICIT_CAST
	typeMatrix[LTID_TINYINT][LTID_FLOAT] = IMPLICIT_CAST
	typeMatrix[LTID_TINYINT][LTID_HUGEINT] = IMPLICIT_CAST
	typeMatrix[LTID_TINYINT][LTID_BIGINT] = IMPLICIT_CAST
	typeMatrix[LTID_TINYINT][LTID_INTEGER] = IMPLICIT_CAST
	typeMatrix[LTID_TINYINT][LTID_SMALLINT] = IMPLICIT_CAST
	typeMatrix[LTID_TINYINT][LTID_DECIMAL] = IMPLICIT_CAST
	typeMatrix[LTID_TINYINT][LTID_VARCHAR] = IMPLICIT_CAST
	//ubigint
	typeMatrix[LTID_UBIGINT][LTID_BOOLEAN] = EXPLICIT_CAST
	typeMatrix[LTID_UBIGINT][LTID_BIT] = EXPLICIT_CAST
	typeMatrix[LTID_UBIGINT][LTID_DOUBLE] = IMPLICIT_CAST
	typeMatrix[LTID_UBIGINT][LTID_FLOAT] = IMPLICIT_CAST
	typeMatrix[LTID_UBIGINT][LTID_HUGEINT] = IMPLICIT_CAST
	typeMatrix[LTID_UBIGINT][LTID_BIGINT] = EXPLICIT_CAST
	typeMatrix[LTID_UBIGINT][LTID_INTEGER] = EXPLICIT_CAST
	typeMatrix[LTID_UBIGINT][LTID_SMALLINT] = EXPLICIT_CAST
	typeMatrix[LTID_UBIGINT][LTID_TINYINT] = EXPLICIT_CAST
	typeMatrix[LTID_UBIGINT][LTID_UINTEGER] = EXPLICIT_CAST
	typeMatrix[LTID_UBIGINT][LTID_USMALLINT] = EXPLICIT_CAST
	typeMatrix[LTID_UBIGINT][LTID_UTINYINT] = EXPLICIT_CAST
	typeMatrix[LTID_UBIGINT][LTID_DECIMAL] = IMPLICIT_CAST
	typeMatrix[LTID_UBIGINT][LTID_VARCHAR] = IMPLICIT_CAST
	//uinteger
	typeMatrix[LTID_UINTEGER][LTID_BOOLEAN] = EXPLICIT_CAST
	typeMatrix[LTID_UINTEGER][LTID_BIT] = EXPLICIT_CAST
	typeMatrix[LTID_UINTEGER][LTID_DOUBLE] = IMPLICIT_CAST
	typeMatrix[LTID_UINTEGER][LTID_FLOAT] = IMPLICIT_CAST
	typeMatrix[LTID_UINTEGER][LTID_HUGEINT] = IMPLICIT_CAST
	typeMatrix[LTID_UINTEGER][LTID_BIGINT] = IMPLICIT_CAST
	typeMatrix[LTID_UINTEGER][LTID_INTEGER] = EXPLICIT_CAST
	typeMatrix[LTID_UINTEGER][LTID_SMALLINT] = EXPLICIT_CAST
	typeMatrix[LTID_UINTEGER][LTID_TINYINT] = EXPLICIT_CAST
	typeMatrix[LTID_UINTEGER][LTID_UBIGINT] = IMPLICIT_CAST
	typeMatrix[LTID_UINTEGER][LTID_USMALLINT] = EXPLICIT_CAST
	typeMatrix[LTID_UINTEGER][LTID_UTINYINT] = EXPLICIT_CAST
	typeMatrix[LTID_UINTEGER][LTID_DECIMAL] = IMPLICIT_CAST
	typeMatrix[LTID_UINTEGER][LTID_VARCHAR] = IMPLICIT_CAST
	//usmallint
	typeMatrix[LTID_USMALLINT][LTID_BOOLEAN] = EXPLICIT_CAST
	typeMatrix[LTID_USMALLINT][LTID_BIT] = EXPLICIT_CAST
	typeMatrix[LTID_USMALLINT][LTID_DOUBLE] = IMPLICIT_CAST
	typeMatrix[LTID_USMALLINT][LTID_FLOAT] = IMPLICIT_CAST
	typeMatrix[LTID_USMALLINT][LTID_HUGEINT] = IMPLICIT_CAST
	typeMatrix[LTID_USMALLINT][LTID_BIGINT] = IMPLICIT_CAST
	typeMatrix[LTID_USMALLINT][LTID_INTEGER] = IMPLICIT_CAST
	typeMatrix[LTID_USMALLINT][LTID_SMALLINT] = EXPLICIT_CAST
	typeMatrix[LTID_USMALLINT][LTID_TINYINT] = EXPLICIT_CAST
	typeMatrix[LTID_USMALLINT][LTID_UBIGINT] = IMPLICIT_CAST
	typeMatrix[LTID_USMALLINT][LTID_UINTEGER] = IMPLICIT_CAST
	typeMatrix[LTID_USMALLINT][LTID_UTINYINT] = EXPLICIT_CAST
	typeMatrix[LTID_USMALLINT][LTID_DECIMAL] = IMPLICIT_CAST
	typeMatrix[LTID_USMALLINT][LTID_VARCHAR] = IMPLICIT_CAST
	//utinyint
	typeMatrix[LTID_UTINYINT][LTID_BOOLEAN] = EXPLICIT_CAST
	typeMatrix[LTID_UTINYINT][LTID_BIT] = EXPLICIT_CAST
	typeMatrix[LTID_UTINYINT][LTID_DOUBLE] = IMPLICIT_CAST
	typeMatrix[LTID_UTINYINT][LTID_FLOAT] = IMPLICIT_CAST
	typeMatrix[LTID_UTINYINT][LTID_HUGEINT] = IMPLICIT_CAST
	typeMatrix[LTID_UTINYINT][LTID_BIGINT] = IMPLICIT_CAST
	typeMatrix[LTID_UTINYINT][LTID_INTEGER] = IMPLICIT_CAST
	typeMatrix[LTID_UTINYINT][LTID_SMALLINT] = IMPLICIT_CAST
	typeMatrix[LTID_UTINYINT][LTID_TINYINT] = EXPLICIT_CAST
	typeMatrix[LTID_UTINYINT][LTID_UBIGINT] = IMPLICIT_CAST
	typeMatrix[LTID_UTINYINT][LTID_UINTEGER] = IMPLICIT_CAST
	typeMatrix[LTID_UTINYINT][LTID_USMALLINT] = IMPLICIT_CAST
	typeMatrix[LTID_UTINYINT][LTID_DECIMAL] = IMPLICIT_CAST
	typeMatrix[LTID_UTINYINT][LTID_VARCHAR] = IMPLICIT_CAST
	//decimal
	typeMatrix[LTID_DECIMAL][LTID_BOOLEAN] = EXPLICIT_CAST
	typeMatrix[LTID_DECIMAL][LTID_DOUBLE] = IMPLICIT_CAST
	typeMatrix[LTID_DECIMAL][LTID_FLOAT] = IMPLICIT_CAST
	typeMatrix[LTID_DECIMAL][LTID_HUGEINT] = EXPLICIT_CAST
	typeMatrix[LTID_DECIMAL][LTID_BIGINT] = EXPLICIT_CAST
	typeMatrix[LTID_DECIMAL][LTID_INTEGER] = EXPLICIT_CAST
	typeMatrix[LTID_DECIMAL][LTID_SMALLINT] = EXPLICIT_CAST
	typeMatrix[LTID_DECIMAL][LTID_VARCHAR] = IMPLICIT_CAST
	//uuid
	typeMatrix[LTID_UUID][LTID_VARCHAR] = IMPLICIT_CAST
}

var ltypes = []LTypeId{
	LTID_INVALID,
	LTID_NULL,
	LTID_UNKNOWN,
	LTID_ANY,
	LTID_USER,
	LTID_BOOLEAN,
	LTID_TINYINT,
	LTID_SMALLINT,
	LTID_INTEGER,
	LTID_BIGINT,
	LTID_DATE,
	LTID_TIME,
	LTID_TIMESTAMP_SEC,
	LTID_TIMESTAMP_MS,
	LTID_TIMESTAMP,
	LTID_TIMESTAMP_NS,
	LTID_DECIMAL,
	LTID_FLOAT,
	LTID_DOUBLE,
	LTID_CHAR,
	LTID_VARCHAR,
	LTID_BLOB,
	LTID_INTERVAL,
	LTID_UTINYINT,
	LTID_USMALLINT,
	LTID_UINTEGER,
	LTID_UBIGINT,
	LTID_TIMESTAMP_TZ,
	LTID_TIME_TZ,
	LTID_BIT,
	LTID_HUGEINT,
	LTID_POINTER,
	LTID_VALIDITY,
	LTID_UUID,
	LTID_STRUCT,
	LTID_LIST,
	LTID_MAP,
	LTID_TABLE,
	LTID_ENUM,
	LTID_AGGREGATE_STATE,
	LTID_LAMBDA,
	LTID_UNION,
}

func init() {
	fillTypeMatrix()
}

func Test_type(t *testing.T) {

}
