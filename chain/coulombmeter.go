package chain

import (
	"context"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
)

type MSGReport struct {
	DeviceID   int64      `json:"deviceId"`
	Success    bool       `json:"success"`
	Timestamp  string     `json:"timestamp"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	TotalValue  float64 `json:"Total_Value"`
	TotalValue1 float64 `json:"Total_Value1"`
	TotalValue2 float64 `json:"Total_Value2"`
	TotalValue3 float64 `json:"Total_Value3"`
	TotalValue4 float64 `json:"Total_Value4"`
	UA          float64 `json:"UA"`
	UB          float64 `json:"UB"`
	UC          float64 `json:"UC"`
	IA          float64 `json:"IA"`
	IB          float64 `json:"IB"`
	IC          float64 `json:"IC"`
	P           float64 `json:"P"`
	PA          float64 `json:"PA"`
	PB          float64 `json:"PB"`
	PC          float64 `json:"PC"`
	Q           float64 `json:"Q"`
	QA          float64 `json:"QA"`
	QB          float64 `json:"QB"`
	QC          float64 `json:"QC"`
	PF          float64 `json:"PF"`
	PFA         float64 `json:"PFA"`
	PFB         float64 `json:"PFB"`
	PFC         float64 `json:"PFC"`
	S           float64 `json:"s"`
}

func MqttServer() {
	mqttString := "tcp://106.12.12.197:1884"

	opts := MQTT.NewClientOptions()
	opts.AddBroker(mqttString)
	opts.SetClientID("info-collector")
	//opts.SetConnectTimeout(time.Duration(60) * time.Second)
	opts.SetCleanSession(true)
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic := "/report-property"
	c.Subscribe(topic, 0, CronReportHandler)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	c.Disconnect(250)
}

func CronReportHandler(client MQTT.Client, msg MQTT.Message) {
	log.Printf("CronReportHandler        [%s]\n", msg.Topic())
	//log.Printf("%s\n", msg.Payload())

	var message MSGReport
	if err := json.Unmarshal(msg.Payload(), &message); err != nil {
		log.Println("json unmarshal  error:", err)
		return
		// panic)
	}
	log.Println(message)
	coulombmeterGateway := message.DeviceID
	var addr = str2bytes("上海市虹口区海伦路111号11栋1号1单元101室")
	//tx, err = dLeaseHold.RegisterDevice(factory, addr, deviceAuth.From, factory.From, big.NewInt(coulombmeterGateway).Bytes())
	tx, _ = dLeaseHold.RegisterDevice(auths["factory"], addr, auths["device"].From, auths["factory"].From, big.NewInt(coulombmeterGateway).Bytes())
	err = GethCli.CheckTx(context.Background(), tx)
	fmt.Println("RegisterDevice", tx.Hash().String())
	if err != nil {
		log.Println(err)
		TotalValue := message.Properties.TotalValue
		tx, err = dLeaseHold.UploadLogs(auths["device"], big.NewInt(int64(TotalValue*1000)))
	}

	GethCli.CheckTx(context.Background(), tx)
	receipt, _ = GethCli.TransactionReceipt(context.Background(), tx.Hash())
	fmt.Println("UploadLogs", tx.Hash().String())
	//fmt.Printf("coulombmeter    %v\n",receipt)

	log.Println(message)
}
