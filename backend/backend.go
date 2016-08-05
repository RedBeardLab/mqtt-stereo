package backend

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/urfave/cli"
)

func createUrl(url string, port int) string {
	return fmt.Sprintf("tcp://%s:%d", url, port)
}

type message struct {
	Time    time.Time `json:string`
	Topic   string    `json:string`
	Payload string    `json:string`
}

func newRecoder(r *os.File, exitSignalCh chan os.Signal) chan MQTT.Message {
	c := make(chan MQTT.Message)
	go func() {
		w := bufio.NewWriter(r)
		for {
			select {
			case toSave := <-c:
				{
					now := time.Now()
					m := message{now,
						toSave.Topic(),
						string(toSave.Payload())}
					mStr, err := json.Marshal(m)
					if err != nil {
						log.Print("Problems with a message: topic %s, payload: %s", toSave.Topic(), toSave.Payload())
					} else {
						fmt.Fprintf(w, "%+v\n", string(mStr))
					}
				}

			case <-exitSignalCh:
				{
					fmt.Println("Flushing to disk")
					w.Flush()
					r.Close()
					os.Exit(1)
					return
				}
			}
		}
	}()
	return c
}

func StartRecording(c *cli.Context) {
	url := createUrl(c.GlobalString("url"), c.GlobalInt("port"))
	file, err := os.Create(c.GlobalString("record"))
	if err != nil {
		log.Fatal(err)
	}
	topic := c.GlobalString("topic")

	opts := MQTT.NewClientOptions()
	opts.AddBroker(url)

	receiver := MQTT.NewClient(opts)
	if token := receiver.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	exitSignalCh := make(chan os.Signal)
	signal.Notify(exitSignalCh, os.Interrupt)
	signal.Notify(exitSignalCh, syscall.SIGTERM)

	recoderCh := newRecoder(file, exitSignalCh)
	f := func(receiver MQTT.Client, msg MQTT.Message) {
		recoderCh <- msg
	}

	if token := receiver.Subscribe(topic, 1, f); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	for {
	}

}

func PlayBack(c *cli.Context) {
	url := createUrl(c.GlobalString("url"), c.GlobalInt("port"))

	opts := MQTT.NewClientOptions()
	opts.AddBroker(url)

	sender := MQTT.NewClient(opts)
	if token := sender.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	playLoop := c.Bool("loop")
	loopCount := 0
	for playLoop || loopCount == 0 {
		loopCount++
		file, err := os.Open(c.GlobalString("record"))
		if err != nil {
			log.Fatal(err)
		}
		reader := bufio.NewScanner(file)
		fastForward := c.Bool("ff")
		var message message
		var previousTime time.Time
		firstRound := true
		for reader.Scan() {
			json.Unmarshal([]byte(reader.Text()), &message)
			if !fastForward {
				if !firstRound {
					toWait := message.Time.Sub(previousTime)
					previousTime = message.Time
					time.Sleep(toWait)
				} else {
					firstRound = false
					previousTime = message.Time
				}
			}
			fmt.Println(message.Payload)
			if token := sender.Publish(message.Topic, 1, false, message.Payload); token.Wait() && token.Error() != nil {
				fmt.Println(token.Error())
			}
		}
	}

}
