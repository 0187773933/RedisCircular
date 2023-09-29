package main

import (
	"fmt"
	"context"
	redis "github.com/redis/go-redis/v9"
	circular_set "github.com/0187773933/RedisCircular/v1/set"
)

var DB *redis.Client

func SetupDB() {
	DB = redis.NewClient( &redis.Options{
		Addr: "localhost:6379" ,
		Password: "" ,
		DB: 6 ,
	})
	var ctx = context.Background()
	ping_result , err := DB.Ping( ctx ).Result()
	fmt.Printf( "DB Connected : Ping = %s\n" , ping_result )
	if err != nil { panic( err ) }
}

func main() {

	SetupDB()
	key := "testing-circular-set"

	DB.Del( context.Background() , key )
	DB.Del( context.Background() , key + ".INDEX" )
	circular_set.Add( DB , key , "A" )
	circular_set.Add( DB , key , "B" )
	circular_set.Add( DB , key , "C" )
	circular_set.Add( DB , key , "D" )
	circular_set.Add( DB , key , "E" )
	circular_set.Add( DB , key , "F" )

	fmt.Printf( "\nForward\n" )
	for i := 0 ; i < 8 ; i++ {
		fmt.Println( circular_set.Next( DB , key ) )
	}

	fmt.Printf( "\nReverse\n" )
	for i := 0 ; i < 9 ; i++ {
		fmt.Println( circular_set.Previous( DB , key ) )
	}

	circular_set.Add( DB , key , "G" )
	circular_set.Add( DB , key , "H" )
	circular_set.Add( DB , key , "I" )
	circular_set.Add( DB , key , "J" )
	circular_set.Add( DB , key , "K" )
	circular_set.Add( DB , key , "L" )

	fmt.Printf( "\nForward\n" )
	for i := 0 ; i < 8 ; i++ {
		fmt.Println( circular_set.Next( DB , key ) )
	}

	fmt.Printf( "\nReverse\n" )
	for i := 0 ; i < 9 ; i++ {
		fmt.Println( circular_set.Previous( DB , key ) )
	}

}