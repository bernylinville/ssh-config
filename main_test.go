package main_test

import (
	"os"
	"path"
	"testing"

	Main "github.com/soulteary/ssh-config"
	Cmd "github.com/soulteary/ssh-config/cmd"
	Fn "github.com/soulteary/ssh-config/internal/fn"
)

func TestProcess(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf("TestProcess() error = %v", err)
	}

	jsonContent, err := os.ReadFile(path.Join(pwd, "./testdata/main-test.json"))
	if err != nil {
		t.Errorf("TestProcess() error = %v", err)
	}

	sshContent, err := os.ReadFile(path.Join(pwd, "./testdata/main-test.cfg"))
	if err != nil {
		t.Errorf("TestProcess() error = %v", err)
	}

	yamlContent, err := os.ReadFile(path.Join(pwd, "./testdata/main-test.yaml"))
	if err != nil {
		t.Errorf("TestProcess() error = %v", err)
	}

	tests := []struct {
		name     string
		fileType string
		input    string
		args     Cmd.Args
		want     []byte
	}{
		{
			name:     "YAML to SSH",
			fileType: "YAML",
			input:    string(yamlContent),
			args:     Cmd.Args{ToSSH: true},
			want:     sshContent,
		},
		{
			name:     "JSON to YAML",
			fileType: "JSON",
			input:    string(jsonContent),
			args:     Cmd.Args{ToYAML: true},
			want:     yamlContent,
		},
		{
			name:     "TEXT to JSON",
			fileType: "TEXT",
			input:    string(sshContent),
			args:     Cmd.Args{ToJSON: true},
			want:     jsonContent,
		},
		{
			name:     "Empty",
			fileType: "TEXT",
			input:    "",
			args:     Cmd.Args{},
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Main.Process(tt.fileType, tt.input, tt.args)
			if tt.name == "TEXT to JSON" {
				gotData := Fn.GetJSONData(string(got))
				wantData := Fn.GetJSONData(string(tt.want))
				if len(gotData) != len(wantData) {
					t.Errorf("Process() = %v, want %v", len(gotData), len(wantData))
				}

				for gotItem := range gotData {
					for wantItem := range wantData {
						if gotData[gotItem].Name == wantData[wantItem].Name {
							for key := range gotData[gotItem].Data {
								if gotData[gotItem].Data[key] != wantData[wantItem].Data[key] {
									t.Errorf("Process() = %v, want %v", gotData[gotItem].Data[key], wantData[wantItem].Data[key])
								}
							}
						}
					}
				}

			} else {
				if string(got) != string(tt.want) {
					t.Errorf("Process() = %v, want %v", string(got), string(tt.want))
				}
			}
		})
	}
}
