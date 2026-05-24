import { createSlice, PayloadAction, createAsyncThunk } from '@reduxjs/toolkit'
import { jobApi, Job } from '@/api/jobs'
import { PaginatedData } from '@/api/auth'

interface JobState {
  jobs: Job[]
  currentJob: Job | null
  loading: boolean
  error: string | null
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

const initialState: JobState = {
  jobs: [],
  currentJob: null,
  loading: false,
  error: null,
  pagination: {
    page: 1,
    page_size: 10,
    total: 0,
    total_pages: 0,
  },
}

export const searchJobs = createAsyncThunk(
  'jobs/search',
  async (params?: Record<string, string>) => {
    const response = await jobApi.search(params)
    return response.data
  }
)

export const getJobById = createAsyncThunk(
  'jobs/getById',
  async (id: number) => {
    const response = await jobApi.getById(id)
    return response.data
  }
)

const jobSlice = createSlice({
  name: 'jobs',
  initialState,
  reducers: {
    setJobs: (state, action: PayloadAction<Job[]>) => {
      state.jobs = action.payload
    },
    setCurrentJob: (state, action: PayloadAction<Job | null>) => {
      state.currentJob = action.payload
    },
    clearJobs: (state) => {
      state.jobs = []
      state.pagination = initialState.pagination
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(searchJobs.pending, (state) => {
        state.loading = true
      })
      .addCase(searchJobs.fulfilled, (state, action: PayloadAction<PaginatedData<Job>>) => {
        state.loading = false
        state.jobs = action.payload.data
        state.pagination = action.payload.pagination
      })
      .addCase(searchJobs.rejected, (state) => {
        state.loading = false
      })
      .addCase(getJobById.pending, (state) => {
        state.loading = true
      })
      .addCase(getJobById.fulfilled, (state, action) => {
        state.loading = false
        state.currentJob = action.payload
      })
      .addCase(getJobById.rejected, (state) => {
        state.loading = false
      })
  },
})

export const { setJobs, setCurrentJob, clearJobs } = jobSlice.actions

export const selectJobs = (state: { jobs: JobState }) => state.jobs

export default jobSlice.reducer
