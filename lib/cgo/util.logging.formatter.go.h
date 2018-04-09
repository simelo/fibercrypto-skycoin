typedef struct{
    GoString_ InfoLevelStyle;
    GoString_ WarnLevelStyle;
    GoString_ ErrorLevelStyle;
    GoString_ FatalLevelStyle;
    GoString_ PanicLevelStyle;
    GoString_ DebugLevelStyle;
    GoString_ PrefixStyle;
    GoString_ TimestampStyle;
    GoString_ CallContextStyle;
    GoString_ CriticalStyle;
}ColorScheme;
typedef struct{
    Handle InfoLevelColor;
    Handle WarnLevelColor;
    Handle ErrorLevelColor;
    Handle FatalLevelColor;
    Handle PanicLevelColor;
    Handle DebugLevelColor;
    Handle PrefixColor;
    Handle TimestampColor;
    Handle CallContextColor;
    Handle CriticalColor;
}compiledColorScheme;
