# MQTT Stereo

A simple MQTT utility to record and play back MQTT streams.

## Obtain

The software is written in Go(lang), if your go environment is set the fastest way to get started is:

``` bash
go get github.com/RedBeardLab/mqtt-stereo
```

if you would like to use this utility but you don't want to set up all the go environment, please open an issues and we will provide the binary for your architecture.

## Usage

The binary should already provide some information, hopefully enough to get started.

There are two modality, play and record.

mqtt-stereo need a file where is read and write the mqtt messages.

## Getting started

There are few parameters that can be set:

 * record, what file use to record and to play the messages (default to "mqtt-record.txt")
 * topic, what topic to register or to play (default to "/#")
 * url, the url of the MQTT broker (default to "localhost")
 * port, in which port the broker is serving (default to 1883)

Finally you need to specify if you want to record or to play a MQTT stream, you can do it, respectively, with the command "record" (aliases "rec" or "r") and the command "play" (alias "p")

## Example

To record an MQTT stream from a remote host, suppose "mqtt-example.com", with serve MQTT on the port 1567 but we are interested only on the topic "/IoT/#" and we want to save everything on the file "IoT_log.txt" you should run:

``` bash
mqtt-stereo --url mqtt-example.com --port 1567 --topic /IoT/# --record IoT_log.txt record
```

After enough messages are received you can exit with Ctrl-C

Now to play back all those messages to another host "iot-mqtt.com" that serve on the port 1884 you should run:

```bash
mqtt-stereo --url iot-mqtt.com --port 1884 --record IoT_log.txt play
```

This command will play all the messages only once respecting the time deltas between the messages.

If for debug reason you want to play all the messages in loop, multiple times, you can add loop subcommand.

```bash
mqtt-stereo --url iot-mqtt.com --port 1884 --record IoT_log.txt play --loop
```

Finally if for benchmark reason you want to run all the messages as fast as possible you can use the fast-forward reason.

```bash
mqtt-stereo --url iot-mqtt.com --port 1884 --record IoT_log.txt play --ff
```

Of course the `--loop` and the `--ff` subcommand can be used together.

## Motivation

At [RedBeardLab.tech](redbeardlab.tech) we are getting a lot of exposure and experience about IoT and MQTT, but we didn't found any suitable tool for debugging and record MQTT messages and mqtt-stereo is our solution to this problem.

Of course this is a very simple implementation, lacking many features, feel free to contribute or open an issue.
