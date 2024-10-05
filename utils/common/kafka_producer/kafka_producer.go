package kafka_producer

import (
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"sync"
	"sync/atomic"
	"tzgit.kaixinxiyou.com/utils/common/jsonEncoder"

	"strings"
	"time"
	"tzgit.kaixinxiyou.com/utils/common/log"
)

type Kafka_Producer struct {
	pProducer    sarama.AsyncProducer
	kafka_iplist string
	input        chan *ProducerMessage
	inFlight     sync.WaitGroup
	closeTag     int32
	sendCount    int64
	successCount int64
	errCount     int64
}
type Kafka_Producer_interface interface {
	AsyncProducer(vjson string)
	Close()
	Kafka_insertDeviceRamLog(data map[interface{}]interface{}, Operation, Table string)
}

type ProducerMessage struct {
	Topic    string
	Value    []byte
	closeTag uint
}

func New_Kafka(iplist string, channelBufferSize int32) (*Kafka_Producer, error) {
	outres := &Kafka_Producer{
		kafka_iplist: iplist,
	}
	var config *sarama.Config
	config = sarama.NewConfig()
	config.Producer.Return.Successes = true //必须有这个选项
	config.Producer.Timeout = 5 * time.Second
	config.Version = sarama.V2_0_0_0
	if channelBufferSize <= 0 {
		channelBufferSize = 8
	}
	outres.input = make(chan *ProducerMessage, channelBufferSize)
	p, err := sarama.NewAsyncProducer(strings.Split(outres.kafka_iplist, ","), config)
	if err != nil {
		return nil, err
	}

	//必须有这个匿名函数内容
	go func(p sarama.AsyncProducer, a *Kafka_Producer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					_, err := err.Msg.Value.Encode()
					if err != nil {
						log.Warn("AsyncProducer   %s", err)
					} else {

					}
					a.errCount++
				}
				//a.inFlight.Done()
				//a.successCount++

			case aa := <-success:
				if aa != nil {
					//log.Debug("%s", aa.Value)
					a.successCount++
				}

			}
		}
	}(p, outres)
	outres.pProducer = p
	go withRecover(outres.dispatcher)
	return outres, nil
}

// asyncProducer 异步生产者
// 并发量大时，必须采用这种方式
func (kfk *Kafka_Producer) AsyncProducer(d map[string]interface{}, topic string) error {

	if topic == "" {
		return errors.New("topic == 空")
	}
	data, err := json.Marshal(d)
	if err != nil {
		log.Error("%v", data)
		data = []byte{}
	}
	return kfk.AsyncProducerBytes(data, topic)
}
func (kfk *Kafka_Producer) AsyncProducerBytes(d []byte, topic string) error {
	//now := time.Now()
	if atomic.LoadInt32(&kfk.closeTag) == 1 {
		return errors.New("已经关闭")
	}
	if topic == "" {
		return errors.New("topic == 空")
	}
	msg := &ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(d),
	}
	var err error
	select {
	case kfk.input <- msg:
	default:
		err = errors.New("pProducer.Input() full")
	}
	return err
}

func (kfk *Kafka_Producer) AsyncProducerJson(d jsonEncoder.ObjectEncoder, topic string) error {
	//now := time.Now()
	if atomic.LoadInt32(&kfk.closeTag) == 1 {
		return errors.New("已经关闭")
	}
	if topic == "" {
		return errors.New("topic == 空")
	}
	//因为d.Marshal 需要回收，因此需要再这里创建新的对象
	msg := &ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(d.Clone()),
	}

	jsonEncoder.Put(d)
	var err error
	select {
	case kfk.input <- msg:
	default:
		err = errors.New("pProducer.Input() full")
	}
	return err
}

// singleton
// dispatches messages by topic
var c = 0

func (kfk *Kafka_Producer) dispatcher() {
	shuttingDown := false
	for msg := range kfk.input {
		if msg == nil {
			log.Warn("Something tried to send a nil message, it was ignored.")
			continue
		}
		if msg.closeTag != 0 {
			shuttingDown = true
			log.Debug("c %v", c)
			continue
		}
		if shuttingDown {
			continue
		}
		msg1 := &sarama.ProducerMessage{
			Topic: msg.Topic,
			Value: sarama.ByteEncoder(msg.Value),
		}
		kfk.sendCount++
		kfk.pProducer.Input() <- msg1
	}

}
func (kfk *Kafka_Producer) Close(timeoutMs int64) error {
	//kfk.pProducer.Close()
	startCloseTime := time.Now()
	atomic.StoreInt32(&kfk.closeTag, 1)
	withRecover(kfk.shutdown)
	for {
		if kfk.successCount+kfk.errCount == kfk.sendCount {
			return nil
		}
		if time.Since(startCloseTime) > time.Duration(timeoutMs)*time.Millisecond {
			return errors.New("The producer timeout closes, and some of the cached data may not be sent properly")
		}
		time.Sleep(100 * time.Millisecond)
	}

}
func (kfk *Kafka_Producer) shutdown() {
	kfk.input <- &ProducerMessage{closeTag: 1}
	kfk.inFlight.Wait()
}
func withRecover(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Warn("%v", err)
		}
	}()
	fn()
}
