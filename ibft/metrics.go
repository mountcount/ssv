package ibft

import (
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/bloxapp/ssv/utils/format"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"strconv"
)

var (
	metricsDecidedSigners = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ssv:validator:ibft_decided_signers",
		Help: "The highest decided sequence number",
	}, []string{"lambda", "pubKey", "seq", "nodeId"})
)

func init() {
	if err := prometheus.Register(metricsDecidedSigners); err != nil {
		log.Println("could not register prometheus collector")
	}
}

// ReportDecided reports on a decided message
func ReportDecided(pk string, msg *proto.SignedMessage) {
	_, role := format.IdentifierUnformat(string(msg.Message.GetLambda()))
	for _, nodeID := range msg.SignerIds {
		metricsDecidedSigners.WithLabelValues(
			role, pk,
			strconv.FormatUint(msg.Message.SeqNumber, 10),
			strconv.FormatUint(nodeID, 10)).Set(1)
	}
}
