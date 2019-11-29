package toolkit

import (
	"github.com/gocql/gocql"
	std "github.com/PKUJohnson/solar/std"
)

type CassandraHelp struct {
	Session       *gocql.Session
	ClusterConfig *gocql.ClusterConfig
}

var cassandraHelper *CassandraHelp

func InitNewCassandraHelp(config std.ConfigCassandra) *CassandraHelp {
	cluster := gocql.NewCluster(config.Hosts...)
	cluster.Keyspace = config.Keyspace
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()

	cassandraHelper := &CassandraHelp{
		Session:       session,
		ClusterConfig: cluster,
	}

	return cassandraHelper
}

func GetCassandraHelp() *CassandraHelp {
	return cassandraHelper
}

// Close closes the cassandra session
func (c *CassandraHelp) Close() {
	c.Session.Close()
}
