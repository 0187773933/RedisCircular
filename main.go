package main

import (
	"fmt"
	"context"
	redis "github.com/redis/go-redis/v9"
	// circular_list "github.com/0187773933/RedisCircular/v1/list"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

var DB *redis.Client

func SetupDB() {
	DB = redis.NewClient( &redis.Options{
		Addr: "localhost:6379" ,
		Password: "" ,
		DB: 4 ,
	})
	var ctx = context.Background()
	ping_result , err := DB.Ping( ctx ).Result()
	fmt.Printf( "DB Connected : PING = %s\n" , ping_result )
	if err != nil { panic( err ) }
}

func TestSet() {
	var ctx = context.Background()
	key := "REDIS_CIRCULAR.TEST.SET"
	key_index := key + ".INDEX"
	DB.Del( ctx , key )
	DB.Del( ctx , key_index )
	for i := 0 ; i < 5 ; i++ {
		circular_set.Add( DB , key , fmt.Sprintf( "asdf-%d" , ( i + 1 ) ) )
	}
	one := circular_set.Next( DB , key ) // should be asdf-1
	two := circular_set.Next( DB , key ) // should be asdf-2
	three := circular_set.Next( DB , key ) // should be asdf-3
	fmt.Println( one , two , three )
	circular_set.Remove( DB , key , "asdf-3" )
	current := circular_set.Current( DB , key ) // should be asdf-4
	next_one := circular_set.Next( DB , key ) // should be asdf-5
	next_two := circular_set.Next( DB , key ) // should be asdf-1
	fmt.Println( current , next_one , next_two )
	circular_set.Add( DB , key , "asdf-6" )
	current = circular_set.Current( DB , key ) // should be asdf-4
	next_one = circular_set.Next( DB , key ) // should be asdf-5
	next_two = circular_set.Next( DB , key ) // should be asdf-1
	next_three := circular_set.Next( DB , key ) // should be asdf-6
	next_four := circular_set.Next( DB , key ) // should be asdf-6
	next_five := circular_set.Next( DB , key ) // should be asdf-6
	fmt.Println( current , next_one , next_two , next_three , next_four , next_five )

	// one := circular_set.Previous( DB , key ) // should be asdf-5
	// two := circular_set.Previous( DB , key ) // should be asdf-4
	// three := circular_set.Previous( DB , key ) // should be asdf-3

}

func main() {
	SetupDB()
	// TestList()
	TestSet()
}