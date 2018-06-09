/**
 * Intermediate representation of a UxOut for sorting and spend choosing.
 */
typedef struct {
	cipher__SHA256   Hash;     ///< Hash of underlying UxOut.
	GoInt64_ BkSeq;           ///< Block height corresponding to the
                            ///< moment balance calculation is performed at.
	cipher__Address  Address;  ///< Account holder address.
	GoInt64_ Coins;           ///< Coins amount (e.g. in SKY).
	GoInt64_ Hours;           ///< Balance of Coin Hours generated by underlying UxOut, depending on UxOut's head time.
} wallet__UxBalance;

/**
 * Internal representation of a Skycoin wallet.
 */
typedef struct {
	GoMap_ Meta;        ///< Records items that are not deterministic, like filename, lable, wallet type, secrets, etc.
	GoSlice_ Entries;   ///< Entries field stores the address entries that are deterministically generated from seed.
} wallet__Wallet;
