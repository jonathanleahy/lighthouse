package network

import "errors"

type (
	Audit struct {
		Domain   string
		Action   string
		DomainID DomainIDSupplier
		Ignore   bool
	}

	DomainIDSupplier func(responseParsed bool, responseBody []byte, responseCode int) string
)

func PreProcessedDomainID(domainID string) DomainIDSupplier {
	return func(responseParsed bool, responseBody []byte, responseCode int) string {
		return domainID
	}
}

func (a *Audit) validate() error {
	if a.Ignore {
		return nil
	}

	if len(a.Domain) == 0 {
		return errors.New(missingAuditDomain)
	}

	if len(a.Action) == 0 {
		return errors.New(missingAuditAction)
	}

	if a.DomainID == nil {
		return errors.New(missingAuditDomainID)
	}

	return nil
}

