typedef struct{
    GoSlice_  Blocks;
} visor__ReadableBlocks;
typedef struct{
    GoSlice_  Txns;
} visor__TransactionResults;
typedef struct{
    visor__TransactionStatus Status;
    GoUint64_ Time;
    visor__ReadableTransaction Transaction;
} visor__TransactionResult;