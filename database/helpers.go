package database

/*
func getPolicies(hook string, token model.Token) ([]string, error) {
	policies := []string{}
	rawPolicies := map[string]interface{}{}
	if err := json.Unmarshal(token.Policies, &rawPolicies); err != nil {
		return policies, err
	}
	val, ok := rawPolicies[hook]
	if ok {
		for _, v := range val.([]interface{}) {
			policies = append(policies, v.(string))
		}
	}
	return policies, nil
}

func getCapabilities(hook string, rule model.Rule) ([]string, error) {
	caps := []string{}
	rawCap := map[string]interface{}{}
	if err := json.Unmarshal(rule.Capabilities, &rawCap); err != nil {
		return caps, err
	}
	val, ok := rawCap[hook]
	if ok {
		for _, v := range val.([]interface{}) {
			caps = append(caps, v.(string))
		}
	}
	return caps, nil
}

func bindPoliciesAndRules(policiesr []model.Policie, rulesr []model.Rule) ([]model.PolicieResponse, error) {

	policier := []model.PolicieResponse{}
	pop := map[string]*model.PolicieResponse{}

	for _, vv := range policiesr {
		// copy policies struct for policies response
		pol := model.PolicieResponse{}
		_ = copier.Copy(&pol, &vv)
		pop[vv.Name] = &pol

		// check rules
		for _, v := range rulesr {
			if v.PoliceId == vv.Id {
				rule := model.RuleResponse{}
				caps, err := getCapabilities("cap", v)
				if err != nil {
					return []model.PolicieResponse{}, err
				}
				// copy rules struct for rules reponse
				if err := copier.Copy(&rule, &v); err != nil {
					return []model.PolicieResponse{}, err
				}
				rule.Capabilities = caps
				pop[vv.Name].Rules = append(pop[vv.Name].Rules, rule)
			}
		}
	}

	// back to the slice
	for _, v := range pop {
		policier = append(policier, *v)
	}

	return policier, nil
}

func bindTokenPolicies(token model.Token) (model.TokenResponse, error) {

	tokenc := model.TokenResponse{}

	tokenPolicies, err := getPolicies("policies", token)
	if err != nil {
		return model.TokenResponse{}, err
	}
	if err := copier.Copy(&tokenc, &token); err != nil {
		return model.TokenResponse{}, err
	}
	tokenc.Policies = tokenPolicies

	return tokenc, nil
}
*/
