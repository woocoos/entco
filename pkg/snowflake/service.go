package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"github.com/tsingsun/woocoo/pkg/conf"
	"os"
	"strconv"
)

var (
	defaultNode *snowflake.Node
)

func init() {
	node := 1
	ns := os.Getenv("SNOWFLAKE_NODE_ID")
	if ns != "" {
		n, _ := strconv.Atoi(ns)
		if n > 0 {
			node = n
		}
	}
	defaultNode, _ = snowflake.NewNode(int64(node))
}

func SetDefaultNode(cnf *conf.Configuration) (err error) {
	if nb := cnf.Int("nodeBits"); nb > 0 {
		snowflake.NodeBits = uint8(nb)
	}
	if sb := cnf.Int("stepBits"); sb > 0 {
		snowflake.StepBits = uint8(sb)
	}
	nid := cnf.Int("nodeID")
	if nid <= 0 {
		// try get from env
		if ns := os.Getenv("SNOWFLAKE_NODE_ID"); ns != "" {
			nid, err = strconv.Atoi(ns)
			if err != nil {
				return err
			}
		}
		if nid <= 0 {
			nid = 1
		}
	}
	defaultNode, err = snowflake.NewNode(int64(nid))
	return err
}

func New() snowflake.ID {
	return defaultNode.Generate()
}
