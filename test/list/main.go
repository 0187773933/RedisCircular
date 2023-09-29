package main

import (
	"fmt"
	"context"
	redis "github.com/redis/go-redis/v9"
	circular_list "github.com/0187773933/RedisCircular/v1/list"
)

// https://github.com/andymccurdy/redis-py/blob/1f857f0053606c23cb3f1abd794e3efbf6981e09/tests/test_commands.py
// https://github.com/ceberous/redis-manager-utils/blob/master/BaseClass.js
// https://github.com/48723247842/RedisCirclularList/blob/master/redis_circular_list/__init__.py
// https://redis.io/commands/sadd
// https://pkg.go.dev/builtin#error
// https://pkg.go.dev/github.com/go-redis/redis/v8#BoolCmd.Err

func main() {

	DB := redis.NewClient( &redis.Options{
		Addr: "localhost:6379" ,
		DB: 3 ,
		Password: "" ,
	})
	defer DB.Close()

	key := "testing-circular-list"

	var ctx = context.Background()

	DB.Del( ctx , key )
	DB.Del( ctx , key + ".INDEX" )

	DB.RPush( ctx , key , "1" )
	DB.RPush( ctx , key , "2" )
	DB.RPush( ctx , key , "3" )
	DB.RPush( ctx , key , "4" )
	DB.RPush( ctx , key , "5" )
	DB.RPush( ctx , key , "6" )

	current , _ := circular_list.Current( DB , key )
	fmt.Printf( "Current = %s\n"  , current )

	fmt.Printf( "Previous = %s\n"  , circular_list.Previous( DB , key ) )
	fmt.Printf( "Previous = %s\n"  , circular_list.Previous( DB , key ) )
	fmt.Printf( "Previous = %s\n"  , circular_list.Previous( DB , key ) )

	fmt.Printf( "Next = %s\n"  , circular_list.Next( DB , key ) )
	fmt.Printf( "Next = %s\n"  , circular_list.Next( DB , key ) )
	fmt.Printf( "Next = %s\n"  , circular_list.Next( DB , key ) )

	current , _ = circular_list.Current( DB , key )
	fmt.Printf( "Current = %s\n"  , current )

}