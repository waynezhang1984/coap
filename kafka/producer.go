package kafka

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	//	"time"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/tools/tls"
	"github.com/coapprocess/configure"
	"github.com/rcrowley/go-metrics"
)

func ProduceMsg(cp configure.ProducerCfg, msg string) {
	if cp.BrokerList == "" {
		printUsageErrorAndExit("no -brokers specified. Alternatively, set the KAFKA_PEERS environment variable")
	}

	if cp.Topic == "" {
		printUsageErrorAndExit("no -topic specified")
	}

	if cp.Verbose {
		sarama.Logger = logger
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	if cp.TlsEnabled {
		tlsConfig, err := tls.NewConfig(cp.TlsClientCert, cp.TlsClientKey)
		if err != nil {
			printErrorAndExit(69, "Failed to create TLS config: %s", err)
		}

		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Config.InsecureSkipVerify = cp.TlsSkipVerify
	}

	switch cp.Partitioner {
	case "":
		if cp.Partition >= 0 {
			config.Producer.Partitioner = sarama.NewManualPartitioner
		} else {
			config.Producer.Partitioner = sarama.NewHashPartitioner
		}
	case "hash":
		config.Producer.Partitioner = sarama.NewHashPartitioner
	case "random":
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	case "manual":
		config.Producer.Partitioner = sarama.NewManualPartitioner
		if cp.Partition == -1 {
			printUsageErrorAndExit("-partition is required when partitioning manually")
		}
	default:
		printUsageErrorAndExit(fmt.Sprintf("Partitioner %s not supported.", cp.Partitioner))
	}

	message := &sarama.ProducerMessage{Topic: cp.Topic, Partition: int32(cp.Partition)}

	if cp.Key != "" {
		message.Key = sarama.StringEncoder(cp.Key)
	}

	if msg != "" {
		message.Value = sarama.StringEncoder(msg)
	} else if stdinAvailable() {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			printErrorAndExit(66, "Failed to read data from the standard input: %s", err)
		}
		message.Value = sarama.ByteEncoder(bytes)
	} else {
		printUsageErrorAndExit("-value is required, or you have to provide the value on stdin")
	}

	producer, err := sarama.NewSyncProducer(strings.Split(cp.BrokerList, ","), config)
	if err != nil {
		printErrorAndExit(69, "Failed to open Kafka producer: %s", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Println("Failed to close Kafka producer cleanly:", err)
		}
	}()
	//fmt.Println("send time: ", time.Now())
	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		printErrorAndExit(69, "Failed to produce message: %s", err)
	} else if !cp.Silent {
		fmt.Printf("topic=%s\tpartition=%d\toffset=%d\n", cp.Topic, partition, offset)
	}
	if cp.ShowMetrics {
		metrics.WriteOnce(config.MetricRegistry, os.Stderr)
	}
}

func stdinAvailable() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
