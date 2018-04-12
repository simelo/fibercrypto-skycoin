
#ifndef SKYSTRUCTS_H
#define SKYSTRUCTS_H

/**
 * Go 8-bit signed integer values.
 */
typedef signed char GoInt8_;
/**
 * Go 8-bit unsigned integer values.
 */
typedef unsigned char GoUint8_;
/**
 * Go 16-bit signed integer values.
 */
typedef short GoInt16_;
/**
 * Go 16-bit unsigned integer values.
 */
typedef unsigned short GoUint16_;
/**
 * Go 32-bit signed integer values.
 */
typedef int GoInt32_;
/**
 * Go 32-bit unsigned integer values.
 */
typedef unsigned int GoUint32_;
/**
 * Go 64-bit signed integer values.
 */
typedef long long GoInt64_;
/**
 * Go 64-bit unsigned integer values.
 */
typedef unsigned long long GoUint64_;
/**
 * Go integer values aligned to the word size of the underlying architecture.
 */
typedef GoInt64_ GoInt_;
/**
 * Go unsigned integer values aligned to the word size of the underlying
 * architecture.
 */
typedef GoUint64_ GoUint_;
/**
 * Architecture-dependent type representing instances Go `uintptr` type.
 * Used as a generic representation of pointer types.
 */
typedef __SIZE_TYPE__ GoUintptr_;
/**
 * Go single precision 32-bits floating point values.
 */
typedef float GoFloat32_;
/**
 * Go double precision 64-bits floating point values.
 */
typedef double GoFloat64_;
/**
 * Instances of Go `complex` type.
 */
typedef float _Complex GoComplex64_;
/**
 * Instances of Go `complex` type.
 */
typedef double _Complex GoComplex128_;
typedef short bool;
typedef GoUint32_ error;
typedef GoUint32_ Handle;

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt._
*/
typedef char _check_for_64_bit_pointer_matchingGoInt[sizeof(void*)==64/8 ? 1:-1];

/**
 * Instances of Go `string` type.
 */
typedef struct {
  const char *p;    ///< Pointer to string characters buffer.
  GoInt_ n;         ///< String size not counting trailing `\0` char
                    ///< if at all included.
} GoString_;
/**
 * Instances of Go `map` type.
 */
typedef void *GoMap_;
/**
 * Instances of Go `chan` channel types.
 */
typedef void *GoChan_;
<<<<<<< HEAD
typedef struct { void *t; void *v; } GoInterface_;
typedef struct { void *data; GoInt_ len; GoInt_ cap; } GoSlice_;

#include "skytypes.gen.h"

typedef struct {
	GoMap_ Meta;
	GoSlice_ Entries;
} Wallet;

// TODO: Remove declarations below since they should generated and included by skytypes.gen.h

/*

=======
/**
 * Instances of Go interface types.
 */
typedef struct {
  void *t;      ///< Pointer to the information of the concrete Go type
                ///< bound to this interface reference.
  void *v;      ///< Pointer to the data corresponding to the value 
                ///< bound to this interface type.
} GoInterface_;
/**
 * Instances of Go slices
 */
typedef struct {
  void *data;   ///< Pointer to buffer containing slice data.
  GoInt_ len;   ///< Number of items stored in slice buffer
  GoInt_ cap;   ///< Maximum number of items that fits in this slice
                ///< considering allocated memory and item type's
                ///< size.
} GoSlice_;

/**
 * TODO: Document
 */
>>>>>>> remotes/github/olemis_t992_libskycoin_tests
typedef unsigned char Ripemd160[20];

/**
 * Addresses of SKY accounts
 */
typedef struct {
	unsigned char Version;  ///< Address version identifier.
                          ///< Used to differentiate testnet
                          ///< vs mainnet addresses, for instance.
	Ripemd160 Key;          ///< Address hash identifier.
} Address;

/**
 * Public key, 33-bytes long.
 */
typedef unsigned char PubKey[33];
/**
 * Container type suitable for storing a variable number of
 * public keys.
 */
typedef GoSlice_ PubKeySlice;
/**
 * Secret key, 32 bytes long.
 */
typedef unsigned char SecKey[32];
/**
 * Integrity checksum, 4-bytes long.
 */
typedef unsigned char Checksum[4];

/**
 * Structure used to specify amounts transferred in a transaction.
 */
typedef struct {
	GoString_ Addr; ///< Sender / receipient address.
	GoInt64_ Coins; ///< Amount transferred (e.g. measured in SKY)
} SendAmount;

/**
 * Memory handles returned back to the caller and manipulated
 * internally by API functions. Usually used to avoid type dependencies
 * with internal implementation types.
 */
typedef GoInt64_ Handle;

/**
 * Hash obtained using SHA256 algorithm, 32 bytes long.
 */
typedef unsigned char SHA256[32];
/**
 * Hash signed using a secret key, 65 bytes long.
 */
typedef unsigned char Sig[65];

/**
 * Skycoin transaction output.
 *
 * Instances are integral part of transactions included in blocks.
 */
typedef struct {
	Address Address;  ///< Receipient address.
	GoInt64_ Coins;   ///< Amount sent to the receipient address.
	GoInt64_ Hours;   ///< TODO: Document TransactionOutput.Hours
} TransactionOutput;

/**
 * Skycoin transaction.
 *
 * Instances of this struct are included in blocks.
 */
typedef struct {
	GoInt32_ Length;    ///< TODO: Document Transaction.Length
	GoInt8_  Type;      ///< TODO: Document Transaction.Type
	SHA256  InnerHash;  ///< TODO: Document Transaction.InnerHash

	GoSlice_ Sigs;      ///< TODO: Document Transaction.Sigs
	GoSlice_ In;        ///< TODO: Document Transaction.In
	GoSlice_ Out;       ///< TODO: Document Transaction.Out
} Transaction;

<<<<<<< HEAD

=======
/**
 * Internal representation of a Skycoin wallet.
 */
typedef struct {
	GoMap_ Meta;        ///< TODO: Document Wallet.Meta
	GoSlice_ Entries;   ///< TODO: Document Wallet.Entries
} Wallet;
>>>>>>> remotes/github/olemis_t992_libskycoin_tests

/**
 * Wallet entry.
 */
typedef struct {
	Address Address;    ///< Wallet address.
	PubKey  Public;     ///< Public key used to generate address.
	SecKey  Secret;     ///< Secret key used to generate address.
} Entry;

/**
 * TODO: Document UxBalance
 */
typedef struct {
	SHA256   Hash;      ///< TODO: Document
	GoInt64_ BkSeq;     ///< Block height corresponding to the
                      ///< moment balance calculation is performed at.
	Address  Address;   ///< Account holder address.
	GoInt64_ Coins;     ///< Coins amount (e.g. in SKY).
	GoInt64_ Hours;     ///< TODO: Document UxBalance.Hours
} UxBalance;

*/

#endif

