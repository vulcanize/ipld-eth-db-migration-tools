# migration-tools
Tools for migrating a [v2 ipld-eth-db](https://github.com/vulcanize/ipld-eth-db/releases/tag/v0.2.1) schema DB to a
[v3 ipld-eth-db](https://github.com/vulcanize/ipld-eth-db/releases/tag/v0.3.2) schema DB.

Reads data from a v2 database, transforming them into v3 DB models, and writing them to a v3 database.

Can be configured to work over a subset of the tables, and over specific block ranges.

While processing headers, state, or state accounts it checks for gaps in the data and writes these out to a file. It only checks for
gaps for these tables, as every other table can (in theory) be empty for a given range. Whereas even a block without any transactions will
produce a header and state trie and state account updates (for the miner's reward).

## Usage
`./migration-tools migrate --config={path_to_toml_config_file}`

Example TOML config:

```toml
[migrator]
    ranges = [
        [0, 1000]
    ]
    start = 0 # $MIGRATION_START
    stop = 1000 # $MIGRATION_STOP
    tableNames = [ # $MIGRATION_TABLE_NAMES
        "headers",
        "transactions",
        "storage"
    ]
    workersPerTable = 1 # $MIGRATION_WORKERS_PER_TABLE
    autoRange = false # $MIGRATION_AUTO_RANGE
    segmentSize = 10000 # $MIGRATION_AUTO_RANGE_SEGMENT_SIZE
    segmentOffset = 0 # $TRANSFER_SEGMENT_OFFSET
    maxPage = 0 # $TRANSFER_MAX_PAGE

[log]
    file = "path/to/log/file" # $LOGRUS_FILE
    level = "info" # $LOGRUS_LEVEL
    readGapsDir = "path/to/read/gaps/dir" # $LOG_READ_GAPS_DIR
    writeGapsDir = "path/to/write/gaps/dir" # $LOG_WRITE_GAPS_DIR

[v2]
    databaseName = "vulcanize_public_v2" # $OLD_DATABASE_NAME
    databaseHostName = "localhost" # $OLD_DATABASE_HOSTNAME
    databasePort = "5432" # $OLD_DATABASE_PORT
    databaseUser = "postgres" # $OLD_DATABASE_USER
    databasePassword = "" # $OLD_DATABASE_PASSWORD
    databaseMaxIdleConns = 50 # $OLD_DATABASE_MAX_IDLE_CONNECTIONS
    databaseMaxOpenConns = 100 # $OLD_DATABASE_MAX_OPEN_CONNECTIONS
    databaseMaxConnLifetime = 0 # $OLD_DATABASE_MAX_CONN_LIFETIME

[v3]
    databaseName = "vulcanize_public_v2" # $NEW_DATABASE_NAME
    databaseHostName = "localhost" # $NEW_DATABASE_HOSTNAME
    databasePort = "5432" # $NEW_DATABASE_PORT
    databaseUser = "postgres" # $NEW_DATABASE_USER
    databasePassword = "" # $NEW_DATABASE_PASSWORD
    databaseMaxIdleConns = 50 # $NEW_DATABASE_MAX_IDLE_CONNECTIONS
    databaseMaxOpenConns = 100 # $NEW_DATABASE_MAX_OPEN_CONNECTIONS
    databaseMaxConnLifetime = 0 # $NEW_DATABASE_MAX_CONN_LIFETIME
```

The command can be configured through the linked TOML file as shown above, through ENV variable bindings, or through CLI flags.
The names of the ENV variables are listed in the comments next to their corresponding TOML binding. For a list of the CLI flags, run: 

`./migration-tools migrate --help`

The precedence of configuration is ENV > CLI > TOML. For example, if a parameter is configured by both ENV variable
and in the TOML file, the ENV variable value is the one used.

The tableNames options are:

public.nodes  
eth.header_cids  
eth.uncles_cids  
eth.transaction_cids  
eth.access_list_elements  
eth.receipt_cids  
eth.log_cids  
eth.state_cids  
eth.state_accounts  
eth.storage_cids  
eth.log_cids.repair

public.blocks should be migrated using pg_dump and COPY FROM using a foreign table to handle unique constraint conflicts on INSERT.



