package setup

func FixateOnceChecks(enabled bool) {
	setBypass.Do(func() {
		bypass = enabled
	})
}

func Once(
	domain string,
	what ...string,
) {
	for i := len(what) - 1; i >= 0; i-- {
		EnsureOpen(domain, what[:i]...)
	}

	key := toKey(domain, what...)

	locks.mu.Lock()
	defer locks.mu.Unlock()

	check(key, unlocked)
	lock(key)
}

func EnsureOpen(
	domain string,
	what ...string,
) {
	key := toKey(domain, what...)

	locks.mu.Lock()
	defer locks.mu.Unlock()

	check(key, unlocked)
}

func EnsureDone(
	domain string,
	what ...string,
) {
	key := toKey(domain, what...)

	locks.mu.Lock()
	defer locks.mu.Unlock()

	check(key, locked)
}
