package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/oauth2"

	units "github.com/docker/go-units"
	"github.com/genuinetools/pkg/cli"
	"github.com/genuinetools/releases/version"
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/sirupsen/logrus"
)

var (
	port     int
	interval time.Duration

	token  string
	enturl string
	orgs   stringSlice
	nouser bool

	updateReleaseBody bool

	debug bool
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

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "releases"
	p.Description = "Server to show latest GitHub Releases for a set of repositories"

	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.IntVar(&port, "port", 8080, "port for the server to listen on")
	p.FlagSet.IntVar(&port, "p", 8080, "port for the server to listen on")
	p.FlagSet.DurationVar(&interval, "interval", time.Hour, "interval on which to refetch release data")

	p.FlagSet.StringVar(&token, "token", os.Getenv("GITHUB_TOKEN"), "GitHub API token (or env var GITHUB_TOKEN)")
	p.FlagSet.StringVar(&enturl, "url", "", "GitHub Enterprise URL")
	p.FlagSet.Var(&orgs, "orgs", "organizations to include")
	p.FlagSet.BoolVar(&nouser, "nouser", false, "do not include your user")

	p.FlagSet.BoolVar(&updateReleaseBody, "update-release-body", false, "update the body message for the release as well")

	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		if len(token) < 1 {
			return fmt.Errorf("GitHub token cannot be empty")
		}

		if nouser && orgs == nil {
			return fmt.Errorf("no organizations provided")
		}
		return nil
	}

	// Set the main program action.
	p.Action = func(ctx context.Context, args []string) error {
		ticker := time.NewTicker(interval)

		// On ^C, or SIGTERM handle exit.
		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt)
		signal.Notify(signals, syscall.SIGTERM)
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		go func() {
			for sig := range signals {
				cancel()
				ticker.Stop()
				logrus.Infof("Received %s, exiting.", sig.String())
				os.Exit(0)
			}
		}()

		// Create the http client.
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)

		// Create the HTTP cache.
		cachePath := "/tmp/cache"
		if err := os.MkdirAll(cachePath, 0777); err != nil {
			logrus.Fatal(err)
		}
		cache := diskcache.New(cachePath)
		tr := httpcache.NewTransport(cache)
		c := &http.Client{Transport: tr}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, c)

		// Create the github client.
		tc := oauth2.NewClient(ctx, ts)
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
			b   bytes.Buffer
			err error
		)

		// Fetch new data and render the template every interval sequence.
		b, err = run(ctx, client, affiliation)
		if err != nil {
			logrus.Warn(err)
		}
		go func() {
			for range ticker.C {
				bt, err := run(ctx, client, affiliation)
				if err != nil {
					logrus.Warn(err)
				} else {
					b = bt
				}
			}
		}()

		// Setup the server.
		mux := http.NewServeMux()

		// Define wildcard/root handler.
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, b.String())
		})

		logrus.Infof("Starting server on port %d...", port)
		logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))

		return nil
	}

	// Run our program.
	p.Run()
}

type release struct {
	Repository          *github.Repository
	Release             *github.RepositoryRelease
	BinaryName          string
	BinaryURL           string
	BinarySHA256        string
	BinaryMD5           string
	BinaryDownloadCount int
	BinarySince         string
}

func run(ctx context.Context, client *github.Client, affiliation string) (bytes.Buffer, error) {
	var (
		page     = 1
		perPage  = 100
		releases = []release{}
		b        bytes.Buffer
		err      error
	)

	logrus.Info("Getting repositories...")
	releases, err = getRepositories(ctx, client, page, perPage, affiliation, releases)
	if err != nil {
		if v, ok := err.(*github.RateLimitError); ok {
			logrus.Warnf("%s Limit: %d; Remaining: %d; Retry After: %s", v.Message, v.Rate.Limit, v.Rate.Remaining, time.Until(v.Rate.Reset.Time).String())
			return b, nil
		}

		logrus.Warnf("getting repositories failed: %v", err)
	}

	// Parse the template.
	logrus.Info("Executing template...")
	t := template.Must(template.New("").Parse(tmpl))
	w := io.Writer(&b)

	// Execute the template.
	err = t.Execute(w, releases)
	return b, err
}

func getRepositories(ctx context.Context, client *github.Client, page, perPage int, affiliation string, releases []release) ([]release, error) {
	opt := &github.RepositoryListOptions{
		Visibility:  "public",
		Affiliation: affiliation,
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}
	repos, resp, err := client.Repositories.List(ctx, "", opt)
	if err != nil {
		return releases, err
	}

	for _, repo := range repos {
		// Skip it if it's archived.
		if repo.GetArchived() {
			continue
		}

		logrus.Debugf("Handling repo %s...", *repo.FullName)
		r, err := handleRepo(ctx, client, repo)
		if err != nil {
			return releases, err
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
	opt := &github.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	releases, resp, err := client.Repositories.ListReleases(ctx, repo.GetOwner().GetLogin(), repo.GetName(), opt)
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusForbidden || err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, err
		}

		// Skip it because there is no release.
		return nil, nil
	}
	if err != nil || len(releases) < 1 {
		return nil, err
	}

	rl := release{
		Repository: repo,
	}
	// Get information about the binary assets.
	arch := "linux-amd64"
	for i := 0; i < len(releases); i++ {
		r := releases[i]

		isLatest := false
		if rl.Release == nil && !r.GetDraft() {
			// If this is the latest release and it's not a draft make it the one
			// to return
			rl.Release = r
			isLatest = true
		}

		// This holds data like os -> arch -> release and we will use it for rendering our
		// release body template.
		allReleases := map[string]map[string]release{}

		// Iterate over the assets.
		for _, asset := range r.Assets {
			rl.BinaryDownloadCount += asset.GetDownloadCount()

			if !strings.Contains(asset.GetName(), ".") {
				// We know we are on a binary and not a hashsum.
				suffix := strings.SplitN(strings.TrimPrefix(asset.GetName(), repo.GetName()+"-"), "-", 2)
				if len(suffix) == 2 {
					// Add this to our overall releases map.
					osn := suffix[0]
					arch := suffix[1]

					// Prefill the map to avoid a panic.
					if _, ok := allReleases[osn]; !ok {
						allReleases[osn] = map[string]release{}
					}

					tr, ok := allReleases[osn][arch]
					if !ok {
						allReleases[osn][arch] = release{
							BinaryURL:  asset.GetBrowserDownloadURL(),
							BinaryName: asset.GetName(),
							Repository: repo,
						}
					} else {
						tr.BinaryURL = asset.GetBrowserDownloadURL()
						tr.BinaryName = asset.GetName()
						allReleases[osn][arch] = tr
					}
				}
			}

			if strings.HasSuffix(asset.GetName(), ".sha256") {
				// We know we are on a sha256sum.
				suffix := strings.SplitN(strings.TrimSuffix(strings.TrimPrefix(asset.GetName(), repo.GetName()+"-"), ".sha256"), "-", 2)
				if len(suffix) == 2 {
					// Add this to our overall releases map.
					osn := suffix[0]
					arch := suffix[1]

					c, err := getReleaseAssetContent(ctx, client, repo, asset.GetID())
					if err != nil {
						return nil, err
					}

					// Prefill the map to avoid a panic.
					if _, ok := allReleases[osn]; !ok {
						allReleases[osn] = map[string]release{}
					}

					tr, ok := allReleases[osn][arch]
					if !ok {
						allReleases[osn][arch] = release{
							BinarySHA256: c,
							Repository:   repo,
						}
					} else {
						tr.BinarySHA256 = c
						allReleases[osn][arch] = tr
					}
				}
			}

			if isLatest && strings.HasSuffix(asset.GetName(), arch) {
				rl.BinaryURL = asset.GetBrowserDownloadURL()
				rl.BinaryName = asset.GetName()
				rl.BinarySince = units.HumanDuration(time.Since(asset.GetCreatedAt().Time))
			}

			if isLatest && strings.HasSuffix(asset.GetName(), arch+".sha256") {
				c, err := getReleaseAssetContent(ctx, client, repo, asset.GetID())
				if err != nil {
					return nil, err
				}
				rl.BinarySHA256 = c
			}

			if isLatest && strings.HasSuffix(asset.GetName(), arch+".md5") {
				c, err := getReleaseAssetContent(ctx, client, repo, asset.GetID())
				if err != nil {
					return nil, err
				}
				rl.BinaryMD5 = c
			}
		}

		if updateReleaseBody {
			// Do this in a go routine we don't really care if it fails.
			go func(repo *github.Repository, r *github.RepositoryRelease, releases map[string]map[string]release) {
				if err := updateRelease(ctx, client, repo, r, releases); err != nil {
					logrus.Warn(err)
				}
			}(repo, r, allReleases)
		}
	}

	return &rl, nil
}

func updateRelease(ctx context.Context, client *github.Client, repo *github.Repository, r *github.RepositoryRelease, releases map[string]map[string]release) error {
	var (
		b bytes.Buffer
	)

	// Parse the template.
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}
	t := template.Must(template.New("").Funcs(funcMap).Delims("<<", ">>").Parse(releaseTmpl))
	w := io.Writer(&b)

	// Execute the template.
	if err := t.Execute(w, releases); err != nil {
		return err
	}

	s := b.String()
	r.Body = &s
	r.Name = r.TagName

	// Send the new body to GitHub to update the release.
	logrus.Debugf("Updating release for %s -> %s...", repo.GetFullName(), r.GetTagName())
	_, resp, err := client.Repositories.EditRelease(ctx, repo.GetOwner().GetLogin(), repo.GetName(), r.GetID(), r)
	if resp.StatusCode == http.StatusForbidden {
		return nil
	}
	return err
}

func getReleaseAssetContent(ctx context.Context, client *github.Client, repo *github.Repository, id int64) (string, error) {
	body, redirectURL, err := client.Repositories.DownloadReleaseAsset(ctx, repo.GetOwner().GetLogin(), repo.GetName(), id)
	if err != nil {
		return "", err
	}
	if body == nil && len(redirectURL) > 0 {
		resp, err := http.Get(redirectURL)
		if err != nil {
			return "", fmt.Errorf("getting redirect url %s failed: %v", redirectURL, err)
		}
		body = resp.Body
	}
	if body == nil {
		return "", errors.New("body for asset was nil")
	}
	defer body.Close()

	b, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	return strings.Split(string(b), " ")[0], nil
}

func in(a stringSlice, s string) bool {
	for _, b := range a {
		if b == s {
			return true
		}
	}
	return false
}
