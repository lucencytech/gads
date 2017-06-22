package v201705

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func getTestConfig() AuthConfig {

	creds := Credentials{
		Config: OAuthConfigArgs{
			ClientID:     os.Getenv("ADWORDS_CLIENT_ID"),
			ClientSecret: os.Getenv("ADWORDS_CLIENT_SECRET"),
		},
		Token: OAuthTokenArgs{
			AccessToken:  os.Getenv("ADWORDS_ACCESS_TOKEN"),
			RefreshToken: os.Getenv("ADWORDS_REFRESH_TOKEN"),
		},
		Auth: Auth{
			CustomerId:     os.Getenv("ADWORDS_TEST_ACCOUNT"),
			DeveloperToken: os.Getenv("ADWORDS_DEVELOPER_TOKEN"),
			PartialFailure: true,
		},
	}

	authconf, _ := NewCredentialsFromParams(creds)
	return authconf
}

func TestSandboxCriteria(t *testing.T) {
	config := getTestConfig()

	campaigns, _, err := NewCampaignService(&config.Auth).Get(Selector{
		Fields: []string{"Id", "Name", "CampaignId"},
	})

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(campaigns)
	campaign := campaigns[0].Id

	adgroups, _, err := NewAdGroupService(&config.Auth).Get(Selector{
		Fields: []string{"Id", "Name"},
		Predicates: []Predicate{
			Predicate{
				Field:    "CampaignId",
				Operator: "EQUALS",
				Values:   []string{strconv.FormatInt(campaign, 10)},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	adgroup, err := func() (*AdGroup, error) {
		for _, a := range adgroups {
			if a.Name == "sidecar-test-adgroup" {
				return &a, nil
			}
		}
		return nil, fmt.Errorf("missing test adgroup\n")
	}()
	if err != nil {
		t.Fatal(err)
	}
	/*
		query := fmt.Sprintf("SELECT * WHERE AdGroupId = %d", adgroup.Id)

		crits, _, err := NewAdGroupCriterionService(&config.Auth).Query(query)
	*/
	crits, _, err := NewAdGroupCriterionService(&config.Auth).Get(Selector{
		Fields: []string{"AdGroupId", "BidModifier", "CriterionUse", "ParentCriterionId", "CriteriaType", "CaseValue", "Id", "BiddingStrategyType", "CpcBid", "BiddingStrategyId"},
		Predicates: []Predicate{
			Predicate{
				Field:    "AdGroupId",
				Operator: "EQUALS",
				Values:   []string{strconv.FormatInt(adgroup.Id, 10)},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	//rootCriterion

	root, rest, toremove := func() (ProductPartition, []BiddableAdGroupCriterion, *BiddableAdGroupCriterion) {
		var root ProductPartition
		var rest []BiddableAdGroupCriterion
		var toremove *BiddableAdGroupCriterion

		for i := 0; i < len(crits); i++ {
			crit, _ := crits[i].(BiddableAdGroupCriterion)
			part := crit.Criterion.(ProductPartition)

			if part.ParentCriterionId == 0 {
				root = part
			} else if part.Dimension.Value == "agi" {
				//	part.Dimension.TypeAttr = "ProductBrand"
				crit.Criterion = part
				toremove = &crit
			} else {
				crit.Criterion = part
				fmt.Printf("CRIT:  %#v\n%#v\n", crit, *crit.BiddingStrategyConfiguration)
				rest = append(rest, crit)
			}
		}
		return root, rest, toremove
	}()

	fmt.Printf("ROOT:  %#v\n", root)

	/*
		removes := AdGroupCriterions{}
		for _, x := range rest {
			removes = append(removes, BiddableAdGroupCriterion{
				AdGroupId: x.AdGroupId,
				Criterion: ProductPartition{
					Id: x.Criterion.(ProductPartition).Id,
				},
			})
		}
	*/

	if toremove != nil {
		aops := AdGroupCriterionOperations{
			"REMOVE": AdGroupCriterions{*toremove},
		}

		res, err := NewAdGroupCriterionService(&config.Auth).Mutate(aops)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(res)
	}

	toadd := rest[0]
	part := toadd.Criterion.(ProductPartition)
	part.Dimension.Value = "agi"
	part.Id = 0
	toadd.Criterion = part
	//toadd.BiddingStrategyConfiguration = nil

	aops := AdGroupCriterionOperations{
		"ADD": AdGroupCriterions{
			toadd,
		},
	}

	res, err := NewAdGroupCriterionService(&config.Auth).Mutate(aops)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestSandboxSharedEntity(t *testing.T) {
	config := getTestConfig()

	campaigns, n, err := NewCampaignService(&config.Auth).Get(Selector{
		Fields: []string{"Id", "Name"},
	})

	fmt.Println(campaigns, n, err)

	campaign := campaigns[0].Id

	res, n, err := NewAdGroupService(&config.Auth).Get(Selector{
		Fields: []string{"Id", "Name"},
		Predicates: []Predicate{
			Predicate{
				Field:    "CampaignId",
				Operator: "EQUALS",
				Values:   []string{strconv.FormatInt(campaign, 10)},
			},
		},
	})

	fmt.Println(res, n, err)

	/*
		err = NewSharedSetService(&config.Auth).Mutate([]SharedSetOperation{
			{"ADD", SharedSet{Name: "sharedset", Type: "NEGATIVE_KEYWORDS"}},
		})

	*/

	sharedsets, n, err := NewSharedSetService(&config.Auth).Get(Selector{
		Fields: []string{"SharedSetId", "Name", "Type"},
	})

	if err != nil {
		t.Error("sharedset: ", err)
	}

	fmt.Println(sharedsets)

	sharedset := sharedsets[0].Id

	err = NewSharedCriterionService(&config.Auth).Mutate([]SharedCriterionOperation{
		{"ADD", SharedCriterion{
			SharedSetId: sharedset,
			Negative:    true,
			Criterion: KeywordCriterion{
				MatchType: "PHRASE",
				Text:      "bbbb",
			},
		}},
	})

	if err != nil {
		t.Error(err)
	}

	err = NewCampaignSharedSetService(&config.Auth).Mutate([]CampaignSharedSetOperation{
		{"REMOVE", CampaignSharedSet{CampaignId: campaign, SharedSetId: sharedset}},
	})

	if err != nil {
		t.Error(err)
	}

	err = NewCampaignSharedSetService(&config.Auth).Mutate([]CampaignSharedSetOperation{
		{"ADD", CampaignSharedSet{CampaignId: campaign, SharedSetId: sharedset}},
	})

	if err != nil {
		t.Error(err)
	}

	sharedsetcrits, _, err := NewSharedCriterionService(&config.Auth).Get(Selector{
		Fields: []string{"SharedSetId", "Negative"},
		Predicates: []Predicate{
			Predicate{
				Field:    "SharedSetId",
				Operator: "EQUALS",
				Values:   []string{strconv.FormatInt(sharedset, 10)},
			},
		},
	})
	if err != nil {
		t.Error(err)
	}

	fmt.Println(sharedsetcrits)

	ss, _, err := NewCampaignSharedSetService(&config.Auth).Get(Selector{
		Fields: []string{"SharedSetId", "CampaignId", "SharedSetName"},
	})

	if err != nil {
		t.Error(err)
	}

	fmt.Println(ss)
}
