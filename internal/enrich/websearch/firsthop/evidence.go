package firsthop

import (
	"strings"

	"osint/internal/enrich/model"
	"osint/internal/enrich/provenance"
	"osint/internal/enrich/websearch/discovery"
)

func modelEvidenceFromPage(
	input model.Input,
	page pageEvidence,
) model.Evidence {
	source :=
		"WebSearch FirstHop"

	requestedURL :=
		strings.TrimSpace(
			page.RequestedURL,
		)

	pageURL :=
		strings.TrimSpace(
			page.URL,
		)

	if page.Redirected &&
		requestedURL != "" {

		source +=
			" | Requested: " +
				requestedURL
	}

	if pageURL != "" {

		source +=
			" | Page: " +
				pageURL
	}

	if page.CrossSite {

		source +=
			" | Cross-site redirect"
	}

	subject :=
		strings.TrimSpace(
			page.Subject.Value,
		)

	if subject != "" {

		source +=
			" | Subject: " +
				subject
	}

	value :=
		model.Evidence{
			Anchor: provenance.Anchor(
				input,
			),

			Subjects: []model.EvidenceSubject{
				page.Subject,
			},

			Source: source,

			Text: strings.TrimSpace(
				page.Text,
			),

			URL: pageURL,
		}

	value.ID =
		provenance.ID(
			"page",
			value,
		)

	return value
}

func modelEvidenceFromSearchSubject(
	input model.Input,
	candidate discovery.SearchResult,
	subject model.EvidenceSubject,
) model.Evidence {
	value :=
		discovery.SearchEvidence(
			input,
			candidate,
		)

	value.Subjects =
		[]model.EvidenceSubject{
			subject,
		}

	subjectName :=
		strings.TrimSpace(
			subject.Value,
		)

	if subjectName != "" {

		value.Source +=
			" | Subject: " +
				subjectName
	}

	value.ID =
		provenance.ID(
			"search",
			value,
		)

	return value
}
