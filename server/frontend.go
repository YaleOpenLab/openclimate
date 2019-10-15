package server
// this scrypt relates to the endpoints and routes related to front end. It is a temporary file for front end demo purpose

import (
	// "encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	cc20 "github.com/Varunram/essentials/chacha20poly1305"
	"github.com/Varunram/essentials/ipfs"
	erpc "github.com/Varunram/essentials/rpc"
	"github.com/Varunram/essentials/utils"
	"github.com/YaleOpenLab/openclimate/blockchain"
	"github.com/YaleOpenLab/openclimate/database"
	"github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
)

func frontendFns() {
	getNationStates()
	getMultiNationals()
	getNationStateId()
	getMultiNationalId()
	getActorId()
	getEarthStatus()
	getActors()
	postFiles()
	postRegister()
	postLogin()
	getFiles()
	addLike()
	searchForEntity()
}

func getId(w http.ResponseWriter, r *http.Request) (string, error) {
	var id string
	err := erpc.CheckGet(w, r)
	if err != nil {
		return id, errors.New("request not get")
	}

	urlParams := strings.Split(r.URL.String(), "/")

	if len(urlParams) < 3 {
		return id, errors.New("no id provided, quitting")
	}

	id = urlParams[2]
	return id, nil
}

func getPutId(w http.ResponseWriter, r *http.Request) (string, error) {
	var id string
	err := erpc.CheckPut(w, r)
	if err != nil {
		log.Println(err)
		return id, errors.New("request not get")
	}

	urlParams := strings.Split(r.URL.String(), "/")

	if len(urlParams) < 3 {
		return id, errors.New("no id provided, quitting")
	}

	id = urlParams[2]
	return id, nil
}

func getNationStates() {
	http.HandleFunc("/nation-states", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		nationStates, err := database.RetrieveAllCountries()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		erpc.MarshalSend(w, nationStates)
	})
}

func getMultiNationals() {
	http.HandleFunc("/multinationals", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		multinationals, err := database.RetrieveAllMultiNationals()
		erpc.MarshalSend(w, multinationals)
	})
}

func getNationStateId() {
	http.HandleFunc("/nation-states/", func(w http.ResponseWriter, r *http.Request) {
		strID, err := getId(w, r)
		if err != nil {
			return
		}

		id, err := strconv.Atoi(strID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		nationState, err := database.RetrieveCountry(id)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		pledges, err := nationState.GetPledges()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		results := make(map[string]interface{})
		results["name"] = nationState.Name
		results["full_name"] = nationState.Name
		results["description"] = nationState.Description
		results["pledges"] = pledges
		results["accountability"] = nationState.Accountability

		erpc.MarshalSend(w, results)
	})
}

func getMultiNationalId() {
	http.HandleFunc("/multinationals/", func(w http.ResponseWriter, r *http.Request) {
		strID, err := getId(w, r)
		if err != nil {
			return
		}

		id, err := strconv.Atoi(strID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		multinational, err := database.RetrieveCompany(id)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		pledges, err := multinational.GetPledges()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		results := make(map[string]interface{})
		results["name"] = multinational.Name
		results["full_name"] = multinational.Name
		results["description"] = multinational.Description
		results["pledges"] = pledges
		results["accountability"] = multinational.Accountability
		results["locations"] = multinational.Locations

		erpc.MarshalSend(w, results)
	})
}

type NationState struct {
	Name        string
	Pledges     []database.Pledge
	Subnational []Subnational
}

type Subnational struct {
	Name    string
	Pledges []database.Pledge
	Assets  []database.Asset
}

func getActorId() {
	http.HandleFunc("/actors/", func(w http.ResponseWriter, r *http.Request) {
		strID, err := getId(w, r)
		if err != nil {
			return
		}

		id, err := utils.ToInt(strID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		company, err := database.RetrieveCompany(id)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		urlParams := strings.Split(r.URL.String(), "/")
		if len(urlParams) < 4 {
			log.Println("insufficient amount of params")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		choice := urlParams[3]

		switch choice {
		case "dashboard":
			pledges, err := company.GetPledges()
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			results := make(map[string]interface{})
			results["full_name"] = company.Name
			results["description"] = company.Description
			results["locations"] = company.Locations
			results["accountability"] = company.Accountability
			results["pledges"] = pledges

			results["direct_emissions"], err = getDirectEmissionsActorId(strID)
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}

			results["mitigation_outcomes"], err = getMitigationOutcomesActorId(strID)
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}

			results["direct_emissions"], err = getWindAndSolarActorId(strID)
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}

			results["disclosure_settings"], err = getDisclosureSettingsActorId(strID)
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}

			results["weighted_score"], err = getWeightedScoreActorId(strID)
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
		// end of dashboard case
		case "nation-states":
			nationStates, err := getActorIdNationStates(company, w, r)
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			erpc.MarshalSend(w, nationStates)
		// end of nation states case
		case "review":
			results := make(map[string]interface{})
			results["certificates"] = company.Certificates
			results["climate_reports"] = company.ClimateReports

			var err error
			results["emissions"], err = blockchain.RetrieveActorEmissions(id)
			if err != nil {
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}
			results["reductions"], err = blockchain.RetrieveActorEmissions(id)
			if err != nil {
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}
			erpc.MarshalSend(w, results)

		// case "manage":
		// 	w.Write([]byte("manage: " + strconv.Itoa(id)))

		case "climate-action-asset":
			if len(urlParams) < 5 {
				log.Println("insufficient amount of params")
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}
			id2 := urlParams[4]
			w.Write([]byte("climate-action-assets ids: " + strconv.Itoa(id) + id2))
		default:
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
	})
}

func getActorIdNationStates(company database.Company, w http.ResponseWriter, r *http.Request) ([]NationState, error) {

	var nationStates []NationState

	countries, err := company.GetCountries()
	if err != nil {
		return nationStates, errors.Wrap(err, "getActorIdNationStates() failed")
	}

	for _, country := range countries {
		var nationState NationState
		states, err := company.GetStates()
		if err != nil {
			return nationStates, errors.Wrap(err, "getActorIdNationStates() failed")
		}

		pledges, err := country.GetPledges()
		if err != nil {
			return nationStates, errors.Wrap(err, "getActorIdNationStates() failed")
		}

		var subnationals []Subnational

		for _, s := range states {
			var subnational Subnational
			pledges, err := s.GetPledges()
			if err != nil {
				return nationStates, errors.Wrap(err, "getActorIdNationStates() failed")
			}
			assets, err := company.GetAssetsByState(s.Name)
			if err != nil {
				return nationStates, errors.Wrap(err, "getActorIdNationStates() failed")
			}

			subnational.Name = s.Name
			subnational.Pledges = pledges
			subnational.Assets = assets
			subnationals = append(subnationals, subnational)
		}

		nationState.Name = country.Name
		nationState.Pledges = pledges
		nationState.Subnational = subnationals
		nationStates = append(nationStates, nationState)
	}

	return nationStates, nil
}

type EarthStatusReturn struct {
	Warminginc               string `json:"warming_in_c"`
	Gtco2left                string `json:"gt_co2_left"`
	Atmosphericco2ppm        string `json:"atmospheric_co2_ppm"`
	Annualglobalemission     string `json:"annual_global_emission"`
	Estimatedbudgetdepletion string `json:"estimated_budget_depletion"`
}

func getEarthStatus() {
	http.HandleFunc("/earth-status", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		var x EarthStatusReturn
		x.Warminginc = "sample"
		x.Gtco2left = "sample"
		x.Atmosphericco2ppm = "sample"
		x.Annualglobalemission = "sample"
		x.Estimatedbudgetdepletion = "sample"

		erpc.MarshalSend(w, x)
	})
}

func getActors() {
	http.HandleFunc("/actors", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		w.Write([]byte("get actors"))
	})
}

func postRegister() {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = r.ParseForm()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		if !checkReqdParams(w, r, "username", "pwhash", "first_name", "last_name", "email", "ein") {
			return
		}

		multiNationalId := r.FormValue("m_id")
		nationId := r.FormValue("n_id")
		var idS string

		var multinational bool
		if multiNationalId == "" {
			if nationId == "" {
				log.Println("both mnc and nation id are nil, quitting")
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}
			multinational = false
			idS = nationId
		} else {
			multinational = true
			idS = multiNationalId
		}

		id, err := utils.ToInt(idS)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		x, _ := database.RetrieveUser(id)

		x.Username = r.FormValue("username")
		x.Pwhash = r.FormValue("pwhash")
		x.FirstName = r.FormValue("first_name")
		x.LastName = r.FormValue("last_name")
		x.Email = r.FormValue("email")
		x.EIN = r.FormValue("ein")
		x.EntityType = "country"

		if multinational {
			x.EntityType = "mnc"
		} else {
			x.EntityType = "country"
		}

		err = x.Save()
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

type PostLoginResponse struct {
	Token string
}

func postLogin() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = r.ParseForm()
		if err != nil {
			log.Println("PARSE FORM ERROR?", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		if !checkReqdPostParams(w, r, "username", "pwhash") {
			return
		}

		username := r.FormValue("username")
		pwhash := r.FormValue("pwhash")

		user, err := database.ValidateUser(username, pwhash)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		token, err := user.GenAccessToken()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var x PostLoginResponse
		x.Token = token
		erpc.MarshalSend(w, x)
	})
}

type postfileReturn struct {
	IpfsHash string
}

func postFiles() {
	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		log.Println("in postfiles endpoint")
		err := erpc.CheckPost(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if !checkReqdPostParams(w, r, "id", "docId", "entity") {
			return
		}

		id := r.FormValue("id")
		docIdString := r.FormValue("docId")
		entity := r.FormValue("entity")

		if entity != "country" && entity != "mnc" && entity != "state" {
			erpc.MarshalSend(w, erpc.StatusBadRequest)
			return
		}

		docId, err := utils.ToInt(docIdString)
		if err != nil {
			erpc.MarshalSend(w, erpc.StatusInternalServerError)
			return
		}

		if docId > 5 || docId < 1 {
			log.Println("invalid doc id, quitting")
			erpc.MarshalSend(w, erpc.StatusBadRequest)
			return
		}

		r.ParseMultipartForm(1 << 21) // max 10MB files
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			log.Println("could not parse form data", err)
			erpc.MarshalSend(w, erpc.StatusBadRequest)
			return
		}

		defer file.Close()

		log.Println("file size: ", fileHeader.Size)
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			erpc.MarshalSend(w, erpc.StatusInternalServerError)
			return
		}

		// encrypt with chacha20 since the file size is variable
		encryptedBytes, err := cc20.Encrypt(fileBytes, globals.IpfsMasterPwd)
		if err != nil {
			erpc.MarshalSend(w, erpc.StatusInternalServerError)
			return
		}

		hash, err := ipfs.IpfsAddBytes(encryptedBytes)
		if err != nil {
			erpc.MarshalSend(w, erpc.StatusInternalServerError)
			return
		}

		switch entity {
		case "country":
			log.Println("storing file against required country")
			idInt, err := utils.ToInt(id)
			if err != nil {
				erpc.MarshalSend(w, erpc.StatusBadRequest)
				return
			}
			x, err := database.RetrieveCountry(idInt)
			if err != nil {
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}
			x.Files = append(x.Files, hash)
			err = x.Save()
			if err != nil {
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
			}
		case "mnc":
			log.Println("storing file against required multinational company")
			idInt, err := utils.ToInt(id)
			if err != nil {
				erpc.MarshalSend(w, erpc.StatusBadRequest)
				return
			}
			x, err := database.RetrieveCompany(idInt)
			if err != nil {
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}
			x.Files = append(x.Files, hash)
			err = x.Save()
			if err != nil {
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
			}
		case "state":
			log.Println("storing file against requried state")
			idInt, err := utils.ToInt(id)
			if err != nil {
				erpc.MarshalSend(w, erpc.StatusBadRequest)
				return
			}
			x, err := database.RetrieveState(idInt)
			if err != nil {
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}
			x.Files = append(x.Files, hash)
			err = x.Save()
			if err != nil {
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
			}
		}

		var pf postfileReturn
		pf.IpfsHash = hash
		erpc.MarshalSend(w, pf)
	})
}

func getFiles() {
	http.HandleFunc("/getfiles", func(w http.ResponseWriter, r *http.Request) {
		log.Println("in getFiles endpoint")

		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if !checkReqdParams(w, r, "hash", "extension") {
			return
		}

		extension := r.URL.Query()["extension"][0]
		hash := r.URL.Query()["hash"][0]

		encryptedFile, err := ipfs.IpfsGetFile(hash, extension)
		if err != nil {
			erpc.MarshalSend(w, erpc.StatusInternalServerError)
			return
		}

		encryptedBytes, err := ioutil.ReadFile(encryptedFile)
		if err != nil {
			erpc.MarshalSend(w, erpc.StatusInternalServerError)
			return
		}

		os.Remove(encryptedFile)

		decryptedBytes, err := cc20.Decrypt(encryptedBytes, globals.IpfsMasterPwd)
		if err != nil {
			erpc.MarshalSend(w, erpc.StatusInternalServerError)
			return
		}

		w.Write(decryptedBytes)
	})
}

func addLike() {
	http.HandleFunc("/like/pledges/", func(w http.ResponseWriter, r *http.Request) {
		strID, err := getPutId(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if !checkReqdPostParams(w, r, "accessToken", "username") {
			return
		}

		accessToken := r.FormValue("accessToken")
		username := r.FormValue("username")

		user, err := database.RetrieveUserByUsername(username)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		if user.AccessToken != accessToken {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		// since the frontend is not expected to pass invalid requests to the liked routes, we don't validate that.
		// this is not expected to be used by any ohter extenral parties, so this is okay I guess.
		user.Liked = append(user.Liked, strID)
		err = user.Save()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, erpc.StatusOK)
	})
}

func addNotVisible() {
	http.HandleFunc("/hide/disclosure-settings/", func(w http.ResponseWriter, r *http.Request) {
		strID, err := getPutId(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if !checkReqdPostParams(w, r, "accessToken", "username") {
			return
		}

		accessToken := r.FormValue("accessToken")
		username := r.FormValue("username")

		user, err := database.RetrieveUserByUsername(username)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		if user.AccessToken != accessToken {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		// since the frontend is not expected to pass invalid requests to the liked routes, we don't validate that.
		// this is not expected to be used by any ohter extenral parties, so this is okay I guess.
		user.NotVisible = append(user.NotVisible, strID)
		err = user.Save()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, erpc.StatusOK)
	})
}

type searchReturn struct {
	Username  string
	Country   string
	Companies []string
	Cities    []string
	Regions   []string
	States    []string
}

func searchForEntity() {
	http.HandleFunc("/actors/search", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if !checkReqdParams(w, r, "search") {
			return
		}

		searchParam := r.URL.Query()["search"][0]

		var country database.Country
		var user database.User
		var companies []database.Company
		var cities []database.City
		var regions []database.Region
		var states []database.State

		user, _ = database.RetrieveUserByUsername(searchParam)

		country, _ = database.RetrieveCountryByName(searchParam)
		companies, _ = database.SearchCompany(searchParam)
		states, _ = database.SearchState(searchParam)
		regions, _ = database.SearchRegion(searchParam)
		cities, _ = database.SearchCity(searchParam)

		var stateNames []string
		var regionNames []string
		var cityNames []string
		var companyNames []string

		for _, state := range states {
			stateNames = append(stateNames, state.Name)
		}

		for _, region := range regions {
			regionNames = append(regionNames, region.Name)
		}

		for _, city := range cities {
			cityNames = append(cityNames, city.Name)
		}

		for _, company := range companies {
			companyNames = append(companyNames, company.Name)
		}
		var x searchReturn

		x.Username = user.Username
		x.Country = country.Name
		x.Companies = companyNames
		x.Cities = cityNames
		x.Regions = regionNames
		x.States = stateNames

		erpc.MarshalSend(w, x)
	})
}
