import { loginUser } from "../api";
import { useState } from "react";
import { useNavigate } from "react-router-dom"

const LoginForm = () => {
    const navigate  = useNavigate()
    const [error, setError] = useState<string | null>(null)

    const onSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        const formData = new FormData(e.target as HTMLFormElement);

        const email = formData.get('email') as string
        const password = formData.get('password') as string

        const response = await loginUser({ email, password })

        if (response?.ok) {
            navigate("/dashboard")
        } else {
            setError("Error occured while logging in")
        }
    }

    return (
        <div>
            <h1>Login</h1>
            <form onSubmit={onSubmit}>
                <div>
                    <label htmlFor="email">Email</label>
                    <input type="email" id="email" />
                </div>
                <div>
                    <label htmlFor="password">Password</label>
                    <input type="password" id="password" />
                </div>
                <button type="submit">Login</button>
            </form>
            {error && <div>{error}</div>}
        </div>
    )
}

export default LoginForm