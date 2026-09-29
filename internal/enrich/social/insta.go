package social

import (
	"bytes"
	"encoding/json"
	"html"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var mentionPattern = regexp.MustCompile(
	`@([A-Za-z0-9][A-Za-z0-9._-]{1,63})`,
)

type instagramStructuredData struct {
	Captions []string

	Tagged []string

	Collaborators []string

	Mentioned []string
}

func parseInstagramPost(
	body []byte,
) (
	string,
	string,
	[]postMention,
	bool,
) {
	document,
		err :=
		goquery.NewDocumentFromReader(
			bytes.NewReader(
				body,
			),
		)

	if err != nil {

		return "",
			"",
			nil,
			false
	}

	title :=
		instagramTitle(
			document,
		)

	var textParts []string

	//
	// Public OpenGraph / Twitter metadata.
	//
	textParts =
		appendUniqueValues(
			textParts,
			instagramDescriptions(
				document,
			)...,
		)

	//
	// Explicit Instagram structures contained in
	// JSON/script data.
	//
	structured :=
		instagramStructured(
			document,
		)

	textParts =
		appendUniqueValues(
			textParts,
			structured.Captions...,
		)

	text :=
		strings.TrimSpace(
			strings.Join(
				textParts,
				"\n",
			),
		)

	//
	// First collect ordinary explicit @mentions
	// from caption/description text.
	//
	mentions :=
		textMentions(
			text,
		)

	seen :=
		mentionSet(
			mentions,
		)

	//
	// Then preserve stronger structured
	// associations.
	//
	//
	// We append readable evidence lines to Text so
	// EvidenceReference offsets continue to point
	// at actual evidence text.
	//
	text,
		mentions =
		appendStructuredAccounts(
			text,
			mentions,
			seen,
			"Tagged account",
			"tagged",
			structured.Tagged,
		)

	text,
		mentions =
		appendStructuredAccounts(
			text,
			mentions,
			seen,
			"Collaborator",
			"collaborator",
			structured.Collaborators,
		)

	text,
		mentions =
		appendStructuredAccounts(
			text,
			mentions,
			seen,
			"Mentioned account",
			"mentioned",
			structured.Mentioned,
		)

	if len(mentions) == 0 {

		return title,
			text,
			nil,
			false
	}

	return title,
		text,
		mentions,
		true
}

func instagramTitle(
	document *goquery.Document,
) string {
	title :=
		metaContent(
			document,
			[]string{
				`meta[property="og:title"]`,
				`meta[name="twitter:title"]`,
			},
		)

	if title != "" {

		return title
	}

	if document == nil {

		return ""
	}

	return strings.TrimSpace(
		document.Find(
			"title",
		).First().Text(),
	)
}

func instagramDescriptions(
	document *goquery.Document,
) []string {
	if document == nil {

		return nil
	}

	var result []string

	for _, selector := range []string{
		`meta[property="og:description"]`,
		`meta[name="description"]`,
		`meta[name="twitter:description"]`,
		`meta[property="twitter:description"]`,
	} {

		document.Find(
			selector,
		).Each(
			func(
				_ int,
				selection *goquery.Selection,
			) {
				value,
					exists :=
					selection.Attr(
						"content",
					)

				if !exists {

					return
				}

				result =
					appendUniqueValues(
						result,
						value,
					)
			},
		)
	}

	return result
}

func instagramStructured(
	document *goquery.Document,
) instagramStructuredData {
	var result instagramStructuredData

	if document == nil {

		return result
	}

	document.Find(
		"script",
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

			value,
				ok :=
				decodeInstagramScript(
					raw,
				)

			if !ok {

				return
			}

			collectInstagramStructured(
				value,
				&result,
			)
		},
	)

	result.Captions =
		uniqueStrings(
			result.Captions,
		)

	result.Tagged =
		uniqueUsernames(
			result.Tagged,
		)

	result.Collaborators =
		uniqueUsernames(
			result.Collaborators,
		)

	result.Mentioned =
		uniqueUsernames(
			result.Mentioned,
		)

	return result
}

func decodeInstagramScript(
	raw string,
) (
	any,
	bool,
) {
	raw =
		strings.TrimSpace(
			html.UnescapeString(
				raw,
			),
		)

	if raw == "" {

		return nil,
			false
	}

	//
	// Normal JSON script.
	//
	var value any

	if json.Unmarshal(
		[]byte(
			raw,
		),
		&value,
	) == nil {

		return value,
			true
	}

	//
	// Some embedded data is wrapped in JavaScript:
	//
	//     something(..., {...});
	//
	// Try the outermost object.
	//
	firstObject :=
		strings.Index(
			raw,
			"{",
		)

	lastObject :=
		strings.LastIndex(
			raw,
			"}",
		)

	if firstObject >= 0 &&
		lastObject >
			firstObject {

		candidate :=
			raw[firstObject : lastObject+1]

		if json.Unmarshal(
			[]byte(
				candidate,
			),
			&value,
		) == nil {

			return value,
				true
		}
	}

	//
	// And the outermost array.
	//
	firstArray :=
		strings.Index(
			raw,
			"[",
		)

	lastArray :=
		strings.LastIndex(
			raw,
			"]",
		)

	if firstArray >= 0 &&
		lastArray >
			firstArray {

		candidate :=
			raw[firstArray : lastArray+1]

		if json.Unmarshal(
			[]byte(
				candidate,
			),
			&value,
		) == nil {

			return value,
				true
		}
	}

	return nil,
		false
}

func collectInstagramStructured(
	value any,
	result *instagramStructuredData,
) {
	if result == nil {

		return
	}

	switch current :=
		value.(type) {

	case map[string]any:

		for key, child := range current {

			normalized :=
				strings.ToLower(
					strings.TrimSpace(
						key,
					),
				)

			switch normalized {

			//
			// Old Instagram post caption graph.
			//
			case "edge_media_to_caption":

				result.Captions =
					append(
						result.Captions,
						collectTextFields(
							child,
							"text",
						)...,
					)

			//
			// Some structured documents expose the
			// caption directly.
			//
			case "caption":

				if text,
					ok :=
					child.(string); ok {

					result.Captions =
						append(
							result.Captions,
							text,
						)

				} else {

					result.Captions =
						append(
							result.Captions,
							collectTextFields(
								child,
								"text",
							)...,
						)
				}

			//
			// Explicitly tagged users.
			//
			case "edge_media_to_tagged_user",
				"tagged_users",
				"tagged_user",
				"usertags",
				"user_tags":

				result.Tagged =
					append(
						result.Tagged,
						collectUsernames(
							child,
						)...,
					)

			//
			// Explicit post collaborators.
			//
			case "coauthor_producers",
				"coauthor_producer",
				"collaborators",
				"collaborator":

				result.Collaborators =
					append(
						result.Collaborators,
						collectUsernames(
							child,
						)...,
					)

			//
			// Explicit mentioned-user containers.
			//
			case "mentioned_users",
				"mentioned_user":

				result.Mentioned =
					append(
						result.Mentioned,
						collectUsernames(
							child,
						)...,
					)
			}

			collectInstagramStructured(
				child,
				result,
			)
		}

	case []any:

		for _, child := range current {

			collectInstagramStructured(
				child,
				result,
			)
		}
	}
}

func collectUsernames(
	value any,
) []string {
	var result []string

	var walk func(
		any,
	)

	walk =
		func(
			current any,
		) {
			switch item :=
				current.(type) {

			case map[string]any:

				for key, child := range item {

					if strings.EqualFold(
						strings.TrimSpace(
							key,
						),
						"username",
					) {

						if username,
							ok :=
							child.(string); ok {

							username =
								accountKey(
									username,
								)

							if username != "" {

								result =
									append(
										result,
										username,
									)
							}
						}
					}

					walk(
						child,
					)
				}

			case []any:

				for _, child := range item {

					walk(
						child,
					)
				}
			}
		}

	walk(
		value,
	)

	return uniqueUsernames(
		result,
	)
}

func collectTextFields(
	value any,
	field string,
) []string {
	var result []string

	var walk func(
		any,
	)

	walk =
		func(
			current any,
		) {
			switch item :=
				current.(type) {

			case map[string]any:

				for key, child := range item {

					if strings.EqualFold(
						strings.TrimSpace(
							key,
						),
						field,
					) {

						if text,
							ok :=
							child.(string); ok {

							text =
								strings.TrimSpace(
									text,
								)

							if text != "" {

								result =
									append(
										result,
										text,
									)
							}
						}
					}

					walk(
						child,
					)
				}

			case []any:

				for _, child := range item {

					walk(
						child,
					)
				}
			}
		}

	walk(
		value,
	)

	return uniqueStrings(
		result,
	)
}

func textMentions(
	text string,
) []postMention {
	matches :=
		mentionPattern.
			FindAllStringSubmatchIndex(
				text,
				-1,
			)

	var result []postMention

	seen :=
		make(
			map[string]struct{},
		)

	for _, match := range matches {

		if len(match) <
			4 {

			continue
		}

		start :=
			match[2]

		end :=
			match[3]

		if start < 0 ||
			end <= start ||
			end >
				len(text) {

			continue
		}

		username :=
			accountKey(
				text[start:end],
			)

		if username == "" {

			continue
		}

		if _, exists :=
			seen[username]; exists {

			continue
		}

		seen[username] =
			struct{}{}

		result =
			append(
				result,
				postMention{
					Username: username,

					Start: start,

					End: end,

					Kind: "mention",
				},
			)
	}

	return result
}

func appendStructuredAccounts(
	text string,
	mentions []postMention,
	seen map[string]struct{},
	label string,
	kind string,
	usernames []string,
) (
	string,
	[]postMention,
) {
	for _, username := range usernames {

		username =
			accountKey(
				username,
			)

		if username == "" {

			continue
		}

		//
		// Preserve the structured association in
		// evidence text even when the caption already
		// mentioned the same account.
		//
		if text != "" {

			text +=
				"\n"
		}

		prefix :=
			label +
				": @"

		start :=
			len(text) +
				len(
					prefix,
				)

		text +=
			prefix +
				username

		end :=
			start +
				len(
					username,
				)

		if _, exists :=
			seen[username]; exists {

			continue
		}

		seen[username] =
			struct{}{}

		mentions =
			append(
				mentions,
				postMention{
					Username: username,

					Start: start,

					End: end,

					Kind: kind,
				},
			)
	}

	return text,
		mentions
}

func mentionSet(
	values []postMention,
) map[string]struct{} {
	result :=
		make(
			map[string]struct{},
		)

	for _, value := range values {

		key :=
			accountKey(
				value.Username,
			)

		if key == "" {

			continue
		}

		result[key] =
			struct{}{}
	}

	return result
}

func metaContent(
	document *goquery.Document,
	selectors []string,
) string {
	if document == nil {

		return ""
	}

	for _, selector := range selectors {

		selection :=
			document.Find(
				selector,
			).First()

		if selection.Length() == 0 {

			continue
		}

		value,
			exists :=
			selection.Attr(
				"content",
			)

		if !exists {

			continue
		}

		value =
			strings.TrimSpace(
				value,
			)

		if value != "" {

			return value
		}
	}

	return ""
}

func appendUniqueValues(
	values []string,
	additions ...string,
) []string {
	seen :=
		make(
			map[string]struct{},
		)

	for _, value := range values {

		key :=
			normalizedText(
				value,
			)

		if key != "" {

			seen[key] =
				struct{}{}
		}
	}

	for _, value := range additions {

		value =
			strings.TrimSpace(
				value,
			)

		key :=
			normalizedText(
				value,
			)

		if key == "" {

			continue
		}

		if _, exists :=
			seen[key]; exists {

			continue
		}

		seen[key] =
			struct{}{}

		values =
			append(
				values,
				value,
			)
	}

	return values
}

func uniqueStrings(
	values []string,
) []string {
	return appendUniqueValues(
		nil,
		values...,
	)
}

func uniqueUsernames(
	values []string,
) []string {
	seen :=
		make(
			map[string]struct{},
		)

	var result []string

	for _, value := range values {

		key :=
			accountKey(
				value,
			)

		if key == "" {

			continue
		}

		if _, exists :=
			seen[key]; exists {

			continue
		}

		seen[key] =
			struct{}{}

		result =
			append(
				result,
				key,
			)
	}

	return result
}

func normalizedText(
	value string,
) string {
	return strings.ToLower(
		strings.Join(
			strings.Fields(
				strings.TrimSpace(
					value,
				),
			),
			" ",
		),
	)
}
