package main

import (
	"obcsdk/lstutil"
)

/*************** Test Objective : Ledger Stress with 2 Clients and 2 Peers *********************
* 
*   1. Connect to a 4 node peer network with security enabled, and deploy a modified version of
*	chaincode_example02 that stores an additional block of data with every transaction
*	Refer to lstutil.go for more details, including parameters and further configuration.
*   2. Invoke TRX_COUNT transactions in parallel, divided among each go client thread
*   3. Check if the total counter value (TRX_COUNT) matches with query on "counter"
* 
*   To use this test script:			go run <testname.go>
*   or (to save output files):			../automation/go_record.sh <testname.go>
* 
***********************************************************************************************/

func main() {
	// args:  testname, # client threads, # peers, total # transactions
	lstutil.RunLedgerStressTest( "LST_2client2peer20K", 2, 2, 20000 )
}
