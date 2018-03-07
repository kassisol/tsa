package errors

import (
	"github.com/juliengk/stack/errors"
)

const (
	// Success indicates no error occurred.
	Success errors.Category = 1000 * iota // 0XXX

	// CertificateError indicates a fault in a certificate.
	CertificateError // 1XXX

	// PrivateKeyError indicates a fault in a private key.
	PrivateKeyError // 2XXX

	// CSRError indicates a problem with CSR parsing
	CSRError // 3XXX

	// RootError indicates a fault in a root.
	RootError // 4XXX

	// IntermediatesError indicates a fault in an intermediate.
	IntermediatesError // 5XXX

	// Serial Number
	SerialError // 6XXX

	// OCSPError indicates a problem with OCSP signing
	OCSPError // 7XXX

	// CertStoreError indicates a problem with the certificate store
	CertStoreError // 8XXX
)

// Parsing errors
const (
	Unknown      errors.Reason = iota // X000
	ReadFailed                        // X001
	DecodeFailed                      // X002
	ParseFailed                       // X003
)

// The following represent certificate non-parsing errors, and must be
// specified along with CertificateError.
const (
	// SelfSigned indicates that a certificate is self-signed and
	// cannot be used in the manner being attempted.
	SelfSigned errors.Reason = 100 * (iota + 1) // Code 11XX

	// VerifyFailed is an X.509 verification failure. The least two
	// significant digits of 12XX is determined as the actual x509
	// error is examined.
	VerifyFailed // Code 12XX

	// BadRequest indicates that the certificate request is invalid.
	BadRequest // Code 13XX

	// MissingSerial indicates that the profile specified
	// 'ClientProvidesSerialNumbers', but the SignRequest did not include a serial
	// number.
	MissingSerial // Code 14XX
)

const (
	certificateInvalid = 10 * (iota + 1) //121X
	unknownAuthority                     //122x
)

// The following represent private-key non-parsing errors, and must be
// specified with PrivateKeyError.
const (
	// Encrypted indicates that the private key is a PKCS #8 encrypted
	// private key. At this time, CFSSL does not support decrypting
	// these keys.
	Encrypted errors.Reason = 100 * (iota + 1) //21XX

	// NotRSAOrECC indicates that they key is not an RSA or ECC
	// private key; these are the only two private key types supported
	// at this time by CFSSL.
	NotRSA //22XX

	// KeyMismatch indicates that the private key does not match
	// the public key or certificate being presented with the key.
	KeyMismatch //23XX

	// GenerationFailed indicates that a private key could not
	// be generated.
	GenerationFailed //24XX

	// Unavailable indicates that a private key mechanism (such as
	// PKCS #11) was requested but support for that mechanism is
	// not available.
	Unavailable
)

// The following are Serial Number related errors, and should be
// specified with SerialError
const (
	IncrementFailed errors.Reason = 100 * (iota + 1) // 61XX
	WriteFailed                                      // X003
)

// The following are OCSP related errors, and should be
// specified with OCSPError
const (
	// IssuerMismatch ocurs when the certificate in the OCSP signing
	// request was not issued by the CA that this responder responds for.
	IssuerMismatch errors.Reason = 100 * (iota + 1) // 71XX

	// InvalidStatus occurs when the OCSP signing requests includes an
	// invalid value for the certificate status.
	InvalidStatus
)

// Certificate persistence related errors specified with CertStoreError
const (
	// InsertionFailed occurs when a SQL insert query failes to complete.
	InsertionFailed = 100 * (iota + 1)
	// RecordNotFound occurs when a SQL query targeting on one unique
	// record failes to update the specified row in the table.
	RecordNotFound

	RecordFound
)

// New returns an error that contains  an error code and message derived from
// the given category, reason. Currently, to avoid confusion, it is not
// allowed to create an error of category Success
func New(category errors.Category, reason errors.Reason) *errors.Error {
	errorCode := int(category) + int(reason)
	var msg string
	switch category {
	case CertificateError:
		switch reason {
		case Unknown:
			msg = "Unknown certificate error"
		case ReadFailed:
			msg = "Failed to read certificate"
		case DecodeFailed:
			msg = "Failed to decode certificate"
		case ParseFailed:
			msg = "Failed to parse certificate"
		case SelfSigned:
			msg = "Certificate is self signed"
		case VerifyFailed:
			msg = "Unable to verify certificate"
		case BadRequest:
			msg = "Invalid certificate request"
		case MissingSerial:
			msg = "Missing serial number in request"
		default:
			msg = errors.DefaultCategoryErrorString("CertificateError", reason)

		}
	case PrivateKeyError:
		switch reason {
		case Unknown:
			msg = "Unknown private key error"
		case ReadFailed:
			msg = "Failed to read private key"
		case DecodeFailed:
			msg = "Failed to decode private key"
		case ParseFailed:
			msg = "Failed to parse private key"
		case Encrypted:
			msg = "Private key is encrypted."
		case NotRSA:
			msg = "Private key algorithm is not RSA"
		case KeyMismatch:
			msg = "Private key does not match public key"
		case GenerationFailed:
			msg = "Failed to new private key"
		case Unavailable:
			msg = "Private key is unavailable"
		default:
			msg = errors.DefaultCategoryErrorString("PrivateKeyError", reason)
		}
	case CSRError:
		switch reason {
		case Unknown:
			msg = "CSR parsing failed due to unknown error"
		case ReadFailed:
			msg = "CSR file read failed"
		case ParseFailed:
			msg = "CSR Parsing failed"
		case DecodeFailed:
			msg = "CSR Decode failed"
		case BadRequest:
			msg = "CSR Bad request"
		default:
			msg = errors.DefaultCategoryErrorString("CSRError", reason)
		}
	case RootError:
		switch reason {
		case Unknown:
			msg = "Unknown root certificate error"
		case ReadFailed:
			msg = "Failed to read root certificate"
		case DecodeFailed:
			msg = "Failed to decode root certificate"
		case ParseFailed:
			msg = "Failed to parse root certificate"
		default:
			msg = errors.DefaultCategoryErrorString("RootError", reason)
		}
	case IntermediatesError:
		switch reason {
		case Unknown:
			msg = "Unknown intermediate certificate error"
		case ReadFailed:
			msg = "Failed to read intermediate certificate"
		case DecodeFailed:
			msg = "Failed to decode intermediate certificate"
		case ParseFailed:
			msg = "Failed to parse intermediate certificate"
		default:
			msg = errors.DefaultCategoryErrorString("IntermediatesError", reason)
		}
	case SerialError:
		switch reason {
		case Unknown:
			msg = "Unknown serial number error"
		case ReadFailed:
			msg = "Failed to read serial number file"
		case IncrementFailed:
			msg = "Failed to increment serial number"
		case WriteFailed:
			msg = "Failed to write to serial number file"
		default:
			msg = errors.DefaultCategoryErrorString("IntermediatesError", reason)
		}
	case OCSPError:
		switch reason {
		case ReadFailed:
			msg = "No certificate provided"
		case IssuerMismatch:
			msg = "Certificate not issued by this issuer"
		case InvalidStatus:
			msg = "Invalid revocation status"
		default:
			msg = errors.DefaultCategoryErrorString("OCSPError", reason)
		}
	case CertStoreError:
		switch reason {
		case Unknown:
			msg = "Certificate store action failed due to unknown error"
		case RecordFound:
			msg = "Certificate already exists and is valid"
		default:
			msg = errors.DefaultCategoryErrorString("CertStoreError", reason)
		}
	default:
		msg = errors.DefaultTypeErrorString(category)
	}

	return &errors.Error{ErrorCode: errorCode, Message: msg}
}
