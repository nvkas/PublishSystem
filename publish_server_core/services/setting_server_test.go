package service

import (
	"publish_server_core/datamodels"
	"testing"
)

func TestExecAllCommand(t *testing.T) {
	type args struct {
		project datamodels.Project
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test1",args{datamodels.Project{User:datamodels.User{ProjectWorkPath:"/home/go/src/"},WarehouseName:"shehao/publish_test",GitAddress:"http://server.spacej.tech:11300/shehao/publish_test.git"}},false},
		{"test2",args{datamodels.Project{WarehouseName:"",GitAddress:""}},true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExecAllCommand(tt.args.project); (err != nil) != tt.wantErr {
				t.Errorf("ExecAllCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
