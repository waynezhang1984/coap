package kafka

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"

	//	"time"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/tools/tls"
	"github.com/coapprocess/configure"

	//	"github.com/devicemanager/mqttclient"
	"github.com/coapprocess/msginfo"
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func ConsumeMsg(cc configure.ConsumerCfg) {
	if cc.BrokerList == "" {
		printUsageErrorAndExit("You have to provide -brokers as a comma-separated list, or set the KAFKA_PEERS environment variable.")
	}

	if cc.Topic == "" {
		printUsageErrorAndExit("-topic is required")
	}

	if cc.Verbose {
		sarama.Logger = logger
	}

	var initialOffset int64
	switch cc.Offset {
	case "oldest":
		initialOffset = sarama.OffsetOldest
	case "newest":
		initialOffset = sarama.OffsetNewest
	default:
		printUsageErrorAndExit("-offset should be `oldest` or `newest`")
	}

	config := sarama.NewConfig()
	if cc.TlsEnabled {
		tlsConfig, err := tls.NewConfig(cc.TlsClientCert, cc.TlsClientKey)
		if err != nil {
			printErrorAndExit(69, "Failed to create TLS config: %s", err)
		}

		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Config.InsecureSkipVerify = cc.TlsSkipVerify
	}

	c, err := sarama.NewConsumer(strings.Split(cc.BrokerList, ","), config)
	if err != nil {
		printErrorAndExit(69, "Failed to start consumer: %s", err)
	}

	partitionList, err := getPartitions(c, cc)
	if err != nil {
		printErrorAndExit(69, "Failed to get the list of partitions: %s", err)
	}

	var (
		messages = make(chan *sarama.ConsumerMessage, cc.BufferSize)
		closing  = make(chan struct{})
		wg       sync.WaitGroup
	)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Kill, os.Interrupt)
		<-signals
		logger.Println("Initiating shutdown of consumer...")
		close(closing)
	}()

	for _, partition := range partitionList {
		pc, err := c.ConsumePartition(cc.Topic, partition, initialOffset)
		if err != nil {
			printErrorAndExit(69, "Failed to start consumer for partition %d: %s", partition, err)
		}

		go func(pc sarama.PartitionConsumer) {
			<-closing
			pc.AsyncClose()
		}(pc)

		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for message := range pc.Messages() {
				messages <- message
			}
		}(pc)
	}

	go func() {
		for msg := range messages {
			//fmt.Println("recv time: ", time.Now())
			//fmt.Printf("Value:\t%s\n", string(msg.Value))

			msg := msginfo.MsgData{
				PubType: "kafka",
				Data:    msg.Value,
			}

			msginfo.GetDataChan() <- msg

			//mqttclient.MqttPub(configure.GetMqttCfg(),string(msg.Value)
		}
	}()

	wg.Wait()
	logger.Println("Done consuming topic", cc.Topic)
	close(messages)

	if err := c.Close(); err != nil {
		logger.Println("Failed to close consumer: ", err)
	}
}

func getPartitions(c sarama.Consumer, cc configure.ConsumerCfg) ([]int32, error) {
	if cc.Partitions == "all" {
		return c.Partitions(cc.Topic)
	}

	tmp := strings.Split(cc.Partitions, ",")
	var pList []int32
	for i := range tmp {
		val, err := strconv.ParseInt(tmp[i], 10, 32)
		if err != nil {
			return nil, err
		}
		pList = append(pList, int32(val))
	}

	return pList, nil
}

func printErrorAndExit(code int, format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", fmt.Sprintf(format, values...))
	fmt.Fprintln(os.Stderr)
	os.Exit(code)
}

func printUsageErrorAndExit(format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", fmt.Sprintf(format, values...))
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Available command line options:")
	flag.PrintDefaults()
	os.Exit(64)
}
