package internal

import "fmt"

func MigrateStores(from Store, to Store) error {

	fromKeys, err := from.GetAllKeys()
	if err != nil {
		return err
	}
	if len(fromKeys) == 0 {
		fmt.Println("ℹ️  No data to migrate")
		return nil
	}

	for _, key := range fromKeys {
		v, err := from.Get(key)
		if err != nil {
			return err
		}
		err = to.Set(key, v)
		if err != nil {
			return err
		}
	}

	return nil
}
