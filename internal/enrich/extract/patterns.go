package extract

import "regexp"

var usernamePattern = regexp.MustCompile(
	`@[A-Za-z0-9._-]{1,64}`,
)

var emailPattern = regexp.MustCompile(
	`(?i)[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}`,
)

var phonePattern = regexp.MustCompile(
	`(?:\+?\d[\d\s()./-]{5,}\d)`,
)

var phoneDigitGroupPattern = regexp.MustCompile(
	`\d+`,
)

var dateLikePhonePattern = regexp.MustCompile(
	`^\d{1,4}\s*[./-]\s*\d{1,2}\s*[./-]\s*\d{1,4}$`,
)

var yearRangePhonePattern = regexp.MustCompile(
	`^(?:19|20)\d{2}\s*[-–—]\s*(?:19|20)\d{2}$`,
)

// Search snippets and public listings frequently
// expose structured information as labelled text.
//
// Include non-location labels too so they act as
// boundaries:
//
// Address: ADAMAS Zipcode: 84801 City: MILOS Phone: ...
//
// Without Phone being recognised as a boundary,
// the City value would incorrectly become:
//
// MILOS Phone: +30 ...
var evidenceFieldPattern = regexp.MustCompile(
	`(?i)\b(location|address|city|country|zip(?:\s*code)?|postcode|postal\s+code|phone|mobile|fax|e-?mail|contact|website|web)\s*:\s*`,
)
