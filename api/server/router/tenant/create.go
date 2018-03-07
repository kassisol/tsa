package tenant

import (
	"fmt"
	"net/http"
	"os"

	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/errors"
	"github.com/juliengk/go-cert/helpers"
	"github.com/juliengk/go-cert/pkix"
	"github.com/juliengk/stack/jsonapi"
	apierr "github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	valid "github.com/kassisol/tsa/pkg/validation"
	"github.com/labstack/echo"
)

func panicIfEmpty(name, value string) {
	if len(value) == 0 {
		panic(fmt.Sprintf("%s cannot be empty", name))
	}
}

func CreateHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		e := apierr.New(apierr.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}
	defer s.End()

	defer func() {
		if r := recover(); r != nil {
			if len(cfg.TenantPath) > 0 {
				os.RemoveAll(cfg.TenantPath)
			}

			err := fmt.Sprintf("%s", r)
			res := jsonapi.NewErrorResponse(1000, err)

			api.JSON(c, http.StatusUnprocessableEntity, res)
		}
	}()

	// Get POST data
	tnt := new(types.Tenant)

	if err := c.Bind(tnt); err != nil {
		r := jsonapi.NewErrorResponse(1000, "Cannot unmarshal JSON")

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	// Validation
	// DV - Name

	// DV - Groups

	// DV - CA Type
	if err := valid.IsValidCAType(tnt.CA.Type); err != nil {
		panic(err)
	}

	// DV - Country
	if err = valid.IsValidCountry(tnt.CA.Country); err != nil {
		panic(err)
	}

	// DV - State
	panicIfEmpty("State", tnt.CA.State)

	// DV - Locality
	panicIfEmpty("Locality", tnt.CA.Locality)

	// DV - Organization
	panicIfEmpty("Organization", tnt.CA.Organization)

	// DV - Organizational Unit
	panicIfEmpty("Organizational", tnt.CA.OrganizationalUnit)

	if err := helpers.IsValidCAOrgUnit(tnt.CA.OrganizationalUnit); err != nil {
		panic(err)
	}

	// Save data to DB
	caou := helpers.UpdateOrgUnitLabel(tnt.CA.OrganizationalUnit)
	cacn := helpers.UpdateCommonNameLabel(tnt.CA.Type, tnt.CA.OrganizationalUnit)

	if err = s.AddTenant(tnt.Name, tnt.AuthGroups, tnt.CA.Type, tnt.CA.Duration, tnt.CA.Expire, tnt.CA.Country, tnt.CA.State, tnt.CA.Locality, tnt.CA.Organization, caou, cacn); err != nil {
		panic(err)
	}

	filters := map[string]string{
		"name": tnt.Name,
	}
	i := s.ListTenants(filters)

	// Initialize CA
	if err := cfg.Tenant(tnt.Name); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	caSubject := pkix.NewSubject(tnt.CA.Country, tnt.CA.State, tnt.CA.Locality, tnt.CA.Organization, caou, cacn)

	caNDN := pkix.NewDNSNames()
	caNE := pkix.NewEmails()
	caIP := pkix.NewIPs()

	caAltnames := pkix.NewSubjectAltNames(*caNDN, *caNE, *caIP)

	caDate := ca.CreateDate(tnt.CA.Duration)
	caSN := 1

	caTemplate, err := ca.CreateTemplate(true, caSubject, caAltnames, caDate, caSN, "")
	if err != nil {
		panic(err)
	}

	if _, err = ca.InitCA(cfg.TenantPath, caTemplate); err != nil {
		panic(err)
	}

	ca.CreateCRLFile(cfg.CA.CRLFile)

	// Response
	response := jsonapi.NewSuccessResponse(i)

	return api.JSON(c, http.StatusCreated, response)
}
