package backend

type Backender interface {
	New(serialNumber int, expireDate, filename, dn string)
	UpdateStatus(serialNumber int, status string)
	Revoke(serialNumber int, date string, reason int)

	List(filter map[string]string) []CertificateResult
	Count(status string) int

	Exists(dn string) bool

	End()
}
