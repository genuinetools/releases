package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/oauth2"

	"github.com/genuinetools/releases/version"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

const (
	// BANNER is what is printed for help/info output.
	BANNER = `          _
 _ __ ___| | ___  __ _ ___  ___  ___
| '__/ _ \ |/ _ \/ _` + "`" + ` / __|/ _ \/ __|
| | |  __/ |  __/ (_| \__ \  __/\__ \
|_|  \___|_|\___|\__,_|___/\___||___/

 Server to show latest GitHub Releases for a set of repositories.
 Version: %s
 Build: %s

`
)

var (
	token  string
	enturl string
	orgs   stringSlice
	nouser bool
	dryrun bool

	debug bool
	vrsn  bool
)

// stringSlice is a slice of strings
type stringSlice []string

// implement the flag interface for stringSlice
func (s *stringSlice) String() string {
	return fmt.Sprintf("%s", *s)
}
func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func init() {
	// parse flags
	flag.StringVar(&token, "token", os.Getenv("GITHUB_TOKEN"), "GitHub API token (or env var GITHUB_TOKEN)")
	flag.StringVar(&enturl, "url", "", "GitHub Enterprise URL")
	flag.Var(&orgs, "orgs", "organizations to include")
	flag.BoolVar(&nouser, "nouser", false, "do not include your user")
	flag.BoolVar(&dryrun, "dry-run", false, "do not change branch settings just print the changes that would occur")

	flag.BoolVar(&vrsn, "version", false, "print version and exit")
	flag.BoolVar(&vrsn, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&debug, "d", false, "run in debug mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, version.VERSION, version.GITCOMMIT))
		flag.PrintDefaults()
	}

	flag.Parse()

	if vrsn {
		fmt.Printf("releases version %s, build %s", version.VERSION, version.GITCOMMIT)
		os.Exit(0)
	}

	// set log level
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if token == "" {
		usageAndExit("GitHub token cannot be empty.", 1)
	}

	if nouser && orgs == nil {
		usageAndExit("no organizations provided", 1)
	}
}

func main() {
	// On ^C, or SIGTERM handle exit.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for sig := range c {
			logrus.Infof("Received %s, exiting.", sig.String())
			os.Exit(0)
		}
	}()

	ctx := context.Background()

	// Create the http client.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Create the github client.
	client := github.NewClient(tc)
	if enturl != "" {
		var err error
		client.BaseURL, err = url.Parse(enturl + "/api/v3/")
		if err != nil {
			logrus.Fatal(err)
		}
	}

	// Affiliation must be set before we add the user to the "orgs".
	affiliation := "owner,collaborator"
	if len(orgs) > 0 {
		affiliation += ",organization_member"
	}

	if !nouser {
		// Get the current user
		user, _, err := client.Users.Get(ctx, "")
		if err != nil {
			if v, ok := err.(*github.RateLimitError); ok {
				logrus.Fatalf("%s Limit: %d; Remaining: %d; Retry After: %s", v.Message, v.Rate.Limit, v.Rate.Remaining, time.Until(v.Rate.Reset.Time).String())
			}

			logrus.Fatal(err)
		}
		username := *user.Login
		// add the current user to orgs
		orgs = append(orgs, username)
	}

	var (
		page     = 1
		perPage  = 100
		releases = []release{}
		err      error
	)

	logrus.Debugf("Getting repositories...")
	releases, err = getRepositories(ctx, client, page, perPage, affiliation, releases)
	if err != nil {
		if v, ok := err.(*github.RateLimitError); ok {
			logrus.Fatalf("%s Limit: %d; Remaining: %d; Retry After: %s", v.Message, v.Rate.Limit, v.Rate.Remaining, time.Until(v.Rate.Reset.Time).String())
		}

		logrus.Fatal(err)
	}
}

type release struct {
	Repository *github.Repository
	Release    *github.RepositoryRelease
}

func getRepositories(ctx context.Context, client *github.Client, page, perPage int, affiliation string, releases []release) ([]release, error) {
	opt := &github.RepositoryListOptions{
		Affiliation: affiliation,
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}
	repos, resp, err := client.Repositories.List(ctx, "", opt)
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		logrus.Debugf("Handling repo %s...", *repo.FullName)
		r, err := handleRepo(ctx, client, repo)
		if err != nil {
			return nil, err
		}
		if r != nil {
			releases = append(releases, *r)
		}
	}

	// Return early if we are on the last page.
	if page == resp.LastPage || resp.NextPage == 0 {
		return releases, nil
	}

	page = resp.NextPage
	return getRepositories(ctx, client, page, perPage, affiliation, releases)
}

// handleRepo will return nil error if the user does not have access to something.
func handleRepo(ctx context.Context, client *github.Client, repo *github.Repository) (*release, error) {
	if !in(orgs, repo.GetOwner().GetLogin()) {
		// return early
		return nil, nil
	}

	r, _, err := client.Repositories.GetLatestRelease(ctx, repo.GetOwner().GetLogin(), repo.GetName())
	if err != nil || r == nil {
		return nil, err
	}

	return &release{
		Repository: repo,
		Release:    r,
	}, nil
}

func in(a stringSlice, s string) bool {
	for _, b := range a {
		if b == s {
			return true
		}
	}
	return false
}

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
