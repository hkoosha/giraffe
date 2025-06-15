package setup

func FixateOnceChecks(enabled bool) {
	setBypass.Do(func() {
		bypass = enabled
	})
}

func Once(
	domain string,
	service string,
	what ...string,
) {
	key := toKey(domain, service, what...)

	locks.mu.Lock()
	defer locks.mu.Unlock()

	check(key, unlocked)
	lock(key)
}

func EnsureOpen(
	domain string,
	service string,
	what ...string,
) {
	key := toKey(domain, service, what...)

	locks.mu.Lock()
	defer locks.mu.Unlock()

	check(key, unlocked)
}

func EnsureDone(
	domain string,
	service string,
	what ...string,
) {
	key := toKey(domain, service, what...)

	locks.mu.Lock()
	defer locks.mu.Unlock()

	check(key, locked)
}
