import request from '@/utils/request'
import type { ApiResponse, Product, ProductRecord, ProductSale, PageResult, PageParams } from '@/types'

export const getProducts = (params: PageParams & { category?: string; low_stock?: boolean }) => {
  return request.get<ApiResponse<PageResult<Product>>>('/products', { params })
}

export const getAllProducts = () => {
  return request.get<ApiResponse<Product[]>>('/products/all')
}

export const getProduct = (id: number) => {
  return request.get<ApiResponse<Product>>(`/products/${id}`)
}

export const createProduct = (data: any) => {
  return request.post<ApiResponse<Product>>('/products', data)
}

export const updateProduct = (id: number, data: any) => {
  return request.put<ApiResponse<Product>>(`/products/${id}`, data)
}

export const deleteProduct = (id: number) => {
  return request.delete<ApiResponse<null>>(`/products/${id}`)
}

export const addStock = (data: { product_id: number; quantity: number; operator_id?: number; remark?: string }) => {
  return request.post<ApiResponse<null>>('/products/add-stock', data)
}

export const deductStock = (data: { product_id: number; quantity: number; appointment_id?: number; operator_id?: number; remark?: string }) => {
  return request.post<ApiResponse<null>>('/products/deduct-stock', data)
}

export const getProductRecords = (params: PageParams & { product_id?: number; change_type?: string }) => {
  return request.get<ApiResponse<PageResult<ProductRecord>>>('/products/records/list', { params })
}

export const getLowStockProducts = () => {
  return request.get<ApiResponse<Product[]>>('/products/low-stock')
}

export const saleProduct = (data: { customer_id: number; product_id: number; quantity: number; pay_method: string; operator_id?: number }) => {
  return request.post<ApiResponse<ProductSale>>('/products/sale', data)
}

export const getProductSales = (params: PageParams & { customer_id?: number }) => {
  return request.get<ApiResponse<PageResult<ProductSale>>>('/products/sales/list', { params })
}

export const stockTake = (data: Record<number, number>) => {
  return request.post<ApiResponse<null>>('/products/stock-take', data)
}
