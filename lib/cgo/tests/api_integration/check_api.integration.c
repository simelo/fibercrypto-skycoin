#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

#define NODE_ADDRESS "SKYCOIN_NODE_HOST"
#define NODE_ADDRESS_DEFAULT "http://127.0.0.1:46420"
#define BUFFER_SIZE 1024
#define STABLE 1

#define STRING_SIZE 128
#define JSON_FILE_SIZE 4096
#define JSON_BIG_FILE_SIZE 32768
#define TEST_DATA_DIR "src/api/integration/testdata/"


char* getNodeAddress(){
	if( STABLE ){
		return NODE_ADDRESS_DEFAULT;
	} else {
		GoString_ nodeAddress;
		memset(&nodeAddress, 0, sizeof(GoString_));
		GoString  nodeEnvName = {NODE_ADDRESS, strlen(NODE_ADDRESS)};
		int result = SKY_cli_Getenv(nodeEnvName, &nodeAddress);
		cr_assert(result == SKY_OK, "Couldn\'t get node address from enviroment");
		registerMemCleanup((void*)nodeAddress.p);
		if( strcmp(nodeAddress.p, "") == 0){
			return NODE_ADDRESS_DEFAULT;
		}
		return (char*)nodeAddress.p;
	}
}

json_value* loadGoldenFile(const char* file){
	char path[STRING_SIZE];
	if(strlen(TEST_DATA_DIR) + strlen(file) < STRING_SIZE){
		strcpy(path, TEST_DATA_DIR);
		strcat(path, file);
		return loadJsonFile(path);
	}
	return NULL;
}

GoString* createGoStringSlice(char** pStrings, int count, GoSlice* slice){
	GoString* goStrings = malloc(sizeof(GoString) * count);
	cr_assert(goStrings != NULL, "Error creating GoString Slice");
	registerMemCleanup( goStrings );
	for(int i = 0; i < count; i++){
		goStrings[i].p = pStrings[i];
		goStrings[i].n = strlen(pStrings[i]);
	}
	slice->data = goStrings;
	slice->len = count;
	slice->cap = count;
	return goStrings;
}

int compareObjectsByHandle(Handle h1, Handle h2){
	GoString_ jsonResult1, jsonResult2;
	int result;
	memset(&jsonResult1, 0, sizeof(GoString_));
	memset(&jsonResult2, 0, sizeof(GoString_));
	
	result = SKY_JsonEncode_Handle(h1, &jsonResult1);
	cr_assert(result == SKY_OK, "Couldn\'t json encode");
	registerMemCleanup((void*)jsonResult1.p);
	
	result = SKY_JsonEncode_Handle(h2, &jsonResult2);
	cr_assert(result == SKY_OK, "Couldn\'t json encode");
	registerMemCleanup((void*)jsonResult2.p);
	
	json_char* json1 = (json_char*)jsonResult1.p;
	json_value* value1 = json_parse(json1, strlen(jsonResult1.p));
	cr_assert(value1 != NULL, "json_parse failed");
	registerJsonFree(value1);
	
	json_char* json2 = (json_char*)jsonResult2.p;
	json_value* value2 = json_parse(json1, strlen(jsonResult2.p));
	cr_assert(value2 != NULL, "json_parse failed");
	registerJsonFree(value2);
	
	int equal = compareJsonValues(value1, value2);
	
	freeRegisteredMemCleanup((void*)jsonResult1.p);
	freeRegisteredMemCleanup((void*)jsonResult2.p);
	freeRegisteredJson(value1);
	freeRegisteredJson(value2);
	return equal;
}

int compareObjectWithGoldenFile(Handle handle, const char* golden_file){
	GoString_ jsonResult;
	int result;
	memset(&jsonResult, 0, sizeof(GoString_));
	
	result = SKY_JsonEncode_Handle(handle, &jsonResult);
	cr_assert(result == SKY_OK, "Couldn\'t json encode");
	registerMemCleanup((void*)jsonResult.p);
	
	json_char* json = (json_char*)jsonResult.p;
	json_value* value = json_parse(json, strlen(jsonResult.p));
	cr_assert(value != NULL, "json_parse failed");
	registerJsonFree(value);
	
	json_value* golden_value = loadGoldenFile(golden_file);
	cr_assert(golden_value != NULL, "loadGoldenFile failed");
	registerJsonFree(golden_value);
	
	int equal = compareJsonValues(value, golden_value);
	
	freeRegisteredJson(value);
	freeRegisteredJson(golden_value);
	freeRegisteredMemCleanup((void*)jsonResult.p);
	
	return equal;
}

Test(api_integration, TestVersion) {
	GoString_ version;
	memset(&version, 0, sizeof(GoString_));
	
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;
	Handle versionDataHandle;
	
	int result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );
	result = SKY_api_Client_Version( &clientHandle, &versionDataHandle );
	cr_assert(result == SKY_OK, "Couldn\'t get version");
	registerHandleClose( versionDataHandle );
	result = SKY_JsonEncode_Handle(versionDataHandle, &version);
	cr_assert(result == SKY_OK, "Couldn\'t json encode version");
	registerMemCleanup((void*)version.p);
	int versionLength = strlen(version.p);
	cr_assert(versionLength > 0, "Invalid version data");
}

Test(api_integration, TestStableCoinSupply) {
	GoString_ jsonResult;
	memset(&jsonResult, 0, sizeof(GoString_));
	
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;
	Handle coinSupplyHandle;
	
	int result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );
	
	result = SKY_api_Client_CoinSupply( &clientHandle, &coinSupplyHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_CoinSupply failed");
	registerHandleClose( coinSupplyHandle );
	
	result = SKY_JsonEncode_Handle(coinSupplyHandle, &jsonResult);
	cr_assert(result == SKY_OK, "Couldn\'t json encode");
	registerMemCleanup((void*)jsonResult.p);
	
	json_char* json = (json_char*)jsonResult.p;
	json_value* value = json_parse(json, strlen(jsonResult.p));
	cr_assert(value != NULL, "json_parse failed");
	registerJsonFree(value);
	
	json_value* json_golden = loadGoldenFile("coinsupply.golden");
	cr_assert(json_golden != NULL, "loadGoldenFile failed");
	registerJsonFree(json_golden);
	
	int equal = compareJsonValues(value, json_golden);
	cr_assert(equal, "Output different than expected");
}

typedef struct{
	char*  	golden_file;
	int		addresses_count;
	char**	addresses;
	int 	hashes_count;
	char**  hashes;
	int		failure;
}test_output;

Test(api_integration, TestStableOutputs) {
	int result;
	GoString_ jsonResult;
	memset(&jsonResult, 0, sizeof(GoString_));
	
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;
	
	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );
	
	char* address_1[] = {
		"ALJVNKYL7WGxFBSriiZuwZKWD4b7fbV1od",
		"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf",
		"qxmeHkwgAMfwXyaQrwv9jq3qt228xMuoT5",
	};
	
	char* hashes_1[] = {
		"9e53268a18f8d32a44b4fb183033b49bebfe9d0da3bf3ef2ad1d560500aa54c6",
		"d91e07318227651129b715d2db448ae245b442acd08c8b4525a934f0e87efce9",
		"01f9c1d6c83dbc1c993357436cdf7f214acd0bfa107ff7f1466d1b18ec03563e",
		"fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f20",
	};
	int test_cases = 3;
	test_output tests[] = {
		{
			"outputs-noargs.golden",
			0, NULL, 0, NULL, 0
		},
		{
			"outputs-addrs.golden",
			3, address_1,
			0, NULL, 0
		},
		{
			"outputs-hashes.golden",
			0, NULL,
			4, hashes_1, 0
		},
	};
	Handle outputHandle;
	GoSlice strings;
	
	for(int i = 0; i < test_cases; i++){
		memset(&strings, 0, sizeof(GoSlice));
		cr_assert(tests[i].addresses_count == 0 || tests[i].hashes_count == 0);
		if(tests[i].addresses_count == 0 && tests[i].hashes_count == 0){
			result = SKY_api_Client_Outputs(&clientHandle, &outputHandle);
		} else if(tests[i].addresses_count > 0){
			createGoStringSlice(tests[i].addresses, 
					tests[i].addresses_count, &strings);
			result = SKY_api_Client_OutputsForAddresses(&clientHandle, 
											strings, &outputHandle);
		} else if(tests[i].hashes_count > 0){
			createGoStringSlice(tests[i].hashes, 
					tests[i].hashes_count, &strings);
			result = SKY_api_Client_OutputsForHashes(&clientHandle, 
											strings, &outputHandle);
		}
		
		if( tests[i].failure ){
			cr_assert(result != SKY_OK, "SKY_api_Client_Outputs should have failed");
			continue;
		}
		cr_assert(result == SKY_OK, "SKY_api_Client_Outputs failed");
		registerHandleClose( outputHandle );
		
		result = SKY_JsonEncode_Handle(outputHandle, &jsonResult);
		cr_assert(result == SKY_OK, "Couldn\'t json encode");
		registerMemCleanup((void*)jsonResult.p);
		
		json_char* json = (json_char*)jsonResult.p;
		json_value* jsonOutput = json_parse(json, strlen(jsonResult.p));
		cr_assert(jsonOutput != NULL, "json_parse failed");
		registerJsonFree(jsonOutput);
		
		json_value* json_golden = loadGoldenFile(tests[i].golden_file);
		cr_assert(json_golden != NULL, "loadGoldenFile failed");
		registerJsonFree(json_golden);
		
		int equal = compareJsonValues(jsonOutput, json_golden);
		cr_assert(equal, "Output different than expected");
	}
}

typedef struct{
	char*  		golden_file;
	char* 		hash;
	GoUint64 	seq;
	int			failure;
}test_block;

Test(api_integration, TestStableBlock) {
	int test_count = 5;
	test_block tests[] = {
		{
			NULL, 
			"80744ec25e6233f40074d35bf0bfdbddfac777869b954a96833cb89f44204444",
			0, 1
		},
		{
			"block-hash.golden",
			"70584db7fb8ab88b8dbcfed72ddc42a1aeb8c4882266dbb78439ba3efcd0458d",
			0, 0,
		},
		{
			"block-hash-genesis.golden",
			"0551a1e5af999fe8fff529f6f2ab341e1e33db95135eef1b2be44fe6981349f3",
			0, 0,
		},
		{
			"block-seq-0.golden",
			NULL,
			0, 0,
		},
		{
			"block-seq-100.golden",
			NULL,
			100, 0,
		},
	};
	
	int result;
	GoString_ jsonResult;
	memset(&jsonResult, 0, sizeof(GoString_));
	
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;
	
	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );
	
	GoString strHash;
	Handle blockHandle;
	
	for(int i = 0; i < test_count; i++){
		if( tests[i].hash != NULL ){
			memset( &strHash, 0, sizeof(GoString) );
			strHash.p = tests[i].hash;
			strHash.n = strlen( tests[i].hash );
			result = SKY_api_Client_BlockByHash(&clientHandle, 
				strHash, &blockHandle);
		} else {
			result = SKY_api_Client_BlockBySeq(&clientHandle,
				tests[i].seq, &blockHandle);
		}
		if( tests[i].failure ){
			cr_assert( result != SKY_OK, "Get Block should have failed" );
			continue;
		}
		cr_assert( result == SKY_OK, "Get Block failed" );
		registerHandleClose( blockHandle );
		
		result = SKY_JsonEncode_Handle(blockHandle, &jsonResult);
		cr_assert(result == SKY_OK, "Couldn\'t json encode");
		registerMemCleanup((void*)jsonResult.p);
		
		json_char* json = (json_char*)jsonResult.p;
		json_value* jsonOutput = json_parse(json, strlen(jsonResult.p));
		cr_assert(jsonOutput != NULL, "json_parse failed");
		registerJsonFree(jsonOutput);
		
		json_value* json_golden = loadGoldenFile(tests[i].golden_file);
		cr_assert(json_golden != NULL, "loadGoldenFile failed");
		registerJsonFree(json_golden);
		
		int equal = compareJsonValues(jsonOutput, json_golden);
		cr_assert(equal, "Output different than expected");
	}
	
	Handle progressHandle;
	printf("Querying every block in the blockchain");
	result = SKY_api_Client_BlockchainProgress(&clientHandle, &progressHandle);
	cr_assert(result == SKY_OK, "SKY_api_Client_BlockchainProgress failed");
	registerHandleClose( progressHandle );
	GoInt64 progress;
	result = SKY_Handle_Progress_GetCurrent( progressHandle, &progress );
	cr_assert(result == SKY_OK, "SKY_Handle_Progress_GetCurrent failed");
	GoInt64 seq, blockSeq;
	Handle prevBlockHandle = 0;
	Handle blockHandle2;
	blockHandle = 0;
	GoString_ hash1, hash2, hash;
	GoString _hash;
	for(seq = 0; seq < progress; seq++){
		result = SKY_api_Client_BlockBySeq(&clientHandle,
				seq, &blockHandle);
		cr_assert( result == SKY_OK, "SKY_api_Client_BlockBySeq failed" );
		registerHandleClose( blockHandle );
		result = SKY_Handle_Block_GetHeadSeq( blockHandle, &blockSeq );
		cr_assert( result == SKY_OK, "SKY_Handle_Block_GetHeadSeq failed" );
		cr_assert( seq == blockSeq, "Incorrect block sequence" );
		if( prevBlockHandle ){
			memset( &hash1, 0, sizeof(GoString_) );
			memset( &hash2, 0, sizeof(GoString_) );
			memset( &hash, 0, sizeof(GoString_) );
			result = SKY_Handle_Block_GetHeadHash( prevBlockHandle, &hash1 );
			cr_assert(result == SKY_OK, "SKY_Handle_Block_GetHeadHash failed");
			registerMemCleanup((void*)hash1.p);
			result = SKY_Handle_Block_GetPreviousBlockHash( blockHandle, &hash2 );
			cr_assert(result == SKY_OK, "SKY_Handle_Block_GetPreviousBlockHash failed");
			registerMemCleanup((void*)hash2.p);
			cr_assert(eq(type(GoString_), hash1, hash2));
			freeRegisteredMemCleanup((void*)hash1.p);
			freeRegisteredMemCleanup((void*)hash1.p);
			result = SKY_Handle_Block_GetHeadHash( blockHandle, &hash );
			registerMemCleanup((void*)hash.p);
			_hash.p = hash.p;
			_hash.n = hash.n;
			result = SKY_api_Client_BlockByHash(&clientHandle, 
				_hash, &blockHandle2);
			cr_assert(result == SKY_OK, "SKY_api_Client_BlockByHash failed");
			registerHandleClose( blockHandle2 );
			
			int equal = compareObjectsByHandle(blockHandle, blockHandle2);
			cr_assert(equal == 1);
			freeRegisteredMemCleanup((void*)hash.p);
			closeRegisteredHandle( blockHandle2 );
		}
		if( prevBlockHandle ){
			closeRegisteredHandle( prevBlockHandle );
		}
		prevBlockHandle = blockHandle;
	}
	if( blockHandle ){
		closeRegisteredHandle( blockHandle );
	}
}

Test(api_integration, TestStableBlockchainMetadata) {
	int result;
	GoString_ jsonResult;
	memset(&jsonResult, 0, sizeof(GoString_));
	
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;
	
	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );
	
	Handle metadataHandle;
	result = SKY_api_Client_BlockchainMetadata( &clientHandle, 
									&metadataHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_BlockchainMetadata failed");
	registerHandleClose( metadataHandle );
	
	int equal = compareObjectWithGoldenFile(metadataHandle, 
									"blockchain-metadata.golden");
	cr_assert(equal, "SKY_api_Client_BlockchainMetadata returned unexpected result");
}

Test(api_integration, TestStableBlockchainProgress) {
	int result;
	GoString_ jsonResult;
	memset(&jsonResult, 0, sizeof(GoString_));
	
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;
	
	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );
	
	Handle progressHandle;
	result = SKY_api_Client_BlockchainProgress( &clientHandle, 
									&progressHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_BlockchainMetadata failed");
	registerHandleClose( progressHandle );
	
	int equal = compareObjectWithGoldenFile(progressHandle, 
									"blockchain-progress.golden");
	cr_assert(equal, "SKY_api_Client_BlockchainProgress returned unexpected result");
}

typedef struct{
	int 	addresses_count;
	char** 	addresses;
	char* 	golden_file;
} test_balance;

Test(api_integration, TestStableBalance) {
	char* addr1[] = {
		"prRXwTcDK24hs6AFxj69UuWae3LzhrsPW9"
	};
	char* addr2[] = {
		"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf"
	};
	char* addr3[] = {
		"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf", 
		"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf"
	};
	char* addr4[] = {
		"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf", 
		"qxmeHkwgAMfwXyaQrwv9jq3qt228xMuoT5"
	};
	int result;
	int tests_count = 4;
	test_balance tests[] = {
		{
			1, addr1, "balance-noaddrs.golden"
		},
		{
			1, addr2, "balance-2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf.golden"
		},
		{
			2, addr3, "balance-2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf.golden"
		},
		{
			2, addr4, "balance-two-addrs.golden"
		},
	};
	
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;
	
	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );
	
	int i;
	GoSlice strings;
	wallet__BalancePair balance;
	for(i = 0; i < tests_count; i++){
		memset( &strings, 0, sizeof(GoSlice) );
		createGoStringSlice( tests[i].addresses, tests[i].addresses_count,
							&strings);
		result = SKY_api_Client_Balance( &clientHandle, 
										strings, &balance );
		cr_assert(result == SKY_OK, "SKY_api_Client_BlockchainMetadata failed");
		json_value* json_golden = loadGoldenFile(tests[i].golden_file);
		cr_assert(json_golden != NULL, "loadGoldenFile failed");
		registerJsonFree(json_golden);
		json_value* value;
		value = get_json_value(json_golden, 
							"confirmed/coins", json_integer);
		cr_assert(value != NULL, "get_json_value confirmed/coins failed");
		cr_assert(value->u.integer == balance.Confirmed.Coins);
		value = get_json_value(json_golden, 
							"confirmed/hours", json_integer);
		cr_assert(value != NULL, "get_json_value confirmed/hours failed");
		cr_assert(value->u.integer == balance.Confirmed.Hours);
		value = get_json_value(json_golden, 
							"predicted/coins", json_integer);
		cr_assert(value != NULL, "get_json_value predicted/coins failed");
		cr_assert(value->u.integer == balance.Predicted.Coins);
		value = get_json_value(json_golden, 
							"predicted/hours", json_integer);
		cr_assert(value != NULL, "get_json_value predicted/hours failed");
		cr_assert(value->u.integer == balance.Predicted.Hours);
	}
}


Test(api_integration, TestStableUxOut) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;
	
	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );
	
	char* golden_file = "uxout.golden";
	char* pUxId = "fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f20";
	GoString strUxId = {pUxId, strlen(pUxId)};
	Handle uxOutHandle;
	result = SKY_api_Client_UxOut( &clientHandle, strUxId, &uxOutHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_UxOut failed");
	registerHandleClose( uxOutHandle );
	
	int equal = compareObjectWithGoldenFile(uxOutHandle, golden_file);
	cr_assert(equal, "SKY_api_Client_UxOut returned unexpected result");
}

typedef struct{
	char* 	address;
	char*	golden_file;
	int 	failure;
}test_address_ux_out;

Test(api_integration, TestStableAddressUxOuts) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;
	
	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );
	int tests_count = 3; 
	test_address_ux_out tests[] = {
		{"", NULL, 1},
		{"prRXwTcDK24hs6AFxj69UuWae3LzhrsPW9", "uxout-noaddr.golden", 0},
		{"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf", "uxout-addr.golden", 0},
		
	};
	GoString addr;
	for(int i = 0; i < tests_count; i++){
		memset(&addr, 0, sizeof(GoString));
		addr.p = tests[i].address;
		addr.n = strlen( tests[i].address );
		Handle outHandle;
		result = SKY_api_Client_AddressUxOuts( &clientHandle, addr, &outHandle );
		if( tests[i].failure ){
			cr_assert(result != SKY_OK, "SKY_api_Client_AddressUxOuts should have failed");
			continue;
		} else {
			cr_assert(result == SKY_OK, "SKY_api_Client_AddressUxOuts failed");
		}
		registerHandleClose( outHandle );
		int equal = compareObjectWithGoldenFile( outHandle, tests[i].golden_file );
		cr_assert( equal == 1 );
	}
}

typedef struct{
	char* golden_file;
	GoUint64 start;
	GoUint64 end;
	int 	 failure;	
}test_blockn;

Handle testBlocks(Client__Handle clientHandle, 
				GoUint64 start, GoUint64 end){
	Handle blocksHandle;
	int result;
	result = SKY_api_Client_Blocks(&clientHandle, start, end, &blocksHandle);
	cr_assert(result == SKY_OK, "SKY_api_Client_Blocks failed");
	registerHandleClose( blocksHandle );
	GoUint64 count = 0;
	result = SKY_Handle_Blocks_GetCount( blocksHandle, &count );
	cr_assert(result == SKY_OK, "SKY_Handle_Blocks_GetCount failed");
	if( start > end ){
		cr_assert(count == 0);
	} else {
		cr_assert(count == end - start + 1);
	}
	GoUint64 i;
	GoString_ hash1, hash2, hash;
	GoString _hash;
	GoUint64 seq;
	int equal;
	for(i = 0; i < count; i++){
		Handle blockHandle = 0, previousBlockHandle = 0;
		Handle blockHandle2;
		result = SKY_Handle_Blocks_GetAt(blocksHandle, i, &blockHandle);
		cr_assert( result == SKY_OK, "Error getting block from blocks handle" );
		registerHandleClose( blockHandle );
		if( i > 0 ){
			memset(&hash1, 0, sizeof(GoString_));
			memset(&hash2, 0, sizeof(GoString_));
			memset(&hash, 0, sizeof(GoString_));
			result = SKY_Handle_Blocks_GetAt(blocksHandle, i - 1, 
								&previousBlockHandle);
			cr_assert( result == SKY_OK, "Error getting previous block from blocks handle" );
			registerHandleClose( previousBlockHandle );
			
			result = SKY_Handle_Block_GetHeadHash(previousBlockHandle, &hash1);
			cr_assert( result == SKY_OK, "Error getting previous block hash");
			registerMemCleanup( (void*)hash1.p );
			result = SKY_Handle_Block_GetPreviousBlockHash(blockHandle, &hash2);
			cr_assert( result == SKY_OK, "Error getting previous block hash");
			registerMemCleanup( (void*)hash2.p );
			
			cr_assert(eq(type(GoString_), hash1, hash2));
			freeRegisteredMemCleanup( (void*)hash1.p );
			freeRegisteredMemCleanup( (void*)hash2.p );
			
		}
		
		result = SKY_Handle_Block_GetHeadHash(blockHandle, &hash);
		cr_assert( result == SKY_OK, "Error getting previous block hash");
		registerMemCleanup( (void*)hash.p );
		
		_hash.p = hash.p;
		_hash.n = hash.n;
		result = SKY_api_Client_BlockByHash(&clientHandle, 
			_hash, &blockHandle2);
		cr_assert( result == SKY_OK, "SKY_api_Client_BlockByHash failed");
		registerHandleClose( blockHandle2 );
		
		result = SKY_Handle_Block_GetHeadSeq( blockHandle2, &seq );
		cr_assert( result == SKY_OK, "SKY_Handle_Block_GetHeadSeq failed");
		cr_assert(seq == i + start);
		
		equal = compareObjectsByHandle( blockHandle, blockHandle2 );
		cr_assert( equal == 1);
		
		freeRegisteredMemCleanup( (void*)hash.p );
		closeRegisteredHandle( blockHandle );
		closeRegisteredHandle( blockHandle2 );
		if( previousBlockHandle > 0 )
			closeRegisteredHandle( previousBlockHandle );
	}
	
	return blocksHandle;
}

Test(api_integration, TestStableBlocks) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;
	
	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );
	
	Handle progressHandle;
	result = SKY_api_Client_BlockchainProgress(&clientHandle, &progressHandle);
	cr_assert(result == SKY_OK, "SKY_api_Client_BlockchainProgress failed");
	registerHandleClose( progressHandle );
	GoUint64 lastNBlocks = 10;
	GoUint64 current;
	result = SKY_Handle_Progress_GetCurrent( progressHandle, &current );
	cr_assert(result == SKY_OK, "SKY_Handle_Progress_GetCurrent failed");
	cr_assert( current > lastNBlocks + 1, "Progress current must be greater than 10" );
	int tests_count = 7;
	test_blockn tests[] = {
		{
			"blocks-first-10.golden", 1, 10, 0
		},
		{
			"blocks-last-10.golden", current - lastNBlocks, current, 0
		},
		{
			"blocks-first-1.golden", 1, 1, 0
		},
		{
			"blocks-all.golden", 0, current, 0
		},
		{
			"blocks-end-less-than-start.golden", 10, 9, 0
		},
		{
			NULL, -10, 9, 1
		},
		{
			NULL, 10, -9, 1
		},
	};
	Handle blocksHandle;
	int equal;
	for(int i = 0; i < tests_count; i++){
		if( tests[i].failure ){
			result = SKY_api_Client_Blocks(&clientHandle, 
					tests[i].start, tests[i].end, &blocksHandle);
			cr_assert(result != SKY_OK, "SKY_api_Client_Blocks should have failed");
		} else {
			blocksHandle = testBlocks(clientHandle, 
					tests[i].start, tests[i].end);
			equal = compareObjectWithGoldenFile(blocksHandle, 
										tests[i].golden_file);
			cr_assert(equal == 1, "SKY_api_Client_Blocks returned a value different than expected.");
			closeRegisteredHandle( blocksHandle );
		}
	}
}