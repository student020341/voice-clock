package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"github.com/faiface/beep"
	"github.com/faiface/beep/vorbis"
)

func formatTimeDisplay(t time.Time) string {
	ampm := "AM"
	hour := t.Hour()
	if hour > 12 {
		ampm = "PM"
		hour -= 12
	}
	minute := t.Minute()
	minuteS := strconv.Itoa(minute)
	if minute < 10 {
		minuteS = fmt.Sprintf("0%v", minuteS)
	}
	second := t.Second()
	secondS := strconv.Itoa(second)
	if second < 10 {
		secondS = fmt.Sprintf("0%v", secondS)
	}
	return fmt.Sprintf("%v:%v:%v %v", hour, minuteS, secondS, ampm)
}

func getTimeData(t time.Time) []string {
	var data []string

	ampm := "am"
	hour := t.Hour()
	minute := t.Minute()

	// AM / PM
	if hour >= 12 {
		ampm = "pm"
	}
	if hour > 12 {
		hour -= 12
	}

	// announce
	announce := "current-time"
	if minute == 0 {
		data = append(data, "dink-donk")
		announce = "its"
	}
	data = append(data, announce)

	// hour
	if hour == 0 {
		data = append(data, "12")
	} else {
		data = append(data, strconv.Itoa(hour))
	}

	// minute
	if minute > 0 {
		if minute < 10 {
			data = append(data, "oh")
		}
		data = append(data, strconv.Itoa(minute))
	}

	// oclock?
	if minute == 0 {
		data = append(data, "oclock")
	}

	// ampm
	data = append(data, ampm)

	return data
}

func mapInputToResource(input string) *fyne.StaticResource {
	switch val := input; val {
	case "0":
		return resource0Ogg
	case "1":
		return resource1Ogg
	case "2":
		return resource2Ogg
	case "3":
		return resource3Ogg
	case "4":
		return resource4Ogg
	case "5":
		return resource5Ogg
	case "6":
		return resource6Ogg
	case "7":
		return resource7Ogg
	case "8":
		return resource8Ogg
	case "9":
		return resource9Ogg
	case "10":
		return resource10Ogg
	case "11":
		return resource11Ogg
	case "12":
		return resource12Ogg
	case "13":
		return resource13Ogg
	case "14":
		return resource14Ogg
	case "15":
		return resource15Ogg
	case "16":
		return resource16Ogg
	case "17":
		return resource17Ogg
	case "18":
		return resource18Ogg
	case "19":
		return resource19Ogg
	case "20":
		return resource20Ogg
	case "21":
		return resource21Ogg
	case "22":
		return resource22Ogg
	case "23":
		return resource23Ogg
	case "24":
		return resource24Ogg
	case "25":
		return resource25Ogg
	case "26":
		return resource26Ogg
	case "27":
		return resource27Ogg
	case "28":
		return resource28Ogg
	case "29":
		return resource29Ogg
	case "30":
		return resource30Ogg
	case "31":
		return resource31Ogg
	case "32":
		return resource32Ogg
	case "33":
		return resource33Ogg
	case "34":
		return resource34Ogg
	case "35":
		return resource35Ogg
	case "36":
		return resource36Ogg
	case "37":
		return resource37Ogg
	case "38":
		return resource38Ogg
	case "39":
		return resource39Ogg
	case "40":
		return resource40Ogg
	case "41":
		return resource41Ogg
	case "42":
		return resource42Ogg
	case "43":
		return resource43Ogg
	case "44":
		return resource44Ogg
	case "45":
		return resource45Ogg
	case "46":
		return resource46Ogg
	case "47":
		return resource47Ogg
	case "48":
		return resource48Ogg
	case "49":
		return resource49Ogg
	case "50":
		return resource50Ogg
	case "51":
		return resource51Ogg
	case "52":
		return resource52Ogg
	case "53":
		return resource53Ogg
	case "54":
		return resource54Ogg
	case "55":
		return resource55Ogg
	case "56":
		return resource56Ogg
	case "57":
		return resource57Ogg
	case "58":
		return resource58Ogg
	case "59":
		return resource59Ogg
	case "am":
		return resourceAMOgg
	case "current-time":
		return resourceCurrentTimeOgg
	case "dink-donk":
		return resourceDinkDonkOgg
	case "its":
		return resourceItsOgg
	case "oclock":
		return resourceOclockOgg
	case "oh":
		return resourceOhOgg
	case "pm":
		return resourcePMOgg
	default:
		return resourceDinkDonkOgg
	}
}

func getStreamer(resource *fyne.StaticResource) (beep.StreamSeekCloser, beep.Format, error) {
	// TODO update to go 1.16+ and change ioutil NopCloser to io
	byteRead := bytes.NewReader(resource.StaticContent)
	readCloser := ioutil.NopCloser(byteRead)
	return vorbis.Decode(readCloser)
}

func mapInputToSequence(data []string, doneChannel chan struct{}) (beep.Streamer, error) {
	var parts []beep.Streamer

	// map string input data to resource, map resource stream seek closer
	for _, datum := range data {
		dataResource := mapInputToResource(datum)
		dataStream, format, err := getStreamer(dataResource)
		if err != nil {
			return nil, err
		}

		// might seem unnecessary for single use, but avoids an issue with the bundled audio
		buf := beep.NewBuffer(format)
		buf.Append(dataStream)
		dataStream.Close()

		//
		streamer := buf.Streamer(0, buf.Len())

		//
		parts = append(parts, streamer)
	}

	if doneChannel != nil {
		parts = append(parts, beep.Callback(func() {
			doneChannel <- struct{}{}
		}))
	}

	return beep.Seq(parts...), nil
}
