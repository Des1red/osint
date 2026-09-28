package variants

import (
	"net/url"
	"strings"
	"sync"
)

const variantWorkers = 3

type LookupFunc func(
	username string,
) PlatformResults

type variantJob struct {
	Index int

	Candidate Candidate
}

type variantJobResult struct {
	Index int

	Result Result

	Found bool
}

type platformSeen struct {
	GitHub map[string]struct{}

	GitLab map[string]struct{}

	Reddit map[string]struct{}

	Facebook map[string]struct{}

	Instagram map[string]struct{}

	Twitter map[string]struct{}

	TikTok map[string]struct{}
}

func Search(
	candidates []Candidate,
	root PlatformResults,
	lookup LookupFunc,
) []Result {
	if len(candidates) == 0 {
		return nil
	}

	jobs :=
		make(
			chan variantJob,
		)

	results :=
		make(
			chan variantJobResult,
			len(candidates),
		)

	workerCount :=
		variantWorkers

	if len(candidates) <
		workerCount {

		workerCount =
			len(candidates)
	}

	var workers sync.WaitGroup

	for i := 0; i < workerCount; i++ {

		workers.Add(
			1,
		)

		go func() {
			defer workers.Done()

			for job := range jobs {

				platformResults :=
					lookup(
						job.Candidate.Username,
					)

				result :=
					Result{
						Username: job.Candidate.Username,

						Match: job.Candidate.Match,

						Source: job.Candidate.Source,

						Platforms: platformResults,
					}

				results <- variantJobResult{
					Index: job.Index,

					Result: result,

					Found: platformResults.FoundAny(),
				}
			}
		}()
	}

	go func() {
		for index, candidate := range candidates {

			jobs <- variantJob{
				Index: index,

				Candidate: candidate,
			}
		}

		close(
			jobs,
		)

		workers.Wait()

		close(
			results,
		)
	}()

	ordered :=
		make(
			[]*Result,
			len(candidates),
		)

	for result := range results {

		if !result.Found {
			continue
		}

		value :=
			result.Result

		ordered[result.Index] =
			&value
	}

	seen :=
		newPlatformSeen()

	seedPlatformSeen(
		seen,
		root,
	)

	var found []Result

	for _, result := range ordered {

		if result == nil {
			continue
		}

		deduplicatePlatforms(
			seen,
			&result.Platforms,
		)

		if !result.Platforms.FoundAny() {
			continue
		}

		found =
			append(
				found,
				*result,
			)
	}

	return found
}

func newPlatformSeen() *platformSeen {
	return &platformSeen{
		GitHub: make(
			map[string]struct{},
		),

		GitLab: make(
			map[string]struct{},
		),

		Reddit: make(
			map[string]struct{},
		),

		Facebook: make(
			map[string]struct{},
		),

		Instagram: make(
			map[string]struct{},
		),

		Twitter: make(
			map[string]struct{},
		),

		TikTok: make(
			map[string]struct{},
		),
	}
}

func seedPlatformSeen(
	seen *platformSeen,
	result PlatformResults,
) {
	if result.GitHub.Found {

		addSeen(
			seen.GitHub,
			result.GitHub.ProfileURL,
			result.GitHub.Source,
		)
	}

	if result.GitLab.Found {

		addSeen(
			seen.GitLab,
			result.GitLab.ProfileURL,
			result.GitLab.Source,
		)
	}

	if result.Reddit.Found {

		addSeen(
			seen.Reddit,
			result.Reddit.ProfileURL,
			result.Reddit.Source,
		)
	}

	if result.Facebook.Found {

		addSeen(
			seen.Facebook,
			result.Facebook.ProfileURL,
			result.Facebook.Source,
		)
	}

	if result.Instagram.Found {

		addSeen(
			seen.Instagram,
			result.Instagram.ProfileURL,
			result.Instagram.Source,
		)
	}

	if result.Twitter.Found {

		addSeen(
			seen.Twitter,
			result.Twitter.ProfileURL,
			result.Twitter.Source,
		)
	}

	if result.TikTok.Found {

		addSeen(
			seen.TikTok,
			result.TikTok.ProfileURL,
			result.TikTok.Source,
		)
	}
}

func deduplicatePlatforms(
	seen *platformSeen,
	result *PlatformResults,
) {
	if result.GitHub.Found {

		if duplicateProfile(
			seen.GitHub,
			result.GitHub.ProfileURL,
			result.GitHub.Source,
		) {

			result.GitHub.Found =
				false
		}
	}

	if result.GitLab.Found {

		if duplicateProfile(
			seen.GitLab,
			result.GitLab.ProfileURL,
			result.GitLab.Source,
		) {

			result.GitLab.Found =
				false
		}
	}

	if result.Reddit.Found {

		if duplicateProfile(
			seen.Reddit,
			result.Reddit.ProfileURL,
			result.Reddit.Source,
		) {

			result.Reddit.Found =
				false
		}
	}

	if result.Facebook.Found {

		if duplicateProfile(
			seen.Facebook,
			result.Facebook.ProfileURL,
			result.Facebook.Source,
		) {

			result.Facebook.Found =
				false
		}
	}

	if result.Instagram.Found {

		if duplicateProfile(
			seen.Instagram,
			result.Instagram.ProfileURL,
			result.Instagram.Source,
		) {

			result.Instagram.Found =
				false
		}
	}

	if result.Twitter.Found {

		if duplicateProfile(
			seen.Twitter,
			result.Twitter.ProfileURL,
			result.Twitter.Source,
		) {

			result.Twitter.Found =
				false
		}
	}

	if result.TikTok.Found {

		if duplicateProfile(
			seen.TikTok,
			result.TikTok.ProfileURL,
			result.TikTok.Source,
		) {

			result.TikTok.Found =
				false
		}
	}
}

func duplicateProfile(
	seen map[string]struct{},
	profileURL string,
	source string,
) bool {
	key :=
		profileKey(
			profileURL,
		)

	if key == "" {

		key =
			profileKey(
				source,
			)
	}

	if key == "" {
		return false
	}

	if _, exists :=
		seen[key]; exists {

		return true
	}

	seen[key] =
		struct{}{}

	return false
}

func addSeen(
	seen map[string]struct{},
	profileURL string,
	source string,
) {
	key :=
		profileKey(
			profileURL,
		)

	if key == "" {

		key =
			profileKey(
				source,
			)
	}

	if key == "" {
		return
	}

	seen[key] =
		struct{}{}
}

func profileKey(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return ""
	}

	parsed, err :=
		url.Parse(
			value,
		)

	if err != nil {

		return strings.ToLower(
			strings.TrimSuffix(
				value,
				"/",
			),
		)
	}

	host :=
		strings.ToLower(
			parsed.Hostname(),
		)

	path :=
		strings.ToLower(
			strings.Trim(
				parsed.Path,
				"/",
			),
		)

	if host == "" {
		return ""
	}

	if path == "" {
		return host
	}

	return host +
		"/" +
		path
}
