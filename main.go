package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	secret := os.Getenv("CHANNEL_SECRET")
	if secret == "" {
		secret = "[待填]"
	}

	token := os.Getenv("CHANNEL_TOKEN")
	if token == "" {
		token = "[待填]"
	}
	bot, err := linebot.New(
		secret,
		token,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					switch {
					case strings.HasPrefix(message.Text, "布版,"):
						content := strings.Split(message.Text, ",")
						if len(content) == 5 {
							err := DeployTOOctopus(content[1], content[2], content[3], content[4])
							if err != nil {
								if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do(); err != nil {
									log.Print(err)
									return
								}
								return
							}
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("deploy complete!")).Do(); err != nil {
								log.Print(err)
								return
							}
						}
					case strings.HasPrefix(message.Text, "布版求救"):
						str := `
						範例: 

						布版,ag-slots,stage,ht-aladdin,1.210913.3
						
						Space:

						ag-ocean
						ag-slots

						Env: 

						dev
						stage
						sandbox
						uat
						pgs-dev
						prod-aceclub
						prod-08
						pgs-prod
						`
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(str)).Do(); err != nil {
							log.Print(err)
						}
					}

					// case *linebot.StickerMessage:
					// 	replyMessage := fmt.Sprintf(
					// 		"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					// 		log.Print(err)
					// 	}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":8066", nil); err != nil {
		log.Fatal(err)
	}
}

func DeployTOOctopus(spaceName, environmentName, projectName, releaseVersion string) (err error) {

	apiURL, err := url.Parse("[待填]")
	if err != nil {
		err = fmt.Errorf("apiURL err: %v ", err)
		return
	}

	apikey := "[待填]"

	fmt.Printf("  %s,%s\n", projectName, releaseVersion)

	// Get reference to space
	space, err := GetSpace(apiURL, apikey, spaceName)
	if err != nil {
		err = fmt.Errorf("get space err: %v ", err)
		return
	}

	// Create client object
	var client *octopusdeploy.Client
	client, err = OctopusAuth(apiURL, apikey, space.ID)
	if err != nil {
		err = fmt.Errorf("get auth err: %v ", err)
		return
	}

	// Get project
	var project *octopusdeploy.Project
	project, err = GetProject(client, apiURL, apikey, space, projectName)
	if err != nil {
		err = fmt.Errorf("get project err: %v ", err)
		return
	}

	// Get environment
	var environment *octopusdeploy.Environment
	environment, err = GetEnvironment(client, apiURL, apikey, space, environmentName)
	if err != nil {
		err = fmt.Errorf("get env err: %v ", err)
		return
	}

	// Get project releases
	var projectReleases []interface{}
	projectReleases, err = GetProjectReleases(apiURL, apikey, space, project)
	if err != nil {
		err = fmt.Errorf("get release err: %v ", err)
		return
	}

	// Loop through releases
	for i := 0; i < len(projectReleases); i++ {
		projectRelease := projectReleases[i].(map[string]interface{})

		// Delete release
		if projectRelease["Version"].(string) == releaseVersion {
			// Create deployment object
			deployment := octopusdeploy.NewDeployment(environment.ID, projectRelease["Id"].(string))

			// Issue deployment
			_, err = client.Deployments.Add(deployment)
			if err != nil {
				err = fmt.Errorf("deploy err: %v ", err)
				return
			}
		}
	}
	return
}

func OctopusAuth(octopusURL *url.URL, apikey, space string) (client *octopusdeploy.Client, err error) {
	client, err = octopusdeploy.NewClient(nil, octopusURL, apikey, space)
	if err != nil {
		return
	}

	return
}

func GetSpace(octopusURL *url.URL, apikey string, spaceName string) (space *octopusdeploy.Space, err error) {
	var client *octopusdeploy.Client
	client, err = OctopusAuth(octopusURL, apikey, "")
	if err != nil {
		return
	}
	spaceQuery := octopusdeploy.SpacesQuery{
		Name: spaceName,
	}

	// Get specific space object
	spaces, err := client.Spaces.Get(spaceQuery)
	if err != nil {
		return
	}

	for _, tmpSpace := range spaces.Items {
		if tmpSpace.Name == spaceName {
			space = tmpSpace
			return
		}
	}

	err = fmt.Errorf("not found space")
	return
}

func GetProjectReleases(
	octopusURL *url.URL,
	apikey string,
	space *octopusdeploy.Space,
	project *octopusdeploy.Project,
) (
	result []interface{},
	err error,
) {
	// Define api endpoint
	projectReleasesEndoint := octopusURL.String() + "/api/" + space.ID + "/projects/" + project.ID + "/releases"

	// Create http client
	httpClient := &http.Client{}
	skipAmount := 0

	// Make request
	var request *http.Request
	request, err = http.NewRequest("GET", projectReleasesEndoint, nil)
	if err != nil {
		return
	}
	request.Header.Set("X-Octopus-apikey", apikey)
	response, err := httpClient.Do(request)
	if err != nil {
		return
	}

	// Get response
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	var releasesJson interface{}
	err = json.Unmarshal(responseData, &releasesJson)
	if err != nil {
		return
	}

	// Map the returned data
	returnedReleases := releasesJson.(map[string]interface{})
	// Returns the list of items, translate it to a map
	returnedItems := returnedReleases["Items"].([]interface{})
	//make(map[string][]octopusdeploy.PropertyValue)

	for {
		// check to see if there's more to get
		fltItemsPerPage := returnedReleases["ItemsPerPage"].(float64)
		itemsPerPage := int(fltItemsPerPage)
		if len(returnedReleases["Items"].([]interface{})) == itemsPerPage {
			// Increment skip accoumt
			skipAmount += len(returnedReleases["Items"].([]interface{}))

			// Make request
			queryString := request.URL.Query()
			queryString.Set("skip", strconv.Itoa(skipAmount))
			request.URL.RawQuery = queryString.Encode()
			response, err := httpClient.Do(request)
			if err != nil {
				break
			}

			responseData, err = ioutil.ReadAll(response.Body)
			if err != nil {
				break
			}
			var releasesJson interface{}
			err = json.Unmarshal(responseData, &releasesJson)
			if err != nil {
				break
			}

			returnedReleases = releasesJson.(map[string]interface{})
			returnedItems = append(returnedItems, returnedReleases["Items"].([]interface{})...)
		} else {
			// err = fmt.Errorf("items per page not found")
			break
		}
	}

	result = returnedItems
	return
}

func GetProject(
	client *octopusdeploy.Client,
	octopusURL *url.URL,
	apikey string,
	space *octopusdeploy.Space,
	projectName string,
) (
	project *octopusdeploy.Project,
	err error,
) {
	projectsQuery := octopusdeploy.ProjectsQuery{
		Name: projectName,
	}

	// Get specific project object
	projects, err := client.Projects.Get(projectsQuery)

	if err != nil {
		return
	}

	for i := range projects.Items {
		if projects.Items[i].Name == projectName {
			project = projects.Items[i]
			return
		}
	}

	err = fmt.Errorf("project not found")
	return
}

func GetEnvironment(
	client *octopusdeploy.Client,
	octopusURL *url.URL,
	apikey string,
	space *octopusdeploy.Space,
	environmentName string,
) (env *octopusdeploy.Environment, err error) {

	// Get environment
	environmentsQuery := octopusdeploy.EnvironmentsQuery{
		Name: environmentName,
	}
	environments, err := client.Environments.Get(environmentsQuery)
	if err != nil {
		return
	}

	// Loop through results
	for i := range environments.Items {
		if environments.Items[i].Name == environmentName {
			env = environments.Items[i]
			return
		}
	}

	err = fmt.Errorf("env not found")
	return
}
