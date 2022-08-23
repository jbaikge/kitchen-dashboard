package dashboard

import "fmt"

type Lock struct {
	Title  string `json:"title"`
	Key    string `json:"key"`
	Locked bool   `json:"locked"`
}

type LockState struct {
	State string
}

func GetLocks(config Config) (locks []Lock, err error) {
	ha := NewHomeAssistant(config)

	locks = make([]Lock, len(config.Locks))
	for i, lockConfig := range config.Locks {
		key := fmt.Sprintf("lock.%s", lockConfig.Key)
		state := new(LockState)
		if err = ha.GetState(key, state); err != nil {
			err = fmt.Errorf("problem getting lock state: %w", err)
			return
		}

		locks[i].Title = lockConfig.Title
		locks[i].Key = lockConfig.Key
		locks[i].Locked = state.State == "locked"
	}
	return
}
