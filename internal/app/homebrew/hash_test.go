package homebrew

import "testing"

func TestHomebrew_buildURLHash(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "correct",
			url:     "https://github.com/NoUseFreak/letitgo/releases/download/0.1.4/darwin_amd64.zip",
			want:    "33f779d32301c252d63bdf2099ffcea370b58fd37a075c670c2ee2359670e21d",
			wantErr: false,
		},
		{
			name:    "non-existing",
			url:     "https://github.com/NoUseFreak/letitgo/releases/download/0.0.0/darwin_amd64.zip",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Homebrew{}
			got, err := h.buildURLHash(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Homebrew.buildURLHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Homebrew.buildURLHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
