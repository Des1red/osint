package rdap

import (
	"encoding/json"
	"strings"
)

type rdapVCard struct {
	Name         string
	Kind         string
	Organization string
	Title        string
	Role         string

	Emails    []string
	Phones    []string
	Addresses []string
	URLs      []string
}

func parseRDAPEntities(input []rdapEntityResponse) []RDAPEntity {
	entities := make([]RDAPEntity, 0, len(input))

	for _, raw := range input {
		card := parseVCard(raw.VCardArray)

		entity := RDAPEntity{
			Handle: raw.Handle,
			Roles:  raw.Roles,
			Status: raw.Status,
			Port43: raw.Port43,

			Name:         card.Name,
			Kind:         card.Kind,
			Organization: card.Organization,
			Title:        card.Title,
			Role:         card.Role,

			Emails:    card.Emails,
			Phones:    card.Phones,
			Addresses: card.Addresses,
			URLs:      card.URLs,

			PublicIDs: raw.PublicIDs,
			Events:    raw.Events,
			Remarks:   raw.Remarks,
			Links:     raw.Links,

			Entities: parseRDAPEntities(raw.Entities),
		}

		entities = append(entities, entity)
	}

	return entities
}

func parseVCard(raw json.RawMessage) rdapVCard {
	var result rdapVCard

	if len(raw) == 0 {
		return result
	}

	var card []any

	if err := json.Unmarshal(raw, &card); err != nil {
		return result
	}

	if len(card) < 2 {
		return result
	}

	properties, ok := card[1].([]any)
	if !ok {
		return result
	}

	for _, item := range properties {
		property, ok := item.([]any)
		if !ok || len(property) < 4 {
			continue
		}

		name, ok := property[0].(string)
		if !ok {
			continue
		}

		value := property[3]

		switch strings.ToLower(name) {
		case "fn":
			result.Name = joinValues(value, " ")

		case "kind":
			result.Kind = joinValues(value, " ")

		case "org":
			result.Organization = joinValues(value, " ")

		case "title":
			result.Title = joinValues(value, " ")

		case "role":
			result.Role = joinValues(value, " ")

		case "email":
			result.Emails = append(
				result.Emails,
				flattenValues(value)...,
			)

		case "tel":
			result.Phones = append(
				result.Phones,
				flattenValues(value)...,
			)

		case "adr":
			address := joinValues(value, ", ")
			if address != "" {
				result.Addresses = append(
					result.Addresses,
					address,
				)
			}

		case "url":
			result.URLs = append(
				result.URLs,
				flattenValues(value)...,
			)
		}
	}

	return result
}

func flattenValues(value any) []string {
	switch v := value.(type) {
	case string:
		v = strings.TrimSpace(v)

		if v == "" {
			return nil
		}

		return []string{v}

	case []any:
		var values []string

		for _, item := range v {
			values = append(
				values,
				flattenValues(item)...,
			)
		}

		return values
	}

	return nil
}

func joinValues(value any, separator string) string {
	return strings.Join(
		flattenValues(value),
		separator,
	)
}
