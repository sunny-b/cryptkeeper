package commands

// 	envBytes, err := base64.StdEncoding.DecodeString(diff)
// 	if err != nil {
// 		return nil, err
// 	}

// 	cfg := make(config.Env)
// 	err = json.Unmarshal(envBytes, &cfg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return cfg, nil
// }

// func calculateDiff(old, new config.Env) map[string]string {
// 	diff := make(map[string]string)

// 	for k, v := range old {
// 		var prefix string
// 		newV, ok := new[k]
// 		if !ok {
// 			prefix = "-"
// 		} else if v != newV {
// 			prefix = "~"
// 		} else {
// 			prefix = "="
// 		}

// 		diff[prefix+k] = newV
// 	}

// 	for k, v := range new {
// 		var prefix string
// 		_, ok := old[k]
// 		if !ok {
// 			prefix = "+"
// 		}

// 		if prefix != "" {
// 			diff[prefix+k] = v
// 		}
// 	}

// 	return diff
// }

// func encodeEnv(env config.Env) (string, error) {
// 	b, err := json.Marshal(env)
// 	if err != nil {
// 		return "", err
// 	}

// 	return base64.StdEncoding.EncodeToString(b), nil
// }

// func exportAll(env config.Env) map[string]string {
// 	diff := make(map[string]string)

// 	for k, v := range env {
// 		diff["+"+k] = v
// 	}

// 	return diff
// }

// func unsetAll(env config.Env) map[string]string {
// 	diff := make(map[string]string)

// 	for k, v := range env {
// 		diff["-"+k] = v
// 	}

// 	return diff
// }

// func fetchDiff(export bool, diffEnvKey string, currentEnv config.Env) config.Env {

// 	diff := make(map[string]string)

// 	if !envVarExists(diffEnvKey) {
// 		if !export {
// 			return nil
// 		}

// 		diff = exportAll(currentEnv)
// 	} else if !export {
// 		diff = unsetAll(currentEnv)
// 	} else {
// 		oldEnv, err := envFromDiff(os.Getenv(diffEnvKey))
// 		if err != nil {
// 			diff = exportAll(currentEnv)
// 		}

// 		if oldEnv != nil {
// 			log.WithField("oldEnv", oldEnv).Debug("oldEnv")
// 			if reflect.DeepEqual(currentEnv, oldEnv) {
// 				return nil
// 			}

// 			diff = calculateDiff(oldEnv, currentEnv)
// 		}
// 	}

// 	return diff
// }
