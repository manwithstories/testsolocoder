import { request } from './index'
import type { Product, ProductImage, Brand, ApiResponse } from '@/types'

export const productApi = {
  listProducts(params?: {
    page?: number
    page_size?: number
    category?: string
    brand_id?: number
    status?: string
    sort_by?: string
    keyword?: string
  }) {
    return request.get<{ list: Product[]; total: number }>('/products', { params })
  },

  getProduct(id: number) {
    return request.get<Product>(`/products/${id}`)
  },

  createProduct(data: {
    title: string
    description: string
    category: string
    brand_id?: number
    brand_name?: string
    original_price?: number
    price: number
    condition?: string
    color?: string
    size?: string
    material?: string
    stock: number
  }) {
    return request.post<Product>('/products', data)
  },

  updateProduct(id: number, data: Record<string, any>) {
    return request.put<Product>(`/products/${id}`, data)
  },

  deleteProduct(id: number) {
    return request.delete(`/products/${id}`)
  },

  updateProductStatus(id: number, status: string) {
    return request.patch(`/products/${id}/status`, { status })
  },

  uploadImages(id: number, formData: FormData) {
    return request.post<ProductImage[]>(`/products/${id}/images`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  getProductImages(id: number) {
    return request.get<ProductImage[]>(`/products/${id}/images`)
  },

  deleteProductImages(id: number) {
    return request.delete(`/products/${id}/images`)
  },

  listMyProducts(params?: {
    page?: number
    page_size?: number
    status?: string
  }) {
    return request.get<{ list: Product[]; total: number }>('/products/seller/my', { params })
  },

  createBrand(data: {
    name: string
    name_cn?: string
    logo?: string
    country?: string
    description?: string
    category: string
  }) {
    return request.post<Brand>('/products/brands', data)
  },

  listBrands(params?: { category?: string }) {
    return request.get<Brand[]>('/products/brands/list', { params })
  },

  getBrand(id: number) {
    return request.get<Brand>(`/products/brands/${id}`)
  }
}

export default productApi
