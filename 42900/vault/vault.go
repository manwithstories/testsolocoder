package vault

import (
	"database/sql"
	"errors"
	"fmt"
	"passman/crypto"
	"passman/storage"
)

const verifierPlaintext = "passman-vault-verifier-v1"

type Manager struct {
	store       *storage.Store
	currentVault *storage.Vault
	masterKey   []byte
}

func NewManager(store *storage.Store) *Manager {
	return &Manager{
		store: store,
	}
}

func (m *Manager) CreateVault(name, masterPassword string) (*storage.Vault, error) {
	_, err := m.store.GetVaultByName(name)
	if err == nil {
		return nil, fmt.Errorf("vault '%s' already exists", name)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	salt, err := crypto.GenerateSalt()
	if err != nil {
		return nil, err
	}

	key := crypto.DeriveKey(masterPassword, salt)
	encData, err := crypto.Encrypt([]byte(verifierPlaintext), key)
	if err != nil {
		return nil, err
	}
	verifier := encData.Encode()

	vault, err := m.store.CreateVault(name, salt, verifier)
	if err != nil {
		return nil, err
	}

	return vault, nil
}

func (m *Manager) OpenVault(name, masterPassword string) error {
	vault, err := m.store.GetVaultByName(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("vault '%s' not found", name)
		}
		return err
	}

	key := crypto.DeriveKey(masterPassword, vault.Salt)

	encData, err := crypto.DecodeEncodedData(vault.Verifier)
	if err != nil {
		return errors.New("invalid vault data")
	}

	plaintext, err := crypto.Decrypt(encData, key)
	if err != nil {
		return errors.New("invalid master password")
	}

	if string(plaintext) != verifierPlaintext {
		return errors.New("invalid master password")
	}

	m.currentVault = vault
	m.masterKey = key

	if err := m.store.SetConfig("current_vault", name); err != nil {
		return err
	}

	return nil
}

func (m *Manager) CloseVault() {
	m.currentVault = nil
	m.masterKey = nil
}

func (m *Manager) GetCurrentVault() *storage.Vault {
	return m.currentVault
}

func (m *Manager) GetMasterKey() []byte {
	return m.masterKey
}

func (m *Manager) ListVaults() ([]*storage.Vault, error) {
	return m.store.ListVaults()
}

func (m *Manager) DeleteVault(name string) error {
	vault, err := m.store.GetVaultByName(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("vault '%s' not found", name)
		}
		return err
	}

	if m.currentVault != nil && m.currentVault.ID == vault.ID {
		m.CloseVault()
		if err := m.store.DeleteConfig("current_vault"); err != nil {
			return err
		}
	}

	return m.store.DeleteVault(vault.ID)
}

func (m *Manager) ChangeMasterPassword(oldPassword, newPassword string) error {
	if m.currentVault == nil {
		return errors.New("no vault is open")
	}

	oldKey := crypto.DeriveKey(oldPassword, m.currentVault.Salt)

	encData, err := crypto.DecodeEncodedData(m.currentVault.Verifier)
	if err != nil {
		return errors.New("invalid vault data")
	}

	plaintext, err := crypto.Decrypt(encData, oldKey)
	if err != nil {
		return errors.New("invalid old master password")
	}

	if string(plaintext) != verifierPlaintext {
		return errors.New("invalid old master password")
	}

	entries, err := m.store.ListEntries(m.currentVault.ID)
	if err != nil {
		return err
	}

	newSalt, err := crypto.GenerateSalt()
	if err != nil {
		return err
	}

	newKey := crypto.DeriveKey(newPassword, newSalt)

	for _, entry := range entries {
		oldEncData, err := crypto.DecodeEncodedData(entry.Password)
		if err != nil {
			return err
		}

		plaintext, err := crypto.Decrypt(oldEncData, oldKey)
		if err != nil {
			return err
		}

		newEncData, err := crypto.Encrypt(plaintext, newKey)
		if err != nil {
			return err
		}

		entry.Password = newEncData.Encode()
		if err := m.store.UpdateEntry(entry); err != nil {
			return err
		}
	}

	newEncVerifier, err := crypto.Encrypt([]byte(verifierPlaintext), newKey)
	if err != nil {
		return err
	}
	newVerifier := newEncVerifier.Encode()

	if err := m.store.UpdateVaultSaltAndVerifier(m.currentVault.ID, newSalt, newVerifier); err != nil {
		return err
	}

	m.currentVault.Salt = newSalt
	m.currentVault.Verifier = newVerifier
	m.masterKey = newKey

	return nil
}

func (m *Manager) GetLastActiveVault() (string, error) {
	name, err := m.store.GetConfig("current_vault")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return name, nil
}

func (m *Manager) EncryptPassword(password string) (string, error) {
	if m.masterKey == nil {
		return "", errors.New("no vault is open")
	}

	encData, err := crypto.Encrypt([]byte(password), m.masterKey)
	if err != nil {
		return "", err
	}

	return encData.Encode(), nil
}

func (m *Manager) DecryptPassword(encodedPassword string) (string, error) {
	if m.masterKey == nil {
		return "", errors.New("no vault is open")
	}

	encData, err := crypto.DecodeEncodedData(encodedPassword)
	if err != nil {
		return "", err
	}

	plaintext, err := crypto.Decrypt(encData, m.masterKey)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
