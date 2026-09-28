package username

import (
	"fmt"
	"strings"

	"osint/internal/information/output"
	"osint/internal/information/platforms"
	"osint/internal/knowledge"
)

type compactPlatform struct {
	Name string

	Fields []output.TreeField
}

func printVariants(
	results []knowledge.VariantResult,
) {
	if len(results) == 0 {
		return
	}

	fmt.Fprintln(
		output.Writer(),
	)

	fmt.Fprintln(
		output.Writer(),
		"Username Variants",
	)

	fmt.Fprintln(
		output.Writer(),
		"-----------------",
	)

	for index, variant := range results {

		if output.Full() {

			printFullVariant(
				variant,
			)

			continue
		}

		printCompactVariant(
			variant,
			index ==
				len(results)-1,
		)
	}
}

func printCompactVariant(
	variant knowledge.VariantResult,
	last bool,
) {
	label :=
		fmt.Sprintf(
			"[%s] %s",
			matchLabel(
				variant.Match,
			),
			variant.Username,
		)

	output.TreeItem(
		"",
		last,
		label,
	)

	prefix :=
		output.TreeChildPrefix(
			"",
			last,
		)

	platformResults :=
		compactVariantPlatforms(
			variant,
		)

	if len(platformResults) == 0 {

		output.TreeValue(
			prefix,
			true,
			"Source",
			variant.Source,
		)

		return
	}

	output.TreeValue(
		prefix,
		false,
		"Source",
		variant.Source,
	)

	for index, platform := range platformResults {

		platformLast :=
			index ==
				len(platformResults)-1

		output.TreeItem(
			prefix,
			platformLast,
			platform.Name,
		)

		output.TreeFields(
			output.TreeChildPrefix(
				prefix,
				platformLast,
			),
			platform.Fields,
		)
	}
}

func compactVariantPlatforms(
	variant knowledge.VariantResult,
) []compactPlatform {
	var result []compactPlatform

	if variant.Platforms.GitHub.Found {

		result =
			append(
				result,
				compactPlatform{
					Name: "GitHub",

					Fields: []output.TreeField{
						{
							Name: "Username",

							Value: variant.Platforms.GitHub.Username,
						},
						{
							Name: "Name",

							Value: variant.Platforms.GitHub.Name,
						},
						{
							Name: "Profile",

							Value: variant.Platforms.GitHub.ProfileURL,
						},
					},
				},
			)
	}

	if variant.Platforms.GitLab.Found {

		extra :=
			variant.Platforms.GitLab.State

		if variant.Platforms.GitLab.Locked {

			if extra != "" {
				extra +=
					", "
			}

			extra +=
				"locked"
		}

		result =
			append(
				result,
				compactPlatform{
					Name: "GitLab",

					Fields: []output.TreeField{
						{
							Name: "Username",

							Value: variant.Platforms.GitLab.Username,
						},
						{
							Name: "Name",

							Value: variant.Platforms.GitLab.Name,
						},
						{
							Name: "State",

							Value: extra,
						},
						{
							Name: "Profile",

							Value: variant.Platforms.GitLab.ProfileURL,
						},
					},
				},
			)
	}

	if variant.Platforms.Reddit.Found {

		result =
			append(
				result,
				compactPlatform{
					Name: "Reddit",

					Fields: []output.TreeField{
						{
							Name: "Username",

							Value: variant.Platforms.Reddit.Username,
						},
						{
							Name: "Karma",

							Value: fmt.Sprintf(
								"%d",
								variant.Platforms.Reddit.TotalKarma,
							),
						},
						{
							Name: "Profile",

							Value: variant.Platforms.Reddit.ProfileURL,
						},
					},
				},
			)
	}

	if variant.Platforms.Facebook.Found {

		result =
			append(
				result,
				compactPlatform{
					Name: "Facebook",

					Fields: []output.TreeField{
						{
							Name: "Username",

							Value: variant.Platforms.Facebook.Username,
						},
						{
							Name: "Name",

							Value: variant.Platforms.Facebook.Name,
						},
						{
							Name: "Access",

							Value: accessState(
								variant.Platforms.Facebook.Accessible,
								variant.Platforms.Facebook.LoginRequired,
							),
						},
						{
							Name: "Profile",

							Value: variant.Platforms.Facebook.ProfileURL,
						},
					},
				},
			)
	}

	if variant.Platforms.Instagram.Found {

		result =
			append(
				result,
				compactPlatform{
					Name: "Instagram",

					Fields: []output.TreeField{
						{
							Name: "Username",

							Value: variant.Platforms.Instagram.Username,
						},
						{
							Name: "Name",

							Value: cleanInstagramName(
								variant.Platforms.Instagram.Name,
							),
						},
						{
							Name: "Access",

							Value: accessState(
								variant.Platforms.Instagram.Accessible,
								variant.Platforms.Instagram.LoginRequired,
							),
						},
						{
							Name: "Profile",

							Value: variant.Platforms.Instagram.ProfileURL,
						},
					},
				},
			)
	}

	if variant.Platforms.Twitter.Found {

		result =
			append(
				result,
				compactPlatform{
					Name: "Twitter/X",

					Fields: []output.TreeField{
						{
							Name: "Username",

							Value: variant.Platforms.Twitter.Username,
						},
						{
							Name: "Name",

							Value: cleanTwitterName(
								variant.Platforms.Twitter.Name,
							),
						},
						{
							Name: "Access",

							Value: accessState(
								variant.Platforms.Twitter.Accessible,
								variant.Platforms.Twitter.LoginRequired,
							),
						},
						{
							Name: "Profile",

							Value: variant.Platforms.Twitter.ProfileURL,
						},
					},
				},
			)
	}

	if variant.Platforms.TikTok.Found {

		result =
			append(
				result,
				compactPlatform{
					Name: "TikTok",

					Fields: []output.TreeField{
						{
							Name: "Username",

							Value: variant.Platforms.TikTok.Username,
						},
						{
							Name: "Name",

							Value: variant.Platforms.TikTok.Name,
						},
						{
							Name: "Access",

							Value: accessState(
								variant.Platforms.TikTok.Accessible,
								variant.Platforms.TikTok.LoginRequired,
							),
						},
						{
							Name: "Profile",

							Value: variant.Platforms.TikTok.ProfileURL,
						},
					},
				},
			)
	}

	return result
}

func printFullVariant(
	variant knowledge.VariantResult,
) {
	fmt.Fprintln(
		output.Writer(),
	)

	fmt.Fprintf(
		output.Writer(),
		"Candidate: %s\n",
		variant.Username,
	)

	fmt.Fprintf(
		output.Writer(),
		"Username Match: %s\n",
		matchLabel(
			variant.Match,
		),
	)

	fmt.Fprintf(
		output.Writer(),
		"Candidate Source: %s\n",
		variant.Source,
	)

	if variant.Platforms.GitHub.Found {

		platforms.PrintGitHub(
			variant.Platforms.GitHub,
		)
	}

	if variant.Platforms.GitLab.Found {

		platforms.PrintGitLab(
			variant.Platforms.GitLab,
		)
	}

	if variant.Platforms.Reddit.Found {

		platforms.PrintReddit(
			variant.Platforms.Reddit,
		)
	}

	if variant.Platforms.Facebook.Found {

		platforms.PrintFacebook(
			variant.Platforms.Facebook,
		)
	}

	if variant.Platforms.Instagram.Found {

		platforms.PrintInstagram(
			variant.Platforms.Instagram,
		)
	}

	if variant.Platforms.Twitter.Found {

		platforms.PrintTwitter(
			variant.Platforms.Twitter,
		)
	}

	if variant.Platforms.TikTok.Found {

		platforms.PrintTikTok(
			variant.Platforms.TikTok,
		)
	}
}

func accessState(
	accessible bool,
	loginRequired bool,
) string {
	switch {

	case accessible:
		return "accessible"

	case loginRequired:
		return "login required"

	default:
		return ""
	}
}

func cleanInstagramName(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return ""
	}

	//
	// Standard Instagram title:
	//
	// Thanos (@thanos) • Instagram photos and videos
	//
	if index :=
		strings.Index(
			value,
			" (@",
		); index > 0 {

		return strings.TrimSpace(
			value[:index],
		)
	}

	//
	// Empty display-name title:
	//
	// (@thanos1) • Instagram photos and videos
	//
	// This contains no useful display name.
	//
	if strings.HasPrefix(
		value,
		"(@",
	) {
		return ""
	}

	//
	// Another Instagram form can be:
	//
	// itsthanos • Instagram photos and videos
	//
	suffix :=
		" • Instagram photos and videos"

	if strings.HasSuffix(
		value,
		suffix,
	) {

		name :=
			strings.TrimSpace(
				strings.TrimSuffix(
					value,
					suffix,
				),
			)

		if strings.HasPrefix(
			name,
			"@",
		) {
			return ""
		}

		return name
	}

	return value
}

func cleanTwitterName(
	value string,
) string {
	value =
		strings.TrimSpace(
			value,
		)

	if value == "" {
		return ""
	}

	if index :=
		strings.Index(
			value,
			" (@",
		); index > 0 {

		return strings.TrimSpace(
			value[:index],
		)
	}

	return strings.TrimSuffix(
		value,
		" on X",
	)
}
