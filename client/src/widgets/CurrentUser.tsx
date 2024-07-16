import { AuthContext } from "../providers/AuthProvider";
import { useContext } from "react";

const CurrentUser = () => {
    const { user } = useContext(AuthContext);

    return (
        <section className="mb-4">
            <h3 className="mb-2">Welcome {user?.username} 👋</h3>
            <div style={{ color: "#6dff69", border: "1px solid #6dff69", padding: "8px", margin: "8px 0"}}>
                <div>ID: {user?.id}</div>
                <div>Username: {user?.username}</div>
                <div>Email: {user?.email}</div>
            </div>
        </section>
    );
}

export default CurrentUser;