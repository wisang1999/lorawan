package lorawan

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMHDR(t *testing.T) {
	Convey("Given an empty MHDR", t, func() {
		var h MHDR
		Convey("Then MarshalBinary returns []byte{0}", func() {
			b, err := h.MarshalBinary()
			So(err, ShouldBeNil)
			So(b, ShouldResemble, []byte{0})
		})

		Convey("Given MType=Proprietary, Major=LoRaWANR1", func() {
			h.MType = Proprietary
			h.Major = LoRaWANR1
			Convey("Then MarshalBinary returns []byte{224}", func() {
				b, err := h.MarshalBinary()
				So(err, ShouldBeNil)
				So(b, ShouldResemble, []byte{224})
			})
		})

		Convey("Given a slice []byte{224}", func() {
			b := []byte{224}
			Convey("Then UnmarshalBinary returns a MHDR with MType=Proprietary, Major=LoRaWANR1", func() {
				err := h.UnmarshalBinary(b)
				So(err, ShouldBeNil)
				So(h, ShouldResemble, MHDR{MType: Proprietary, Major: LoRaWANR1})
			})
		})
	})
}

func TestPHYPayload(t *testing.T) {
	Convey("Given an empty PHYPayload with empty MACPayload", t, func() {
		p := PHYPayload{MACPayload: &MACPayload{}}

		Convey("Then MarshalBinary returns []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}", func() {
			b, err := p.MarshalBinary()
			So(err, ShouldBeNil)
			So(b, ShouldResemble, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		})

		Convey("Given MHDR(MType=JoinAccept, Major=LoRaWANR1), MACPayload(FHDR(DevAddr=[4]byte{1, 2, 3, 4})), MIC=[4]byte{4, 3, 2, 1}", func() {
			p.MHDR.MType = JoinAccept
			p.MHDR.Major = LoRaWANR1
			p.MACPayload = &MACPayload{
				FHDR: FHDR{
					DevAddr: DevAddr([4]byte{1, 2, 3, 4}),
				},
			}
			p.MIC = [4]byte{4, 3, 2, 1}

			Convey("Given the NwkSKey []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}", func() {
				nwkSKey := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

				Convey("Then ValidateMIC returns false", func() {
					v, err := p.ValidateMIC(nwkSKey)
					So(err, ShouldBeNil)
					So(v, ShouldBeFalse)

				})

				Convey("calculateMIC returns []byte{0x7b, 0x14, 0x9d, 0x8a}", func() {
					mic, err := p.calculateMIC(nwkSKey)
					So(err, ShouldBeNil)
					So(mic, ShouldResemble, []byte{0x7b, 0x14, 0x9d, 0x8a})
				})

				Convey("Given the MIC is []byte{0x7b, 0x14, 0x9d, 0x8a}", func() {
					p.MIC = [4]byte{0x7b, 0x14, 0x9d, 0x8a}

					Convey("Then ValidateMIC returns true", func() {
						v, err := p.ValidateMIC(nwkSKey)
						So(err, ShouldBeNil)
						So(v, ShouldBeTrue)
					})
				})
			})

			Convey("Then MarshalBinary returns []byte{32, 1, 2, 3, 4, 0, 0, 0, 4, 3, 2, 1}", func() {
				b, err := p.MarshalBinary()
				So(err, ShouldBeNil)
				So(b, ShouldResemble, []byte{32, 1, 2, 3, 4, 0, 0, 0, 4, 3, 2, 1})
			})
		})

		Convey("Given the slice of bytes with an invalid size", func() {
			b := make([]byte, 4)
			Convey("Then UnmarshalBinary returns an error", func() {
				err := p.UnmarshalBinary(b)
				So(err, ShouldResemble, errors.New("lorawan: at least 5 bytes needed to decode PHYPayload"))
			})
		})

		Convey("Given the slice []byte{32, 1, 2, 3, 4, 0, 0, 0, 4, 3, 2, 1}", func() {
			b := []byte{64, 1, 2, 3, 4, 0, 0, 0, 4, 3, 2, 1}
			Convey("Then UnmarshalBinary does not return an error", func() {
				err := p.UnmarshalBinary(b)
				So(err, ShouldBeNil)

				Convey("Then MHDR=(MType=UnconfirmedDataUp, Major=LoRaWANR1)", func() {
					So(p.MHDR, ShouldResemble, MHDR{MType: UnconfirmedDataUp, Major: LoRaWANR1})
				})
				Convey("Then MACPayload(FHDR(DevAddr=[4]byte{1, 2, 3, 4}))", func() {
					So(p.MACPayload, ShouldResemble, &MACPayload{FHDR: FHDR{DevAddr: DevAddr([4]byte{1, 2, 3, 4})}})
				})
				Convey("Then MIC=[4]byte{4, 3, 2, 1}", func() {
					So(p.MIC, ShouldResemble, [4]byte{4, 3, 2, 1})
				})
			})
		})
	})
}

func TestPHYPayloadJoinRequest(t *testing.T) {
	Convey("Given an empty PHYPayload with empty JoinRequestPayload", t, func() {
		p := PHYPayload{MACPayload: &JoinRequestPayload{}}
		Convey("Then MarshalBinary returns []byte with 23 0x00 bytes", func() {
			exp := make([]byte, 23)
			b, err := p.MarshalBinary()
			So(err, ShouldBeNil)
			So(b, ShouldResemble, exp)
		})

		Convey("Given MHDR=(MType=JoinRequest, Major=LoRaWANR1), MACPayload=JoinRequestPayload(AppEUI=[8]byte{1, 1, 1, 1, 1, 1, 1, 1}, DevEUI=[8]byte{2, 2, 2, 2, 2, 2, 2, 2} and DevNonce=[2]byte{3, 3}), MIC=[4]byte{4, 5, 6, 7}", func() {
			p.MHDR = MHDR{MType: JoinRequest, Major: LoRaWANR1}
			p.MACPayload = &JoinRequestPayload{
				AppEUI:   [8]byte{1, 1, 1, 1, 1, 1, 1, 1},
				DevEUI:   [8]byte{2, 2, 2, 2, 2, 2, 2, 2},
				DevNonce: [2]byte{3, 3},
			}
			p.MIC = [4]byte{4, 5, 6, 7}

			Convey("Then MarshalBinary returns []byte{0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3,  4, 5, 6, 7}", func() {
				b, err := p.MarshalBinary()
				So(err, ShouldBeNil)
				So(b, ShouldResemble, []byte{0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 4, 5, 6, 7})
			})
		})

		Convey("Given the slice []byte{0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 4, 5, 6, 7}", func() {
			b := []byte{0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 4, 5, 6, 7}

			Convey("Then UnmarshalBinary does not return an error", func() {
				err := p.UnmarshalBinary(b)
				So(err, ShouldBeNil)

				Convey("Then MHDR=(MType=JoinRequest, Major=LoRaWANR1)", func() {
					So(p.MHDR, ShouldResemble, MHDR{MType: JoinRequest, Major: LoRaWANR1})
				})
				Convey("Then MACPayload=JoinRequestPayload(AppEUI=[8]byte{1, 1, 1, 1, 1, 1, 1, 1}, DevEUI=[8]byte{2, 2, 2, 2, 2, 2, 2, 2} and DevNonce=[2]byte{3, 3})", func() {
					So(p.MACPayload, ShouldResemble, &JoinRequestPayload{
						AppEUI:   [8]byte{1, 1, 1, 1, 1, 1, 1, 1},
						DevEUI:   [8]byte{2, 2, 2, 2, 2, 2, 2, 2},
						DevNonce: [2]byte{3, 3},
					})
				})
				Convey("MIC=[4]byte{4, 5, 6, 7}", func() {
					So(p.MIC, ShouldResemble, [4]byte{4, 5, 6, 7})
				})
			})
		})
	})
}

func TestPHYPayloadJoinAccept(t *testing.T) {
	Convey("Given an empty PHYPayload with empty JoinAcceptPayload", t, func() {
		p := PHYPayload{MACPayload: &JoinAcceptPayload{}}
		Convey("Then MarshalBinary returns []byte with 17 0x00", func() {
			exp := make([]byte, 17)
			b, err := p.MarshalBinary()
			So(err, ShouldBeNil)
			So(b, ShouldResemble, exp)
		})

		Convey("Given MHDR=(MType=JoinAccept, Major=LoRaWANR1), MACPayload=JoinAcceptPayload(AppNonce=[3]byte{1, 1, 1}, NetID=[3]byte{2, 2, 2}, DevAddr=[4]byte{1, 2, 3, 4}, DLSettings=(RX2DataRate=1, RX1DRoffset=2), RXDelay=7), MIC=[4]byte{8, 9 , 10, 11}", func() {
			p.MHDR = MHDR{MType: JoinAccept, Major: LoRaWANR1}
			p.MACPayload = &JoinAcceptPayload{
				AppNonce:   [3]byte{1, 1, 1},
				NetID:      [3]byte{2, 2, 2},
				DevAddr:    DevAddr([4]byte{1, 2, 3, 4}),
				DLSettings: DLsettings{RX2DataRate: 1, RX1DRoffset: 2},
				RXDelay:    7,
			}
			p.MIC = [4]byte{8, 9, 10, 11}

			// no encryption and invalid MIC
			Convey("Then MarshalBinary returns []byte{32, 1, 1, 1, 2, 2, 2, 1, 2, 3, 4, 33, 7, 8, 9, 10, 11}", func() {
				b, err := p.MarshalBinary()
				So(err, ShouldBeNil)
				So(b, ShouldResemble, []byte{32, 1, 1, 1, 2, 2, 2, 1, 2, 3, 4, 33, 7, 8, 9, 10, 11})
			})

			Convey("Given AppKey []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}", func() {
				appKey := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

				Convey("Then ValidateMIC returns false", func() {
					v, err := p.ValidateMIC(appKey)
					So(err, ShouldBeNil)
					So(v, ShouldBeFalse)
				})

				Convey("Given SetMIC is called", func() {
					err := p.SetMIC(appKey)
					So(err, ShouldBeNil)

					// todo: validate if this mic is actually valid
					Convey("Then SetMIC sets the MIC to [4]byte{27, 61, 112, 162}", func() {
						So(p.MIC, ShouldResemble, [4]byte{27, 61, 112, 162})
					})

					Convey("Given EncryptMACPayload is called", func() {
						err := p.EncryptMACPayload(appKey)
						So(err, ShouldBeNil)

						Convey("Then MACPayload should be of type *DataPayload", func() {
							dp, ok := p.MACPayload.(*DataPayload)
							So(ok, ShouldBeTrue)

							Convey("Then DataPayload Bytes equals []byte{234, 201, 51, 48, 151, 50, 166, 172, 136, 105, 14, 81, 71, 167, 87, 205}", func() {
								So(dp.Bytes, ShouldResemble, []byte{234, 201, 51, 48, 151, 50, 166, 172, 136, 105, 14, 81, 71, 167, 87, 205})
							})

							Convey("Then marshalBinary returns []byte{32, 234, 201, 51, 48, 151, 50, 166, 172, 136, 105, 14, 81, 71, 167, 87, 205, 27, 61, 112, 162}", func() {
								b, err := p.MarshalBinary()
								So(err, ShouldBeNil)
								So(b, ShouldResemble, []byte{32, 234, 201, 51, 48, 151, 50, 166, 172, 136, 105, 14, 81, 71, 167, 87, 205, 27, 61, 112, 162})
							})
						})
					})
				})

				Convey("Given the MIC is [4]byte{27, 61, 112, 162}", func() {
					p.MIC = [4]byte{27, 61, 112, 162}
					Convey("Then ValidateMIC returns true", func() {
						v, err := p.ValidateMIC(appKey)
						So(err, ShouldBeNil)
						So(v, ShouldBeTrue)
					})
				})
			})
		})

		Convey("Given the slice []byte{32, 234, 201, 51, 48, 151, 50, 166, 172, 136, 105, 14, 81, 71, 167, 87, 205, 27, 61, 112, 162}", func() {
			b := []byte{32, 234, 201, 51, 48, 151, 50, 166, 172, 136, 105, 14, 81, 71, 167, 87, 205, 27, 61, 112, 162}

			Convey("Then UnmarshalBinary does not return an error", func() {
				err := p.UnmarshalBinary(b)
				So(err, ShouldBeNil)

				Convey("Then MHDR=(MType=JoinAccept, Major=LoRaWANR1)", func() {
					So(p.MHDR, ShouldResemble, MHDR{MType: JoinAccept, Major: LoRaWANR1})
				})

				Convey("Then MACPayload is of type *DataPayload", func() {
					dp, ok := p.MACPayload.(*DataPayload)
					So(ok, ShouldBeTrue)

					Convey("Then Bytes equals []byte{234, 201, 51, 48, 151, 50, 166, 172, 136, 105, 14, 81, 71, 167, 87, 205}", func() {
						So(dp.Bytes, ShouldResemble, []byte{234, 201, 51, 48, 151, 50, 166, 172, 136, 105, 14, 81, 71, 167, 87, 205})
					})

					Convey("Given AppKey []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}", func() {
						appKey := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

						Convey("Given DecryptMACPayload is called", func() {
							err := p.DecryptMACPayload(appKey)
							So(err, ShouldBeNil)

							Convey("Then MACPayload is of type *JoinAcceptPayload", func() {
								ja, ok := p.MACPayload.(*JoinAcceptPayload)
								So(ok, ShouldBeTrue)

								Convey("Then MACPayload=JoinAcceptPayload(AppNonce=[3]byte{1, 1, 1}, NetID=[3]byte{2, 2, 2}, DevAddr=[4]byte{1, 2, 3, 4}, DLSettings=(RX2DataRate=1, RX1DRoffset=2), RXDelay=7", func() {
									So(ja, ShouldResemble, &JoinAcceptPayload{
										AppNonce:   [3]byte{1, 1, 1},
										NetID:      [3]byte{2, 2, 2},
										DevAddr:    DevAddr([4]byte{1, 2, 3, 4}),
										DLSettings: DLsettings{RX2DataRate: 1, RX1DRoffset: 2},
										RXDelay:    7,
									})
								})

								Convey("Then ValidateMIC returns true", func() {
									v, err := p.ValidateMIC(appKey)
									So(err, ShouldBeNil)
									So(v, ShouldBeTrue)
								})
							})
						})
					})
				})

				Convey("Then MIC=[4]byte{27, 61, 112, 162}", func() {
					So(p.MIC, ShouldResemble, [4]byte{27, 61, 112, 162})
				})
			})
		})
	})
}

func ExampleNewPHYPayload() {
	nwkSKey := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	appSKey := []byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	// uplink and downlink messages are (un)marshalled and encrypted / decrypted
	// differently
	uplink := true

	macPayload := NewMACPayload(uplink)
	macPayload.FHDR = FHDR{
		DevAddr: DevAddr([4]byte{1, 2, 3, 4}),
		FCtrl: FCtrl{
			ADR:       false,
			ADRACKReq: false,
			ACK:       false,
		},
		FCnt:  0,
		FOpts: []MACCommand{}, // you can leave this out when there is no MAC command to send
	}
	macPayload.FPort = 10
	macPayload.FRMPayload = []Payload{&DataPayload{Bytes: []byte{1, 2, 3, 4}}}

	if err := macPayload.EncryptFRMPayload(appSKey); err != nil {
		panic(err)
	}

	payload := NewPHYPayload(uplink)
	payload.MHDR = MHDR{
		MType: ConfirmedDataUp,
		Major: LoRaWANR1,
	}
	payload.MACPayload = macPayload

	if err := payload.SetMIC(nwkSKey); err != nil {
		panic(err)
	}

	bytes, err := payload.MarshalBinary()
	if err != nil {
		panic(err)
	}

	fmt.Println(bytes)

	// Output:
	// [128 1 2 3 4 0 0 0 10 59 85 197 241 187 130 109 49]
}

func ExampleNewPHYPayload_joinRequest() {
	uplink := true
	appKey := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	payload := NewPHYPayload(uplink)
	payload.MHDR = MHDR{
		MType: JoinRequest,
		Major: LoRaWANR1,
	}
	payload.MACPayload = &JoinRequestPayload{
		AppEUI:   [8]byte{1, 1, 1, 1, 1, 1, 1, 1},
		DevEUI:   [8]byte{2, 2, 2, 2, 2, 2, 2, 2},
		DevNonce: [2]byte{3, 3},
	}

	if err := payload.SetMIC(appKey); err != nil {
		panic(err)
	}

	bytes, err := payload.MarshalBinary()
	if err != nil {
		panic(err)
	}

	fmt.Println(bytes)

	// Output:
	// [0 1 1 1 1 1 1 1 1 2 2 2 2 2 2 2 2 3 3 9 185 123 50]
}

func ExampleNewPHYPayload_joinAcceptSend() {
	uplink := false
	appKey := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	payload := NewPHYPayload(uplink)
	payload.MHDR = MHDR{
		MType: JoinAccept,
		Major: LoRaWANR1,
	}
	payload.MACPayload = &JoinAcceptPayload{
		AppNonce:   [3]byte{1, 1, 1},
		NetID:      [3]byte{2, 2, 2},
		DevAddr:    DevAddr([4]byte{1, 2, 3, 4}),
		DLSettings: DLsettings{RX2DataRate: 0, RX1DRoffset: 0},
		RXDelay:    0,
	}
	// set the MIC before encryption
	if err := payload.SetMIC(appKey); err != nil {
		panic(err)
	}
	if err := payload.EncryptMACPayload(appKey); err != nil {
		panic(err)
	}

	bytes, err := payload.MarshalBinary()
	if err != nil {
		panic(err)
	}

	fmt.Println(bytes)

	// Output:
	// [32 64 253 162 88 11 45 30 206 20 214 140 149 191 32 154 238 227 185 68 130]
}

func ExampleNewPHYPayload_joinAcceptReceive() {
	uplink := false
	appKey := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	bytes := []byte{32, 171, 84, 244, 227, 34, 30, 148, 118, 211, 1, 33, 90, 24, 50, 81, 139, 128, 229, 23, 154}

	payload := NewPHYPayload(uplink)
	if err := payload.UnmarshalBinary(bytes); err != nil {
		panic(err)
	}

	if err := payload.DecryptMACPayload(appKey); err != nil {
		panic(err)
	}

	_, ok := payload.MACPayload.(*JoinAcceptPayload)
	if !ok {
		panic("*JoinAcceptPayload expected")
	}

	v, err := payload.ValidateMIC(appKey)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)

	// Output:
	// true
}
