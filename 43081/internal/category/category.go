package category

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"finance-tracker/internal/logger"
	"finance-tracker/internal/models"
	"finance-tracker/internal/storage"
)

func Create(name, categoryType string, parentID *int64) (*models.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("category name cannot be empty")
	}

	categoryType = strings.ToLower(strings.TrimSpace(categoryType))
	if categoryType != "income" && categoryType != "expense" {
		return nil, fmt.Errorf("invalid category type: %s, must be 'income' or 'expense'", categoryType)
	}

	if parentID != nil {
		parent, err := GetByID(*parentID)
		if err != nil {
			return nil, fmt.Errorf("parent category not found: %w", err)
		}
		if parent.Type != categoryType {
			return nil, fmt.Errorf("child category type must match parent type")
		}
	}

	now := time.Now()
	var parentVal interface{}
	if parentID != nil {
		parentVal = *parentID
	} else {
		parentVal = nil
	}

	result, err := storage.DB.Exec(`
		INSERT INTO categories (name, type, parent_id, created_at)
		VALUES (?, ?, ?, ?)
	`, name, categoryType, parentVal, now)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("category '%s' already exists", name)
		}
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	id, _ := result.LastInsertId()
	category := &models.Category{
		ID:        id,
		Name:      name,
		Type:      categoryType,
		ParentID:  parentID,
		CreatedAt: now,
	}

	logger.Info("Category created: %s (ID: %d, type: %s)", name, id, categoryType)
	return category, nil
}

func GetByID(id int64) (*models.Category, error) {
	var category models.Category
	var parentID sql.NullInt64

	err := storage.DB.QueryRow(`
		SELECT id, name, type, parent_id, created_at
		FROM categories WHERE id = ?
	`, id).Scan(&category.ID, &category.Name, &category.Type, &parentID, &category.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if parentID.Valid {
		category.ParentID = &parentID.Int64
	}

	return &category, nil
}

func GetByNameAndType(name, categoryType string) (*models.Category, error) {
	var category models.Category
	var parentID sql.NullInt64

	err := storage.DB.QueryRow(`
		SELECT id, name, type, parent_id, created_at
		FROM categories WHERE name = ? AND type = ? AND parent_id IS NULL
	`, name, categoryType).Scan(&category.ID, &category.Name, &category.Type, &parentID, &category.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

func List(categoryType string) ([]models.Category, error) {
	query := `SELECT id, name, type, parent_id, created_at FROM categories`
	args := []interface{}{}

	if categoryType != "" {
		query += ` WHERE type = ?`
		args = append(args, categoryType)
	}
	query += ` ORDER BY parent_id IS NULL DESC, name`

	rows, err := storage.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		var parentID sql.NullInt64
		err := rows.Scan(&category.ID, &category.Name, &category.Type, &parentID, &category.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		if parentID.Valid {
			category.ParentID = &parentID.Int64
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func ListTree(categoryType string) ([]models.Category, error) {
	all, err := List(categoryType)
	if err != nil {
		return nil, err
	}

	type node struct {
		cat      *models.Category
		children []*node
	}

	nodeMap := make(map[int64]*node)
	for i := range all {
		nodeMap[all[i].ID] = &node{cat: &all[i]}
	}

	var rootNodes []*node
	for i := range all {
		cat := &all[i]
		n := nodeMap[cat.ID]
		if cat.ParentID == nil {
			rootNodes = append(rootNodes, n)
		} else {
			parent := nodeMap[*cat.ParentID]
			if parent != nil {
				parent.children = append(parent.children, n)
			}
		}
	}

	var buildTree func(n *node) models.Category
	buildTree = func(n *node) models.Category {
		result := *n.cat
		result.Children = nil
		for _, child := range n.children {
			result.Children = append(result.Children, buildTree(child))
		}
		return result
	}

	var roots []models.Category
	for _, n := range rootNodes {
		roots = append(roots, buildTree(n))
	}

	return roots, nil
}

func Delete(id int64) error {
	category, err := GetByID(id)
	if err != nil {
		return err
	}

	var childCount int
	err = storage.DB.QueryRow(`SELECT COUNT(*) FROM categories WHERE parent_id = ?`, id).Scan(&childCount)
	if err != nil {
		return fmt.Errorf("failed to check child categories: %w", err)
	}
	if childCount > 0 {
		return fmt.Errorf("cannot delete category '%s' with %d child categories", category.Name, childCount)
	}

	var transCount int
	err = storage.DB.QueryRow(`SELECT COUNT(*) FROM transactions WHERE category_id = ?`, id).Scan(&transCount)
	if err != nil {
		return fmt.Errorf("failed to check category transactions: %w", err)
	}
	if transCount > 0 {
		return fmt.Errorf("cannot delete category '%s' with %d transactions", category.Name, transCount)
	}

	_, err = storage.DB.Exec(`DELETE FROM categories WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	logger.Info("Category deleted: %s (ID: %d)", category.Name, id)
	return nil
}

func Exists(id int64, categoryType string) bool {
	var count int
	err := storage.DB.QueryRow(`SELECT COUNT(*) FROM categories WHERE id = ? AND type = ?`, id, categoryType).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}
