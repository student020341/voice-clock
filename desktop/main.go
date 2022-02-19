package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func initializeSpeaker() error {
	testResource := mapInputToResource("0")
	tempStreamer, format, err := getStreamer(testResource)
	if err != nil {
		return err
	}

	tempStreamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	return err
}

func getHalloweenStreamer() (beep.Streamer, error) {
	donkStream, format, err := getStreamer(resourceDinkDonkOgg)
	if err != nil {
		return nil, err
	}

	announceStream, _, err := getStreamer(resourceCurrentTimeOgg)
	if err != nil {
		return nil, err
	}

	numberStream, _, err := getStreamer(resource6Ogg)
	if err != nil {
		return nil, err
	}

	//
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	peepoS := beep.ResampleRatio(4, 0.66, announceStream)
	// there was an issue looping these, problem with byte? fake read closer?
	// building new streamer from buffer works!
	monkaW := beep.ResampleRatio(4, 0.50, numberStream)

	buf1 := beep.NewBuffer(format)
	buf1.Append(donkStream)
	buf1.Append(peepoS)
	donkStream.Close()
	announceStream.Close()

	buf2 := beep.NewBuffer(format)
	buf2.Append(monkaW)
	numberStream.Close()

	clipChunk1 := buf1.Streamer(0, buf1.Len())
	clipChunk2 := beep.Loop(3, buf2.Streamer(0, buf2.Len()))

	return beep.Seq(clipChunk1, clipChunk2), nil
}

func doOctoberAudio() {

	hstreamer, err := getHalloweenStreamer()
	if err != nil {
		fmt.Printf("error: %+v", err)
		return
	}

	speaker.Play(hstreamer)
}

func doAudio(channel chan struct{}) {
	data := getTimeData(time.Now())
	seq, err := mapInputToSequence(data, channel)
	if err != nil {
		fmt.Printf("error: %+v", err)
		return
	}

	speaker.Play(seq)
}

func main() {
	fmt.Println("initializing speaker...")
	initializeSpeaker()
	autoHour := false // say time automatically on the hour
	playing := false
	doneChannel := make(chan struct{})

	fmt.Println("building windows...")
	a := app.New()
	a.SetIcon(resourceIconPng)
	w := a.NewWindow("Crump Clock")
	timeLabel := widget.NewLabel("--")

	// port input
	portInput := widget.NewEntry()
	portInput.SetText("2005")
	startServerButton := widget.NewButton("Start Server", nil)
	startServerButton.OnTapped = func() {
		if _, err := strconv.Atoi(portInput.Text); err == nil {
			startServerButton.Disable()
			portInput.Disable()
			portStr := fmt.Sprintf(":%s", portInput.Text)
			startServerButton.SetText(fmt.Sprintf("listening on %s", portStr))
			// start simple server, don't block
			go func() {
				http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
					// allow interrupt
					playing = true
					speaker.Clear()
					doAudio(doneChannel)
					fmt.Fprintln(w, "ok :)")
				})
				http.ListenAndServe(portStr, nil)
			}()
		} else {
			dialog.ShowError(errors.New(fmt.Sprintf("'%s' is not a number, yo", portInput.Text)), w)
		}
	}

	// update time label, do tick stuff
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				currentTime := time.Now()
				timeLabel.SetText(formatTimeDisplay(currentTime))

				// announce new hour
				if autoHour && currentTime.Minute() == 0 && currentTime.Second() == 0 {
					seq, err := mapInputToSequence(getTimeData(currentTime), doneChannel)
					if err == nil {
						// disable speak button for a moment and clear speaker
						playing = true
						speaker.Clear()
						speaker.Play(seq)
					}
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-doneChannel:
				playing = false
			}
		}
	}()

	w.SetContent(container.NewVBox(
		timeLabel,
		widget.NewButton("Speak", func() {
			if !playing {
				playing = true
				doAudio(doneChannel)
			}
		}),
		widget.NewCheck("Announce Hour", func(value bool) {
			autoHour = value
		}),
		widget.NewLabel("Listen Port"),
		portInput,
		startServerButton,
	))

	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}
