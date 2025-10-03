package internal

import "fmt"

func MigrateStores(from, to Store) error {
	keys, err := from.GetAllKeys()
	if err != nil {
		return fmt.Errorf("Error getting keys: %v\n", err)
	}

	for _, key := range keys {
		value, err := from.Get(key)
		if err != nil {
			return fmt.Errorf("Error getting value for %s: %v\n", key, err)
		}

		err = to.Set(key, value)
		if err != nil {

			return fmt.Errorf("Error setting value for %s: %v\n", key, err)
		}
	}

	return nil
}
