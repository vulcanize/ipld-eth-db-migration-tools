// VulcanizeDB
// Copyright © 2022 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package public_blocks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/migration-tools/pkg/util"
)

const DefaultV2FDWTableName = "v2db_public_blocks"

const (
	getMaxPagePgStr   = `SELECT MAX(ctid)::TEXT FROM public.blocks`
	transferPagePgStr = "INSERT INTO public.blocks(key, data)" +
		"SELECT key, data FROM $1 WHERE ctid = ANY (ARRAY(SELECT ('($2,' || s.i || ')')::tid " +
		"FROM generate_series(0, current_setting('block_size')::int/4) AS s(i)))" +
		"ON CONFLICT (key) DO NOTHING"
)

func GetMaxPage(db *sqlx.DB) (uint64, error) {
	var ctidTuple string
	if err := db.Get(ctidTuple, getMaxPagePgStr); err != nil {
		return 0, err
	}
	tupleSplit := strings.Split(ctidTuple, ",")
	if len(tupleSplit) != 2 {
		return 0, fmt.Errorf("incorrect format for the tuple returned by getMaxPagePgStr: %s", ctidTuple)
	}
	maxPageStr := strings.TrimPrefix(tupleSplit[0], "(")
	return strconv.ParseUint(maxPageStr, 10, 64)
}

func GetPageSegments(maxPage, segmentSize uint64) [][2]uint64 {
	numSegments := maxPage / segmentSize
	remainder := maxPage % segmentSize
	if remainder != 0 {
		numSegments = numSegments + 1
	}
	segments := make([][2]uint64, numSegments)
	currentPage := uint64(0)
	for i := uint64(0); i <= numSegments; i++ {
		segments[i] = [2]uint64{currentPage, currentPage + segmentSize - 1}
		currentPage = currentPage + segmentSize
	}
	return segments
}

func TransferPages(db *sqlx.DB, fdwTableName string, firstPage, lastPage uint64) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			util.Rollback(tx)
			panic(p)
		} else if err != nil {
			util.Rollback(tx)
		}
	}()
	for i := firstPage; i <= lastPage; i++ {
		_, err = tx.Exec(transferPagePgStr, fdwTableName, i)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
