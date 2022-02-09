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

package csv

// WriteCSVStr provides explicit typing for write csv statements
type WriteCSVStr string

const (
	CSVWriteIPLDsStr WriteCSVStr = "%s, \\x%x\n"

	CSVWriteEthUnclesStr WriteCSVStr = "%s, %s, %s, %s, %s, %s\n"

	CSVWriteEthTransactionsStr WriteCSVStr = "%s, %s, %s, %s, %s, %d, %s, \\x%x, %d, %s\n"

	CSVWriteEthStorageStr WriteCSVStr = "%s, \\x%x, %s, %s, \\x%x, %d, %t, %s\n"

	CSVWriteEthStateStr WriteCSVStr = "%s, %s, %s, \\x%x, %d, %t, %s\n"

	CSVWriteNodesStr WriteCSVStr = "%s, %s, %s, %s, %d\n"

	CSVWriteEthReceiptsStr WriteCSVStr = "%s, %s, %s, %s, %s, %s, %d, %s\n"

	CSVWriteEthLogsStr WriteCSVStr = "%s, %s, %s, %s, %d, %s, %s, %s, %s, \\x%x\n"

	CSVWriteEthHeadersStr WriteCSVStr = "%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, \\x%x, %d, %s, %d, %s\n"

	CSVWriteEthAccountsStr WriteCSVStr = "%s, \\x%x, %s, %d, \\x%x, '%s\n"

	CSVWriteAccessListElementsStr WriteCSVStr = "%s, %d, %s, %s\n"
)
