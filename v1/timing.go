package oakacs

// It's also important to be aware that using subtle.ConstantTimeCompare() can leak information about username and password length. To prevent this, we should hash both the provided and expected username and password values using a fast cryptographic hash function like SHA-256 before comparing them. This ensures that both the provided and expected values that we are comparing are equal in length and prevents subtle.ConstantTimeCompare() itself from returning early.
