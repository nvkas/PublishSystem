package test

import (
	"Publish_System/services"
	"reflect"
	"testing"
)

func TestDealWithData(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
		//正常情况
		{"1",args{data:[]byte(`{
  "ref": "refs/heads/master",
  "before": "fc5075afcb02fcf9d356d4b1a23348506a522947",
  "after": "fc5075afcb02fcf9d356d4b1a23348506a522947",
  "compare_url": "",
  "commits": [
    {
      "id": "fc5075afcb02fcf9d356d4b1a23348506a522947",
      "message": "123\n",
      "url": "http://localhost:3000/ch/git_hook_test/commit/fc5075afcb02fcf9d356d4b1a23348506a522947",
      "author": {
        "name": "miracle",
        "email": "asd471331055@vip.qq.com",
        "username": "ch"
      },
      "committer": {
        "name": "miracle",
        "email": "asd471331055@vip.qq.com",
        "username": "ch"
      },
      "added": [
        "test.txt"
      ],
      "removed": [],
      "modified": [],
      "timestamp": "0001-01-01T00:00:00Z"
    }
  ],
  "repository": {
    "id": 1,
    "owner": {
      "id": 1,
      "username": "ch",
      "login": "ch",
      "full_name": "",
      "email": "asd471331055@vip.qq.com",
      "avatar_url": "https://secure.gravatar.com/avatar/bdafa8c50b079c7bf7a647de881e3b0b?d=identicon"
    },
    "name": "git_hook_test",
    "full_name": "ch/git_hook_test",
    "description": "",
    "private": false,
    "fork": false,
    "parent": null,
    "empty": false,
    "mirror": false,
    "size": 12288,
    "html_url": "http://localhost:3000/ch/git_hook_test",
    "ssh_url": "miracle@localhost:ch/git_hook_test.git",
    "clone_url": "http://localhost:3000/ch/git_hook_test.git",
    "website": "",
    "stars_count": 0,
    "forks_count": 0,
    "watchers_count": 1,
    "open_issues_count": 0,
    "default_branch": "master",
    "created_at": "2019-02-17T16:52:44+08:00",
    "updated_at": "2019-02-17T17:00:23+08:00"
  },
  "pusher": {
    "id": 1,
    "username": "ch",
    "login": "ch",
    "full_name": "",
    "email": "asd471331055@vip.qq.com",
    "avatar_url": "https://secure.gravatar.com/avatar/bdafa8c50b079c7bf7a647de881e3b0b?d=identicon"
  },
  "sender": {
    "id": 1,
    "username": "ch",
    "login": "ch",
    "full_name": "",
    "email": "asd471331055@vip.qq.com",
    "avatar_url": "https://secure.gravatar.com/avatar/bdafa8c50b079c7bf7a647de881e3b0b?d=identicon"
  }
}`)}, map[string]interface{}{"repository_name":"ch/git_hook_test","pusher":"ch"}},

		//异常情况
		{"2",args{data:[]byte(`{"test":"test"}`)}, map[string]interface{}{"repository_name":nil,"pusher":nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.DealWithData(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DealWithData() = %v, want %v", got, tt.want)
			}
		})
	}
}
