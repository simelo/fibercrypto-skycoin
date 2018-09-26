%include "arrays_csharp.i"
%include cpointer.i
// %pointer_functions(cipher_PubKey, cipher_PubKeyp);
// %pointer_functions(cipher_SecKey, cipher_SecKeyp);
// %pointer_functions(cipher__Ripemd160, cipher__Ripemd160p);
// %pointer_functions(cipher_Sig, cipher_Sigp);
%pointer_functions(GoSlice, GoSlicep);
%pointer_functions(_GoString_, GoStringp);
%pointer_functions(int, intp);
%pointer_functions(Transaction__Handle, Transaction__Handlep);
// %pointer_functions(byte, bytep);

%inline %{
#include "json.h"
	//Define function SKY_handle_close to avoid including libskycoin.h
void SKY_handle_close(Handle p0);

int MEMPOOLIDX = 0;
void *MEMPOOL[1024 * 256];

int JSONPOOLIDX = 0;
json_value* JSON_POOL[128];

int HANDLEPOOLIDX = 0;
Handle HANDLE_POOL[128];

typedef struct {
	Client__Handle client;
	WalletResponse__Handle wallet;
} wallet_register;

int WALLETPOOLIDX = 0;
wallet_register WALLET_POOL[64];

int stdout_backup;
int pipefd[2];

void * registerMemCleanup(void *p) {
	int i;
	for (i = 0; i < MEMPOOLIDX; i++) {
		if(MEMPOOL[i] == NULL){
			MEMPOOL[i] = p;
			return p;
		}
	}
	MEMPOOL[MEMPOOLIDX++] = p;
	return p;
}

void freeRegisteredMemCleanup(void *p){
	int i;
	for (i = 0; i < MEMPOOLIDX; i++) {
		if(MEMPOOL[i] == p){
			free(p);
			MEMPOOL[i] = NULL;
			break;
		}
	}
}

int registerJsonFree(void *p){
	int i;
	for (i = 0; i < JSONPOOLIDX; i++) {
		if(JSON_POOL[i] == NULL){
			JSON_POOL[i] = p;
			return i;
		}
	}
	JSON_POOL[JSONPOOLIDX++] = p;
	return JSONPOOLIDX-1;
}

void freeRegisteredJson(void *p){
	int i;
	for (i = 0; i < JSONPOOLIDX; i++) {
		if(JSON_POOL[i] == p){
			JSON_POOL[i] = NULL;
			json_value_free( (json_value*)p );
			break;
		}
	}
}

int registerWalletClean(Client__Handle clientHandle,
						WalletResponse__Handle walletHandle){
	int i;
	for (i = 0; i < WALLETPOOLIDX; i++) {
		if(WALLET_POOL[i].wallet == 0 && WALLET_POOL[i].client == 0){
			WALLET_POOL[i].wallet = walletHandle;
			WALLET_POOL[i].client = clientHandle;
			return i;
		}
	}
	WALLET_POOL[WALLETPOOLIDX].wallet = walletHandle;
	WALLET_POOL[WALLETPOOLIDX].client = clientHandle;
	return WALLETPOOLIDX++;
}

int registerHandleClose(Handle handle){
	int i;
	for (i = 0; i < HANDLEPOOLIDX; i++) {
		if(HANDLE_POOL[i] == 0){
			HANDLE_POOL[i] = handle;
			return i;
		}
	}
	HANDLE_POOL[HANDLEPOOLIDX++] = handle;
	return HANDLEPOOLIDX - 1;
}

void closeRegisteredHandle(Handle handle){
	int i;
	for (i = 0; i < HANDLEPOOLIDX; i++) {
		if(HANDLE_POOL[i] == handle){
			HANDLE_POOL[i] = 0;
			SKY_handle_close(handle);
			break;
		}
	}
}

void cleanupWallet(Client__Handle client, WalletResponse__Handle wallet){
	int result;
	GoString_ strWalletDir;
	GoString_ strFileName;
	memset(&strWalletDir, 0, sizeof(GoString_));
	memset(&strFileName, 0, sizeof(GoString_));


	result = SKY_api_Handle_Client_GetWalletDir(client, &strWalletDir);
	if( result != 0 ){
		return;
	}
	result = SKY_api_Handle_Client_GetWalletFileName(wallet, &strFileName);
	if( result != 0 ){
		free( (void*)strWalletDir.p );
		return;
	}
	char fullPath[128];
	if( strWalletDir.n + strFileName.n < 126){
		strcpy( fullPath, strWalletDir.p );
		if( fullPath[0] == 0 || fullPath[strlen(fullPath) - 1] != '/' )
			strcat(fullPath, "/");
		strcat( fullPath, strFileName.p );
		result = unlink( fullPath );
		if( strlen(fullPath) < 123 ){
			strcat( fullPath, ".bak" );
			result = unlink( fullPath );
		}
	}
	GoString str = { strFileName.p, strFileName.n };
	result = SKY_api_Client_UnloadWallet( client, str );
	GoString strFullPath = { fullPath, strlen(fullPath) };
	free( (void*)strWalletDir.p );
	free( (void*)strFileName.p );
}

void cleanRegisteredWallet(
			Client__Handle client,
			WalletResponse__Handle wallet){

	int i;
	for (i = 0; i < WALLETPOOLIDX; i++) {
		if(WALLET_POOL[i].wallet == wallet && WALLET_POOL[i].client == client){
			WALLET_POOL[i].wallet = 0;
			WALLET_POOL[i].client = 0;
			cleanupWallet( client, wallet );
			return;
		}
	}
}

void cleanupMem() {
	int i;

	for (i = 0; i < WALLETPOOLIDX; i++) {
		if(WALLET_POOL[i].client != 0 && WALLET_POOL[i].wallet != 0){
			cleanupWallet( WALLET_POOL[i].client, WALLET_POOL[i].wallet );
		}
	}

  void **ptr;
  for (i = MEMPOOLIDX, ptr = MEMPOOL; i; --i) {
	if( *ptr )
		free(*ptr);
	ptr++;
  }
  for (i = JSONPOOLIDX, ptr = (void*)JSON_POOL; i; --i) {
	if( *ptr )
		json_value_free(*ptr);
	ptr++;
  }
  for (i = 0; i < HANDLEPOOLIDX; i++) {
	  if( HANDLE_POOL[i] )
		SKY_handle_close(HANDLE_POOL[i]);
  }
}


void setup(void) {
	srand ((unsigned int) time (NULL));
}

void teardown(void) {
	cleanupMem();
}

// TODO: Move to libsky_io.c
void fprintbuff(FILE *f, void *buff, size_t n) {
  unsigned char *ptr = (unsigned char *) buff;
  fprintf(f, "[ ");
  for (; n; --n, ptr++) {
    fprintf(f, "%02d ", *ptr);
  }
  fprintf(f, "]");
}

int parseBoolean(const char* str, int length){
	int result = 0;
	if(length == 1){
		result = str[0] == '1' || str[0] == 't' || str[0] == 'T';
	} else {
		result = strncmp(str, "true", length) == 0 ||
			strncmp(str, "True", length) == 0 ||
			strncmp(str, "TRUE", length) == 0;
	}
	return result;
}

void toGoString(GoString_ *s, GoString *r){
GoString * tmp = r;

  *tmp = (*(GoString *) s);
}

int copySlice(GoSlice_* pdest, GoSlice_* psource, int elem_size){
  pdest->len = psource->len;
  pdest->cap = psource->len;
  int size = pdest->len * elem_size;
  pdest->data = malloc(size);
	if( pdest->data == NULL )
		return 1;
  registerMemCleanup( pdest->data );
  memcpy(pdest->data, psource->data, size );
	return 0;
}



int concatSlices(GoSlice_* slice1, GoSlice_* slice2, int elem_size, GoSlice_* result){
	int size1 = slice1->len;
	int size2 = slice2->len;
	int size = size1 + size2;
	if (size <= 0)
		return 1;
	void* data = malloc(size * elem_size);
	if( data == NULL )
		return 1;
	registerMemCleanup( data );
	result->data = data;
	result->len = size;
	result->cap = size;
	char* p = data;
	if(size1 > 0){
		memcpy( p, slice1->data, size1 * elem_size );
		p += (elem_size * size1);
	}
	if(size2 > 0){
		memcpy( p, slice2->data, size2 * elem_size );
	}
	return 0;
}
    void parseJsonMetaData(char *metadata, int *n, int *r, int *p, int *keyLen)
{
	*n = *r = *p = *keyLen = 0;
	int length = strlen(metadata);
	int openingQuote = -1;
	const char *keys[] = {"n", "r", "p", "keyLen"};
	int keysCount = 4;
	int keyIndex = -1;
	int startNumber = -1;
	for (int i = 0; i < length; i++)
	{
		if (metadata[i] == '\"')
		{
			startNumber = -1;
			if (openingQuote >= 0)
			{
				keyIndex = -1;
				metadata[i] = 0;
				for (int k = 0; k < keysCount; k++)
				{
					if (strcmp(metadata + openingQuote + 1, keys[k]) == 0)
					{
						keyIndex = k;
					}
				}
				openingQuote = -1;
			}
			else
			{
				openingQuote = i;
			}
		}
		else if (metadata[i] >= '0' && metadata[i] <= '9')
		{
			if (startNumber < 0)
				startNumber = i;
		}
		else if (metadata[i] == ',')
		{
			if (startNumber >= 0)
			{
				metadata[i] = 0;
				int number = atoi(metadata + startNumber);
				startNumber = -1;
				if (keyIndex == 0)
					*n = number;
				else if (keyIndex == 1)
					*r = number;
				else if (keyIndex == 2)
					*p = number;
				else if (keyIndex == 3)
					*keyLen = number;
			}
		}
		else
		{
			startNumber = -1;
		}
	}
}

int cutSlice(GoSlice_* slice, int start, int end, int elem_size, GoSlice_* result){
	int size = end - start;
	if( size <= 0)
		return 1;
	void* data = malloc(size * elem_size);
	if( data == NULL )
		return 1;
	registerMemCleanup( data );
	result->data = data;
	result->len = size;
	result->cap = size;
	char* p = slice->data;
	p += (elem_size * start);
	memcpy( data, p, elem_size * size );
	return 0;
}

coin__Transaction* makeEmptyTransaction(Transaction__Handle* handle){
  int result;
  coin__Transaction* ptransaction = NULL;
  result  = SKY_coin_Create_Transaction(handle);
   registerHandleClose(*handle);
  result = SKY_coin_GetTransactionObject( *handle, &ptransaction );
    return ptransaction;
}

int makeAddress(cipher__Address* paddress){
  cipher__PubKey pubkey;
  cipher__SecKey seckey;
  cipher__Address address;
  int result;

  result = SKY_cipher_GenerateKeyPair(&pubkey, &seckey);
  if(result != 0){
	  return 1;
  }

  result = SKY_cipher_AddressFromPubKey( &pubkey, paddress );
    if(result != 0){
	  return 1;
  }
  return result;
}
    %}