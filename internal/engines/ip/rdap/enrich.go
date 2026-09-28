package rdap

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

func enrichEntities(
	client *http.Client,
	entities []RDAPEntity,
	source string,
) []RDAPEntity {
	sourceURL, err := url.Parse(source)
	if err != nil {
		return entities
	}

	authority := sourceURL.Host

	cache := make(map[string]RDAPEntity)
	stack := make(map[string]bool)

	for i := range entities {
		entities[i] = enrichEntity(
			client,
			entities[i],
			authority,
			cache,
			stack,
		)
	}

	return entities
}

func enrichEntity(
	client *http.Client,
	entity RDAPEntity,
	authority string,
	cache map[string]RDAPEntity,
	stack map[string]bool,
) RDAPEntity {
	self := entitySelfLink(
		entity.Links,
		authority,
	)

	if self != "" {
		if cached, ok := cache[self]; ok {
			entity = mergeEntity(
				entity,
				cached,
			)
		}

		// Prevent recursive entity loops.
		if stack[self] {
			return entity
		}

		if _, exists := cache[self]; !exists {
			enriched, err := fetchEntity(
				client,
				self,
			)

			if err == nil {
				cache[self] = enriched

				entity = mergeEntity(
					entity,
					enriched,
				)
			}
		}

		stack[self] = true
	}

	for i := range entity.Entities {
		entity.Entities[i] = enrichEntity(
			client,
			entity.Entities[i],
			authority,
			cache,
			stack,
		)
	}

	if self != "" {
		delete(stack, self)

		// Store the fully enriched version.
		cache[self] = entity
	}

	return entity
}

func fetchEntity(
	client *http.Client,
	target string,
) (RDAPEntity, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		target,
		nil,
	)
	if err != nil {
		return RDAPEntity{}, err
	}

	req.Header.Set(
		"Accept",
		"application/rdap+json",
	)

	req.Header.Set(
		"User-Agent",
		"OSINT-Master/1.0",
	)

	resp, err := client.Do(req)
	if err != nil {
		return RDAPEntity{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RDAPEntity{}, &rdapHTTPError{
			Status: resp.Status,
		}
	}

	var data rdapEntityResponse

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&data); err != nil {
		return RDAPEntity{}, err
	}

	entities := parseRDAPEntities(
		[]rdapEntityResponse{data},
	)

	if len(entities) == 0 {
		return RDAPEntity{}, nil
	}

	return entities[0], nil
}

func entitySelfLink(
	links []RDAPLink,
	authority string,
) string {
	for _, link := range links {
		if !strings.EqualFold(
			link.Rel,
			"self",
		) {
			continue
		}

		if !strings.EqualFold(
			link.Type,
			"application/rdap+json",
		) {
			continue
		}

		target, err := url.Parse(link.Href)
		if err != nil {
			continue
		}

		if target.Scheme != "https" {
			continue
		}

		if !strings.EqualFold(
			target.Host,
			authority,
		) {
			continue
		}

		return target.String()
	}

	return ""
}

func mergeEntity(
	current RDAPEntity,
	enriched RDAPEntity,
) RDAPEntity {
	if enriched.Handle != "" {
		current.Handle = enriched.Handle
	}

	if enriched.Name != "" {
		current.Name = enriched.Name
	}

	if enriched.Kind != "" {
		current.Kind = enriched.Kind
	}

	if enriched.Organization != "" {
		current.Organization = enriched.Organization
	}

	if enriched.Title != "" {
		current.Title = enriched.Title
	}

	if enriched.Role != "" {
		current.Role = enriched.Role
	}

	if enriched.Port43 != "" {
		current.Port43 = enriched.Port43
	}

	current.Roles = mergeStrings(
		current.Roles,
		enriched.Roles,
	)

	current.Status = mergeStrings(
		current.Status,
		enriched.Status,
	)

	current.Emails = mergeStrings(
		current.Emails,
		enriched.Emails,
	)

	current.Phones = mergeStrings(
		current.Phones,
		enriched.Phones,
	)

	current.Addresses = mergeStrings(
		current.Addresses,
		enriched.Addresses,
	)

	current.URLs = mergeStrings(
		current.URLs,
		enriched.URLs,
	)

	if len(enriched.PublicIDs) > 0 {
		current.PublicIDs = enriched.PublicIDs
	}

	if len(enriched.Events) > 0 {
		current.Events = enriched.Events
	}

	if len(enriched.Remarks) > 0 {
		current.Remarks = enriched.Remarks
	}

	if len(enriched.Links) > 0 {
		current.Links = enriched.Links
	}

	if len(enriched.Entities) > 0 {
		current.Entities = enriched.Entities
	}

	return current
}

func mergeStrings(
	current []string,
	additional []string,
) []string {
	seen := make(map[string]bool)

	var result []string

	for _, value := range current {
		if value == "" || seen[value] {
			continue
		}

		seen[value] = true
		result = append(
			result,
			value,
		)
	}

	for _, value := range additional {
		if value == "" || seen[value] {
			continue
		}

		seen[value] = true
		result = append(
			result,
			value,
		)
	}

	return result
}

type rdapHTTPError struct {
	Status string
}

func (e *rdapHTTPError) Error() string {
	return "RDAP entity lookup failed: " + e.Status
}
