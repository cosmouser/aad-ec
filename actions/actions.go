package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/cosmouser/aad-ec/config"
	log "github.com/sirupsen/logrus"
)

// CurrentToken holds the current JWT access token
var CurrentToken *AccessResponse

func init() {
	checkToken()
}

func checkToken() {
	var err error
	if CurrentToken == nil {
		CurrentToken = &AccessResponse{ExpiresOn: "0"}
	}
	expDate, err := strconv.Atoi(CurrentToken.ExpiresOn)
	if err != nil {
		log.Fatal(err)
	}
	if time.Unix(int64(expDate), 0).Before(time.Now()) {
		CurrentToken, err = RequestToken()
		if err != nil {
			log.Error(err)
		}
	}
}

// RequestToken sends a request to Microsoft's Graph API and returns a JWT
func RequestToken() (accResp *AccessResponse, err error) {
	accResp = new(AccessResponse)
	var reqClient = &http.Client{
		Timeout: time.Second * 10,
	}
	formValues := url.Values{}
	formValues.Set("resource", config.C.Resource)
	formValues.Set("client_id", config.C.ClientID)
	formValues.Set("client_secret", config.C.ClientSecret)
	formValues.Set("grant_type", config.C.GrantType)
	br := strings.NewReader(formValues.Encode())
	resURI := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", config.C.Tenant)
	req, err := http.NewRequest("POST", resURI, br)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := reqClient.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(accResp); err != nil {
		log.Error(err)
		return
	}
	return
}

func getPlans(userPrincipalName string) ([]byte, error) {
	var reqClient = &http.Client{
		Timeout: time.Second * 10,
	}
	reqURI := fmt.Sprintf("%s/%s/users/%s/assignedPlans",
		config.C.Resource,
		config.C.APIVersion,
		userPrincipalName,
	)
	req, err := http.NewRequest("GET", reqURI, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+CurrentToken.AccessToken)
	resp, err := reqClient.Do(req)
	if err != nil {
		return []byte(err.Error()), err
	}
	defer resp.Body.Close()

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(err.Error()), err
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode >= 400 || resp.StatusCode <= 509 {
			errv := &ErrorResponse{}
			err = json.Unmarshal(output, errv)
			if err != nil {
				return []byte(err.Error()), err
			}
			return []byte(errv.Error.Message), errors.New(errv.Error.Code)
		}
		return output, errors.New("request was not successful")
	}
	user := &APResponse{}
	err = json.Unmarshal(output, user)

	assignedPlans := []AssignedPlan{}
	for _, v := range user.Value {
		if v.CapabilityStatus == "Enabled" {
			assignedPlans = append(assignedPlans, v)
		}
	}
	apiResponse := &GetPlansResponse{assignedPlans}
	return json.Marshal(apiResponse)
}

// APIHandler is the request handler for calls to the API
func APIHandler(w http.ResponseWriter, r *http.Request) {
	checkToken()
	ra := r.Header.Get("X-Real-IP")
	if ra == "" {
		ra = r.RemoteAddr
	}
	rURI, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		log.Fatal(err)
	}
	query := rURI.Query()
	switch version := query.Get("version"); version {
	case "0.2":
		switch r.Method {
		case "GET":
			if err := checkmail.ValidateFormat(query.Get("uid")); err != nil {
				log.WithFields(log.Fields{
					"remoteAddr": ra,
					"uid":        query.Get("uid"),
				}).Error(err)
				http.Error(w, err.Error(), 503)
				return

			}
			data, err := getPlans(query.Get("uid"))
			if err != nil {
				if err.Error() == "Request_ResourceNotFound" {
					log.WithFields(log.Fields{
						"remoteAddr": ra,
						"uid":        query.Get("uid"),
						"message":    string(data),
					}).Error(err.Error())
					http.Error(w, err.Error(), 404)
					return
				}
				log.WithFields(log.Fields{
					"remoteAddr": ra,
					"uid":        query.Get("uid"),
					"message":    string(data),
				}).Error(err.Error())
				http.Error(w, "request was not successful", 500)
				return

			}
			w.Header().Set("Content-Type", "application/json")
			log.WithFields(log.Fields{
				"remoteAddr": ra,
				"uid":        query.Get("uid"),
			}).Info("Response delivered")
			w.Write(data)
		}

	default:
		w.Write([]byte(`Error: invalid version or no version specified`))
	}
	return
}

// IndexHandler serves the home page with the form for submitting internal requests to aad-ec
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var data = struct {
		ExternalURL string
	}{
		ExternalURL: config.C.ExternalURL,
	}
	if err := indexTemplate.Execute(w, &data); err != nil {
		log.Error(err)
	}
	return
}

var indexTemplate = template.Must(template.New("1").Parse(`<!DOCTYPE html>
<html>
  <head>
	<title>AzureAD Entitlement Checker</title>
	<link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous"> 
    <script src="https://code.jquery.com/jquery-3.3.1.min.js" integrity="sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=" crossorigin="anonymous"></script>  </head>
  <body class="container-fluid">
	<div class="row">
	  <div class="col-sm-4 col-sm-offset-4">
		<h1>AzureAD Entitlement Checker</h1>
		<form id="checker" method="GET" action="/ece/getPlans">
		  <div class="form-group"><span class="label label-default">v0.2</span> <label for="uid">Email of user:</label>
		  <input type="text" class="form-control" id="uid" name="uid" value="cosmo@ucsc.edu"></div>
		  <input type="hidden" id="version" name="version" value="0.2">
		  <button type="submit" class="btn btn-primary">Submit</button>
		</form>
		<br>
		<div class="well" id="results">Results will show here</div>
		<script>
		  $("#checker").submit(function(e) {
			var form = $(this);
		    $.ajax({
			  type: "GET",
			  url: '{{.ExternalURL}}/ece/getPlans?' + form.serialize(),
			  success: function(data) {
				  var out = "";
				  for(var i in data.assignedPlans) {
					  out += data.assignedPlans[i].service + "<br>";
				  }
				$("#results").html(out);
			  },
			  error: function(jqxhr) {
				$("#results").text(jqxhr.responseText);
			  },
			  beforeSend: function(xhr, settings) {
				xhr.setRequestHeader('Accept', 'application/json');
			  }
			});
			e.preventDefault();
		  });
	  </script>
	  </div>
	</div>
  </body>
</html>
`))
