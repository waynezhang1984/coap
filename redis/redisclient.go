package redisclient

import (
	"os"
	"runtime"

	"github.com/coapprocess/configure"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}

func ConnRedis() (c redis.Conn, err error) {
	url := configure.GetRedisUrl()
	db, err := redis.Dial("tcp", url)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "redisclient.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("Connect to redis error: ", err)

		return
	}
	return db, err
}

func SaveRedisString(key, value string) error {

	url := configure.GetRedisUrl()
	c, err := redis.Dial("tcp", url)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "redisclient.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("Connect to redis error: ", err)

		return err
	}
	defer c.Close()
	_, err = c.Do("SET", key, value)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "redisclient.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("redis set failed: ", err)
	}

	return err
}

func ReadRedisString(key string) (string, error) {

	url := configure.GetRedisUrl()
	c, err := redis.Dial("tcp", url)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "redisclient.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("Connect to redis error: ", err)

		return "", err
	}
	defer c.Close()
	value, err := redis.String(c.Do("GET", key))
	if err != nil {
		//fmt.Println("redis get failed:", err)
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "redisclient.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("redis get failed: ", err)
		return "", err
	}
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "redisclient.go",
		"line":     line,
		"func":     f.Name(),
		"value":    value,
	}).Info("ReadRedisByte")

	return value, nil
}

func DelRedis(key string) error {

	url := configure.GetRedisUrl()
	c, err := redis.Dial("tcp", url)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "redisclient.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("Connect to redis error: ", err)

		return err
	}
	defer c.Close()

	_, err = c.Do("DEL", key)
	if err != nil {
		//fmt.Println("redis delete failed:", err)
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "redisclient.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("redis delete failed: ", err)
	}
	return err
}

func SaveRedisByte(key string, value []byte) error {

	url := configure.GetRedisUrl()
	c, err := redis.Dial("tcp", url)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "redisclient.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("Connect to redis error: ", err)

		return err
	}

	defer c.Close()
	_, err = c.Do("SET", key, value)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "redisclient.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("redis set failed: ", err)
	}

	return err
}

func ReadRedisByte(key string) ([]byte, error) {

	url := configure.GetRedisUrl()
	c, err := redis.Dial("tcp", url)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "redisclient.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("Connect to redis error: ", err)

		return []byte{}, err
	}
	defer c.Close()
	value, err := redis.Bytes(c.Do("GET", key))
	if err != nil {
		//fmt.Println("redis get failed:", err)
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "redisclient.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("redis get failed: ", err)
		return []byte{}, err
	}
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "redisclient.go",
		"line":     line,
		"func":     f.Name(),
		"value":    value,
	}).Info("ReadRedisByte")

	return value, nil
}
