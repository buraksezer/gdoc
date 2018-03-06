package cmd

import (
	"fmt"

	bolt "github.com/coreos/bbolt"
	"github.com/spf13/cobra"
)

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
			logger.Fatal(err)
		}
		db, err := bolt.Open(dpath, 0600, nil)
		if err != nil {
			logger.Fatalf("failed to open db: %v", err)
		}
		defer func() {
			if derr := db.Close(); derr != nil {
				logger.Println("failed to close boltdb: ", derr)
			}
		}()

		name, pkg := args[0], args[1]
		err = db.Update(func(tx *bolt.Tx) error {
			b, berr := tx.CreateBucketIfNotExists([]byte("alias"))
			if berr != nil {
				return fmt.Errorf("create bucket: %s", berr)
			}
			return b.Put([]byte(name), []byte(pkg))
		})
		if err != nil {
			logger.Fatalf("failed to get aliases: %v", err)
		}
	},
}

var delAliasCmd = &cobra.Command{
	Use:   "del [pkgname]",
	Short: "Delete an alias from database",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dpath, err := findGdocDB()
		if err != nil {
			logger.Fatal(err)
		}
		db, err := bolt.Open(dpath, 0600, nil)
		if err != nil {
			logger.Fatalf("failed to open db: %v", err)
		}
		defer func() {
			if derr := db.Close(); derr != nil {
				logger.Println("failed to close boltdb: ", derr)
			}
		}()

		name := args[0]
		err = db.Update(func(tx *bolt.Tx) error {
			b, berr := tx.CreateBucketIfNotExists([]byte("alias"))
			if berr != nil {
				return fmt.Errorf("create bucket: %s", berr)
			}
			return b.Delete([]byte(name))
		})
		if err != nil {
			logger.Fatalf("failed to get aliases: %v", err)
		}
	},
}

var listAliasCmd = &cobra.Command{
	Use:   "list",
	Short: "List available aliases",
	Run: func(cmd *cobra.Command, args []string) {
		dpath, err := findGdocDB()
		if err != nil {
			logger.Fatal(err)
		}
		db, err := bolt.Open(dpath, 0600, nil)
		if err != nil {
			logger.Fatalf("failed to open db: %v", err)
		}
		defer func() {
			if derr := db.Close(); derr != nil {
				logger.Println("failed to close boltdb: ", derr)
			}
		}()

		err = db.Update(func(tx *bolt.Tx) error {
			b, berr := tx.CreateBucketIfNotExists([]byte("alias"))
			if berr != nil {
				return fmt.Errorf("create bucket: %s", berr)
			}
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("name: %s, path: %s\n", k, v)
			}
			return nil
		})
		if err != nil {
			logger.Fatalf("failed to get aliases: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
	aliasCmd.AddCommand(setAliasCmd)
	aliasCmd.AddCommand(delAliasCmd)
	aliasCmd.AddCommand(listAliasCmd)
}
