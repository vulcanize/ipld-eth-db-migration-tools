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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/vulcanize/migration-tools/pkg/sql"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff/indexer/database/sql/postgres"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	migration_tools "github.com/vulcanize/migration-tools/pkg"
	"github.com/vulcanize/migration-tools/pkg/eth_headers"
	"github.com/vulcanize/migration-tools/pkg/eth_logs"
	"github.com/vulcanize/migration-tools/pkg/eth_receipts"
	"github.com/vulcanize/migration-tools/pkg/eth_transactions"
	"github.com/vulcanize/migration-tools/pkg/public_blocks"
)

type benchMarkRepairCondition struct {
	NumOfWorkers int
	SegmentSize  uint64
}

var benchMarkConditions = []benchMarkRepairCondition{
	{
		NumOfWorkers: 1,
		SegmentSize:  100,
	},
	{
		NumOfWorkers: 5,
		SegmentSize:  100,
	},
	{
		NumOfWorkers: 10,
		SegmentSize:  100,
	},
	{
		NumOfWorkers: 50,
		SegmentSize:  100,
	},

	{
		NumOfWorkers: 1,
		SegmentSize:  1000,
	},
	{
		NumOfWorkers: 5,
		SegmentSize:  1000,
	},
	{
		NumOfWorkers: 10,
		SegmentSize:  1000,
	},
	{
		NumOfWorkers: 50,
		SegmentSize:  1000,
	},

	{
		NumOfWorkers: 1,
		SegmentSize:  10000,
	},
	{
		NumOfWorkers: 5,
		SegmentSize:  10000,
	},
	{
		NumOfWorkers: 10,
		SegmentSize:  10000,
	},
}

var v3BenchmarkDBConfig = postgres.Config{
	Hostname:     "localhost",
	Port:         5432,
	DatabaseName: "vulcanize_v3_benchmark",
	Username:     "postgres",
	Password:     "",
}

var _ = Describe("Benchmark log repair", Serial, func() {
	It("benchmarks the repair of log data", Serial, Label("measurement"), func() {
		for i, con := range benchMarkConditions {
			fmt.Printf("----------setting up benchmark %d (workers: %d, segment size: %d, blocks: %d, txs: %d, rcts: %d, logs: %d, logs missing IPLDs: %d, IPLDs: %d)----------\r\n",
				i, con.NumOfWorkers, con.SegmentSize, numHeaders, numTransactions, numReceipts, numLogs, numLogsMissingIPLDs, numIPLDs)
			start := time.Now()
			setupBenchMarkDB()
			fmt.Printf("----------benchmark %d setup finished (total time: %s)----------\r\n", i, time.Now().Sub(start).String())
			blockRanges, err := migration_tools.DetectAndSegmentRangeByChunkSize(v3BenchmarkDBConfig, con.SegmentSize)
			Expect(err).ToNot(HaveOccurred())
			conf := &migration_tools.Config{
				ReadDB:          v3BenchmarkDBConfig,
				WriteDB:         v3BenchmarkDBConfig,
				WorkersPerTable: con.NumOfWorkers,
			}
			migrator, err = migration_tools.NewMigrator(context.Background(), conf)
			Expect(err).ToNot(HaveOccurred())
			wg := new(sync.WaitGroup)

			fmt.Printf("----------running benchmark %d----------\r\n", i)
			start = time.Now()
			rangeChan := make(chan [2]uint64)
			readGapsChan, writeGapsChan, doneChan, quitChan, errChan := migrator.Migrate(wg, migration_tools.EthLogsRepair, rangeChan)

			wg.Add(1)
			go func() {
				defer wg.Done()
				for _, blockRange := range blockRanges {
					select {
					case <-doneChan:
						return
					default:
						rangeChan <- blockRange
					}
				}
				close(quitChan)
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case readGap := <-readGapsChan:
						fmt.Printf("Migrator %s table read gap: %v\r\n", migration_tools.EthLogsRepair, readGap)
					case writeGap := <-writeGapsChan:
						fmt.Printf("Migrator %s table write gap: %v\r\n", migration_tools.EthLogsRepair, writeGap)
					case err := <-errChan:
						fmt.Printf("Migrator %s table error: %v\r\n", migration_tools.EthLogsRepair, err)
					case <-doneChan:
						return
					}
				}
			}()
			wg.Wait()
			fmt.Printf("----------finished benchmark %d (run time: %s, workers: %d, segment size: %d, blocks: %d, txs: %d, rcts: %d, logs: %d, logs missing IPLDs: %d, IPLDs: %d)----------\r\n",
				i, time.Now().Sub(start).String(), con.NumOfWorkers, con.SegmentSize, numHeaders, numTransactions, numReceipts, numLogs, numLogsMissingIPLDs, numIPLDs)
			time.Sleep(30 * time.Second)
			_, err = sqlxDB.Exec("ALTER TABLE eth.log_cids ADD CONSTRAINT log_cids_leaf_mh_key_fkey " +
				"FOREIGN KEY (leaf_mh_key) REFERENCES public.blocks (key) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED")
			Expect(err).ToNot(HaveOccurred())
			clearBenchMarkDB()
		}
	})
})

var (
	numBatches  = 200
	batchSize   = 500
	gapStepSize = 2

	numHeaders        = numBatches * batchSize
	numHeaderPerBatch = batchSize

	numTransactionsPerHeader = 3
	numTransactions          = numHeaders * numTransactionsPerHeader
	numTransactionPerBatch   = numHeaderPerBatch * numTransactionsPerHeader

	numReceipts         = numTransactions
	numReceiptsPerBatch = numTransactionPerBatch

	numReceiptsPerHeader = numTransactionsPerHeader
	numLogs              = numReceipts * numReceiptsPerHeader
	numLogsPerBatch      = numReceiptsPerBatch * numReceiptsPerHeader

	numLogsMissingIPLDs         = (numHeaders / gapStepSize) * numReceiptsPerHeader // every other header is missing an IPLD for one of three logs in each receipt
	numLogsMissingIPLDsPerBatch = (numHeaderPerBatch / gapStepSize) * numReceiptsPerHeader

	numIPLDs         = numHeaders + numTransactions + numReceipts + (numLogs - numLogsMissingIPLDs)
	numIPLDsPerBatch = numHeaderPerBatch + numTransactionPerBatch + numReceiptsPerBatch + (numLogsPerBatch - numLogsMissingIPLDsPerBatch)
)

const benchMarkTarPath = "bulk_benchmark_data.tar"

var tarExistsFlag bool

func setupBenchMarkDB() {
	if _, err := os.Stat(benchMarkTarPath); errors.Is(err, os.ErrNotExist) {
		fmt.Print("----------setting up benchmark database by generating and writing data----------\r\n")
		tarExistsFlag = false
		setupBenchmarkDBFromGeneratedData()
	} else {
		fmt.Printf("----------setting up benchmark database from .tar %s----------\r\n", benchMarkTarPath)
		tarExistsFlag = true
		setupBenchmarkDBFromTar()
	}
}

func clearBenchMarkDB() {
	if tarExistsFlag {
		teardownBenchmarkDB()
	} else {
		dropBenchmarkData()
	}
}

func setupBenchmarkDBFromTar() {
	cmd := exec.Command("createdb", "vulcanize_v3_benchmark")
	err = cmd.Run()
	Expect(err).ToNot(HaveOccurred())

	connStr := v3BenchmarkDBConfig.DbConnectionString()
	sqlxDB, err = sqlx.Connect("postgres", connStr)
	Expect(err).ToNot(HaveOccurred())
	cmd = exec.Command("pg_restore", "bulk_benchmark_data.tar", "--dbname=vulcanize_v3_benchmark")
	err = cmd.Run()
	Expect(err).ToNot(HaveOccurred())
}

func teardownBenchmarkDB() {
	migrator.Close()
	sqlxDB.Close()
	cmd := exec.Command("dropdb", "vulcanize_v3_benchmark")
	err = cmd.Run()
	Expect(err).ToNot(HaveOccurred())
}

func setupBenchmarkDBFromGeneratedData() {
	connStr := v3BenchmarkDBConfig.DbConnectionString()
	sqlxDB, err = sqlx.Connect("postgres", connStr)
	Expect(err).ToNot(HaveOccurred())
	writer := sql.NewWriter(sqlxDB)

	_, err := sqlxDB.Exec("ALTER TABLE eth.log_cids DROP CONSTRAINT IF EXISTS log_cids_leaf_mh_key_fkey")
	Expect(err).ToNot(HaveOccurred())

	_, err = sqlxDB.Exec(writeNodeStr,
		"mockName",
		"mockGenesisBlock",
		1,
		baseHeader.NodeID,
		1337)
	Expect(err).ToNot(HaveOccurred())

	var lastHeader *eth_headers.HeaderModelV3
	for k := 0; k < numBatches; k++ {
		headers := make([]eth_headers.HeaderModelV3, numHeaderPerBatch)
		transactions := make([]eth_transactions.TransactionModelV3, 0, numTransactionPerBatch)
		receipts := make([]eth_receipts.ReceiptModelV3, 0, numReceiptsPerBatch)
		logs := make([]eth_logs.LogModelV3, 0, numLogsPerBatch) // This is the limiting factor in batching: most records + 10 fields; 65535 is the max args Postgres supports per write conn
		iplds := make([]public_blocks.IPLDModel, 0, numIPLDsPerBatch)
		for i := 0; i < numHeaderPerBatch; i++ {
			header, headerIPLD := generateHeader(lastHeader)
			lastHeader = &header
			headers[i] = header
			iplds = append(iplds, headerIPLD)

			txs, txIPLDs := generateTxs(header.BlockHash)
			transactions = append(transactions, txs...)
			iplds = append(iplds, txIPLDs...)

			rcts, rctIPLDs := generateRcts(txs)
			receipts = append(receipts, rcts...)
			iplds = append(iplds, rctIPLDs...)

			ls, logIPLDs := generateLogs(i, rcts)
			logs = append(logs, ls...)
			iplds = append(iplds, logIPLDs...)
		}

		err = writer.Write(sql.PgWriteIPLDsStr, iplds)
		Expect(err).ToNot(HaveOccurred())

		err = writer.Write(sql.PgWriteEthHeadersStr, headers)
		Expect(err).ToNot(HaveOccurred())

		err = writer.Write(sql.PgWriteEthTransactionsStr, transactions)
		Expect(err).ToNot(HaveOccurred())

		err = writer.Write(sql.PgWriteEthReceiptsStr, receipts)
		Expect(err).ToNot(HaveOccurred())

		err = writer.Write(sql.PgWriteEthLogsStr, logs)
		Expect(err).ToNot(HaveOccurred())
	}
}

func dropBenchmarkData() {
	tearDownSQLXDB(sqlxDB)
	migrator.Close()
	sqlxDB.Close()
}

var baseHeader = eth_headers.HeaderModelV3{
	BlockNumber:     "0",
	BlockHash:       randomHash().String(),
	ParentHash:      randomHash().String(),
	CID:             "mockCID",
	MhKey:           "",
	TotalDifficulty: "10000000",
	NodeID:          "mockNodeID",
	Reward:          "10000000",
	StateRoot:       randomHash().String(),
	UncleRoot:       randomHash().String(),
	TxRoot:          randomHash().String(),
	RctRoot:         randomHash().String(),
	Bloom:           []byte{1, 2, 3, 4, 5},
	Timestamp:       uint64(time.Now().UnixNano()),
	TimesValidated:  1,
	Coinbase:        randomAddress().String(),
}

func generateHeader(parentHeader *eth_headers.HeaderModelV3) (eth_headers.HeaderModelV3, public_blocks.IPLDModel) {
	if parentHeader == nil {
		baseHeader.MhKey = keccak256ToMhKey(common.Hex2Bytes(baseHeader.BlockHash))
		return baseHeader, public_blocks.IPLDModel{
			Key:  baseHeader.MhKey,
			Data: randomBytes(),
		}
	}
	header := baseHeader
	height, _ := strconv.ParseInt(parentHeader.BlockNumber, 10, 64)
	header.BlockNumber = strconv.Itoa(int(height + 1))
	header.ParentHash = parentHeader.BlockHash
	header.BlockHash = randomHash().String()
	header.MhKey = keccak256ToMhKey(common.Hex2Bytes(header.BlockHash))
	return header, public_blocks.IPLDModel{
		Key:  header.MhKey,
		Data: randomBytes(),
	}
}

var baseTx = eth_transactions.TransactionModelV3{
	HeaderID: "",
	Index:    0,
	TxHash:   "",
	CID:      "mockCID",
	MhKey:    "",
	Dst:      randomAddress().String(),
	Src:      randomAddress().String(),
	Data:     randomBytes(),
	Type:     0,
	Value:    "1000",
}

func generateTxs(headerHash string) ([]eth_transactions.TransactionModelV3, []public_blocks.IPLDModel) {
	txs := make([]eth_transactions.TransactionModelV3, 3)
	iplds := make([]public_blocks.IPLDModel, 3)
	for i := 0; i < 3; i++ {
		tx := baseTx
		tx.TxHash = randomHash().String()
		tx.HeaderID = headerHash
		tx.Index = int64(i)
		tx.MhKey = keccak256ToMhKey(common.Hex2Bytes(tx.TxHash))
		txs[i] = tx
		iplds[i] = public_blocks.IPLDModel{
			Key:  tx.MhKey,
			Data: randomBytes(),
		}
	}
	return txs, iplds
}

var baseRct = eth_receipts.ReceiptModelV3{
	TxID:         "",
	LeafCID:      "mockCID",
	LeafMhKey:    "",
	PostStatus:   0,
	PostState:    randomHash().String(),
	Contract:     randomAddress().String(),
	ContractHash: randomHash().String(),
	LogRoot:      randomHash().String(),
}

func generateRcts(txs []eth_transactions.TransactionModelV3) ([]eth_receipts.ReceiptModelV3, []public_blocks.IPLDModel) {
	rcts := make([]eth_receipts.ReceiptModelV3, len(txs))
	iplds := make([]public_blocks.IPLDModel, len(txs))
	for i, tx := range txs {
		rct := baseRct
		rct.TxID = tx.TxHash
		rct.LeafMhKey = keccak256ToMhKey(randomHash().Bytes())
		rcts[i] = rct
		iplds[i] = public_blocks.IPLDModel{
			Key:  rct.LeafMhKey,
			Data: randomBytes(),
		}
	}
	return rcts, iplds
}

var baseLog = eth_logs.LogModelV3{
	ReceiptID: "",
	LeafCID:   "mockCID",
	LeafMhKey: "",
	Address:   "",
	Index:     0,
	Data:      []byte{},
	Topic0:    "",
	Topic1:    "",
	Topic2:    "",
	Topic3:    "",
}

func generateLogs(i int, rcts []eth_receipts.ReceiptModelV3) ([]eth_logs.LogModelV3, []public_blocks.IPLDModel) {
	logs := make([]eth_logs.LogModelV3, len(rcts)*3)
	iplds := make([]public_blocks.IPLDModel, 0, len(rcts)*2)
	index := 0
	for _, rct := range rcts {
		for j := 0; j < 3; j++ {
			gethLog := &types.Log{
				Address: randomAddress(),
				Topics: []common.Hash{
					randomHash(),
				},
				Data: randomBytes(),
			}
			logRLP, err := rlp.EncodeToBytes(gethLog)
			Expect(err).ToNot(HaveOccurred())
			log := baseLog
			log.ReceiptID = rct.TxID
			log.Index = int64(index)
			log.Address = gethLog.Address.String()
			log.Topic0 = gethLog.Topics[0].String()
			log.Data = gethLog.Data
			log.LeafMhKey = keccak256ToMhKey(crypto.Keccak256(logRLP))
			logs[index] = log
			index++
			if i%2 == 0 && j == 2 { // every other block we exclude one of the three logs for each receipt
				continue
			}
			iplds = append(iplds, public_blocks.IPLDModel{
				Key:  log.LeafMhKey,
				Data: logRLP,
			})
		}
	}
	return logs, iplds
}
