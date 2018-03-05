package cmd

import (
	"fmt"
	"log"

	bolt "github.com/coreos/bbolt"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Set, delete or list package aliases",
}

var setAliasCmd = &cobra.Command{
	Use:   "set [pkgname] [pkgpath]",
	Short: "Set alias to easy access package documentation",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dpath, err := findGdocDB()
		if err != nil {
			log.Fatal("error:", err)
		}
		db, err := bolt.Open(dpath, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		name, pkg := args[0], args[1]
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("alias"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return b.Put([]byte(name), []byte(pkg))
		})
	},
}

var delAliasCmd = &cobra.Command{
	Use:   "del [pkgname]",
	Short: "Delete an alias from database",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dpath, err := findGdocDB()
		if err != nil {
			log.Fatal("error:", err)
		}
		db, err := bolt.Open(dpath, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		name := args[0]
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("alias"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return b.Delete([]byte(name))
		})
	},
}

var listAliasCmd = &cobra.Command{
	Use:   "list",
	Short: "List available aliases",
	Run: func(cmd *cobra.Command, args []string) {
		dpath, err := findGdocDB()
		if err != nil {
			log.Fatal("error:", err)
		}
		db, err := bolt.Open(dpath, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("alias"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("name: %s, path: %s\n", k, v)
			}
			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(aliasCmd)
	aliasCmd.AddCommand(setAliasCmd)
	aliasCmd.AddCommand(delAliasCmd)
	aliasCmd.AddCommand(listAliasCmd)
}
