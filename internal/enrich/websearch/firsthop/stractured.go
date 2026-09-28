package firsthop

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type employmentEvidence struct {
	Title string

	Organization string

	StartDate string

	EndDate string

	Current bool

	Summary string
}

type educationEvidence struct {
	School string

	Degrees []string

	Majors []string

	StartDate string

	EndDate string

	Summary string
}

type structuredPageEvidence struct {
	References []string

	URLs []string

	Organizations []string

	Locations []string

	Employment []employmentEvidence

	Education []educationEvidence
}

type structuredCollector struct {
	result structuredPageEvidence

	references map[string]struct{}

	urls map[string]struct{}

	organizations map[string]struct{}

	locations map[string]struct{}

	employment map[string]struct{}

	education map[string]struct{}
}

func structuredEvidence(
	document *goquery.Document,
	fullName string,
	pageURL string,
) structuredPageEvidence {
	if document == nil {
		return structuredPageEvidence{}
	}

	collector :=
		newStructuredCollector()

	firstName,
		lastName,
		hasIdentity :=
		fullNameTokens(
			fullName,
		)

	allowStandaloneOrganization :=
		standaloneOrganizationPage(
			pageURL,
		)

	document.Find(
		`script[type="application/ld+json"]`,
	).Each(
		func(
			_ int,
			selection *goquery.Selection,
		) {
			raw :=
				strings.TrimSpace(
					selection.Text(),
				)

			if raw == "" {
				return
			}

			var value any

			if err :=
				json.Unmarshal(
					[]byte(
						raw,
					),
					&value,
				); err != nil {

				return
			}

			walkStructured(
				value,
				firstName,
				lastName,
				hasIdentity,
				allowStandaloneOrganization,
				collector,
			)
		},
	)

	return collector.result
}

func newStructuredCollector() *structuredCollector {
	return &structuredCollector{
		references: make(
			map[string]struct{},
		),

		urls: make(
			map[string]struct{},
		),

		organizations: make(
			map[string]struct{},
		),

		locations: make(
			map[string]struct{},
		),

		employment: make(
			map[string]struct{},
		),

		education: make(
			map[string]struct{},
		),
	}
}

func walkStructured(
	value any,
	firstName string,
	lastName string,
	hasIdentity bool,
	allowStandaloneOrganization bool,
	collector *structuredCollector,
) {
	switch node :=
		value.(type) {

	case []any:

		for _, child := range node {

			walkStructured(
				child,
				firstName,
				lastName,
				hasIdentity,
				allowStandaloneOrganization,
				collector,
			)
		}

	case map[string]any:

		if schemaType(
			node["@type"],
			"Person",
		) &&
			hasIdentity &&
			containsNameTokens(
				entityName(
					node,
				),
				firstName,
				lastName,
			) {

			collectPersonNode(
				node,
				collector,
			)
		}

		if allowStandaloneOrganization &&
			organizationType(
				node["@type"],
			) {

			collectOrganizationNode(
				node,
				collector,
			)
		}

		for _, child := range node {

			walkStructured(
				child,
				firstName,
				lastName,
				hasIdentity,
				allowStandaloneOrganization,
				collector,
			)
		}
	}
}

func collectPersonNode(
	node map[string]any,
	collector *structuredCollector,
) {
	collectContactFields(
		node,
		collector,
	)

	collectURLFields(
		node,
		collector,
	)

	collectAddress(
		node["address"],
		collector,
	)

	collectPlace(
		node["homeLocation"],
		collector,
	)

	collectPlace(
		node["workLocation"],
		collector,
	)

	jobTitle :=
		firstTextValue(
			node["jobTitle"],
		)

	employmentBefore :=
		len(
			collector.result.Employment,
		)

	collectEmploymentValue(
		node["worksFor"],
		jobTitle,
		collector,
	)

	if len(
		collector.result.Employment,
	) == employmentBefore &&
		jobTitle != "" {

		collector.addEmployment(
			employmentEvidence{
				Title: jobTitle,
			},
		)
	}

	collectOrganizationRelationship(
		node["affiliation"],
		collector,
	)

	collectOrganizationRelationship(
		node["memberOf"],
		collector,
	)

	collectEducationValue(
		node["alumniOf"],
		collector,
	)
}

func collectOrganizationNode(
	node map[string]any,
	collector *structuredCollector,
) {
	name :=
		entityName(
			node,
		)

	if name != "" {
		collector.addOrganization(
			name,
		)
	}

	collectContactFields(
		node,
		collector,
	)

	collectURLFields(
		node,
		collector,
	)

	collectAddress(
		node["address"],
		collector,
	)

	collectPlace(
		node["location"],
		collector,
	)
}

func collectContactFields(
	node map[string]any,
	collector *structuredCollector,
) {
	for _, value := range textValues(
		node["email"],
	) {

		collector.addReference(
			value,
		)
	}

	for _, value := range textValues(
		node["telephone"],
	) {

		collector.addReference(
			value,
		)
	}

	collectContactPoint(
		node["contactPoint"],
		collector,
	)
}

func collectContactPoint(
	value any,
	collector *structuredCollector,
) {
	switch node :=
		value.(type) {

	case []any:

		for _, child := range node {

			collectContactPoint(
				child,
				collector,
			)
		}

	case map[string]any:

		for _, email := range textValues(
			node["email"],
		) {

			collector.addReference(
				email,
			)
		}

		for _, phone := range textValues(
			node["telephone"],
		) {

			collector.addReference(
				phone,
			)
		}

		collectURLFields(
			node,
			collector,
		)
	}
}

func collectURLFields(
	node map[string]any,
	collector *structuredCollector,
) {
	for _, value := range textValues(
		node["url"],
	) {

		collector.addURL(
			value,
		)
	}

	for _, value := range textValues(
		node["sameAs"],
	) {

		collector.addURL(
			value,
		)
	}
}

func collectEmploymentValue(
	value any,
	defaultTitle string,
	collector *structuredCollector,
) {
	switch node :=
		value.(type) {

	case nil:

		return

	case string:

		organization :=
			strings.TrimSpace(
				node,
			)

		if organization == "" {
			return
		}

		collector.addOrganization(
			organization,
		)

		collector.addEmployment(
			employmentEvidence{
				Title: defaultTitle,

				Organization: organization,
			},
		)

	case []any:

		for _, child := range node {

			collectEmploymentValue(
				child,
				defaultTitle,
				collector,
			)
		}

	case map[string]any:

		//
		// Some JSON-LD represents employment
		// through a Role object.
		//
		if schemaType(
			node["@type"],
			"Role",
		) {

			title :=
				firstTextValue(
					node["roleName"],
				)

			if title == "" {
				title =
					defaultTitle
			}

			startDate :=
				firstTextValue(
					node["startDate"],
				)

			endDate :=
				firstTextValue(
					node["endDate"],
				)

			organizations :=
				organizationNames(
					node["worksFor"],
				)

			if len(organizations) == 0 {

				organizations =
					organizationNames(
						node["memberOf"],
					)
			}

			if len(organizations) == 0 {

				organizations =
					organizationNames(
						node["affiliation"],
					)
			}

			if len(organizations) == 0 {

				if title != "" {

					collector.addEmployment(
						employmentEvidence{
							Title: title,

							StartDate: startDate,

							EndDate: endDate,
						},
					)
				}

				return
			}

			for _, organization := range organizations {

				collector.addOrganization(
					organization,
				)

				collector.addEmployment(
					employmentEvidence{
						Title: title,

						Organization: organization,

						StartDate: startDate,

						EndDate: endDate,
					},
				)
			}

			return
		}

		if organizationType(
			node["@type"],
		) {

			collectOrganizationNode(
				node,
				collector,
			)
		}

		organization :=
			entityName(
				node,
			)

		if organization == "" {
			return
		}

		collector.addOrganization(
			organization,
		)

		collector.addEmployment(
			employmentEvidence{
				Title: defaultTitle,

				Organization: organization,
			},
		)
	}
}

func collectOrganizationRelationship(
	value any,
	collector *structuredCollector,
) {
	for _, organization := range organizationNames(
		value,
	) {

		collector.addOrganization(
			organization,
		)
	}
}

func collectEducationValue(
	value any,
	collector *structuredCollector,
) {
	switch node :=
		value.(type) {

	case nil:

		return

	case string:

		school :=
			strings.TrimSpace(
				node,
			)

		if school == "" {
			return
		}

		collector.addEducation(
			educationEvidence{
				School: school,
			},
		)

	case []any:

		for _, child := range node {

			collectEducationValue(
				child,
				collector,
			)
		}

	case map[string]any:

		if schemaType(
			node["@type"],
			"Role",
		) {

			schools :=
				organizationNames(
					node["alumniOf"],
				)

			if len(schools) == 0 {

				schools =
					organizationNames(
						node["memberOf"],
					)
			}

			for _, school := range schools {

				collector.addEducation(
					educationEvidence{
						School: school,

						StartDate: firstTextValue(
							node["startDate"],
						),

						EndDate: firstTextValue(
							node["endDate"],
						),
					},
				)
			}

			return
		}

		school :=
			entityName(
				node,
			)

		if school == "" {
			return
		}

		collector.addEducation(
			educationEvidence{
				School: school,
			},
		)
	}
}

func collectAddress(
	value any,
	collector *structuredCollector,
) {
	switch node :=
		value.(type) {

	case []any:

		for _, child := range node {

			collectAddress(
				child,
				collector,
			)
		}

	case map[string]any:

		var parts []string

		add :=
			func(
				value string,
			) {
				value =
					strings.TrimSpace(
						value,
					)

				if value == "" {
					return
				}

				for _, existing := range parts {

					if strings.EqualFold(
						existing,
						value,
					) {

						return
					}
				}

				parts =
					append(
						parts,
						value,
					)
			}

		add(
			firstTextValue(
				node["addressLocality"],
			),
		)

		add(
			firstTextValue(
				node["addressRegion"],
			),
		)

		add(
			countryValue(
				node["addressCountry"],
			),
		)

		if len(parts) > 0 {

			collector.addLocation(
				strings.Join(
					parts,
					", ",
				),
			)
		}
	}
}

func collectPlace(
	value any,
	collector *structuredCollector,
) {
	switch node :=
		value.(type) {

	case string:

		collector.addLocation(
			node,
		)

	case []any:

		for _, child := range node {

			collectPlace(
				child,
				collector,
			)
		}

	case map[string]any:

		name :=
			entityName(
				node,
			)

		if name != "" {

			collector.addLocation(
				name,
			)
		}

		collectAddress(
			node["address"],
			collector,
		)
	}
}

func organizationNames(
	value any,
) []string {
	var result []string

	switch node :=
		value.(type) {

	case string:

		node =
			strings.TrimSpace(
				node,
			)

		if node != "" {

			result =
				append(
					result,
					node,
				)
		}

	case []any:

		for _, child := range node {

			result =
				append(
					result,
					organizationNames(
						child,
					)...,
				)
		}

	case map[string]any:

		if schemaType(
			node["@type"],
			"Role",
		) {

			result =
				append(
					result,
					organizationNames(
						node["worksFor"],
					)...,
				)

			result =
				append(
					result,
					organizationNames(
						node["memberOf"],
					)...,
				)

			result =
				append(
					result,
					organizationNames(
						node["affiliation"],
					)...,
				)

			return uniqueStructuredStrings(
				result,
			)
		}

		name :=
			entityName(
				node,
			)

		if name != "" {

			result =
				append(
					result,
					name,
				)
		}
	}

	return uniqueStructuredStrings(
		result,
	)
}

func entityName(
	node map[string]any,
) string {
	if node == nil {
		return ""
	}

	name :=
		firstTextValue(
			node["name"],
		)

	if name != "" {
		return name
	}

	return firstTextValue(
		node["legalName"],
	)
}

func organizationType(
	value any,
) bool {
	types :=
		[]string{
			"Organization",
			"Corporation",
			"LocalBusiness",
			"Store",
			"ProfessionalService",
			"EducationalOrganization",
			"CollegeOrUniversity",
			"School",
		}

	for _, expected := range types {

		if schemaType(
			value,
			expected,
		) {

			return true
		}
	}

	return false
}

func schemaType(
	value any,
	expected string,
) bool {
	switch current :=
		value.(type) {

	case string:

		return strings.EqualFold(
			strings.TrimSpace(
				current,
			),
			expected,
		)

	case []any:

		for _, child := range current {

			if schemaType(
				child,
				expected,
			) {

				return true
			}
		}
	}

	return false
}

func textValues(
	value any,
) []string {
	var result []string

	switch current :=
		value.(type) {

	case string:

		current =
			strings.TrimSpace(
				current,
			)

		if current != "" {

			result =
				append(
					result,
					current,
				)
		}

	case []any:

		for _, child := range current {

			result =
				append(
					result,
					textValues(
						child,
					)...,
				)
		}

	case map[string]any:

		name :=
			entityName(
				current,
			)

		if name != "" {

			result =
				append(
					result,
					name,
				)
		}
	}

	return uniqueStructuredStrings(
		result,
	)
}

func firstTextValue(
	value any,
) string {
	values :=
		textValues(
			value,
		)

	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func countryValue(
	value any,
) string {
	switch current :=
		value.(type) {

	case string:

		return strings.TrimSpace(
			current,
		)

	case map[string]any:

		return entityName(
			current,
		)
	}

	return ""
}

func standaloneOrganizationPage(
	value string,
) bool {
	parsed, err :=
		url.Parse(
			value,
		)

	if err != nil {
		return false
	}

	host :=
		strings.ToLower(
			strings.TrimSpace(
				parsed.Hostname(),
			),
		)

	host =
		strings.TrimPrefix(
			host,
			"www.",
		)

	host =
		strings.TrimPrefix(
			host,
			"m.",
		)

	blocked :=
		[]string{
			"instagram.com",
			"facebook.com",
			"linkedin.com",
			"twitter.com",
			"x.com",
			"tiktok.com",
			"reddit.com",
			"github.com",
			"gitlab.com",
			"youtube.com",
		}

	for _, domain := range blocked {

		if host == domain ||
			strings.HasSuffix(
				host,
				"."+domain,
			) {

			return false
		}
	}

	return true
}

func (
	collector *structuredCollector,
) addReference(
	value string,
) {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return
	}

	lower :=
		strings.ToLower(
			value,
		)

	if strings.HasPrefix(
		lower,
		"mailto:",
	) {

		value =
			value[len("mailto:"):]

		if index :=
			strings.Index(
				value,
				"?",
			); index >= 0 {

			value =
				value[:index]
		}
	}

	if strings.HasPrefix(
		lower,
		"tel:",
	) {

		value =
			value[len("tel:"):]
	}

	if decoded,
		err :=
		url.QueryUnescape(
			value,
		); err == nil {

		value =
			decoded
	}

	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return
	}

	key :=
		strings.ToLower(
			value,
		)

	if _, exists :=
		collector.references[key]; exists {

		return
	}

	collector.references[key] =
		struct{}{}

	collector.result.References =
		append(
			collector.result.References,
			value,
		)
}

func (
	collector *structuredCollector,
) addURL(
	value string,
) {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return
	}

	parsed, err :=
		url.Parse(
			value,
		)

	if err != nil {
		return
	}

	scheme :=
		strings.ToLower(
			parsed.Scheme,
		)

	if scheme != "http" &&
		scheme != "https" {

		return
	}

	if parsed.Hostname() == "" {
		return
	}

	parsed.Fragment =
		""

	value =
		parsed.String()

	key :=
		strings.ToLower(
			value,
		)

	if _, exists :=
		collector.urls[key]; exists {

		return
	}

	collector.urls[key] =
		struct{}{}

	collector.result.URLs =
		append(
			collector.result.URLs,
			value,
		)
}

func (
	collector *structuredCollector,
) addOrganization(
	value string,
) {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return
	}

	key :=
		strings.ToLower(
			value,
		)

	if _, exists :=
		collector.organizations[key]; exists {

		return
	}

	collector.organizations[key] =
		struct{}{}

	collector.result.Organizations =
		append(
			collector.result.Organizations,
			value,
		)
}

func (
	collector *structuredCollector,
) addLocation(
	value string,
) {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return
	}

	key :=
		strings.ToLower(
			value,
		)

	if _, exists :=
		collector.locations[key]; exists {

		return
	}

	collector.locations[key] =
		struct{}{}

	collector.result.Locations =
		append(
			collector.result.Locations,
			value,
		)
}

func (
	collector *structuredCollector,
) addEmployment(
	value employmentEvidence,
) {
	value.Title =
		strings.TrimSpace(
			value.Title,
		)

	value.Organization =
		strings.TrimSpace(
			value.Organization,
		)

	if value.Title == "" &&
		value.Organization == "" {

		return
	}

	key :=
		strings.ToLower(
			value.Title +
				":" +
				value.Organization +
				":" +
				value.StartDate +
				":" +
				value.EndDate,
		)

	if _, exists :=
		collector.employment[key]; exists {

		return
	}

	collector.employment[key] =
		struct{}{}

	collector.result.Employment =
		append(
			collector.result.Employment,
			value,
		)
}

func (
	collector *structuredCollector,
) addEducation(
	value educationEvidence,
) {
	value.School =
		strings.TrimSpace(
			value.School,
		)

	if value.School == "" {
		return
	}

	key :=
		strings.ToLower(
			value.School +
				":" +
				value.StartDate +
				":" +
				value.EndDate,
		)

	if _, exists :=
		collector.education[key]; exists {

		return
	}

	collector.education[key] =
		struct{}{}

	collector.result.Education =
		append(
			collector.result.Education,
			value,
		)
}

func uniqueStructuredStrings(
	values []string,
) []string {
	seen :=
		make(
			map[string]struct{},
		)

	var result []string

	for _, value := range values {

		value =
			strings.TrimSpace(
				value,
			)

		if value == "" {
			continue
		}

		key :=
			strings.ToLower(
				value,
			)

		if _, exists :=
			seen[key]; exists {

			continue
		}

		seen[key] =
			struct{}{}

		result =
			append(
				result,
				value,
			)
	}

	return result
}
