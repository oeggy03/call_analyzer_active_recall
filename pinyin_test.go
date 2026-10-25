package main

import "testing"

func TestMarkedPinyinRendersNumberedReadings(t *testing.T) {
	tests := []struct {
		reading string
		want    string
	}{
		{reading: "hei1 wu3", want: "hēi wǔ"},
		{reading: "xia4 xian4", want: "xià xiàn"},
		{reading: "tun2", want: "tún"},
		{reading: "teng2 nuo2", want: "téng nuó"},
		{reading: "li4 yong4 lü4", want: "lì yòng lǜ"},
		{reading: "zhong1 tai2", want: "zhōng tái"},
		{reading: "kai1 ji1", want: "kāi jī"},
		{reading: "da4 ka3", want: "dà kǎ"},
		{reading: "biao1 di4", want: "biāo dì"},
		{reading: "ka3 xing2", want: "kǎ xíng"},
		{reading: "ma5", want: "ma"},
		{reading: "nv3", want: "nǚ"},
		{reading: "nu:3", want: "nǚ"},
		{reading: "lve4", want: "lüè"},
		{reading: "jiu3", want: "jiǔ"},
		{reading: "dui4", want: "duì"},
		{reading: "nǐ hǎo", want: "nǐ hǎo"},
		{reading: "xi1'an1", want: "xī'ān"},
		{reading: "", want: ""},
		{reading: "ci2 0", want: "cí 0"},
	}

	for _, test := range tests {
		if got := markedPinyin(test.reading); got != test.want {
			t.Errorf("markedPinyin(%q) = %q, want %q", test.reading, got, test.want)
		}
	}
}
