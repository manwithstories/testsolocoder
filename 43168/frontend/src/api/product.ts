import request from '@/utils/request'
import type {
  Product,
  ProductCategory,
  ProductFormData,
  ListProductsParams,
  PaginationResult
} from '@/types'

export interface ProductListResult extends PaginationResult<Product> {}

export const listProducts = (params: ListProductsParams) => {
  return request.get<any, ProductListResult>('/products/', { params })
}

export const getProduct = (id: number | string) => {
  return request.get<any, Product>(`/products/${id}`)
}

export const createProduct = (data: ProductFormData) => {
  return request.post<any, Product>('/products/', data)
}

export const updateProduct = (id: number | string, data: ProductFormData) => {
  return request.put<any, Product>(`/products/${id}`, data)
}

export const deleteProduct = (id: number | string) => {
  return request.delete<any, void>(`/products/${id}`)
}

export const listProductCategories = () => {
  return request.get<any, ProductCategory[]>('/products/categories')
}

export const uploadProductImage = (formData: FormData) => {
  return request.post<any, { url: string }>('/products/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export const inquireProduct = (id: number | string, data: { remark?: string }) => {
  return request.post<any, { orderId: number }>(`/products/${id}/inquire`, data)
}
