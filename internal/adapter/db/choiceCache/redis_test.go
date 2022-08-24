package choiceCache

import (
	"errors"
	"testing"
	"time"

	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

var (
	redisServer *miniredis.Miniredis
	redisClient *redis.Client
)

func TestSet(t *testing.T) {
	setUp()
	defer teardown()
	repo := NewChoiceCache(redisClient, logging.GetLogger("debug"))
	type args struct {
		choiceTitle string
		voteTitle   string
		count       int
		expire      time.Duration
	}
	type mockCall func()
	testCases := []struct {
		title   string
		input   args
		isError bool
		mock    mockCall
	}{
		{
			title:   "Success Set()",
			input:   args{choiceTitle: "choice", voteTitle: "vote", count: 1, expire: 1 * time.Second},
			mock:    func() {},
			isError: false,
		},
		{
			title: "reddis internal error and Set() return error",
			input: args{choiceTitle: "choice", voteTitle: "vote", count: 1, expire: 1 * time.Second},
			mock: func() {
				redisServer.SetError("interanl redis error")
			},
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			err := repo.Set(test.input.voteTitle, test.input.choiceTitle, test.input.count, test.input.expire)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	setUp()
	defer teardown()

	repo := NewChoiceCache(redisClient, logging.GetLogger("debug"))
	type mockCall func(voteTitle string, choiceTitle string, count int, expire time.Duration) error
	type args struct {
		choiceTitle string
		voteTitle   string
		count       int
		expire      time.Duration
	}
	testCases := []struct {
		title   string
		input   args
		isError bool
		mock    mockCall
		want    int
	}{
		{
			title:   "Get should find title and return count",
			input:   args{choiceTitle: "choice", voteTitle: "vote", count: 1, expire: 1 * time.Second},
			isError: false,
			mock: func(voteTitle string, choiceTitle string, count int, expire time.Duration) error {
				return repo.Set(voteTitle, choiceTitle, count, expire)
			},
			want: 1,
		},
		{
			title:   "Get doens't find key and should return error",
			input:   args{choiceTitle: "wrong key", voteTitle: "wrong key", count: 1, expire: 1 * time.Second},
			isError: true,
			mock: func(voteTitle string, choiceTitle string, count int, expire time.Duration) error {
				return repo.Set("some vote", "some choice", count, expire)
			},
			want: -1,
		},
		{
			title:   "reddis internal error and Get return error ",
			input:   args{choiceTitle: "choice", voteTitle: "vote", count: 1, expire: 1 * time.Second},
			isError: true,
			mock: func(voteTitle string, choiceTitle string, count int, expire time.Duration) error {
				redisServer.SetError("interanl redis error")
				return errors.New("internal error")
			},
			want: -1,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			_ = test.mock(test.input.voteTitle, test.input.choiceTitle, test.input.count, test.input.expire)
			got, err := repo.Get(test.input.voteTitle, test.input.choiceTitle)
			if !test.isError {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			} else {
				assert.Error(t, err)
			}

		})
	}

}

func setUp() {
	redisServer = mockRedis()
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
}

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()

	if err != nil {
		panic(err)
	}

	return s
}

func teardown() {
	redisServer.Close()
}
