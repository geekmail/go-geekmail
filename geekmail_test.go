package geekmail

var (
	TestConf = Conf{
		GitHubAuth: GitHubAuth{
			Repository: "github.com/geekmail/geekmail-sample",
			Secret:     "password123",
		},
		// Overwrite .APIToken in login_test.go (it's .gitignored). e.g.:
		//
		// package geekmail
		// func init() { TestConf.APIToken = "ABC" }
		APIToken: "APITOKEN",
	}
)
