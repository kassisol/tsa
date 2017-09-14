package system

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/errors"
	"github.com/juliengk/go-cert/helpers"
	"github.com/juliengk/go-cert/pkix"
	"github.com/juliengk/go-utils/validation"
	"github.com/juliengk/stack/jsonapi"
	apierr "github.com/kassisol/tsa/api/errors"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	valid "github.com/kassisol/tsa/pkg/validation"
	"github.com/labstack/echo"
)

var ErrCountryCodeLength = fmt.Errorf("Country should be a 2 letters code")

func CAInitHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		e := errors.New(apierr.DatabaseError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}
	defer s.End()

	if len(s.ListConfigs("ca")) > 0 {
		r := jsonapi.NewErrorResponse(1000, "Initialization already done")

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	defer func() {
		if r := recover(); r != nil {
			os.RemoveAll(cfg.CA.Dir.Root)

			s.RemoveConfig("ca_type", "ALL")
			s.RemoveConfig("ca_duration", "ALL")
			s.RemoveConfig("ca_country", "ALL")
			s.RemoveConfig("ca_state", "ALL")
			s.RemoveConfig("ca_locality", "ALL")
			s.RemoveConfig("ca_org", "ALL")
			s.RemoveConfig("ca_ou", "ALL")
			s.RemoveConfig("ca_cn", "ALL")

			err := fmt.Sprintf("%s", r)
			res := jsonapi.NewErrorResponse(1000, err)

			api.JSON(c, http.StatusUnprocessableEntity, res)
		}
	}()

	// Get POST data
	ci := new(types.CertificationAuthority)

	if err := c.Bind(ci); err != nil {
		r := jsonapi.NewErrorResponse(1000, "Cannot unmarshal JSON")

		return api.JSON(c, http.StatusInternalServerError, r)
	}

	// Validation
	// DV - CA Type
	if err := valid.IsValidCAType(ci.Type); err != nil {
		panic(err)
	}

	// DV - Country
	if len(ci.Country) != 2 {
		panic(ErrCountryCodeLength)
	}

	for _, c := range ci.Country {
		if err := validation.IsUpper(string(c)); err != nil {
			cErr := fmt.Errorf("Country: %v", err)
			panic(cErr)
		}
	}

	// DV - Organizational Unit
	if err := helpers.IsValidCAOrgUnit(ci.OrganizationalUnit); err != nil {
		panic(err)
	}

	// Save datas to DB
	caou := helpers.UpdateOrgUnitLabel(ci.OrganizationalUnit)
	cacn := helpers.UpdateCommonNameLabel(ci.Type, ci.OrganizationalUnit)

	s.AddConfig("ca_type", ci.Type)
	s.AddConfig("ca_duration", strconv.Itoa(ci.Duration))
	s.AddConfig("ca_country", ci.Country)
	s.AddConfig("ca_state", ci.State)
	s.AddConfig("ca_locality", ci.Locality)
	s.AddConfig("ca_org", ci.Organization)
	s.AddConfig("ca_ou", caou)
	s.AddConfig("ca_cn", cacn)

	// Initialize CA
	caSubject := pkix.NewSubject(ci.Country, ci.State, ci.Locality, ci.Organization, caou, cacn)

	caNDN := pkix.NewDNSNames()
	caNE := pkix.NewEmails()
	caIP := pkix.NewIPs()

	caAltnames := pkix.NewSubjectAltNames(*caNDN, *caNE, *caIP)

	caDate := ca.CreateDate(ci.Duration)
	caSN := 1

	caTemplate, err := ca.CreateTemplate(true, caSubject, caAltnames, caDate, caSN, "")
	if err != nil {
		panic(err)
	}

	if _, err = ca.InitCA(cfg.App.Dir.Root, caTemplate); err != nil {
		panic(err)
	}

	ca.CreateCRLFile(cfg.CA.CRLFile)

	// Response
	response := jsonapi.NewSuccessResponse(ci)

	return api.JSON(c, http.StatusCreated, response)
}
