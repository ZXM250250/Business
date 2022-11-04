package tools

import (
	"github.com/garyburd/redigo/redis"
)

func RedisConfPipline(pip ...func(c redis.Conn)) error {
	c, err := redis.Dial("tcp", "101.43.4.210:6379")
	if err != nil {
		return err
	}
	defer c.Close()
	for _, f := range pip {
		f(c)
	}
	c.Flush()
	return nil
}

var con redis.Conn

func init() {
	var err error
	con, err = redis.Dial("tcp", "101.43.4.210:6379")
	if err != nil {
		panic(err)
		return
	}

}
func GetRedisCon() redis.Conn {
	return con
}

func RedisConfDo(commandName string, args ...interface{}) (interface{}, error) {
	return con.Do(commandName, args...)
}
