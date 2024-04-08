import Dashboard from '../pages/Dashboard';
import Login from '../pages/Login';
import Registrer from '../pages/Register';
import { createBrowserRouter, RouterProvider as BaseRouterProvider } from 'react-router-dom';


const router = createBrowserRouter([
    {
        path: '/login',
        element: <Login />
    },
    {
        path: '/register',
        element: <Registrer />
    },
    {
        path: '/dashboard',
        element: <Dashboard />
    }
]);

const Loader = () => {
    return (
        <div>Loading...</div>
    )

}

const RouterProvider = () => {
    return (
        <BaseRouterProvider router={router} fallbackElement={<Loader />} />
    )
}

export default RouterProvider