import JoinRoom from '../pages/JoinRoom';
import Room from '../pages/Room';
import Login from '../pages/Login';
import Registrer from '../pages/Register';
import { createBrowserRouter, RouterProvider as BaseRouterProvider } from 'react-router-dom';


export const router = createBrowserRouter([
    {
        path: '/login',
        element: <Login />
    },
    {
        path: '/register',
        element: <Registrer />
    },
    {
        path: '/joinRoom',
        element: <JoinRoom />
    },
    {
        path: '/room/:id',
        element: <Room />
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