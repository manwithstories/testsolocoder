package entry

import (
	"errors"
	"passman/storage"
	"passman/vault"
	"strings"
)

type Manager struct {
	vaultManager *vault.Manager
	store        *storage.Store
}

func NewManager(vaultManager *vault.Manager, store *storage.Store) *Manager {
	return &Manager{
		vaultManager: vaultManager,
		store:        store,
	}
}

func (m *Manager) Add(name, username, password, url, tags string) (*storage.Entry, error) {
	currentVault := m.vaultManager.GetCurrentVault()
	if currentVault == nil {
		return nil, errors.New("no vault is open")
	}

	encryptedPassword, err := m.vaultManager.EncryptPassword(password)
	if err != nil {
		return nil, err
	}

	entry := &storage.Entry{
		VaultID:  currentVault.ID,
		Name:     name,
		Username: username,
		Password: encryptedPassword,
		URL:      url,
		Tags:     tags,
	}

	return m.store.CreateEntry(entry)
}

func (m *Manager) Get(id int64) (*storage.Entry, string, error) {
	entry, err := m.store.GetEntry(id)
	if err != nil {
		return nil, "", err
	}

	decryptedPassword, err := m.vaultManager.DecryptPassword(entry.Password)
	if err != nil {
		return nil, "", err
	}

	return entry, decryptedPassword, nil
}

func (m *Manager) List() ([]*storage.Entry, error) {
	currentVault := m.vaultManager.GetCurrentVault()
	if currentVault == nil {
		return nil, errors.New("no vault is open")
	}

	return m.store.ListEntries(currentVault.ID)
}

func (m *Manager) Search(nameQuery, tagQuery string) ([]*storage.Entry, error) {
	currentVault := m.vaultManager.GetCurrentVault()
	if currentVault == nil {
		return nil, errors.New("no vault is open")
	}

	return m.store.SearchEntries(currentVault.ID, nameQuery, tagQuery)
}

func (m *Manager) Update(id int64, name, username, password, url, tags string) error {
	entry, err := m.store.GetEntry(id)
	if err != nil {
		return err
	}

	if name != "" {
		entry.Name = name
	}
	if username != "" {
		entry.Username = username
	}
	if password != "" {
		encryptedPassword, err := m.vaultManager.EncryptPassword(password)
		if err != nil {
			return err
		}
		entry.Password = encryptedPassword
	}
	if url != "" {
		entry.URL = url
	}
	if tags != "" {
		entry.Tags = tags
	}

	return m.store.UpdateEntry(entry)
}

func (m *Manager) Delete(id int64) error {
	return m.store.DeleteEntry(id)
}

func (m *Manager) DecryptEntryPassword(entry *storage.Entry) (string, error) {
	return m.vaultManager.DecryptPassword(entry.Password)
}

func (m *Manager) GetAllTags() ([]string, error) {
	currentVault := m.vaultManager.GetCurrentVault()
	if currentVault == nil {
		return nil, errors.New("no vault is open")
	}

	entries, err := m.store.ListEntries(currentVault.ID)
	if err != nil {
		return nil, err
	}

	tagMap := make(map[string]bool)
	for _, entry := range entries {
		if entry.Tags == "" {
			continue
		}
		tags := strings.Split(entry.Tags, ",")
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tagMap[tag] = true
			}
		}
	}

	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}

	return tags, nil
}
