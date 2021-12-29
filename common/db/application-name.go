package db

import "net/url"

func addApplication(connstr string, application string) (string, error) {
	url, err := url.Parse(connstr)
	if err != nil {
		return "", err
	}
	q := url.Query()
	q.Add("application_name", application)
	url.RawQuery = q.Encode()
	return url.String(), nil
}
