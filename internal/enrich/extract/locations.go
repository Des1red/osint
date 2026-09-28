package extract

import (
	"strings"

	"osint/internal/enrich/model"
)

type labeledField struct {
	Name string

	Value string
}

func extractLocations(
	input model.Input,
) []model.LocationReference {
	var locations []model.LocationReference

	//
	// Explicit structured location evidence.
	//
	for _, value := range input.Locations {

		location :=
			strings.TrimSpace(
				value.Location,
			)

		if location == "" {
			continue
		}

		locations =
			append(
				locations,
				model.LocationReference{
					Location: location,

					Source: value.Source,

					Evidence: wholeEvidenceReference(
						value.EvidenceID,
					),
				},
			)
	}

	//
	// Explicitly labelled location evidence
	// appearing inside free text.
	//
	//
	// Examples:
	//
	// Location: Greece
	//
	// Address: ADAMAS
	// Zipcode: 84801
	// City: MILOS
	//
	// City: ATHENS
	//
	for _, value := range input.Text {

		found :=
			locationsFromText(
				value.Text,
			)

		for _, location := range found {

			locations =
				append(
					locations,
					model.LocationReference{
						Location: location,

						Source: value.Source,

						Evidence: wholeEvidenceReference(
							value.EvidenceID,
						),
					},
				)
		}
	}

	return uniqueLocations(
		locations,
	)
}

func locationsFromText(
	value string,
) []string {
	fields :=
		extractLabeledFields(
			value,
		)

	if len(fields) == 0 {
		return nil
	}

	var locations []string

	for index, field := range fields {

		switch field.Name {

		case "location",
			"city",
			"country":

			if validLocationValue(
				field.Value,
			) {

				locations =
					append(
						locations,
						field.Value,
					)
			}

		case "address":

			address :=
				addressFromFields(
					fields,
					index,
				)

			if validLocationValue(
				address,
			) {

				locations =
					append(
						locations,
						address,
					)
			}
		}
	}

	return uniqueStrings(
		locations,
	)
}

func extractLabeledFields(
	value string,
) []labeledField {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return nil
	}

	matches :=
		evidenceFieldPattern.FindAllStringSubmatchIndex(
			value,
			-1,
		)

	if len(matches) == 0 {
		return nil
	}

	var result []labeledField

	for index, match := range matches {

		if len(match) < 4 {
			continue
		}

		name :=
			normalizeFieldName(
				value[match[2]:match[3]],
			)

		start :=
			match[1]

		end :=
			len(value)

		if index+1 <
			len(matches) {

			end =
				matches[index+1][0]
		}

		if start >= end {
			continue
		}

		fieldValue :=
			cleanLocationFieldValue(
				value[start:end],
			)

		if fieldValue == "" {
			continue
		}

		result =
			append(
				result,
				labeledField{
					Name: name,

					Value: fieldValue,
				},
			)
	}

	return result
}

func addressFromFields(
	fields []labeledField,
	index int,
) string {
	if index < 0 ||
		index >= len(fields) {

		return ""
	}

	address :=
		strings.TrimSpace(
			fields[index].Value,
		)

	if address == "" {
		return ""
	}

	parts :=
		[]string{
			address,
		}

	//
	// A common public listing format is:
	//
	// Address: ADAMAS
	// Zipcode: 84801
	// City: MILOS
	//
	// Preserve those adjacent fields as one
	// address instead of losing the postcode
	// and city context.
	//
	for next := index + 1; next < len(fields) &&
		next <= index+3; next++ {

		field :=
			fields[next]

		switch field.Name {

		case "zipcode":

			if field.Value != "" {

				parts =
					append(
						parts,
						field.Value,
					)
			}

		case "city",
			"country":

			if validLocationValue(
				field.Value,
			) {

				parts =
					append(
						parts,
						field.Value,
					)
			}

		default:

			return strings.Join(
				parts,
				", ",
			)
		}
	}

	return strings.Join(
		parts,
		", ",
	)
}

func cleanLocationFieldValue(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return ""
	}

	for _, separator := range []string{
		" · ",
		" | ",
		"\n",
		"\r",
	} {

		if index :=
			strings.Index(
				value,
				separator,
			); index >= 0 {

			value =
				value[:index]
		}
	}

	value =
		strings.Trim(
			strings.TrimSpace(
				value,
			),
			`"'“”‘’.,;:()[]{}-`,
		)

	return strings.Join(
		strings.Fields(
			value,
		),
		" ",
	)
}

func normalizeFieldName(
	value string,
) string {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	value =
		strings.ReplaceAll(
			value,
			" ",
			"",
		)

	value =
		strings.ReplaceAll(
			value,
			"-",
			"",
		)

	switch value {

	case "zip",
		"zipcode",
		"postcode",
		"postalcode":

		return "zipcode"

	case "location":
		return "location"

	case "address":
		return "address"

	case "city":
		return "city"

	case "country":
		return "country"

	case "phone":
		return "phone"

	case "mobile":
		return "mobile"

	case "fax":
		return "fax"

	case "email":
		return "email"

	case "contact":
		return "contact"

	case "website",
		"web":

		return "website"
	}

	return value
}

func validLocationValue(
	value string,
) bool {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return false
	}

	if len(value) > 120 {
		return false
	}

	if len(
		strings.Fields(
			value,
		),
	) > 12 {

		return false
	}

	if onlyDigits(
		value,
	) {
		return false
	}

	return true
}

func onlyDigits(
	value string,
) bool {
	found :=
		false

	for _, character := range value {

		switch {

		case character >= '0' &&
			character <= '9':

			found =
				true

		case character == ' ',
			character == '-':

		default:
			return false
		}
	}

	return found
}

func uniqueLocations(
	values []model.LocationReference,
) []model.LocationReference {
	indexes :=
		make(
			map[string]int,
		)

	var result []model.LocationReference

	for _, value := range values {

		location :=
			strings.TrimSpace(
				value.Location,
			)

		if location == "" {
			continue
		}

		key :=
			strings.ToLower(
				location +
					":" +
					value.Source,
			)

		if index,
			exists :=
			indexes[key]; exists {

			result[index].Evidence =
				mergeEvidence(
					result[index].Evidence,
					value.Evidence,
				)

			continue
		}

		indexes[key] =
			len(result)

		result =
			append(
				result,
				value,
			)
	}

	return result
}
