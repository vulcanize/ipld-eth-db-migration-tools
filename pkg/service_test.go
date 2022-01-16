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
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff/indexer/database/sql"
	"github.com/ethereum/go-ethereum/statediff/indexer/database/sql/postgres"
	"github.com/ethereum/go-ethereum/statediff/indexer/interfaces"
	"github.com/ethereum/go-ethereum/statediff/indexer/node"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	migration_tools "github.com/vulcanize/migration-tools/pkg"
)

func writeV2SQL(sqlDB sql.Database) {
	ind, err = sql.NewStateDiffIndexer(context.Background(), migration_tools.TestConfig, sqlDB)
	Expect(err).ToNot(HaveOccurred())
	var tx interfaces.Batch
	tx, err = ind.PushBlock(
		mockBlock,
		migration_tools.MockReceipts,
		migration_tools.MockBlock.Difficulty())
	Expect(err).ToNot(HaveOccurred())
	defer func() {
		err := tx.Submit(err)
		Expect(err).ToNot(HaveOccurred())
	}()
	for _, n := range migration_tools.StateDiffs {
		err = ind.PushStateNode(tx, n, mockBlock.Hash().String())
		Expect(err).ToNot(HaveOccurred())
	}
	Expect(tx.(*sql.BatchTx).BlockNumber).To(Equal(migration_tools.BlockNumber.Uint64()))
}

func tearDown() {
	tearDownDatabase(sqlDB)
	err = sqlDB.Close()
	Expect(err).ToNot(HaveOccurred())
}

var _ = Describe("Migration Service", Serial, func() {
	Describe("LogTrie repair", Serial, func() {
		BeforeEach(func() {
			driver, err := postgres.NewSQLXDriver(context.Background(), v3DBConfig, node.Info{})
			Expect(err).ToNot(HaveOccurred())
			sqlDB, err = postgres.NewPostgresDB(driver), nil
			Expect(err).ToNot(HaveOccurred())
			prepDatabase(sqlDB)
			writeV2SQL(sqlDB)
			_, err = sqlDB.Exec(context.Background(), "ALTER TABLE eth.log_cids DROP CONSTRAINT IF EXISTS log_cids_leaf_mh_key_fkey")
			Expect(err).ToNot(HaveOccurred())
			conf := &migration_tools.Config{
				ReadDB:          v3DBConfig,
				WriteDB:         v3DBConfig,
				WorkersPerTable: 1,
			}
			migrator, err = migration_tools.NewMigrator(context.Background(), conf)
			Expect(err).To(BeNil())
		})
		AfterEach(func() {
			_, err := sqlDB.Exec(context.Background(), "ALTER TABLE eth.log_cids ADD CONSTRAINT log_cids_leaf_mh_key_fkey "+
				"FOREIGN KEY (leaf_mh_key) REFERENCES public.blocks (key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED")
			Expect(err).ToNot(HaveOccurred())
			tearDown()
			migrator.Close()
		})

		It("throws no errors on empty range", Serial, Label("test"), func() {
			wg := new(sync.WaitGroup)
			blockRangeChan := make(chan [2]uint64)
			readGaps, writeGaps, doneChan, quitChan, errChan := migrator.Migrate(wg, migration_tools.EthLogsRepair, blockRangeChan)
			rng := [2]uint64{migration_tools.BlockNumber.Uint64(), migration_tools.BlockNumber.Uint64() + 1000}
			blockRangeChan <- rng

			wg.Add(1)
			var readGap [2]uint64
			var writeGap [2]uint64
			go func() {
				defer wg.Done()
				for {
					select {
					case r := <-readGaps:
						readGap = r
					case w := <-writeGaps:
						writeGap = w
					case migrateErr := <-errChan:
						err = migrateErr
					case <-doneChan:
						return
					}
				}
			}()
			close(quitChan)
			wg.Wait()
			Expect(readGap).To(Equal([2]uint64{0, 0}))
			Expect(writeGap).To(Equal([2]uint64{0, 0}))
			Expect(err).ToNot(HaveOccurred())
		})

		It("throws no errors on empty range", Serial, Label("test"), func() {
			wg := new(sync.WaitGroup)
			blockRangeChan := make(chan [2]uint64)
			readGaps, writeGaps, doneChan, quitChan, errChan := migrator.Migrate(wg, migration_tools.EthLogsRepair, blockRangeChan)
			rng := [2]uint64{0, 1000}
			blockRangeChan <- rng

			wg.Add(1)
			var readGap [2]uint64
			var writeGap [2]uint64
			go func() {
				defer wg.Done()
				for {
					select {
					case r := <-readGaps:
						readGap = r
					case w := <-writeGaps:
						writeGap = w
					case migrateErr := <-errChan:
						err = migrateErr
					case <-doneChan:
						return
					}
				}
			}()
			close(quitChan)
			wg.Wait()
			Expect(readGap).To(Equal([2]uint64{0, 0}))
			Expect(writeGap).To(Equal([2]uint64{0, 0}))
			Expect(err).ToNot(HaveOccurred())
		})

		It("repairs missing log IPLDs", Serial, Label("test"), func() {
			// explicitly check the IPLDs are correct
			type logIPLD struct {
				Index   int    `db:"index"`
				Address string `db:"address"`
				Data    []byte `db:"data"`
				Topic0  string `db:"topic0"`
				Topic1  string `db:"topic1"`
			}
			rctPgStr := `SELECT receipt_cids.leaf_cid FROM eth.receipt_cids, eth.transaction_cids, eth.header_cids
				WHERE receipt_cids.tx_id = transaction_cids.tx_hash
				AND transaction_cids.header_id = header_cids.block_hash
				AND header_cids.block_number = $1
				ORDER BY transaction_cids.index`
			logPgStr := `SELECT log_cids.index, log_cids.address, log_cids.topic0, log_cids.topic1, data FROM eth.log_cids
    				INNER JOIN eth.receipt_cids ON (log_cids.rct_id = receipt_cids.tx_id)
					INNER JOIN public.blocks ON (log_cids.leaf_mh_key = blocks.key)
					WHERE receipt_cids.leaf_cid = $1 ORDER BY eth.log_cids.index ASC`
			rcts := make([]string, 0)
			err = sqlDB.Select(context.Background(), &rcts, rctPgStr, migration_tools.BlockNumber.Uint64())
			Expect(err).To(BeNil())
			Expect(len(rcts)).To(Equal(5))

			for i := range rcts {
				results := make([]logIPLD, 0)
				err = sqlDB.Select(context.Background(), &results, logPgStr, rcts[i])
				Expect(err).To(BeNil())

				expectedLogs := migration_tools.MockReceipts[i].Logs
				Expect(len(results)).To(Equal(len(expectedLogs)))

				var nodeElements []interface{}
				for idx, r := range results {
					// Decode the log leaf node.
					err = rlp.DecodeBytes(r.Data, &nodeElements)
					Expect(err).To(BeNil())
					if len(nodeElements) == 2 {
						log := new(types.Log)
						rlp.DecodeBytes(nodeElements[1].([]byte), log)
						logRaw, err := rlp.EncodeToBytes(expectedLogs[idx])
						Expect(err).To(BeNil())
						Expect(logRaw).To(Equal(nodeElements[1].([]byte)))
					} else {
						log := new(types.Log)
						rlp.DecodeBytes(r.Data, log)
						logRaw, err := rlp.EncodeToBytes(expectedLogs[idx])
						Expect(err).To(BeNil())
						Expect(logRaw).To(Equal(r.Data))
					}
				}
			}

			// remove the short log IPLDs
			_, err = sqlDB.Exec(context.Background(), `DELETE FROM blocks WHERE key = $1`, migration_tools.ShotLog1MhKey)
			Expect(err).To(BeNil())
			_, err = sqlDB.Exec(context.Background(), `DELETE FROM blocks WHERE key = $1`, migration_tools.ShotLog2MhKey)
			Expect(err).To(BeNil())

			// explicitly check the IPLDs are gone
			rcts2 := make([]string, 0)
			err = sqlDB.Select(context.Background(), &rcts2, rctPgStr, migration_tools.BlockNumber.Uint64())
			Expect(err).To(BeNil())
			Expect(len(rcts2)).To(Equal(5))

			for i := range rcts2 {
				results := make([]logIPLD, 0)
				err = sqlDB.Select(context.Background(), &results, logPgStr, rcts2[i])
				Expect(err).To(BeNil())

				expectedLogs := migration_tools.MockReceipts[i].Logs
				if i == 3 || i == 1 {
					Expect(len(results)).To(Equal(len(expectedLogs) - 1))
				} else {
					Expect(len(results)).To(Equal(len(expectedLogs)))
				}

				var nodeElements []interface{}
				for idx, r := range results {
					// Decode the log leaf node.
					err = rlp.DecodeBytes(r.Data, &nodeElements)
					Expect(err).To(BeNil())
					if len(nodeElements) == 2 {
						logRaw, err := rlp.EncodeToBytes(expectedLogs[idx])
						Expect(err).To(BeNil())
						Expect(logRaw).To(Equal(nodeElements[1].([]byte)))
					} else {
						logRaw, err := rlp.EncodeToBytes(expectedLogs[idx])
						Expect(err).To(BeNil())
						Expect(logRaw).To(Equal(r.Data))
					}
				}
			}

			wg := new(sync.WaitGroup)
			blockRangeChan := make(chan [2]uint64)
			readGaps, writeGaps, doneChan, quitChan, errChan := migrator.Migrate(wg, migration_tools.EthLogsRepair, blockRangeChan)
			rng := [2]uint64{migration_tools.BlockNumber.Uint64(), migration_tools.BlockNumber.Uint64() + 1000}
			blockRangeChan <- rng

			wg.Add(1)
			var readGap [2]uint64
			var writeGap [2]uint64
			go func() {
				defer wg.Done()
				for {
					select {
					case r := <-readGaps:
						readGap = r
					case w := <-writeGaps:
						writeGap = w
					case migrateErr := <-errChan:
						err = migrateErr
					case <-doneChan:
						return
					}
				}
			}()

			close(quitChan)
			wg.Wait()
			Expect(readGap).To(Equal([2]uint64{0, 0}))
			Expect(writeGap).To(Equal([2]uint64{0, 0}))
			Expect(err).ToNot(HaveOccurred())

			// explicitly check the IPLDs are back
			rcts3 := make([]string, 0)
			err = sqlDB.Select(context.Background(), &rcts3, rctPgStr, migration_tools.BlockNumber.Uint64())
			Expect(err).To(BeNil())
			Expect(len(rcts3)).To(Equal(5))

			for i := range rcts {
				results := make([]logIPLD, 0)
				err = sqlDB.Select(context.Background(), &results, logPgStr, rcts[i])
				Expect(err).To(BeNil())
				expectedLogs := migration_tools.MockReceipts[i].Logs
				Expect(len(results)).To(Equal(len(expectedLogs)))

				var nodeElements []interface{}
				for idx, r := range results {
					// Decode the log leaf node.
					err = rlp.DecodeBytes(r.Data, &nodeElements)
					Expect(err).To(BeNil())
					if len(nodeElements) == 2 {
						logRaw, err := rlp.EncodeToBytes(expectedLogs[idx])
						Expect(err).To(BeNil())
						Expect(logRaw).To(Equal(nodeElements[1].([]byte)))
					} else {
						logRaw, err := rlp.EncodeToBytes(expectedLogs[idx])
						Expect(err).To(BeNil())
						Expect(logRaw).To(Equal(r.Data))
					}
				}
			}
		})
	})
})

func prepDatabase(db sql.Database) {
	tx, err := db.Begin(context.Background())
	Expect(err).ToNot(HaveOccurred())

	_, err = tx.Exec(context.Background(), `DELETE FROM eth.header_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.uncle_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.transaction_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.receipt_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.state_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.storage_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.state_accounts`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.access_list_elements`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.log_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM blocks`)
	Expect(err).ToNot(HaveOccurred())
	err = tx.Commit(context.Background())
	Expect(err).ToNot(HaveOccurred())
}

func tearDownDatabase(db sql.Database) {
	tx, err := db.Begin(context.Background())
	Expect(err).ToNot(HaveOccurred())

	_, err = tx.Exec(context.Background(), `DELETE FROM eth.header_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.uncle_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.transaction_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.receipt_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.state_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.storage_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.state_accounts`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.access_list_elements`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM eth.log_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM blocks`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(context.Background(), `DELETE FROM nodes`)
	Expect(err).ToNot(HaveOccurred())
	err = tx.Commit(context.Background())
	Expect(err).ToNot(HaveOccurred())
}
