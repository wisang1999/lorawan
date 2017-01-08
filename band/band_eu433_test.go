package band

import (
	"testing"

	"github.com/brocaar/lorawan"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEU433Band(t *testing.T) {
	Convey("Given the EU 433 band is selected with repeaterCompatible=true", t, func() {
		band, err := GetConfig(EU_433, true, lorawan.DwellTimeNoLimit)
		So(err, ShouldBeNil)

		Convey("Then GetRX1Channel returns the uplink channel", func() {
			for i := 0; i < 3; i++ {
				rx1Chan := band.GetRX1Channel(i)
				So(rx1Chan, ShouldEqual, i)
			}
		})

		Convey("Then GetRX1Frequency returns the uplink frequency", func() {
			for _, c := range band.DownlinkChannels {
				freq, err := band.GetRX1Frequency(c.Frequency)
				So(err, ShouldBeNil)
				So(freq, ShouldEqual, c.Frequency)
			}
		})

		Convey("Then the max payload size (N) is 222 for DR4", func() {
			So(band.MaxPayloadSize[4].N, ShouldEqual, 222)
		})
	})

	Convey("Given the EU 433 band is selected with repeaterCompatible=false", t, func() {
		band, err := GetConfig(EU_433, false, lorawan.DwellTimeNoLimit)
		So(err, ShouldBeNil)

		Convey("Then the max payload size (N) is 242 for DR4", func() {
			So(band.MaxPayloadSize[4].N, ShouldEqual, 242)
		})
	})
}
