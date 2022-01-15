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

package migration_tools

import (
	"context"
	"strconv"

	"github.com/ethereum/go-ethereum/statediff/indexer/database/sql/postgres"
)

// PgReadMinAndMaxBlockNumbers for finding the min and max block height in the DB
const PgReadMinAndMaxBlockNumbers = `SELECT MIN(block_number) min, MAX(block_number) max
									FROM eth.header_cids`

// MinAndMux struct to hold min and max block_number values
type MinAndMux struct {
	Min string `db:"min"`
	Max string `db:"max"`
}

// SegmentRangeByChunkSize splits the provided range up into segments based on the desired size of the segments
func SegmentRangeByChunkSize(chunkSize, start, stop uint64) [][2]uint64 {
	totalRangeSize := stop - start + 1
	numOfChunks := totalRangeSize / chunkSize
	remainder := totalRangeSize % chunkSize

	chunks := make([][2]uint64, numOfChunks)
	for i := uint64(0); i < numOfChunks; i++ {
		chunkStart := start + (i * chunkSize)
		chunkEnd := chunkStart + chunkSize - 1
		chunks[i] = [2]uint64{start, chunkEnd}
		if i == numOfChunks-1 && remainder != 0 {
			chunks = append(chunks, [2]uint64{chunkEnd + 1, chunkEnd + remainder})
		}
	}
	return chunks
}

// DetectAndSegmentRangeByChunkSize finds the min and max block heights in the DB, and breaks the range
// up into segments based on the provided chunk size
func DetectAndSegmentRangeByChunkSize(readConf postgres.Config, chunkSize uint64) ([][2]uint64, error) {
	readDB, err := NewDB(context.Background(), readConf)
	if err != nil {
		return nil, err
	}
	defer readDB.Close()
	minMax := new(MinAndMux)
	if err := readDB.Get(minMax, PgReadMinAndMaxBlockNumbers); err != nil {
		return nil, err
	}
	min, err := strconv.ParseUint(minMax.Min, 10, 64)
	if err != nil {
		return nil, err
	}
	max, err := strconv.ParseUint(minMax.Max, 10, 64)
	if err != nil {
		return nil, err
	}
	return SegmentRangeByChunkSize(chunkSize, min, max), nil
}
