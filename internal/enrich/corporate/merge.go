package corporate

import (
	"strings"

	"osint/internal/enrich/corporate/opencorp"
	identitymatch "osint/internal/enrich/match"
	"osint/internal/enrich/model"
)

func mergeOpenCorp(
	result model.EnrichmentResult,
	fullName string,
	openCorpResult opencorp.Result,
) (
	model.EnrichmentResult,
	int,
) {
	retained :=
		0

	for _, officer := range openCorpResult.Officers {

		//
		// OpenCorporates returns name-search
		// candidates.
		//
		// Keep only records whose officer name
		// actually matches the requested identity.
		//
		if !identitymatch.IdentityPresent(
			fullName,
			officer.Name,
		) {

			continue
		}

		retained++

		result =
			mergeOfficer(
				result,
				officer,
			)
	}

	return result,
		retained
}

func mergeOfficer(
	result model.EnrichmentResult,
	officer opencorp.Officer,
) model.EnrichmentResult {
	source :=
		strings.TrimSpace(
			officer.ProfileURL,
		)

	if source == "" {

		source =
			"OpenCorporates"
	}

	company :=
		strings.TrimSpace(
			officer.Company,
		)

	position :=
		strings.TrimSpace(
			officer.Position,
		)

	jurisdiction :=
		strings.TrimSpace(
			officer.Jurisdiction,
		)

	startDate :=
		strings.TrimSpace(
			officer.StartDate,
		)

	endDate :=
		strings.TrimSpace(
			officer.EndDate,
		)

	//
	// A same-name corporate record is useful
	// evidence, but by itself does not prove that
	// the requested person owns the record.
	//
	// Preserve it as unattributed until stronger
	// identity evidence exists.
	//
	if company != "" {

		result.UnattributedOrganizations =
			append(
				result.UnattributedOrganizations,
				model.OrganizationReference{
					Organization: company,

					Source: source,
				},
			)

		result.UnattributedEmployment =
			append(
				result.UnattributedEmployment,
				model.EmploymentReference{
					Title: position,

					Organization: company,

					StartDate: startDate,

					EndDate: endDate,

					Current: endDate == "",

					Summary: corporateSummary(
						officer,
					),

					Source: source,
				},
			)
	}

	if jurisdiction != "" {

		result.UnattributedLocations =
			append(
				result.UnattributedLocations,
				model.LocationReference{
					Location: jurisdiction,

					Source: source,
				},
			)
	}

	profileURL :=
		strings.TrimSpace(
			officer.ProfileURL,
		)

	if profileURL != "" {

		result.UnattributedLinks =
			append(
				result.UnattributedLinks,
				model.ExternalLink{
					URL: profileURL,

					Category: "Corporate Registry",

					Source: source,
				},
			)
	}

	companyURL :=
		strings.TrimSpace(
			officer.CompanyURL,
		)

	if companyURL != "" {

		result.UnattributedLinks =
			append(
				result.UnattributedLinks,
				model.ExternalLink{
					URL: companyURL,

					Category: "Corporate Registry",

					Source: source,
				},
			)
	}

	return result
}

func corporateSummary(
	officer opencorp.Officer,
) string {
	name :=
		strings.TrimSpace(
			officer.Name,
		)

	if name == "" {

		return "OpenCorporates officer candidate"
	}

	return "OpenCorporates officer candidate for " +
		name
}
