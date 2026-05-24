package repository

import (
	"furniture-platform/internal/model"
)

// BatchGetOptions 批量获取多个产品的选项
func (r *ProductRepository) BatchGetOptions(productIDs []uint) (map[uint][]*model.ProductOption, error) {
	result := make(map[uint][]*model.ProductOption)
	if len(productIDs) == 0 {
		return result, nil
	}
	var options []*model.ProductOption
	if err := r.db.Where("product_id IN ?", productIDs).
		Order("product_id ASC, sort ASC, id ASC").
		Find(&options).Error; err != nil {
		return nil, err
	}
	for _, o := range options {
		result[o.ProductID] = append(result[o.ProductID], o)
	}
	return result, nil
}

// BatchGetImages 批量获取多个产品的图片
func (r *ProductRepository) BatchGetImages(productIDs []uint) (map[uint][]*model.ProductImage, error) {
	result := make(map[uint][]*model.ProductImage)
	if len(productIDs) == 0 {
		return result, nil
	}
	var images []*model.ProductImage
	if err := r.db.Where("product_id IN ?", productIDs).
		Order("product_id ASC, sort ASC, id ASC").
		Find(&images).Error; err != nil {
		return nil, err
	}
	for _, img := range images {
		result[img.ProductID] = append(result[img.ProductID], img)
	}
	return result, nil
}
