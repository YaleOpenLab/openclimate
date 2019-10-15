package database

import (
	"github.com/Varunram/essentials/utils"
	"github.com/pkg/errors"
	"log"
)

func Populate() {
	InitUSStates()
	populateCountries()
	populateRegions()
	populateUSStates()
	populateAvangridCompany()
	populateAvangridAssets()
	populateAdminUsers()
	populateTestUsers()
}

// Test function populating the countries bucket with dummy values
// to test the rpc endpoint for countries
func populateCountries() error {
	log.Println("populating countries", CountryIds)
	for _, country := range CountryIds {
		log.Println("inserting: ", country)
		_, err := NewCountry(country)
		if err != nil {
			return errors.Wrap(err, "Failed to add countries")
		}
	}
	return nil
}

func populateRegions() error {
	_, err := NewRegion("New England", "USA")
	if err != nil {
		return errors.Wrap(err, "Failed populate regions test")
	}
	_, err = NewRegion("Shanghai", "China")
	if err != nil {
		return errors.Wrap(err, "Failed populate regions test")
	}
	_, err = NewRegion("Osaka", "Japan")
	if err != nil {
		return errors.Wrap(err, "Failed populate regions test")
	}
	_, err = NewRegion("Cancun", "Mexico")
	if err != nil {
		return errors.Wrap(err, "Failed populate regions test")
	}
	_, err = NewRegion("Addis Ababa", "Ethiopia")
	if err != nil {
		return errors.Wrap(err, "Failed populate regions test")
	}
	return nil
}

// Test function populating the regions bucket with the US states
func populateUSStates() {
	USStates = []string{"Alabama", "Alaska", "American Samoa", "Arizona", "Arkansas", "California", "Colorado", "Connecticut", "Delaware", "District of Columbia", "Federated States of Micronesia", "Florida", "Georgia", "Guam", "Hawaii", "Idaho", "Illinois", "Indiana", "Iowa", "Kansas", "Kentucky", "Louisiana", "Maine", "Marshall Islands", "Maryland", "Massachusetts", "Michigan", "Minnesota", "Mississippi", "Missouri", "Montana", "Nebraska", "Nevada", "New Hampshire", "New Jersey", "New Mexico", "New York", "North Carolina", "North Dakota", "Northern Mariana Islands", "Ohio", "Oklahoma", "Oregon", "Palau", "Pennsylvania", "Puerto Rico", "Rhode Island", "South Carolina", "South Dakota", "Tennessee", "Texas", "Utah", "Vermont", "Virgin Island", "Virginia", "Washington", "West Virginia", "Wisconsin", "Wyoming"}
	for _, state := range USStates {
		_, err := NewState(state, "USA")
		if err != nil {
			log.Println(errors.Wrap(err, "could not populate US States"))
		}
	}

	// Test Connecticut pledges
	ct, err := RetrieveStateByName("Connecticut", "USA")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = NewPledge("Emissions reduction", 2001, 2050, 80, true, "state", ct.Index)
	if err != nil {
		log.Println(err)
		return
	}
}

func populateAvangridCompany() {

	avangrid, err := NewCompany("Avangrid", "USA")
	if err != nil {
		log.Println(err)
		return
	}

	// Add States
	ct, err := RetrieveStateByName("Connecticut", "USA")
	if err != nil {
		log.Println(err)
		return
	}
	ny, err := RetrieveStateByName("New York", "USA")
	if err != nil {
		log.Println(err)
		return
	}
	ma, err := RetrieveStateByName("Massachusetts", "USA")
	if err != nil {
		log.Println(err)
		return
	}
	err = avangrid.AddStates(ct.Index, ny.Index, ma.Index)
	if err != nil {
		log.Println(err)
		return
	}

	// Add Regions
	ne, err := RetrieveRegionByName("New England", "USA")
	if err != nil {
		log.Println(err)
		return
	}
	err = avangrid.AddRegions(ne.Index)
	if err != nil {
		log.Println(err)
		return
	}

	// Add Countries
	us, err := RetrieveCountryByName("USA")
	if err != nil {
		log.Println(err)
		return
	}
	err = avangrid.AddCountries(us.Index)
	if err != nil {
		log.Println(err)
		return
	}

	// Add Pledges
	pledge, err := NewPledge("reduction", 2015, 2050, 50, true, "company", avangrid.GetID())
	if err != nil {
		log.Println(err)
		return
	}
	err = avangrid.AddPledges(pledge.ID)
	if err != nil {
		log.Println(err)
		return
	}
}

func populateAvangridAssets() {

	avangrid, err := RetrieveCompanyByName("Avangrid", "USA")
	if err != nil {
		log.Println(err)
		return
	}
	bfc, err := NewAsset("Bridgeport 4MW Fuel Cell", avangrid.GetID(), "Bridgeport", "Connecticut", "Gas Fuel Cell")
	if err != nil {
		log.Println(err)
		return
	}

	err = bfc.ReportAssetData(2018, 206, 19164)
	if err != nil {
		log.Println(err)
		return
	}

	updateA := bfc
	updateA.ActionType = []string{"Emissions", "Mitigation"}
	err = UpdateAsset(bfc.Index, updateA)
	if err != nil {
		log.Println(err)
		return
	}

	nhfc, err := NewAsset("New Haven Fuel Cell", avangrid.GetID(), "New Haven", "Connecticut", "Solar Array")
	if err != nil {
		log.Println(err)
		return
	}
	bs, err := NewAsset("Bridgeport Solar 2.2MW", avangrid.GetID(), "Bridgeport", "Connecticut", "Solar Array")
	if err != nil {
		log.Println(err)
		return
	}
	wh, err := NewAsset("Woodbridge High", avangrid.GetID(), "Woodbridge", "Connecticut", "Gas Fuel Cell")
	if err != nil {
		log.Println(err)
		return
	}
	gfc, err := NewAsset("Glastonbury Fuel Cell", avangrid.GetID(), "Glastonbury", "Connecticut", "Gas Fuel Cell")
	if err != nil {
		log.Println(err)
		return
	}

	err = avangrid.AddAssets(bfc.Index, nhfc.Index, bs.Index, wh.Index, gfc.Index)
	if err != nil {
		log.Println(err)
		return
	}
}

func populateAdminUsers() error {
	pwhash := utils.SHA3hash("p")

	_, err := NewUser("amanda", pwhash, "amanda@test.com", "company", "Avangrid", "USA")
	if err != nil {
		log.Println(err, "failed to populate user amanda")
	}

	b, err := NewUser("brian", pwhash, "brian@test.com", "company", "Avangrid", "USA")
	if err != nil {
		return errors.Wrap(err, "failed to populate user brian")
	}
	b.Verified = true
	b.Save()

	// users, err := RetrieveAllUsers()
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(users)

	return nil
}

func populateTestUsers() error {
	pwhash := utils.SHA3hash("a")
	user, err := NewUser("testuser", pwhash, "user@test.com", "country", "USA", "")
	if err != nil {
		return errors.Wrap(err, "failed to create test user in country: USA")
	}
	user.Verified = true
	return user.Save()
}
