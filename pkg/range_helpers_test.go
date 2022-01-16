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

package migration_tools_test

import (
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	migration_tools "github.com/vulcanize/migration-tools/pkg"
)

type chunkSizeTestCase struct {
	Start               uint64
	Stop                uint64
	ChunkSize           uint64
	ExpectedNumOfChunks uint64
}

var chunkSizeTestCases = []chunkSizeTestCase{
	{
		Start:               0,
		Stop:                99,
		ChunkSize:           10,
		ExpectedNumOfChunks: (99 - 0 + 1) / 10,
	},
	{
		Start:               0,
		Stop:                10000,
		ChunkSize:           19,
		ExpectedNumOfChunks: ((10000 + 1) / 19) + 1,
	},
	{
		Start:               0,
		Stop:                10000,
		ChunkSize:           100,
		ExpectedNumOfChunks: ((10000 + 1) / 100) + 1,
	},
	{
		Start:               21534252,
		Stop:                1123425221,
		ChunkSize:           10000,
		ExpectedNumOfChunks: ((1123425221 - 21534252 + 1) / 10000) + 1,
	},
	{
		Start:               21534252,
		Stop:                1123425221,
		ChunkSize:           7712,
		ExpectedNumOfChunks: ((1123425221 - 21534252 + 1) / 7712) + 1,
	},
	{
		Start:               123,
		Stop:                1123,
		ChunkSize:           10,
		ExpectedNumOfChunks: ((1123 - 123 + 1) / 10) + 1,
	},
}

const writeSingleV3Header = `INSERT INTO eth.header_cids (block_number, block_hash, parent_hash, cid, mh_key, td, node_id,
							reward, state_root, uncle_root, tx_root, receipt_root, bloom, timestamp, times_validated, coinbase)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

const writeIPDL = `INSERT INTO public.blocks (key, data) VALUES ($1, $2)`

var _ = Describe("Range Segmenting Helpers", Serial, func() {
	Describe("SegmentRangeByChunkSize", Serial, func() {
		It("segments properly", Serial, func() {
			for _, testCase := range chunkSizeTestCases {
				segments := migration_tools.SegmentRangeByChunkSize(testCase.ChunkSize, testCase.Start, testCase.Stop)
				Expect(uint64(len(segments))).To(Equal(testCase.ExpectedNumOfChunks))
				Expect(segments[0][0]).To(Equal(testCase.Start))
				Expect(segments[0][1]).To(Equal(testCase.Start + testCase.ChunkSize - 1))
				Expect(segments[1][0]).To(Equal(testCase.Start + testCase.ChunkSize))
				Expect(segments[1][1]).To(Equal(testCase.Start + testCase.ChunkSize + testCase.ChunkSize - 1))
				rem := (testCase.Stop - testCase.Start + 1) % testCase.ChunkSize
				if rem == 0 {
					Expect(segments[len(segments)-1][0]).To(Equal(testCase.Stop - testCase.ChunkSize + 1))
				} else {
					Expect(segments[len(segments)-1][0]).To(Equal(testCase.Stop - rem + 1))
				}
				Expect(segments[len(segments)-1][1]).To(Equal(testCase.Stop))
			}
		})
	})
	Describe("DetectAndSegmentRangeByChunkSize", Serial, Label("test"), func() {
		BeforeEach(func() {
			connStr := v3DBConfig.DbConnectionString()
			sqlxDB, err = sqlx.Connect("postgres", connStr)
			Expect(err).ToNot(HaveOccurred())
			tearDownSQLXDB(sqlxDB)
			header := migration_tools.MockHeader

			_, err := sqlxDB.Exec(writeIPDL,
				"mockMhKey",
				[]byte{1, 2, 3, 4})
			Expect(err).ToNot(HaveOccurred())

			_, err = sqlxDB.Exec(writeNodeStr,
				"mockName",
				"mockGenesisBlock",
				1,
				"mockNodeID",
				1337)
			Expect(err).ToNot(HaveOccurred())

			// stop
			_, err = sqlxDB.Exec(writeSingleV3Header,
				header.Number.String(),
				randomHash().String(),
				header.ParentHash.String(),
				"mockCID",
				"mockMhKey",
				header.Difficulty.String(),
				"mockNodeID",
				"1010230213",
				header.Root.String(),
				header.UncleHash.String(),
				header.TxHash.String(),
				header.ReceiptHash.String(),
				header.Bloom.Bytes(),
				header.Time,
				1,
				header.Coinbase.String())
			Expect(err).ToNot(HaveOccurred())

			// start
			_, err = sqlxDB.Exec(writeSingleV3Header,
				"0",
				randomHash().String(),
				header.ParentHash.String(),
				"mockCID",
				"mockMhKey",
				header.Difficulty.String(),
				"mockNodeID",
				"1010230213",
				header.Root.String(),
				header.UncleHash.String(),
				header.TxHash.String(),
				header.ReceiptHash.String(),
				header.Bloom.Bytes(),
				header.Time,
				1,
				header.Coinbase.String())
			Expect(err).ToNot(HaveOccurred())

			// middle
			_, err = sqlxDB.Exec(writeSingleV3Header,
				"1000",
				randomHash().String(),
				header.ParentHash.String(),
				"mockCID",
				"mockMhKey",
				header.Difficulty.String(),
				"mockNodeID",
				"1010230213",
				header.Root.String(),
				header.UncleHash.String(),
				header.TxHash.String(),
				header.ReceiptHash.String(),
				header.Bloom.Bytes(),
				header.Time,
				1,
				header.Coinbase.String())
			Expect(err).ToNot(HaveOccurred())

		})
		AfterEach(func() {
			tearDownSQLXDB(sqlxDB)
		})
		It("segments properly", Serial, func() {
			segments, err := migration_tools.DetectAndSegmentRangeByChunkSize(v3DBConfig, 1000)
			Expect(err).ToNot(HaveOccurred())
			Expect(uint64(len(segments))).To(Equal(((migration_tools.BlockNumber.Uint64() + 1) / 1000) + 1))
			Expect(segments[0][0]).To(Equal(uint64(0)))
			Expect(segments[len(segments)-1][1]).To(Equal(migration_tools.BlockNumber.Uint64()))
		})
	})
})
