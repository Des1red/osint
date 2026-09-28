package enrichment

import (
	"fmt"
	"strings"

	"osint/internal/enrich/model"
	"osint/internal/information/output"
)

func evidenceSummary(
	values []model.EvidenceReference,
) string {
	if len(values) == 0 {

		return ""
	}

	seen :=
		make(
			map[string]struct{},
		)

	var result []string

	for _, value := range values {

		id :=
			strings.TrimSpace(
				value.EvidenceID,
			)

		if id == "" {

			continue
		}

		text :=
			id

		//
		// 0:0 means the whole evidence record
		// rather than one exact text occurrence.
		//
		if value.Start != 0 ||
			value.End != 0 {

			text =
				fmt.Sprintf(
					"%s [%d:%d]",
					id,
					value.Start,
					value.End,
				)
		}

		if _, exists :=
			seen[text]; exists {

			continue
		}

		seen[text] =
			struct{}{}

		result =
			append(
				result,
				text,
			)
	}

	return strings.Join(
		result,
		", ",
	)
}

func appendProvenanceFields(
	fields []output.TreeField,
	source string,
	evidence []model.EvidenceReference,
) []output.TreeField {
	if !output.Full() {

		return fields
	}

	source =
		strings.TrimSpace(
			source,
		)

	if source != "" {

		fields =
			append(
				fields,
				output.TreeField{
					Name: "Discovered From",

					Value: source,
				},
			)
	}

	evidenceText :=
		evidenceSummary(
			evidence,
		)

	if evidenceText != "" {

		fields =
			append(
				fields,
				output.TreeField{
					Name: "Evidence",

					Value: evidenceText,
				},
			)
	}

	return fields
}
