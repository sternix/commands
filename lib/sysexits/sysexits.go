package sysexits

const (
	OK		int = 0		/* successful termination */
	BASE		int = 64	/* base value for error messages */
	USAGE		int = 64	/* command line usage error */
	DATAERR		int = 65	/* data format error */
	NOINPUT		int = 66	/* cannot open input */
	NOUSER		int = 67	/* addressee unknown */
	NOHOST		int = 68	/* host name unknown */
	UNAVAILABLE	int = 69	/* service unavailable */
	SOFTWARE	int = 70	/* internal software error */
	OSERR		int = 71	/* system error (e.g., can't fork) */
	OSFILE		int = 72	/* critical OS file missing */
	CANTCREAT	int = 73	/* can't create (user) output file */
	IOERR		int = 74	/* input/output error */
	TEMPFAIL	int = 75	/* temp failure; user is invited to retry */
	PROTOCOL	int = 76	/* remote error in protocol */
	NOPERM		int = 77	/* permission denied */
	CONFIG		int = 78	/* configuration error */
	MAX		int = 78	/* maximum listed value */
)
