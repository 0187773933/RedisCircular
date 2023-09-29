package list

import (
	"fmt"
	"context"
	"strconv"
	redis "github.com/redis/go-redis/v9"
)

func Current( redis_connection *redis.Client , redis_circular_list_key string ) ( result string , index string ) {
	result = "failed"
	index = "0"
	var ctx = context.Background()
	// 1.) Get Length
	circular_list_length , circular_list_length_error := redis_connection.LLen( ctx , redis_circular_list_key ).Result()
	if circular_list_length_error != nil { panic( circular_list_length_error ) }
	if circular_list_length < 1 { return; }

	// 2.) Get Current Index
	circular_list_key_index_key := fmt.Sprintf( "%s.INDEX" , redis_circular_list_key )
	circular_list_key_index , circular_list_key_index_error := redis_connection.Get( ctx , circular_list_key_index_key ).Result()
	if circular_list_key_index_error != nil {
		_ , initialize_index_key_error := redis_connection.Set( ctx , circular_list_key_index_key , "0" , 0 ).Result()
		if initialize_index_key_error != nil { panic( "Could Not Reinitialize Index Variable" ) }
		circular_list_key_index = "0"
	}
	circular_list_index_int , _ := strconv.ParseInt( circular_list_key_index , 0 , 64 )
	circular_list_index_int_64 := int64( circular_list_index_int )

	// 3.) Get Current Value at Index
	current_in_circle , current_in_circle_error := redis_connection.LIndex( ctx , redis_circular_list_key , circular_list_index_int_64 ).Result()
	if current_in_circle_error != nil { panic( circular_list_key_index_error ) }
	result = current_in_circle
	index = circular_list_key_index
	return
}

func Previous( redis_connection *redis.Client , redis_circular_list_key string ) ( result string ) {

	result = "failed"
	var ctx = context.Background()
	//redis := get_redis_connection( "localhost:6379" , 3 , "" )

	// ==================== 1.) Get Length ===================================================================================================================
	circular_list_length_int_64 , circular_list_length_error := redis_connection.LLen( ctx , redis_circular_list_key ).Result()
	if circular_list_length_error != nil { panic( "Couldn't Get List Length" ) }
	if circular_list_length_int_64 < 1 { panic( "Circular List Length < 1" ) }
	// ==================== 1.) Get Length ===================================================================================================================


	// ==================== 2.) Get Current Index ============================================================================================================
	circular_list_key_index_key := fmt.Sprintf( "%s.INDEX" , redis_circular_list_key )
	circular_list_key_index , circular_list_key_index_error := redis_connection.Get( ctx , circular_list_key_index_key ).Result()

	// Create .INDEX on the circular list key if it doesnt already exist
	if circular_list_key_index_error != nil {
		circular_list_key_index = "0"
		_ , initialize_index_key_error := redis_connection.Set( ctx , circular_list_key_index_key , circular_list_key_index , 0 ).Result()
		if initialize_index_key_error != nil { panic( fmt.Sprintf( "Couldn't SET %s" , circular_list_key_index_key  ) ) }
	}
	circular_list_index_int , _ := strconv.ParseInt( circular_list_key_index , 0 , 64 )
	circular_list_index_int_64 := int64( circular_list_index_int )
	// ==================== 2.) Get Current Index ============================================================================================================


	// ==================== 3.) if index == 0 { Set Index To Last Item In List  } else { Decrement Index } ===================================================
	if circular_list_index_int_64 == 0 {
		circular_list_index_int_64 = ( circular_list_length_int_64 - 1 )
		_ , circular_list_set_result_error := redis_connection.Set( ctx , circular_list_key_index_key , circular_list_index_int_64 , 0 ).Result()
		if circular_list_set_result_error != nil { panic( fmt.Sprintf( "Could not set %s" , circular_list_key_index_key ) ) }
		} else {
			circular_list_index_int_64 = ( circular_list_index_int_64 - 1 )
			_ , circular_list_decr_result_error := redis_connection.Decr( ctx , circular_list_key_index_key ).Result()
			if circular_list_decr_result_error != nil { panic( fmt.Sprintf( "Couldn't DECR %s" , circular_list_key_index_key  ) ) }
		}
	// ==================== 3.) if index == 0 { Set Index To Last Item In List  } else { Decrement Index } ===================================================

	// ==================== 4.) Get Previous in List @ Updated Index =========================================================================================
	previous_in_circle_list , previous_in_circle_list_error := redis_connection.LIndex( ctx , redis_circular_list_key , circular_list_index_int_64 ).Result()
	if previous_in_circle_list_error != nil { panic( fmt.Sprintf( "Couldn't LINDEX %s %d" , circular_list_key_index_key , circular_list_index_int_64 ) ) }
	result = previous_in_circle_list
	// ==================== 4.) Get Previous in List @ Updated Index =========================================================================================

	return
}

func Next( redis_connection *redis.Client , redis_circular_list_key string ) ( result string ) {

	result = "failed"
	var ctx = context.Background()
	//redis := get_redis_connection( "localhost:6379" , 3 , "" )

	// ==================== 1.) Get Length ==================================================================================================================
	circular_list_length_int_64 , circular_list_length_error := redis_connection.LLen( ctx , redis_circular_list_key ).Result()
	if circular_list_length_error != nil { panic( "Couldn't Get List Length" ) }
	if circular_list_length_int_64 < 1 { panic( "Circular List Length < 1" ) }
	// ==================== 1.) Get Length ==================================================================================================================


	// ==================== 2.) Get Current Index ===========================================================================================================
	circular_list_key_index_key := fmt.Sprintf( "%s.INDEX" , redis_circular_list_key )
	circular_list_key_index , circular_list_key_index_error := redis_connection.Get( ctx , circular_list_key_index_key ).Result()

	// Create .INDEX on the circular list key if it doesnt already exist
	if circular_list_key_index_error != nil {
		circular_list_key_index = "0"
		_ , initialize_index_key_error := redis_connection.Set( ctx , circular_list_key_index_key , circular_list_key_index , 0 ).Result()
		if initialize_index_key_error != nil { panic( fmt.Sprintf( "Couldn't SET %s" , circular_list_key_index_key  ) ) }
	}
	circular_list_index_int , _ := strconv.ParseInt( circular_list_key_index , 0 , 64 )
	circular_list_index_int_64 := int64( circular_list_index_int )
	// ==================== 2.) Get Current Index ===========================================================================================================


	// ==================== 3.) if index == ( circular_list_length - 1 ) { Set Index To 0 } else { Increment Index } ========================================
	if circular_list_index_int_64 == ( circular_list_length_int_64 - 1 ) {
		circular_list_index_int_64 = 0
		_ , circular_list_set_result_error := redis_connection.Set( ctx , circular_list_key_index_key , circular_list_index_int_64 , 0 ).Result()
		if circular_list_set_result_error != nil { panic( fmt.Sprintf( "Could not set %s" , circular_list_key_index_key ) ) }
		} else {
			circular_list_index_int_64 = ( circular_list_index_int_64 + 1 )
			_ , circular_list_decr_result_error := redis_connection.Incr( ctx , circular_list_key_index_key ).Result()
			if circular_list_decr_result_error != nil { panic( fmt.Sprintf( "Couldn't DECR %s" , circular_list_key_index_key  ) ) }
		}
	// ==================== 3.) if index == ( circular_list_length - 1 ) { Set Index To 0 } else { Increment Index } ========================================

	// ==================== 4.) Get Next in List @ Updated Index ============================================================================================
	next_in_circle_list , next_in_circle_list_error := redis_connection.LIndex( ctx , redis_circular_list_key , circular_list_index_int_64 ).Result()
	if next_in_circle_list_error != nil { panic( fmt.Sprintf( "Couldn't LINDEX %s %d" , circular_list_key_index_key , circular_list_index_int_64 ) ) }
	result = next_in_circle_list
	// ==================== 4.) Get Next in List @ Updated Index ============================================================================================

	return
}