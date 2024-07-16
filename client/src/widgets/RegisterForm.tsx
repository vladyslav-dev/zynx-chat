import { loginUser, registerUser } from "../api";
import { useState } from "react";
import { useNavigate } from "react-router-dom"

const RegisterForm = () => {
    const navigate = useNavigate()
    const [error, setError] = useState<string | null>(null)

    const onSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        const formData = new FormData(e.target as HTMLFormElement);
        const { username, email, password } = Object.fromEntries(formData.entries())

        const response = await registerUser({ username, email, password })

        if ("id" in response) {
            localStorage.setItem("user", JSON.stringify(response))

            const loginResponse = await loginUser({ email, password });

            if (loginResponse?.ok) {
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
                    <input type="username" id="username" name="username" />
                </div>
                <div>
                    <label htmlFor="email">Email</label>
                    <input type="email" id="email" name="email" />
                </div>
                <div>
                    <label htmlFor="password">Password</label>
                    <input type="password" id="password" name="password" />
                </div>
                <button type="submit">Create account</button>
            </form>
            {error && <div>{error}</div>}
        </div>
    )
}

export default RegisterForm