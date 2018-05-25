package sound

import (
	"fmt"
	"io"

	"github.com/fvdveen/bf"
	"github.com/jonas747/dca"
	"github.com/op/go-logging"
	"gopkg.in/oleiade/lane.v1"
)

var (
	log   = logging.MustGetLogger("commands/sound")
	queue = lane.NewQueue()
	skip  = make(chan interface{})
	stop  = make(chan interface{})
	_     = startQueue()

	playing bool
)

type sound struct {
	ctx     bf.Context
	content [][]byte

	author, name string
	view         bool
}

func encodeSessionToBytes(enc dca.OpusReader) ([][]byte, error) {
	slice := [][]byte{}

loop:
	for {
		f, err := enc.OpusFrame()
		if err == io.EOF {
			break loop
		} else if err != nil {
			return nil, err
		}
		slice = append(slice, f)
	}

	return slice, nil
}

func startQueue() interface{} {
	go playQueue()
	return nil
}

func playQueue() {
	if queue == nil {
		queue = lane.NewQueue()
	}
	for {
		playing = true

		if queue.Head() != nil {
			s := queue.Dequeue()
			snd, ok := s.(sound)
			if !ok {
				playing = false
				continue
			}

			vc, err := snd.ctx.JoinVoiceChannel(false, true)
			if err != nil {
				log.Errorf("Could not join voice channel: %v", err)
				playing = false
				continue
			}

			if err := vc.Speaking(true); err != nil {
				log.Errorf("could not set speaking status: %v", err)
				if err := vc.Disconnect(); err != nil {
					log.Errorf("Could not disconnect: %v", err)
				}
				playing = false
				continue
			}

			if snd.view {
				if err := snd.ctx.SendMessage(fmt.Sprintf("Now playing %s - %s", snd.name, snd.author)); err != nil {
					log.Errorf("Could not send message: %v", err)
				}
			}

		sendSound:
			for _, frame := range snd.content {
				select {
				case <-skip:
					break sendSound
				case <-stop:
					for queue.Head() != nil {
						_ = queue.Dequeue()
					}
					break sendSound
				default:
					vc.OpusSend <- frame
				}
			}

			for queue.Head() != nil {
				s = queue.Dequeue()
				snd, ok := s.(sound)
				if !ok {
					playing = false
					continue
				}

				if snd.view {
					if err := snd.ctx.SendMessage(fmt.Sprintf("Now playing %s - %s", snd.name, snd.author)); err != nil {
						log.Errorf("Could not send message: %v", err)
					}
				}

			sendMoreSound:
				for _, frame := range snd.content {
					select {
					case <-skip:
						break sendMoreSound
					case <-stop:
						for queue.Head() != nil {
							_ = queue.Dequeue()
						}
						break sendMoreSound
					default:
						vc.OpusSend <- frame
					}
				}
			}

			if err := vc.Speaking(false); err != nil {
				log.Errorf("Could not set speaking status: %v", err)
			}

			if err := vc.Disconnect(); err != nil {
				log.Errorf("Could not disconnect: %v", err)
			}
		}

		playing = false
	}
}
