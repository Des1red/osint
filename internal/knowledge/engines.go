package knowledge

import (
	"osint/internal/engines/domain"
	"osint/internal/engines/fullname"
	"osint/internal/engines/ip"
	"osint/internal/engines/username"
)

//
// Engine results.
//

type FullNameResult = fullname.FullNameResult

type IPResult = ip.IPResult

type UsernameResult = username.UsernameResult

type DomainResult = domain.DomainResult
