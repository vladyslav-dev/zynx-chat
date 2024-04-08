import { loginUser, registerUser } from "../api";
import { useState } from "react";
import { useNavigate } from "react-router-dom"

const RegisterForm = () => {
    const navigate = useNavigate()
    const [error, setError] = useState<string | null>(null)

    const onSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        const formData = new FormData(e.target as HTMLFormElement);

        const username = formData.get('username') as string
        const email = formData.get('email') as string
        const password = formData.get('password') as string

        const response = await registerUser({ username, email, password })

        if (response?.ok) {
            const response = await loginUser({ email, password });

            if (response?.ok) {
                navigate("/dashboard")
            } else {
                setError("Error occured while logging in")
            }
        } else {
            setError("Error occured while registering user")
        }
    }

    return (
        <div>
            <h1>Register</h1>
            <form onSubmit={onSubmit}>
                <div>
                    <label htmlFor="username">Username</label>
                    <input type="username" id="username" />
                </div>
                <div>
                    <label htmlFor="email">Email</label>
                    <input type="email" id="email" />
                </div>
                <div>
                    <label htmlFor="password">Password</label>
                    <input type="password" id="password" />
                </div>
                <button type="submit">Create account</button>
            </form>
            {error && <div>{error}</div>}
        </div>
    )
}

export default RegisterForm