package wavdata

import "math"

type Format struct {
	Encoding Encoding
	Bits     uint16
}

type Encoding uint16

// TODO: normalize names
const (
	/* Windows WAVE File Encoding Tags */
	Unknown                 = Encoding(0x0000) // Unknown Format
	PCM                     = Encoding(0x0001) // PCM
	ADPCM                   = Encoding(0x0002) // Microsoft ADPCM Format
	IEEE_Float              = Encoding(0x0003) // IEEE Float
	VSELP                   = Encoding(0x0004) // Compaq Computer's VSELP
	IBM_CSVD                = Encoding(0x0005) // IBM CVSD
	ALAW                    = Encoding(0x0006) // ALAW
	MULAW                   = Encoding(0x0007) // MULAW
	OKI_ADPCM               = Encoding(0x0010) // OKI ADPCM
	DVI_ADPCM               = Encoding(0x0011) // Intel's DVI ADPCM
	MEDIASPACE_ADPCM        = Encoding(0x0012) // Videologic's MediaSpace ADPCM
	SIERRA_ADPCM            = Encoding(0x0013) // Sierra ADPCM
	G723_ADPCM              = Encoding(0x0014) // G.723 ADPCM
	DIGISTD                 = Encoding(0x0015) // DSP Solution's DIGISTD
	DIGIFIX                 = Encoding(0x0016) // DSP Solution's DIGIFIX
	DIALOGIC_OKI_ADPCM      = Encoding(0x0017) // Dialogic OKI ADPCM
	MEDIAVISION_ADPCM       = Encoding(0x0018) // MediaVision ADPCM
	CU_CODEC                = Encoding(0x0019) // HP CU
	YAMAHA_ADPCM            = Encoding(0x0020) // Yamaha ADPCM
	SONARC                  = Encoding(0x0021) // Speech Compression's Sonarc
	TRUESPEECH              = Encoding(0x0022) // DSP Group's True Speech
	ECHOSC1                 = Encoding(0x0023) // Echo Speech's EchoSC1
	AUDIOFILE_AF36          = Encoding(0x0024) // Audiofile AF36
	APTX                    = Encoding(0x0025) // APTX
	AUDIOFILE_AF10          = Encoding(0x0026) // AudioFile AF10
	PROSODY_1612            = Encoding(0x0027) // Prosody 1612
	LRC                     = Encoding(0x0028) // LRC
	AC2                     = Encoding(0x0030) // Dolby AC2
	GSM610                  = Encoding(0x0031) // GSM610
	MSNAUDIO                = Encoding(0x0032) // MSNAudio
	ANTEX_ADPCME            = Encoding(0x0033) // Antex ADPCME
	CONTROL_RES_VQLPC       = Encoding(0x0034) // Control Res VQLPC
	DIGIREAL                = Encoding(0x0035) // Digireal
	DIGIADPCM               = Encoding(0x0036) // DigiADPCM
	CONTROL_RES_CR10        = Encoding(0x0037) // Control Res CR10
	VBXADPCM                = Encoding(0x0038) // NMS VBXADPCM
	ROLAND_RDAC             = Encoding(0x0039) // Roland RDAC
	ECHOSC3                 = Encoding(0x003A) // EchoSC3
	ROCKWELL_ADPCM          = Encoding(0x003B) // Rockwell ADPCM
	ROCKWELL_DIGITALK       = Encoding(0x003C) // Rockwell Digit LK
	XEBEC                   = Encoding(0x003D) // Xebec
	G721_ADPCM              = Encoding(0x0040) // Antex Electronics G.721
	G728_CELP               = Encoding(0x0041) // G.728 CELP
	MSG723                  = Encoding(0x0042) // MSG723
	MPEG                    = Encoding(0x0050) // MPEG Layer 1,2
	RT24                    = Encoding(0x0051) // RT24
	PAC                     = Encoding(0x0051) // PAC
	MPEGLAYER3              = Encoding(0x0055) // MPEG Layer 3
	CIRRUS                  = Encoding(0x0059) // Cirrus
	ESPCM                   = Encoding(0x0061) // ESPCM
	VOXWARE                 = Encoding(0x0062) // Voxware (obsolete)
	CANOPUS_ATRAC           = Encoding(0x0063) // Canopus Atrac
	G726_ADPCM              = Encoding(0x0064) // G.726 ADPCM
	G722_ADPCM              = Encoding(0x0065) // G.722 ADPCM
	DSAT                    = Encoding(0x0066) // DSAT
	DSAT_DISPLAY            = Encoding(0x0067) // DSAT Display
	VOXWARE_BYTE_ALIGNED    = Encoding(0x0069) // Voxware Byte Aligned (obsolete)
	VOXWARE_AC8             = Encoding(0x0070) // Voxware AC8 (obsolete)
	VOXWARE_AC10            = Encoding(0x0071) // Voxware AC10 (obsolete)
	VOXWARE_AC16            = Encoding(0x0072) // Voxware AC16 (obsolete)
	VOXWARE_AC20            = Encoding(0x0073) // Voxware AC20 (obsolete)
	VOXWARE_RT24            = Encoding(0x0074) // Voxware MetaVoice (obsolete)
	VOXWARE_RT29            = Encoding(0x0075) // Voxware MetaSound (obsolete)
	VOXWARE_RT29HW          = Encoding(0x0076) // Voxware RT29HW (obsolete)
	VOXWARE_VR12            = Encoding(0x0077) // Voxware VR12 (obsolete)
	VOXWARE_VR18            = Encoding(0x0078) // Voxware VR18 (obsolete)
	VOXWARE_TQ40            = Encoding(0x0079) // Voxware TQ40 (obsolete)
	SOFTSOUND               = Encoding(0x0080) // Softsound
	VOXWARE_TQ60            = Encoding(0x0081) // Voxware TQ60 (obsolete)
	MSRT24                  = Encoding(0x0082) // MSRT24
	G729A                   = Encoding(0x0083) // G.729A
	MVI_MV12                = Encoding(0x0084) // MVI MV12
	DF_G726                 = Encoding(0x0085) // DF G.726
	DF_GSM610               = Encoding(0x0086) // DF GSM610
	ISIAUDIO                = Encoding(0x0088) // ISIAudio
	ONLIVE                  = Encoding(0x0089) // Onlive
	SBC24                   = Encoding(0x0091) // SBC24
	DOLBY_AC3_SPDIF         = Encoding(0x0092) // Dolby AC3 SPDIF
	ZYXEL_ADPCM             = Encoding(0x0097) // ZyXEL ADPCM
	PHILIPS_LPCBB           = Encoding(0x0098) // Philips LPCBB
	PACKED                  = Encoding(0x0099) // Packed
	RHETOREX_ADPCM          = Encoding(0x0100) // Rhetorex ADPCM
	IRAT                    = Encoding(0x0101) // BeCubed Software's IRAT
	VIVO_G723               = Encoding(0x0111) // Vivo G.723
	VIVO_SIREN              = Encoding(0x0112) // Vivo Siren
	DIGITAL_G723            = Encoding(0x0123) // Digital G.723
	CREATIVE_ADPCM          = Encoding(0x0200) // Creative ADPCM
	CREATIVE_FASTSPEECH8    = Encoding(0x0202) // Creative FastSpeech8
	CREATIVE_FASTSPEECH10   = Encoding(0x0203) // Creative FastSpeech10
	QUARTERDECK             = Encoding(0x0220) // Quarterdeck
	FM_TOWNS_SND            = Encoding(0x0300) // FM Towns Snd
	BTV_DIGITAL             = Encoding(0x0400) // BTV Digital
	VME_VMPCM               = Encoding(0x0680) // VME VMPCM
	OLIGSM                  = Encoding(0x1000) // OLIGSM
	OLIADPCM                = Encoding(0x1001) // OLIADPCM
	OLICELP                 = Encoding(0x1002) // OLICELP
	OLISBC                  = Encoding(0x1003) // OLISBC
	OLIOPR                  = Encoding(0x1004) // OLIOPR
	LH_CODEC                = Encoding(0x1100) // LH Codec
	NORRIS                  = Encoding(0x1400) // Norris
	ISIAUDIO2               = Encoding(0x1401) // ISIAudio
	SOUNDSPACE_MUSICOMPRESS = Encoding(0x1500) // Soundspace Music Compression
	DVM                     = Encoding(0x2000) // DVM
	EXTENSIBLE              = Encoding(0xFFFE) // SubFormat
	DEVELOPMENT             = Encoding(0xFFFF) // Development
)

type Codec struct {
	ReadF32 func(src []byte, dst []float32, sampleCount int) (advance int)
}

var Codecs = map[Format]Codec{
	Format{PCM, 8}: {
		ReadF32: func(src []byte, dst []float32, sampleCount int) int {
			h := 0
			for k := 0; k < sampleCount; k++ {
				v := uint8(src[h])
				dst[k] = float32(v)/128.0 - 1.0
				h += 1
			}
			return h
		},
	},

	Format{PCM, 16}: {
		ReadF32: func(src []byte, dst []float32, sampleCount int) int {
			h := 0
			for k := 0; k < sampleCount; k++ {
				v := int16(src[h]) | int16(src[h+1])<<8
				dst[k] = float32(v) / float32(0x8000)
				h += 16
			}
			return h
		},
	},

	Format{PCM, 24}: {
		ReadF32: func(src []byte, dst []float32, sampleCount int) int {
			h := 0
			for k := 0; k < sampleCount; k++ {
				v := int32(src[h]) | int32(src[h+1])<<8 | int32(src[h+2])<<16
				if v&0x800000 != 0 {
					v |= ^0xffffff
				}
				dst[k] = float32(v) / float32(0x800000)
				h += 3
			}
			return h
		},
	},

	Format{PCM, 32}: {
		ReadF32: func(src []byte, dst []float32, sampleCount int) int {
			h := 0
			for k := 0; k < sampleCount; k++ {
				v := int32(src[h]) | int32(src[h+1])<<8 | int32(src[h+2])<<16 | int32(src[h+3])<<24
				dst[k] = float32(v) / float32(0x80000000)
				h += 4
			}
			return h
		},
	},

	Format{IEEE_Float, 32}: {
		ReadF32: func(src []byte, dst []float32, sampleCount int) int {
			h := 0
			for k := 0; k < sampleCount; k++ {
				v := uint32(src[h]) | uint32(src[h+1])<<8 | uint32(src[h+2])<<16 | uint32(src[h+3])<<24
				dst[k] = math.Float32frombits(v)
				h += 4
			}
			return h
		},
	},
}
