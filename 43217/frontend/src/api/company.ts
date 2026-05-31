import { request } from '@/utils/request'
import type { Company, Department, Employee, CompanyBudget, DepartmentAppointment, PaginationParams, PaginationResponse } from '@/types'

export const registerCompany = (data: any): Promise<{ company: Company; hr_user: any }> => {
  return request.post('/companies/register', data)
}

export const getCompany = (id: number): Promise<Company> => {
  return request.get(`/companies/${id}`)
}

export const updateCompany = (data: any): Promise<void> => {
  return request.put('/companies', data)
}

export const getDepartments = (): Promise<Department[]> => {
  return request.get('/departments')
}

export const addDepartment = (data: any): Promise<Department> => {
  return request.post('/departments', data)
}

export const updateDepartment = (id: number, data: any): Promise<void> => {
  return request.put(`/departments/${id}`, data)
}

export const getEmployees = (params: PaginationParams): Promise<PaginationResponse<Employee>> => {
  return request.get('/employees', { params })
}

export const getEmployee = (id: number): Promise<Employee> => {
  return request.get(`/employees/${id}`)
}

export const addEmployee = (data: any): Promise<Employee> => {
  return request.post('/employees', data)
}

export const updateEmployee = (id: number, data: any): Promise<void> => {
  return request.put(`/employees/${id}`, data)
}

export const setBudget = (data: any): Promise<void> => {
  return request.post('/budgets', data)
}

export const getBudget = (year: number): Promise<CompanyBudget> => {
  return request.get('/budgets', { params: { year } })
}

export const setDepartmentAppointment = (data: any): Promise<void> => {
  return request.post('/department-appointments', data)
}

export const getDepartmentAppointments = (departmentId: number, year: number): Promise<DepartmentAppointment[]> => {
  return request.get(`/departments/${departmentId}/appointments`, { params: { year } })
}
