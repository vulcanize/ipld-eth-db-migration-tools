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

package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	migration_tools "github.com/vulcanize/migration-tools/pkg"
)

var (
	cfgFile        string
	subCommand     string
	logWithCommand log.Entry
)

var rootCmd = &cobra.Command{
	Use:              "migration-tools",
	PersistentPreRun: initFuncs,
}

func Execute() {
	log.Info("----- Starting migrator -----")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initFuncs(cmd *cobra.Command, args []string) {
	viper.BindEnv(migration_tools.TOML_LOGRUS_FILE, migration_tools.LOGRUS_FILE)
	logfile := viper.GetString(migration_tools.TOML_LOGRUS_FILE)
	if logfile != "" {
		file, err := os.OpenFile(logfile,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.Infof("Directing output to %s", logfile)
			log.SetOutput(file)
		} else {
			log.SetOutput(os.Stdout)
			log.Info("Failed to log to file, using default stdout")
		}
	} else {
		log.SetOutput(os.Stdout)
	}
	if err := logLevel(); err != nil {
		log.Fatal("Could not set log level: ", err)
	}
}

func logLevel() error {
	viper.BindEnv(migration_tools.TOML_LOGRUS_LEVEL, migration_tools.LOGRUS_LEVEL)
	lvl, err := log.ParseLevel(viper.GetString(migration_tools.TOML_LOGRUS_LEVEL))
	if err != nil {
		return err
	}
	log.SetLevel(lvl)
	if lvl > log.InfoLevel {
		log.SetReportCaller(true)
	}
	log.Info("Log level set to ", lvl.String())
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file location")

	// log flags
	rootCmd.PersistentFlags().String(migration_tools.CLI_LOG_READ_GAPS_DIR, "./readGaps/", "directory to write out read gaps to")
	rootCmd.PersistentFlags().String(migration_tools.CLI_LOG_WRITE_GAPS_DIR, "./writeGaps/", "directory to write out write gaps to")
	rootCmd.PersistentFlags().String(migration_tools.CLI_LOG_TRANSFER_GAPS_DIR, "./transferGaps/", "directory to write out transfer gaps to")
	rootCmd.PersistentFlags().String(migration_tools.CLI_LOGRUS_LEVEL, log.InfoLevel.String(), "log level (trace, debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().String(migration_tools.CLI_LOGRUS_FILE, "", "file path for logging")

	// old db flags
	rootCmd.PersistentFlags().String(migration_tools.CLI_OLD_DATABASE_NAME, "vulcanize_old", "name for the old database")
	rootCmd.PersistentFlags().String(migration_tools.CLI_OLD_DATABASE_HOSTNAME, "localhost", "hostname for the old database")
	rootCmd.PersistentFlags().String(migration_tools.CLI_OLD_DATABASE_PORT, "5432", "port for the old database")
	rootCmd.PersistentFlags().String(migration_tools.CLI_OLD_DATABASE_USER, "postgres", "username to use with the old database")
	rootCmd.PersistentFlags().String(migration_tools.CLI_OLD_DATABASE_PASSWORD, "", "password to use for the old database")
	rootCmd.PersistentFlags().Int(migration_tools.CLI_OLD_DATABASE_MAX_IDLE_CONNECTIONS, 0, "max idle connections for the old database")
	rootCmd.PersistentFlags().Int(migration_tools.CLI_OLD_DATABASE_MAX_OPEN_CONNECTIONS, 0, "max open connections for the old database")
	rootCmd.PersistentFlags().Duration(migration_tools.CLI_OLD_DATABASE_MAX_CONN_LIFETIME, 0, "max connection lifetime for the old database")

	// new db flags
	rootCmd.PersistentFlags().String(migration_tools.CLI_NEW_DATABASE_NAME, "vulcanize_v3", "name for the new database")
	rootCmd.PersistentFlags().String(migration_tools.CLI_NEW_DATABASE_HOSTNAME, "localhost", "hostname for the new database")
	rootCmd.PersistentFlags().String(migration_tools.CLI_NEW_DATABASE_PORT, "5432", "port for the new database")
	rootCmd.PersistentFlags().String(migration_tools.CLI_NEW_DATABASE_USER, "postgres", "username to use with the new database")
	rootCmd.PersistentFlags().String(migration_tools.CLI_NEW_DATABASE_PASSWORD, "", "password to use for the new database")
	rootCmd.PersistentFlags().Int(migration_tools.CLI_NEW_DATABASE_MAX_IDLE_CONNECTIONS, 0, "max idle connections for the new database")
	rootCmd.PersistentFlags().Int(migration_tools.CLI_NEW_DATABASE_MAX_OPEN_CONNECTIONS, 0, "max open connections for the new database")
	rootCmd.PersistentFlags().Duration(migration_tools.CLI_NEW_DATABASE_MAX_CONN_LIFETIME, 0, "max connection lifetime for the new database")

	// log TOML bindings
	viper.BindPFlag(migration_tools.TOML_LOG_READ_GAPS_DIR, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_LOG_READ_GAPS_DIR))
	viper.BindPFlag(migration_tools.TOML_LOG_WRITE_GAPS_DIR, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_LOG_WRITE_GAPS_DIR))
	viper.BindPFlag(migration_tools.TOML_LOGRUS_LEVEL, rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag(migration_tools.TOML_LOGRUS_FILE, rootCmd.PersistentFlags().Lookup("log-file"))

	// old db TOML bindings
	viper.BindPFlag(migration_tools.TOML_OLD_DATABASE_NAME, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_OLD_DATABASE_NAME))
	viper.BindPFlag(migration_tools.TOML_OLD_DATABASE_HOSTNAME, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_OLD_DATABASE_HOSTNAME))
	viper.BindPFlag(migration_tools.TOML_OLD_DATABASE_PORT, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_NEW_DATABASE_PORT))
	viper.BindPFlag(migration_tools.TOML_OLD_DATABASE_USER, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_NEW_DATABASE_USER))
	viper.BindPFlag(migration_tools.TOML_OLD_DATABASE_PASSWORD, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_OLD_DATABASE_PASSWORD))
	viper.BindPFlag(migration_tools.TOML_OLD_DATABASE_MAX_IDLE_CONNECTIONS, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_OLD_DATABASE_MAX_IDLE_CONNECTIONS))
	viper.BindPFlag(migration_tools.TOML_OLD_DATABASE_MAX_OPEN_CONNECTIONS, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_OLD_DATABASE_MAX_OPEN_CONNECTIONS))
	viper.BindPFlag(migration_tools.TOML_OLD_DATABASE_MAX_CONN_LIFETIME, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_OLD_DATABASE_MAX_CONN_LIFETIME))

	// new db TOML bindings
	viper.BindPFlag(migration_tools.TOML_NEW_DATABASE_NAME, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_NEW_DATABASE_NAME))
	viper.BindPFlag(migration_tools.TOML_NEW_DATABASE_HOSTNAME, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_NEW_DATABASE_HOSTNAME))
	viper.BindPFlag(migration_tools.TOML_NEW_DATABASE_PORT, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_NEW_DATABASE_PORT))
	viper.BindPFlag(migration_tools.TOML_NEW_DATABASE_USER, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_NEW_DATABASE_USER))
	viper.BindPFlag(migration_tools.TOML_NEW_DATABASE_PASSWORD, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_NEW_DATABASE_PASSWORD))
	viper.BindPFlag(migration_tools.TOML_NEW_DATABASE_MAX_IDLE_CONNECTIONS, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_NEW_DATABASE_MAX_IDLE_CONNECTIONS))
	viper.BindPFlag(migration_tools.TOML_NEW_DATABASE_MAX_OPEN_CONNECTIONS, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_NEW_DATABASE_MAX_OPEN_CONNECTIONS))
	viper.BindPFlag(migration_tools.TOML_NEW_DATABASE_MAX_CONN_LIFETIME, rootCmd.PersistentFlags().Lookup(migration_tools.CLI_NEW_DATABASE_MAX_CONN_LIFETIME))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			log.Printf("Using config file: %s", viper.ConfigFileUsed())
		} else {
			log.Fatal(fmt.Sprintf("Couldn't read config file: %s", err.Error()))
		}
	} else {
		log.Warn("No config file passed with --config flag")
	}
}
