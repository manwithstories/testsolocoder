import { Routes, Route, Navigate } from 'react-router-dom';
import { Layout } from './components/Layout';
import { LoginPage } from './pages/Auth/LoginPage';
import { RegisterPage } from './pages/Auth/RegisterPage';
import { DashboardPage } from './pages/Dashboard/DashboardPage';
import { JobListPage } from './pages/Jobs/JobListPage';
import { JobDetailPage } from './pages/Jobs/JobDetailPage';
import { CreateJobPage } from './pages/Jobs/CreateJobPage';
import { MyApplicationsPage } from './pages/Jobs/MyApplicationsPage';
import { ScheduleListPage } from './pages/Schedules/ScheduleListPage';
import { ScheduleBoardPage } from './pages/Schedules/ScheduleBoardPage';
import { CheckInPage } from './pages/CheckIns/CheckInPage';
import { CheckInRecordsPage } from './pages/CheckIns/CheckInRecordsPage';
import { SalaryListPage } from './pages/Salaries/SalaryListPage';
import { SalaryDetailPage } from './pages/Salaries/SalaryDetailPage';
import { EvaluationListPage } from './pages/Evaluations/EvaluationListPage';
import { MyEvaluationsPage } from './pages/Evaluations/MyEvaluationsPage';
import { StatsPage } from './pages/Stats/StatsPage';
import { ProfilePage } from './pages/Profile/ProfilePage';
import { ProtectedRoute } from './components/ProtectedRoute';

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      
      <Route path="/" element={<ProtectedRoute><Layout /></ProtectedRoute>}>
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<DashboardPage />} />
        
        <Route path="jobs">
          <Route index element={<JobListPage />} />
          <Route path="create" element={<CreateJobPage />} />
          <Route path=":id" element={<JobDetailPage />} />
          <Route path="my-applications" element={<MyApplicationsPage />} />
        </Route>
        
        <Route path="schedules">
          <Route index element={<ScheduleListPage />} />
          <Route path="board" element={<ScheduleBoardPage />} />
        </Route>
        
        <Route path="checkins">
          <Route index element={<CheckInPage />} />
          <Route path="records" element={<CheckInRecordsPage />} />
        </Route>
        
        <Route path="salaries">
          <Route index element={<SalaryListPage />} />
          <Route path=":id" element={<SalaryDetailPage />} />
        </Route>
        
        <Route path="evaluations">
          <Route index element={<EvaluationListPage />} />
          <Route path="mine" element={<MyEvaluationsPage />} />
        </Route>
        
        <Route path="stats" element={<StatsPage />} />
        <Route path="profile" element={<ProfilePage />} />
      </Route>
      
      <Route path="*" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  );
}

export default App;
