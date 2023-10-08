package set

import (
	// "fmt"
	"context"
	redis "github.com/redis/go-redis/v9"
)

func Index( rdb *redis.Client , key string ) ( result int ) {
	var ctx = context.Background()
	index , err := rdb.Get( ctx , key ).Int()
	if err == redis.Nil {
		// If the key doesn't exist, initialize it to 0
		rdb.Set( ctx , key , 0 , 0 )
		return 0
	} else if err != nil {
		panic( err )
	}
	return index
}

func Add( rdb *redis.Client , key string , member string ) {
	// Get the current length of the sorted set to use as the score for the new item
	var ctx = context.Background()
	score := rdb.ZCard( ctx , key ).Val()
	// Add the new member to the sorted set with the score
	rdb.ZAddNX( ctx , key , redis.Z{ Score: float64( score ) , Member: member } ).Result()
}

func Current( rdb *redis.Client , key string ) ( result string ) {
	var ctx = context.Background()
	key_index := key + ".INDEX"
	index := Index( rdb , key_index )
	items , err := rdb.ZRange( ctx , key , int64( index ) , int64( index ) ).Result()
	if err != nil { panic( err ) }
	if len( items ) > 0 {
		result = items[ 0 ]
	} else {
		result = ""
	}
	return
}

func Next( rdb *redis.Client , key string ) ( result string ) {
	var ctx = context.Background()

	key_index := key + ".INDEX"
	index := Index( rdb , key_index )

	// Increment the index
	index = ( index + 1 )
	total , err := rdb.ZCard( ctx , key ).Result()
	if err != nil { panic( err ) }
	if int64( index ) >= total {
		index = 0 // Wrap around
		// fmt.Println( "Wrapped around to the beginning!" , index )
	}
	// fmt.Printf("Current Index: %d\n", index)
	rdb.Set( ctx , key_index , index , 0 )

	// Get the item at the current index
	items , err := rdb.ZRange( ctx , key , int64( index ) , int64( index ) ).Result()
	if err != nil { panic( err ) }

	if len( items ) > 0 {
		result = items[ 0 ]
	} else {
		result = ""
	}
	return
}

func Previous( rdb *redis.Client , key string ) ( result string ) {
	var ctx = context.Background()

	key_index := key + ".INDEX"
	index := Index( rdb , key_index )

	// Decrement the index
	index = ( index - 1 )
	if index < 0 {
		total, err := rdb.ZCard( ctx , key ).Result()
		if err != nil { panic( err ) }
		index = ( int( total ) - 1 ) // Wrap around
		// fmt.Println( "Wrapped around to the end!" , index )
	}
	rdb.Set( ctx , key_index , index , 0 )

	// Print the current index
	// fmt.Printf( "Current Index: %d\n" , index )

	// Get the item at the adjusted index
	items , err := rdb.ZRange( ctx , key , int64( index ) , int64( index ) ).Result()
	if err != nil { panic( err ) }

	if len( items ) > 0 {
		result = items[ 0 ]
	} else {
		result = ""
	}
	return
}