import Login from '../../pages/auth/login';
import Register from '../../pages/auth/register';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import AuthProvider from './AuthProvider';
import MainLayout from '../layouts/MainLayout';
import Settings from '../../pages/settings';
import Conversation from '../../pages/conversation';
import Profile from '../../pages/profile';

const RouterProvider = () => {
    return (
        <BrowserRouter>
            <AuthProvider>
                <Routes>
                    <Route path="/login" element={<Login />} />
                    <Route path="/register" element={<Register />} />
                    <Route
                        path="/*"
                        element={
                            <MainLayout>
                                <Routes>
                                    <Route path="/" element={<Conversation />} />
                                    <Route path="settings" element={<Settings />} />
                                    <Route path="profile" element={<Profile />} />
                                    <Route path="*" element={<Navigate to="/" replace />} />
                                </Routes>
                            </MainLayout>
                        }
                    />
                </Routes>
            </AuthProvider>
        </BrowserRouter>
    )
}

export default RouterProvider