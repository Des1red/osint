package knowledge

import (
	"osint/internal/enrich"
)

//
// Shared enrichment.
//

type EnrichmentResult = enrich.EnrichmentResult

type UsernameReference = enrich.UsernameReference

type EmailReference = enrich.EmailReference

type PhoneReference = enrich.PhoneReference

type SocialReference = enrich.SocialReference

type ExternalLink = enrich.ExternalLink

type LocationReference = enrich.LocationReference

type OrganizationReference = enrich.OrganizationReference

type EmploymentReference = enrich.EmploymentReference

type EducationReference = enrich.EducationReference

type PersonReference = enrich.PersonReference

type RelativeSearchResult = enrich.RelativeSearchResult

type RelatedAccountReference = enrich.RelatedAccountReference

//
// Directory results.
//

type DirectoryResult = enrich.DirectoryResult

type DirectoryEntry = enrich.DirectoryEntry

type DirectoryField = enrich.DirectoryField
