package main

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func Test_validPaths(t *testing.T) {
	type args struct {
		conf *config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validPaths(tt.args.conf); (err != nil) != tt.wantErr {
				t.Errorf("validPaths() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setCredentials(t *testing.T) {
	type args struct {
		exportType string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setCredentials(tt.args.exportType)
		})
	}
}

func Test_exportGratka(t *testing.T) {
	type args struct {
		name string
		conf *config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := exportGratka(tt.args.name, tt.args.conf); (err != nil) != tt.wantErr {
				t.Errorf("exportGratka() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
