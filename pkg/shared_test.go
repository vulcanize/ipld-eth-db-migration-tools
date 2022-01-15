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
	"bytes"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff/indexer/database/sql"
	"github.com/ethereum/go-ethereum/statediff/indexer/database/sql/postgres"
	"github.com/ethereum/go-ethereum/statediff/indexer/interfaces"
	"github.com/ethereum/go-ethereum/statediff/indexer/ipld"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	dshelp "github.com/ipfs/go-ipfs-ds-help"
	"github.com/jmoiron/sqlx"
	"github.com/multiformats/go-multihash"
	. "github.com/onsi/gomega"

	migration_tools "github.com/vulcanize/migration-tools/pkg"
)

var (
	sqlxDB                                                *sqlx.DB
	sqlDB                                                 sql.Database
	err                                                   error
	ind                                                   interfaces.StateDiffIndexer
	migrator                                              migration_tools.Migrator
	tx1, tx2, tx3, tx4, tx5, rct1, rct2, rct3, rct4, rct5 []byte
	mockBlock                                             *types.Block
)

var v3DBConfig = postgres.Config{
	Hostname:     "localhost",
	Port:         5432,
	DatabaseName: "vulcanize_testing_v3",
	Username:     "postgres",
	Password:     "",
}

const writeNodeStr = `INSERT INTO public.nodes (client_name, genesis_block, network_id, node_id, chain_id)
						VALUES ($1, $2, $3, $4, $5)`

func init() {
	rand.Seed(time.Now().UnixNano())

	mockBlock = migration_tools.MockBlock
	txs, rcts := migration_tools.MockBlock.Transactions(), migration_tools.MockReceipts

	buf := new(bytes.Buffer)
	txs.EncodeIndex(0, buf)
	tx1 = make([]byte, buf.Len())
	copy(tx1, buf.Bytes())
	buf.Reset()

	txs.EncodeIndex(1, buf)
	tx2 = make([]byte, buf.Len())
	copy(tx2, buf.Bytes())
	buf.Reset()

	txs.EncodeIndex(2, buf)
	tx3 = make([]byte, buf.Len())
	copy(tx3, buf.Bytes())
	buf.Reset()

	txs.EncodeIndex(3, buf)
	tx4 = make([]byte, buf.Len())
	copy(tx4, buf.Bytes())
	buf.Reset()

	txs.EncodeIndex(4, buf)
	tx5 = make([]byte, buf.Len())
	copy(tx5, buf.Bytes())
	buf.Reset()

	rcts.EncodeIndex(0, buf)
	rct1 = make([]byte, buf.Len())
	copy(rct1, buf.Bytes())
	buf.Reset()

	rcts.EncodeIndex(1, buf)
	rct2 = make([]byte, buf.Len())
	copy(rct2, buf.Bytes())
	buf.Reset()

	rcts.EncodeIndex(2, buf)
	rct3 = make([]byte, buf.Len())
	copy(rct3, buf.Bytes())
	buf.Reset()

	rcts.EncodeIndex(3, buf)
	rct4 = make([]byte, buf.Len())
	copy(rct4, buf.Bytes())
	buf.Reset()

	rcts.EncodeIndex(4, buf)
	rct5 = make([]byte, buf.Len())
	copy(rct5, buf.Bytes())
	buf.Reset()

	receiptTrie := ipld.NewRctTrie()

	receiptTrie.Add(0, rct1)
	receiptTrie.Add(1, rct2)
	receiptTrie.Add(2, rct3)
	receiptTrie.Add(3, rct4)
	receiptTrie.Add(4, rct5)

	rctLeafNodes, keys, _ := receiptTrie.GetLeafNodes()

	rctleafNodeCids := make([]cid.Cid, len(rctLeafNodes))
	orderedRctLeafNodes := make([][]byte, len(rctLeafNodes))
	for i, rln := range rctLeafNodes {
		var idx uint

		r := bytes.NewReader(keys[i].TrieKey)
		rlp.Decode(r, &idx)
		rctleafNodeCids[idx] = rln.Cid()
		orderedRctLeafNodes[idx] = rln.RawData()
	}
}

func randomHash() common.Hash {
	hash := make([]byte, 32)
	rand.Read(hash)
	return common.BytesToHash(hash)
}

func randomAddress() common.Address {
	addr := make([]byte, 20)
	rand.Read(addr)
	return common.BytesToAddress(addr)
}

func randomBytes() []byte {
	by := make([]byte, 256)
	rand.Read(by)
	return by
}

func keccak256ToMhKey(h []byte) string {
	buf, err := multihash.Encode(h, multihash.KECCAK_256)
	if err != nil {
		panic(err)
	}
	return blockstore.BlockPrefix.String() + dshelp.MultihashToDsKey(multihash.Multihash(buf)).String()
}

func tearDownSQLXDB(db *sqlx.DB) {
	tx, err := db.Begin()
	Expect(err).ToNot(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM eth.header_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM eth.uncle_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM eth.transaction_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM eth.receipt_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM eth.state_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM eth.storage_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM eth.state_accounts`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM eth.access_list_elements`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM eth.log_cids`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM blocks`)
	Expect(err).ToNot(HaveOccurred())
	_, err = tx.Exec(`DELETE FROM nodes`)
	Expect(err).ToNot(HaveOccurred())
	err = tx.Commit()
	Expect(err).ToNot(HaveOccurred())
}
