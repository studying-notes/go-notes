package client

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/streadway/amqp"
)

var (
	// beatChan = make(chan bool, 16)
	beatTime = time.Second * 30
	pubTime  = time.Second * 16
	tickTime = time.Second * 8

	messageTTL  = int64(time.Hour / time.Millisecond)          // TTL for message in queue 
	queueExpire = int64(time.Hour * 24 * 7 / time.Millisecond) // expire time for unused queue

	errAck     = errors.New("ack")
	errNack    = errors.New("nack")
	errFull    = errors.New("full")
	errCancel  = errors.New("cancel")
	errTimeout = errors.New("timeout")
)

// Client for mq
type Client struct {
	server *broker
	device string
	conn   *amqp.Connection
	ch     *amqp.Channel
	que    amqp.Queue
	msgs   <-chan amqp.Delivery

	destructor sync.Once
	confirm    chan amqp.Confirmation
	pubChan    chan *publishMsg
	ctx        context.Context
	cancel     context.CancelFunc
	onPublish  int32
}

type topic struct {
	chanName  string
	chanType  string
	keyPrefix string
}

type broker struct {
	server   string
	ssl      bool
	vhost    string
	userPass string

	ca   string
	cert string
	key  string

	quePrefix string

	topics []topic
}

// Close the client
func (clt *Client) Close() {
	clt.destructor.Do(func() {
		if clt.conn != nil {
			clt.conn.Close()
			clt.conn = nil
			clt.ch = nil
		}
		clt.cancel()
	})
}

func (clt *Client) connInit(tlsConfig *tls.Config, server *broker, proxyAddr string) (err error) {

	protocol := "amqp"
	if server.ssl {
		protocol = "amqps"
	}
	var conn *amqp.Connection
	if len(proxyAddr) > 0 {
		conn, err = amqp.DialConfig(protocol+"://"+server.userPass+"@"+server.server+"/"+server.vhost,
			amqp.Config{
				Heartbeat:       10 * time.Second,
				TLSClientConfig: tlsConfig,
				Locale:          "en_US",
				Dial: func(network, addr string) (net.Conn, error) {
					return net.Dial("tcp4", proxyAddr)
				}})
	} else {
		conn, err = amqp.DialTLS(protocol+"://"+server.userPass+"@"+server.server+"/"+server.vhost, tlsConfig)
	}

	if err != nil {
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		return
	}
	ch.Confirm(false)
	clt.confirm = make(chan amqp.Confirmation, 16)
	clt.confirm = ch.NotifyPublish(clt.confirm)

	for _, topic := range server.topics {
		err = ch.ExchangeDeclare(
			topic.chanName, // name
			topic.chanType, // type
			true,           // durable
			false,          // auto-deleted
			false,          // internal
			false,          // no-wait
			nil,            // arguments
		)
		if err != nil {
			return
		}
	}

	clt.conn = conn
	clt.ch = ch

	return
}

func (clt *Client) queInit(server *broker, ifFresh bool) (err error) {

	var num int
	ch := clt.ch

	if ifFresh {
		num, err = ch.QueueDelete(
			server.quePrefix+"."+clt.device,
			false,
			false,
			false,
		)
		if err != nil {
			return
		}
		log.Println("[RABBITMQ_CLIENT]", clt.device, "queue deleted with", num, "message purged")
	}

	args := make(amqp.Table)
	args["x-message-ttl"] = messageTTL
	args["x-expires"] = queueExpire
	q, err := ch.QueueDeclare(
		server.quePrefix+"."+clt.device, // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		args,  // arguments
	)
	if err != nil {
		return
	}

	for _, topic := range server.topics {
		err = ch.QueueBind(
			q.Name,
			topic.keyPrefix+"."+clt.device,
			topic.chanName,
			false,
			nil,
		)
		if err != nil {
			return
		}
	}

	clt.que = q
	return
}

type publishMsg struct {
	topicId      int
	keySuffix    string
	msg          []byte
	expire       time.Duration
	startTime    time.Time
	ctx          context.Context
	cancel       context.CancelFunc
	ackErr       error
	// ackChan      chan error
}

func (clt *Client) sendPublish(topicId int, keySuffix string, msg []byte, expire time.Duration) error {
	topic := &clt.server.topics[topicId]
	if expire <= 0 {
		return errors.New("Expiration parameter error")
	}
	
	return clt.ch.Publish(topic.chanName, topic.keyPrefix+"."+keySuffix, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msg,
		Expiration:  fmt.Sprintf("%d", int64(expire/time.Millisecond)),
	})
}

// heartBeat publish beatmsg to topic through mq client.
// check the beatChan to identify whether the client connection is health.
func (clt *Client) heartBeat(topicId int, keySuffix string, beatChan chan bool) bool {
	pubTimer := time.NewTimer(time.Second * 8)
	var beatMsg struct {
		T       string `json:"t"`
		Expired int64  `json:"expired"`
		Msid    string `json:"msid"`
	}

	beatMsg.T = "heartbeat"
	beatMsg.Expired = time.Now().Add(beatTime).Unix()
	uuid := md5.Sum([]byte("heartbeat-" + fmt.Sprintf("%d", beatMsg.Expired) + clt.device))
	beatMsg.Msid = fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
	beatStr, _ := json.Marshal(beatMsg)

	for i := 0; i < 3; i++ {
		pubTimer.Reset(pubTime)
		err := clt.Publish(topicId, true, keySuffix, beatTime, beatStr)
		if err != nil {
			log.Println("[RABBITMQ_CLIENT]", "heartBeat publish error:", err)
			continue
		}
		select {
		case <-beatChan:
			return true
		case <-pubTimer.C:
		}
	}
	return false
}

func newConsumer(parent context.Context, msgProcess func(topic string, message []byte), server string, device string) *Client {

	clt := new(Client)
	clt.ctx, clt.cancel = context.WithCancel(parent)

	clt.server = server
	clt.device = device
	tlsConfig, err := clt.server.getTLSconfig()
	if err != nil {
		log.Println("[RABBITMQ_CLIENT]", "tls ERROR:", err)
		return nil
	}

	err = clt.connInit(tlsConfig, clt.server)
	if err != nil {
		log.Println("[RABBITMQ_CLIENT]", "connInit ERROR:", err)
		return nil
	}
	err = clt.queInit(clt.server, false)
	if err != nil {
		clt.Close()
		err = clt.connInit(tlsConfig, clt.server)
		if err != nil {
			log.Println("[RABBITMQ_CLIENT]", "connInit ERROR:", err)
			return nil
		}
		err = clt.queInit(clt.server, true)
	}
	if err != nil {
		clt.Close()
		log.Println("[RABBITMQ_CLIENT]", "queInit ERROR:", err)
		return nil
	}

	msgs, err := clt.ch.Consume(
		clt.que.Name, // queue
		clt.device,   // consumer
		false,        // auto ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // args
	)
	if err != nil {
		clt.Close()
		log.Println("[RABBITMQ_CLIENT]", "Start consume ERROR:", err)
		return nil
	}

	clt.msgs = msgs
	clt.pubChan = make(chan *publishMsg, 4)

	go func() {
		cc := make(chan *amqp.Error)
		e := <-clt.ch.NotifyClose(cc)
		log.Println("[RABBITMQ_CLIENT]", "channel close error:", e.Error())
		clt.cancel()
	}()

	go func() {
		for d := range msgs {
			msg := d.Body
			msgProcess(d.Exchange, msg)
			d.Ack(false)
		}
	}()

	return clt
}

func (clt *Client) publishProc() {
	ticker := time.NewTicker(tickTime)
	deliveryMap := make(map[uint64]*publishMsg)

	defer func() {
		atomic.AddInt32(&clt.onPublish, -1)
		ticker.Stop()
		for _, msg := range deliveryMap {
			msg.ackErr = errCancel
			msg.cancel()
		}
	}()

	var deliveryTag uint64 = 1
	var ackTag uint64 = 1
	var pMsg *publishMsg
	for {
		select {

		case <-clt.ctx.Done():
			return

		case pMsg = <-clt.pubChan:
			pMsg.startTime = time.Now()
			err := clt.sendPublish(pMsg.topicId, pMsg.keySuffix, pMsg.msg, pMsg.expire)
			if err != nil {
				pMsg.ackErr = err
				pMsg.cancel()
			}
			deliveryMap[deliveryTag] = pMsg
			deliveryTag++

		case c, ok := <-clt.confirm:
			if !ok {
				log.Println("[RABBITMQ_CLIENT]", "client Publish notify channel error")
				return
			}
			pMsg = deliveryMap[c.DeliveryTag]
			// fmt.Println("DeliveryTag:", c.DeliveryTag)
			delete(deliveryMap, c.DeliveryTag)
			if c.Ack {
				pMsg.ackErr = nil
				pMsg.cancel()
			} else {
				pMsg.ackErr = errNack
				pMsg.cancel()
			}
		case <-ticker.C:
			now := time.Now()
			for {
				if len(deliveryMap) == 0 {
					break
				}
				pMsg = deliveryMap[ackTag]
				if pMsg != nil {
					if now.Sub(pMsg.startTime.Add(pubTime)) > 0 {
						pMsg.ackErr = errTimeout
						pMsg.cancel()
						delete(deliveryMap, ackTag)
					} else {
						break
					}
				}
				ackTag++
			}
		}
	}
}

// Publish used to send message to topic in messageQueue.
func (clt *Client) Publish(topicID int, keySuffix string, expire time.Duration, msg []byte) (err error) {

	pMsg := publishMsg{
		topicId:      topicID,
		keySuffix:    keySuffix,
		expire:       expire,
		msg:          msg,
	}

	// 在client中，pubchan是预先建立好的，但是只有在有publish时，才创建publishProc
	// 如果发送过程出现异常导致publishProc退出，此时onPublish被置零，可以再次创建新的publishProc
	// pubChan中可能有残存的msg，如果没有及时新的publishProc启动，则对这些消息的处理是无用的
	// 如果当前client关闭，pubchan不会立刻关闭（等待gc），已经进入发送过程的publish会等待超时
	if atomic.AddInt32(&clt.onPublish, 1) == 1 {
		go clt.publishProc()
	} else {
		atomic.AddInt32(&clt.onPublish, -1)
	}

	timer := time.NewTimer(pubTime)
	defer timer.Stop()

	pMsg.ctx, pMsg.cancel = context.WithCancel(context.Background())
	defer pMsg.cancel()
	select {
	case <-timer.C:
		err = errFull
		return
	case clt.pubChan <- &pMsg:
	}

	timer.Reset(pubTime)
	select {
	case <-pMsg.ctx.Done():
		err = pMsg.ackErr
		break
	case <-timer.C:
		err = errTimeout
		break
	case <-clt.ctx.Done():
		err = errCancel
		break
	}

	return
}

// NewPublisher init a Publisher of rabbitmq client
func NewPublisher(parent context.Context, server string) *Client {
	clt := new(Client)
	clt.ctx, clt.cancel = context.WithCancel(parent)
	clt.server = server
	tlsConfig, err := clt.server.getTLSconfig()
	if err != nil {
		log.Println("[RABBITMQ_CLIENT]", "tls ERROR:", err)
		return nil
	}

	err = clt.connInit(tlsConfig, clt.server)
	if err != nil {
		log.Println("[RABBITMQ_CLIENT]", "connInit ERROR:", err)
		return nil
	}

	return clt
}