package main

import (
	"io"
	"log"
	"os"

	"github.com/maorfr/cain/pkg/cain"

	"github.com/spf13/cobra"
)

func main() {
	cmd := NewRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		log.Fatal("Failed to execute command")
	}
}

// NewRootCmd represents the base command when called without any subcommands
func NewRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cain",
		Short: "",
		Long:  ``,
	}

	out := cmd.OutOrStdout()

	cmd.AddCommand(NewBackupCmd(out))
	cmd.AddCommand(NewRestoreCmd(out))
	cmd.AddCommand(NewSchemaCmd(out))

	return cmd
}

type backupCmd struct {
	namespace string
	selector  string
	container string
	keyspace  string
	bucket    string
	tag       string
	parallel  int
	sum       bool

	out io.Writer
}

// NewBackupCmd performs a backup of a cassandra cluster
func NewBackupCmd(out io.Writer) *cobra.Command {
	b := &backupCmd{out: out}

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "backup cassandra cluster to S3",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cain.Backup(b.namespace, b.selector, b.container, b.keyspace, b.bucket, b.parallel); err != nil {
				log.Fatal(err)
			}
		},
	}
	f := cmd.Flags()

	f.StringVarP(&b.namespace, "namespace", "n", "", "namespace to find cassandra cluster")
	f.StringVarP(&b.selector, "selector", "l", "", "selector to filter on")
	f.StringVarP(&b.container, "container", "c", "cassandra", "container name to backup")
	f.StringVarP(&b.keyspace, "keyspace", "k", "", "keyspace to backup")
	f.StringVarP(&b.bucket, "bucket", "b", "", "bucket to backup to")
	f.IntVarP(&b.parallel, "parallel", "p", 1, "number of files to copy in parallel. set this flag to 0 for full parallelism")

	cmd.MarkFlagRequired("namespace")
	cmd.MarkFlagRequired("selector")
	cmd.MarkFlagRequired("keyspace")
	cmd.MarkFlagRequired("bucket")

	return cmd
}

type restoreCmd struct {
	bucket      string
	namespace   string
	cluster     string
	keyspace    string
	tag         string
	toNamespace string
	selector    string
	container   string
	parallel    int

	out io.Writer
}

// NewRestoreCmd performs a restore from backup of a cassandra cluster
func NewRestoreCmd(out io.Writer) *cobra.Command {
	r := &restoreCmd{out: out}

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "restore cassandra cluster from S3",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cain.Restore(r.bucket, r.namespace, r.cluster, r.keyspace, r.tag, r.toNamespace, r.selector, r.container, r.parallel); err != nil {
				log.Fatal(err)
			}
		},
	}
	f := cmd.Flags()

	f.StringVarP(&r.bucket, "bucket", "b", "", "bucket to restore from")
	f.StringVarP(&r.namespace, "namespace", "n", "", "namespace to restore from")
	f.StringVar(&r.cluster, "cluster", "", "cassandra cluster to restore from")
	f.StringVarP(&r.keyspace, "keyspace", "k", "", "keyspace to restore")
	f.StringVarP(&r.tag, "tag", "t", "", "tag to restore")
	f.StringVar(&r.toNamespace, "to-namespace", "", "namespace to find cassandra cluster")
	f.StringVarP(&r.selector, "selector", "l", "", "selector to filter on")
	f.StringVarP(&r.container, "container", "c", "cassandra", "container name to restore")
	f.IntVarP(&r.parallel, "parallel", "p", 1, "number of files to copy in parallel. set this flag to 0 for full parallelism")

	cmd.MarkFlagRequired("bucket")
	cmd.MarkFlagRequired("namespace")
	cmd.MarkFlagRequired("cluster")
	cmd.MarkFlagRequired("keyspace")
	cmd.MarkFlagRequired("tag")
	cmd.MarkFlagRequired("selector")

	return cmd
}

type schemaCmd struct {
	namespace string
	selector  string
	container string
	sum       bool

	out io.Writer
}

// NewSchemaCmd performs schema related actions
func NewSchemaCmd(out io.Writer) *cobra.Command {
	s := &schemaCmd{out: out}

	cmd := &cobra.Command{
		Use:   "schema",
		Short: "get schema of cassandra cluster",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cain.Schema(s.namespace, s.selector, s.container, s.sum); err != nil {
				log.Fatal(err)
			}
		},
	}
	f := cmd.Flags()

	f.StringVarP(&s.namespace, "namespace", "n", "", "namespace to find cassandra cluster")
	f.StringVarP(&s.selector, "selector", "l", "", "selector to filter on")
	f.StringVarP(&s.container, "container", "c", "cassandra", "container name to restore")
	f.BoolVar(&s.sum, "sum", false, "print only checksum")

	cmd.MarkFlagRequired("namespace")
	cmd.MarkFlagRequired("selector")

	return cmd
}
