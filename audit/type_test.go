package audit

import (
	"testing"
)

func TestUnString(t *testing.T) {
	var tests = []struct {
		input    string
		expected Type
	}{
		{"messageDelete", MessageDelete},
		{"MemberNickname", MemberNickname},
		{"VoiceAudioState", VoiceAudioState},
		{"voiceaudiostate", VoiceAudioState},
		{"invalidtype", Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ans := UnString(tt.input)
			if ans != tt.expected {
				t.Errorf("got %s, want %s", ans.String(), tt.expected.String())
			}
		})
	}
}
