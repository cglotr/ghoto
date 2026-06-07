package googleauth

import "net/http"

type Google_auth interface {
	Get_client() (*http.Client, error)
}
